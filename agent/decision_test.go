package agentmodel

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestStructuredText_AcceptsStringObjectArray(t *testing.T) {
	cases := map[string]string{
		"string": `"My card was charged twice."`,
		"object": `{"message":"hi","order_id":"A-104"}`,
		"array":  `["Hi","My card was charged twice."]`,
	}
	for name, raw := range cases {
		t.Run(name, func(t *testing.T) {
			var text StructuredText
			if err := json.Unmarshal([]byte(raw), &text); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if string(text) != raw {
				t.Errorf("decoded = %s, want %s", text, raw)
			}
			out, err := json.Marshal(text)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}
			if string(out) != raw {
				t.Errorf("marshalled = %s, want passthrough %s", out, raw)
			}
		})
	}
}

func TestStructuredText_Rejects(t *testing.T) {
	cases := map[string]string{
		"null":    `null`,
		"number":  `42`,
		"boolean": `true`,
	}
	for name, raw := range cases {
		t.Run(name, func(t *testing.T) {
			var text StructuredText
			if err := json.Unmarshal([]byte(raw), &text); err == nil {
				t.Errorf("expected error for %s, got nil", name)
			}
		})
	}
}

func TestDecodeDecisionQuestion_Variants(t *testing.T) {
	boolean, err := DecodeDecisionQuestion([]byte(`{"type":"boolean","instructions":"Is it urgent?"}`))
	if err != nil {
		t.Fatalf("boolean: %v", err)
	}
	booleanQuestion, ok := boolean.(BooleanQuestion)
	if !ok || booleanQuestion.Type != DecisionQuestionTypeBoolean || booleanQuestion.Criteria != nil {
		t.Errorf("boolean = %#v, want BooleanQuestion without criteria", boolean)
	}
	if booleanQuestion.QuestionType() != DecisionQuestionTypeBoolean {
		t.Errorf("QuestionType() = %q", booleanQuestion.QuestionType())
	}

	withCriteria, err := DecodeDecisionQuestion([]byte(`{"type":"boolean","instructions":"x","criteria":{"true":"asks for money back","false":{"what":"only reports"}}}`))
	if err != nil {
		t.Fatalf("boolean with criteria: %v", err)
	}
	if criteria := withCriteria.(BooleanQuestion).Criteria; criteria == nil || string(criteria.True) != `"asks for money back"` || string(criteria.False) != `{"what":"only reports"}` {
		t.Errorf("boolean criteria = %#v", criteria)
	}

	choice, err := DecodeDecisionQuestion([]byte(`{"type":"choice","instructions":"What is it about?","criteria":{"billing":"charges","bug":{"what":"defect","examples":["crash"]}}}`))
	if err != nil {
		t.Fatalf("choice: %v", err)
	}
	choiceQuestion, ok := choice.(ChoiceQuestion)
	if !ok || len(choiceQuestion.Criteria) != 2 || string(choiceQuestion.Criteria["bug"]) != `{"what":"defect","examples":["crash"]}` {
		t.Errorf("choice = %#v", choice)
	}

	score, err := DecodeDecisionQuestion([]byte(`{"type":"score","instructions":["How frustrated?","focus on tone"],"criteria":["calm","annoyed","angry"]}`))
	if err != nil {
		t.Fatalf("score: %v", err)
	}
	scoreQuestion, ok := score.(ScoreQuestion)
	if !ok || len(scoreQuestion.Criteria) != 3 || string(scoreQuestion.Criteria[2]) != `"angry"` {
		t.Errorf("score = %#v", score)
	}
	if string(scoreQuestion.Instructions) != `["How frustrated?","focus on tone"]` {
		t.Errorf("structured instructions = %s", scoreQuestion.Instructions)
	}
}

func TestDecodeDecisionQuestion_NullBooleanCriteriaIsAbsent(t *testing.T) {
	question, err := DecodeDecisionQuestion([]byte(`{"type":"boolean","instructions":"x","criteria":null}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if question.(BooleanQuestion).Criteria != nil {
		t.Error("null criteria should decode to nil")
	}
}

func TestDecodeDecisionQuestion_Rejects(t *testing.T) {
	cases := map[string]string{
		"vendor type noul":           `{"type":"noul","instructions":"x"}`,
		"unknown type":               `{"type":"rating","instructions":"x"}`,
		"missing type":               `{"instructions":"x"}`,
		"missing instructions":       `{"type":"boolean"}`,
		"null instructions":          `{"type":"boolean","instructions":null}`,
		"number instructions":        `{"type":"boolean","instructions":3}`,
		"choice criteria array":      `{"type":"choice","instructions":"x","criteria":["a","b"]}`,
		"choice criteria null value": `{"type":"choice","instructions":"x","criteria":{"a":null}}`,
		"score criteria object":      `{"type":"score","instructions":"x","criteria":{"0":"a","1":"b"}}`,
		"score null level":           `{"type":"score","instructions":"x","criteria":["a",null,"c"]}`,
		"boolean stray criteria key": `{"type":"boolean","instructions":"x","criteria":{"true":"a","maybe":"b"}}`,
		"boolean criteria array":     `{"type":"boolean","instructions":"x","criteria":["a","b"]}`,
		"unknown field":              `{"type":"boolean","instruction":"x"}`,
		"not an object":              `"boolean"`,
		"empty":                      ``,
	}
	for name, raw := range cases {
		t.Run(name, func(t *testing.T) {
			if _, err := DecodeDecisionQuestion([]byte(raw)); err == nil {
				t.Errorf("expected error for %q, got nil", name)
			}
		})
	}
}

func TestDecodeDecisionQuestion_ShapeErrorsSpeakWire(t *testing.T) {
	cases := map[string]struct {
		raw  string
		want string
	}{
		"choice criteria array":  {`{"type":"choice","instructions":"x","criteria":["a"]}`, "criteria: must be an object mapping each option name to its description"},
		"score criteria object":  {`{"type":"score","instructions":"x","criteria":{"0":"a"}}`, "criteria: must be an ordered array of level descriptions"},
		"boolean criteria array": {`{"type":"boolean","instructions":"x","criteria":["a"]}`, `criteria: must be an object with "true" and "false" descriptions`},
		"type not a string":      {`{"type":5}`, "type"},
	}
	for name, testCase := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := DecodeDecisionQuestion([]byte(testCase.raw))
			if err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !strings.Contains(err.Error(), testCase.want) {
				t.Errorf("error = %q, want it to contain %q", err.Error(), testCase.want)
			}
			if strings.Contains(err.Error(), "Go struct") || strings.Contains(err.Error(), "agentmodel.") {
				t.Errorf("error leaks Go types: %q", err.Error())
			}
		})
	}
}

func TestDecisionRequest_RoundTrip(t *testing.T) {
	raw := `{"model":{"model_id":"jev-latest"},"state":{"message":"hi"},"questions":{` +
		`"refund":{"type":"boolean","instructions":"x"},` +
		`"topic":{"type":"choice","instructions":"y","criteria":{"a":"A"}},` +
		`"mood":{"type":"score","instructions":"z","criteria":["l0","l1"]}}}`

	var request DecisionRequest
	if err := json.Unmarshal([]byte(raw), &request); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if _, ok := request.Questions["refund"].(BooleanQuestion); !ok {
		t.Errorf("refund = %T, want BooleanQuestion", request.Questions["refund"])
	}
	if _, ok := request.Questions["topic"].(ChoiceQuestion); !ok {
		t.Errorf("topic = %T, want ChoiceQuestion", request.Questions["topic"])
	}
	if _, ok := request.Questions["mood"].(ScoreQuestion); !ok {
		t.Errorf("mood = %T, want ScoreQuestion", request.Questions["mood"])
	}

	out, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	for _, want := range []string{`"type":"boolean"`, `"type":"choice"`, `"type":"score"`, `"state":{"message":"hi"}`} {
		if !strings.Contains(string(out), want) {
			t.Errorf("marshalled request lacks %s: %s", want, out)
		}
	}
	var again DecisionRequest
	if err := json.Unmarshal(out, &again); err != nil {
		t.Fatalf("re-unmarshal: %v", err)
	}
	if len(again.Questions) != 3 {
		t.Errorf("round trip lost questions: %d", len(again.Questions))
	}
}

func TestDecisionQuestions_NamesOffendingKey(t *testing.T) {
	var questions DecisionQuestions
	err := json.Unmarshal([]byte(`{"ok":{"type":"boolean","instructions":"x"},"bad":{"type":"choice","instructions":"y","criteria":["a"]}}`), &questions)
	if err == nil || !strings.Contains(err.Error(), `questions["bad"]`) {
		t.Errorf("error = %v, want one naming questions[\"bad\"]", err)
	}
}

func TestDecisionAnswers_RoundTrip(t *testing.T) {
	raw := `{` +
		`"refund":{"type":"boolean","boolean":false,"probability":0,"confidence":1},` +
		`"topic":{"type":"choice","choice":"a","probabilities":{"a":1},"confidence":1},` +
		`"mood":{"type":"score","score":0.5,"legend":["l0","l1"],"probabilities":[0.5,0.5],"confidence":0}}`

	var answers DecisionAnswers
	if err := json.Unmarshal([]byte(raw), &answers); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	boolean, ok := answers["refund"].(BooleanAnswer)
	if !ok || boolean.Boolean || boolean.Probability != 0 || boolean.Confidence != 1 {
		t.Errorf("refund = %#v", answers["refund"])
	}
	if choice, ok := answers["topic"].(ChoiceAnswer); !ok || choice.Choice != "a" || choice.Probabilities["a"] != 1 {
		t.Errorf("topic = %#v", answers["topic"])
	}
	score, ok := answers["mood"].(ScoreAnswer)
	if !ok || score.Score != 0.5 || len(score.Legend) != 2 || string(score.Legend[1]) != `"l1"` || score.Probabilities[1] != 0.5 {
		t.Errorf("mood = %#v", answers["mood"])
	}

	out, err := json.Marshal(DecisionResult{ModelId: "jev-1.13.0", Answers: answers})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	for _, want := range []string{`"boolean":false`, `"probability":0`, `"confidence":0`, `"legend":["l0","l1"]`, `"probabilities":[0.5,0.5]`} {
		if !strings.Contains(string(out), want) {
			t.Errorf("marshalled result lacks %s: %s", want, out)
		}
	}
}

func TestDecodeDecisionAnswer_RejectsUnknownType(t *testing.T) {
	if _, err := DecodeDecisionAnswer([]byte(`{"type":"noul","noul":0.5}`)); err == nil {
		t.Error("expected error for vendor type noul, got nil")
	}
}
