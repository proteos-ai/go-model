package agentmodel

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// StructuredText is the text shape a decision provider reads: a JSON string, a JSON
// object, or a JSON array — never null, a number, or a bool. It carries the state
// under judgment, a question's instructions, every criteria description, and the
// legend entries echoed back on a score. Structure is passed through verbatim so
// callers can write structured guidance (named fields such as what / not_for /
// examples, backtick paths into the state like `ticket.messages[0].text`).
//
// A defined type over json.RawMessage does not inherit its methods, so both JSON
// methods are re-implemented here: without them the value would marshal as base64.
type StructuredText json.RawMessage

// MarshalJSON emits the raw bytes untouched (an empty value marshals as null).
func (text StructuredText) MarshalJSON() ([]byte, error) {
	if len(text) == 0 {
		return []byte("null"), nil
	}
	return []byte(text), nil
}

// UnmarshalJSON accepts a string, object, or array and rejects everything else.
// encoding/json calls this for a literal null too, so null is rejected here rather
// than silently becoming an empty value.
func (text *StructuredText) UnmarshalJSON(data []byte) error {
	trimmed := bytes.TrimSpace(data)
	if len(trimmed) == 0 {
		return errors.New("structured text must not be empty")
	}
	switch trimmed[0] {
	case '"', '{', '[':
		*text = StructuredText(append((*text)[:0], trimmed...))
		return nil
	default:
		return fmt.Errorf("structured text must be a JSON string, object, or array, got %s", describeJsonValue(trimmed))
	}
}

func describeJsonValue(raw []byte) string {
	switch {
	case bytes.Equal(raw, []byte("null")):
		return "null"
	case bytes.Equal(raw, []byte("true")), bytes.Equal(raw, []byte("false")):
		return "a boolean"
	default:
		return "a number"
	}
}

// DecisionQuestionType discriminates a DecisionQuestion and its DecisionAnswer.
type DecisionQuestionType string

const (
	// DecisionQuestionTypeBoolean asks whether a statement about the state holds
	// and is answered with the calibrated probability that it does. TypeSafe calls
	// this primitive "noul" (short for Bernoulli — the answer is the Bernoulli
	// parameter p, not a bool); the domain keeps the plain word and confines the
	// vendor's to its adapter.
	DecisionQuestionTypeBoolean DecisionQuestionType = "boolean"
	// DecisionQuestionTypeChoice picks one option out of an unordered set.
	DecisionQuestionTypeChoice DecisionQuestionType = "choice"
	// DecisionQuestionTypeScore places the state on an ordered scale of levels.
	DecisionQuestionTypeScore DecisionQuestionType = "score"
)

// DecisionQuestion is one atomic question over the shared state — a tagged union
// decoded by its type (precedent: ToolBinding). Concrete variants: BooleanQuestion,
// ChoiceQuestion, ScoreQuestion. Each carries its own Type so the default marshaller
// emits the discriminator.
type DecisionQuestion interface {
	QuestionType() DecisionQuestionType
}

// BooleanQuestion: "is this statement true?" Criteria optionally spells out what
// makes the statement hold (true) or not (false) when the boundary is subtle.
type BooleanQuestion struct {
	Type         DecisionQuestionType `json:"type"`
	Instructions StructuredText       `json:"instructions"`
	Criteria     *BooleanCriteria     `json:"criteria,omitempty"`
}

// BooleanCriteria describes the two outcomes of a BooleanQuestion.
type BooleanCriteria struct {
	True  StructuredText `json:"true"`
	False StructuredText `json:"false"`
}

func (BooleanQuestion) QuestionType() DecisionQuestionType { return DecisionQuestionTypeBoolean }

// ChoiceQuestion: "which one of these?" Criteria maps each option name to its
// description; options are unordered. Option-count limits are the provider's.
type ChoiceQuestion struct {
	Type         DecisionQuestionType      `json:"type"`
	Instructions StructuredText            `json:"instructions"`
	Criteria     map[string]StructuredText `json:"criteria"`
}

func (ChoiceQuestion) QuestionType() DecisionQuestionType { return DecisionQuestionTypeChoice }

// ScoreQuestion: "where on this scale?" Criteria is the ordered list of level
// descriptions; the level number is the position (0-based). Every level must be
// described — a null level fails to decode. Level-count limits are the provider's.
type ScoreQuestion struct {
	Type         DecisionQuestionType `json:"type"`
	Instructions StructuredText       `json:"instructions"`
	Criteria     []StructuredText     `json:"criteria"`
}

func (ScoreQuestion) QuestionType() DecisionQuestionType { return DecisionQuestionTypeScore }

// NewBooleanQuestion builds a boolean question with its discriminator set.
func NewBooleanQuestion(instructions StructuredText, criteria *BooleanCriteria) BooleanQuestion {
	return BooleanQuestion{Type: DecisionQuestionTypeBoolean, Instructions: instructions, Criteria: criteria}
}

// NewChoiceQuestion builds a choice question with its discriminator set.
func NewChoiceQuestion(instructions StructuredText, criteria map[string]StructuredText) ChoiceQuestion {
	return ChoiceQuestion{Type: DecisionQuestionTypeChoice, Instructions: instructions, Criteria: criteria}
}

// NewScoreQuestion builds a score question with its discriminator set.
func NewScoreQuestion(instructions StructuredText, criteria []StructuredText) ScoreQuestion {
	return ScoreQuestion{Type: DecisionQuestionTypeScore, Instructions: instructions, Criteria: criteria}
}

// DecisionQuestions maps a caller-chosen key to a question; answers come back under
// the same keys. Decodes each value by its type.
type DecisionQuestions map[string]DecisionQuestion

// UnmarshalJSON decodes every question through DecodeDecisionQuestion, naming the
// offending key on failure.
func (questions *DecisionQuestions) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	if raw == nil {
		*questions = nil
		return nil
	}
	decoded := make(DecisionQuestions, len(raw))
	for key, value := range raw {
		question, err := DecodeDecisionQuestion(value)
		if err != nil {
			return fmt.Errorf("questions[%q]: %w", key, err)
		}
		decoded[key] = question
	}
	*questions = decoded
	return nil
}

// DecodeDecisionQuestion decodes one raw question by its "type" into the matching
// variant (mirrors DecodeToolBinding). Unknown fields are rejected so a misspelt
// field or a stray criteria key surfaces instead of being silently dropped.
func DecodeDecisionQuestion(raw json.RawMessage) (DecisionQuestion, error) {
	questionType, err := peekDecisionType(raw)
	if err != nil {
		return nil, err
	}
	switch questionType {
	case DecisionQuestionTypeBoolean:
		var question BooleanQuestion
		if err := decodeStrict(raw, &question); err != nil {
			return nil, describeQuestionDecodeError(questionType, err)
		}
		if err := requireInstructions(question.Instructions); err != nil {
			return nil, err
		}
		return question, nil
	case DecisionQuestionTypeChoice:
		var question ChoiceQuestion
		if err := decodeStrict(raw, &question); err != nil {
			return nil, describeQuestionDecodeError(questionType, err)
		}
		if err := requireInstructions(question.Instructions); err != nil {
			return nil, err
		}
		return question, nil
	case DecisionQuestionTypeScore:
		var question ScoreQuestion
		if err := decodeStrict(raw, &question); err != nil {
			return nil, describeQuestionDecodeError(questionType, err)
		}
		if err := requireInstructions(question.Instructions); err != nil {
			return nil, err
		}
		return question, nil
	default:
		return nil, unknownDecisionType(questionType)
	}
}

// describeQuestionDecodeError rewrites encoding/json's container-shape errors
// (which name Go types) into the wire vocabulary a caller can act on. Errors
// raised by StructuredText and by unknown-field rejection already read well and
// pass through.
func describeQuestionDecodeError(questionType DecisionQuestionType, err error) error {
	var typeError *json.UnmarshalTypeError
	if !errors.As(err, &typeError) {
		return err
	}
	if strings.HasPrefix(typeError.Field, "criteria") {
		return fmt.Errorf("criteria: %s", criteriaShape(questionType))
	}
	if typeError.Field == "" {
		return errors.New("must be a JSON object")
	}
	return fmt.Errorf("%s has the wrong JSON shape", typeError.Field)
}

func criteriaShape(questionType DecisionQuestionType) string {
	switch questionType {
	case DecisionQuestionTypeBoolean:
		return `must be an object with "true" and "false" descriptions`
	case DecisionQuestionTypeChoice:
		return "must be an object mapping each option name to its description"
	case DecisionQuestionTypeScore:
		return "must be an ordered array of level descriptions"
	default:
		return "has the wrong JSON shape"
	}
}

// DecisionAnswer is the calibrated answer to one question — a tagged union decoded
// by its type. Concrete variants: BooleanAnswer, ChoiceAnswer, ScoreAnswer. Every
// variant carries confidence = (n·peak − 1)/(n − 1) over its distribution, so
// gating code can threshold without switching on the type.
type DecisionAnswer interface {
	AnswerType() DecisionQuestionType
}

// BooleanAnswer: Probability is P(true); Boolean is the argmax (probability > 0.5);
// Confidence is |2p − 1| (the general formula at n = 2).
type BooleanAnswer struct {
	Type        DecisionQuestionType `json:"type"`
	Boolean     bool                 `json:"boolean"`
	Probability float64              `json:"probability"`
	Confidence  float64              `json:"confidence"`
}

func (BooleanAnswer) AnswerType() DecisionQuestionType { return DecisionQuestionTypeBoolean }

// ChoiceAnswer: Choice is the highest-probability option; Probabilities is the full
// distribution keyed by option name.
type ChoiceAnswer struct {
	Type          DecisionQuestionType `json:"type"`
	Choice        string               `json:"choice"`
	Probabilities map[string]float64   `json:"probabilities"`
	Confidence    float64              `json:"confidence"`
}

func (ChoiceAnswer) AnswerType() DecisionQuestionType { return DecisionQuestionTypeChoice }

// ScoreAnswer: Score is the probability-weighted mean of the level numbers (so it
// can fall between levels); Legend echoes the level descriptions and Probabilities
// the distribution, both indexed by level (index = level number).
type ScoreAnswer struct {
	Type          DecisionQuestionType `json:"type"`
	Score         float64              `json:"score"`
	Legend        []StructuredText     `json:"legend"`
	Probabilities []float64            `json:"probabilities"`
	Confidence    float64              `json:"confidence"`
}

func (ScoreAnswer) AnswerType() DecisionQuestionType { return DecisionQuestionTypeScore }

// DecisionAnswers maps each question key to its answer. Decodes each value by its
// type so Go consumers of a DecisionResult get the typed variants back.
type DecisionAnswers map[string]DecisionAnswer

// UnmarshalJSON decodes every answer through DecodeDecisionAnswer, naming the
// offending key on failure.
func (answers *DecisionAnswers) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	if raw == nil {
		*answers = nil
		return nil
	}
	decoded := make(DecisionAnswers, len(raw))
	for key, value := range raw {
		answer, err := DecodeDecisionAnswer(value)
		if err != nil {
			return fmt.Errorf("answers[%q]: %w", key, err)
		}
		decoded[key] = answer
	}
	*answers = decoded
	return nil
}

// DecodeDecisionAnswer decodes one raw answer by its "type" into the matching variant.
func DecodeDecisionAnswer(raw json.RawMessage) (DecisionAnswer, error) {
	answerType, err := peekDecisionType(raw)
	if err != nil {
		return nil, err
	}
	switch answerType {
	case DecisionQuestionTypeBoolean:
		var answer BooleanAnswer
		if err := json.Unmarshal(raw, &answer); err != nil {
			return nil, err
		}
		return answer, nil
	case DecisionQuestionTypeChoice:
		var answer ChoiceAnswer
		if err := json.Unmarshal(raw, &answer); err != nil {
			return nil, err
		}
		return answer, nil
	case DecisionQuestionTypeScore:
		var answer ScoreAnswer
		if err := json.Unmarshal(raw, &answer); err != nil {
			return nil, err
		}
		return answer, nil
	default:
		return nil, unknownDecisionType(answerType)
	}
}

// DecisionRequest is a single, stateless decision call: one text-only state and a
// map of atomic questions the provider answers in parallel. It is the calibrated
// counterpart to GenerationRequest — no prose comes back, only typed answers.
//
// Model reuses ModelConfig (the same concept generate uses); only ModelId is
// meaningful for a decision, the sampling knobs are rejected by the domain rather
// than silently ignored.
type DecisionRequest struct {
	Model     ModelConfig       `json:"model"`
	State     StructuredText    `json:"state"`
	Questions DecisionQuestions `json:"questions"`
}

// DecisionResult is the provider's reply: one answer per question key, the model
// that produced them (resolved id, e.g. "jev-1.13.0"), and token usage.
type DecisionResult struct {
	ModelId string          `json:"model_id"`
	Answers DecisionAnswers `json:"answers"`
	Usage   ModelUsage      `json:"usage"`
}

func peekDecisionType(raw json.RawMessage) (DecisionQuestionType, error) {
	if len(bytes.TrimSpace(raw)) == 0 {
		return "", errors.New("must be a JSON object")
	}
	var head struct {
		Type DecisionQuestionType `json:"type"`
	}
	if err := json.Unmarshal(raw, &head); err != nil {
		var typeError *json.UnmarshalTypeError
		if errors.As(err, &typeError) {
			if typeError.Field == "" {
				return "", errors.New("must be a JSON object")
			}
			return "", errors.New("type must be a string (boolean | choice | score)")
		}
		return "", err
	}
	if head.Type == "" {
		return "", errors.New("type is required (boolean | choice | score)")
	}
	return head.Type, nil
}

func unknownDecisionType(decisionType DecisionQuestionType) error {
	return fmt.Errorf("unknown type %q (want boolean | choice | score)", decisionType)
}

func requireInstructions(instructions StructuredText) error {
	if len(instructions) == 0 {
		return errors.New("instructions is required")
	}
	return nil
}

// decodeStrict unmarshals with unknown fields rejected — a typo in a question is a
// caller bug worth a 400, not something to forward.
func decodeStrict(raw json.RawMessage, destination any) error {
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	return decoder.Decode(destination)
}
