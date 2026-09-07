package conversationmodel

import "time"

// SendEligibility is the answer to "may this send go out right now?" — the
// read model behind the dry-run endpoint and the shape a denied Send carries in
// its error details. Reason is the error code the send would fail with
// (sending_window_closed, sending_limit_reached, frequency_cap_reached,
// contact_blocked, contact_opted_out); empty when allowed. EarliestAllowedAt is
// the earliest candidate instant every temporal blocker clears — nil when
// allowed or when the block is permanent (suppression). The first failing
// check wins RuleId/RuleType/ContactId.
type SendEligibility struct {
	IsAllowed         bool            `json:"is_allowed"`
	Reason            string          `json:"reason,omitempty"`
	RuleId            string          `json:"rule_id,omitempty"`
	RuleType          SendingRuleType `json:"rule_type,omitempty"`
	EarliestAllowedAt *time.Time      `json:"earliest_allowed_at,omitempty"`
	ContactId         string          `json:"contact_id,omitempty"`
	ContactAddressId  string          `json:"contact_address_id,omitempty"`
}

// Details projects the eligibility onto the error `details` object clients
// (SDK, MCP, UI) read to schedule a retry.
func (eligibility SendEligibility) Details() map[string]any {
	details := map[string]any{}
	if eligibility.RuleId != "" {
		details["rule_id"] = eligibility.RuleId
	}
	if eligibility.RuleType != "" {
		details["rule_type"] = eligibility.RuleType
	}
	if eligibility.EarliestAllowedAt != nil {
		details["earliest_allowed_at"] = eligibility.EarliestAllowedAt.UTC().Format(time.RFC3339)
	}
	if eligibility.ContactId != "" {
		details["contact_id"] = eligibility.ContactId
	}
	if eligibility.ContactAddressId != "" {
		details["contact_address_id"] = eligibility.ContactAddressId
	}
	return details
}
