package conversationmodel

// SendingLimitPreset is a named bundle of recommended `limit` rule configs for
// one class of sender account (a free LinkedIn account, an established
// mailbox, …). Presets are a static catalog, never persisted: applying one
// expands into concrete SendingRule rows of type limit linked to the chosen
// connections — nothing at evaluation time knows about presets (the lemlist
// "recommended defaults" model).
type SendingLimitPreset struct {
	Key         string `json:"key"`
	Name        string `json:"name"`
	Description string `json:"description"`
	// ConnectorKeys lists the connectors this preset is meant for; the catalog
	// endpoint filters by it and ApplyPreset rejects a mismatch.
	ConnectorKeys []ConnectorKey `json:"connector_keys"`
	// IsRecommended marks the conservative default the UI highlights.
	IsRecommended bool                     `json:"is_recommended"`
	Rules         []SendingLimitPresetRule `json:"rules"`
}

// SendingLimitPresetRule is one limit the preset expands into.
type SendingLimitPresetRule struct {
	RuleType   SendingRuleType `json:"rule_type"`
	RuleConfig LimitRuleConfig `json:"rule_config"`
}
