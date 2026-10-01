package conversationapi

import conversationmodel "go.proteos.ai/model/conversation"

// Normalized returns the request with every server-side default filled in,
// exactly as the row will be stored: evaluator.kind empty ⇒ decide; screening
// omitted on a reason definition with a boolean question ⇒ enabled with the
// service threshold; is_enabled nil ⇒ true. The ONE place these defaults live:
// the service applies it before validating and storing, and `pro module`
// applies it to both the local manifest and the deployed row so an omitted
// field never reads as drift.
func (request CreateConversationTagDefinitionRequest) Normalized() CreateConversationTagDefinitionRequest {
	if request.Evaluator.Kind == "" {
		request.Evaluator.Kind = conversationmodel.ConversationTagEvaluatorDecide
	}
	if request.Evaluator.Kind == conversationmodel.ConversationTagEvaluatorReason &&
		request.Question.Type == conversationmodel.ConversationTagQuestionTypeBoolean && request.Evaluator.Screening == nil {
		request.Evaluator.Screening = &conversationmodel.ConversationTagScreening{IsEnabled: true}
	}
	if request.IsEnabled == nil {
		isEnabled := true
		request.IsEnabled = &isEnabled
	}
	return request
}

// Normalized returns the request with the stored defaults filled in:
// quiet_window_seconds nil ⇒ 90; max_wait_seconds nil ⇒ 600, floored to the
// quiet window so an explicit long quiet window never yields an unsatisfiable
// timing; is_enabled nil ⇒ true. Shared by the service and `pro module` (see
// CreateConversationTagDefinitionRequest.Normalized).
func (request CreateConversationTagSetRequest) Normalized() CreateConversationTagSetRequest {
	if request.QuietWindowSeconds == nil {
		quietWindowSeconds := conversationmodel.ConversationTagSetDefaultQuietWindowSeconds
		request.QuietWindowSeconds = &quietWindowSeconds
	}
	if request.MaxWaitSeconds == nil {
		maxWaitSeconds := max(conversationmodel.ConversationTagSetDefaultMaxWaitSeconds, *request.QuietWindowSeconds)
		request.MaxWaitSeconds = &maxWaitSeconds
	}
	if request.IsEnabled == nil {
		isEnabled := true
		request.IsEnabled = &isEnabled
	}
	return request
}
