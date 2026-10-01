package conversationapi

import (
	"testing"

	conversationmodel "go.proteos.ai/model/conversation"
)

func TestCreateConversationTagDefinitionRequestNormalized(t *testing.T) {
	isDisabled := false
	boolean := conversationmodel.ConversationTagQuestion{Type: conversationmodel.ConversationTagQuestionTypeBoolean}

	omitted := CreateConversationTagDefinitionRequest{Question: boolean}.Normalized()
	if omitted.Evaluator.Kind != conversationmodel.ConversationTagEvaluatorDecide {
		t.Fatalf("an omitted evaluator defaults to decide, got %q", omitted.Evaluator.Kind)
	}
	if omitted.Evaluator.Screening != nil {
		t.Fatal("decide never gets a screening")
	}
	if omitted.IsEnabled == nil || !*omitted.IsEnabled {
		t.Fatal("is_enabled nil defaults to true")
	}

	reason := CreateConversationTagDefinitionRequest{
		Question:  boolean,
		Evaluator: conversationmodel.ConversationTagEvaluator{Kind: conversationmodel.ConversationTagEvaluatorReason},
		IsEnabled: &isDisabled,
	}.Normalized()
	if reason.Evaluator.Screening == nil || !reason.Evaluator.Screening.IsEnabled || reason.Evaluator.Screening.MinProbability != nil {
		t.Fatalf("reason + boolean defaults to screening on with the service threshold, got %+v", reason.Evaluator.Screening)
	}
	if *reason.IsEnabled {
		t.Fatal("an explicit is_enabled is kept")
	}

	choice := CreateConversationTagDefinitionRequest{
		Question:  conversationmodel.ConversationTagQuestion{Type: conversationmodel.ConversationTagQuestionTypeChoice},
		Evaluator: conversationmodel.ConversationTagEvaluator{Kind: conversationmodel.ConversationTagEvaluatorReason},
	}.Normalized()
	if choice.Evaluator.Screening != nil {
		t.Fatal("a choice question is never screened")
	}

	explicit := CreateConversationTagDefinitionRequest{
		Question:  boolean,
		Evaluator: conversationmodel.ConversationTagEvaluator{Kind: conversationmodel.ConversationTagEvaluatorReason, Screening: &conversationmodel.ConversationTagScreening{IsEnabled: false}},
	}.Normalized()
	if explicit.Evaluator.Screening.IsEnabled {
		t.Fatal("an explicit screening is kept")
	}
}

func TestCreateConversationTagSetRequestNormalized(t *testing.T) {
	omitted := CreateConversationTagSetRequest{}.Normalized()
	if *omitted.QuietWindowSeconds != 90 || *omitted.MaxWaitSeconds != 600 || !*omitted.IsEnabled {
		t.Fatalf("omitted timing defaults to 90/600 enabled, got %d/%d %v", *omitted.QuietWindowSeconds, *omitted.MaxWaitSeconds, *omitted.IsEnabled)
	}

	longQuiet := 900
	floored := CreateConversationTagSetRequest{QuietWindowSeconds: &longQuiet}.Normalized()
	if *floored.MaxWaitSeconds != 900 {
		t.Fatalf("an omitted max wait is floored to the quiet window, got %d", *floored.MaxWaitSeconds)
	}

	shortMax := 30
	explicit := CreateConversationTagSetRequest{QuietWindowSeconds: &longQuiet, MaxWaitSeconds: &shortMax}.Normalized()
	if *explicit.MaxWaitSeconds != 30 {
		t.Fatalf("an explicit max wait is kept for validation to reject, got %d", *explicit.MaxWaitSeconds)
	}
}
