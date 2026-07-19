package conversationmodel

// ConversationParticipant is one member of a conversation's roster — the
// complete set of everyone in the thread, INCLUDING the connected account
// itself. It embeds the identity snapshot (ContactRef) and marks the
// account's own entries with IsSelf, so a UI renders "who it's with" as the
// non-self subset without needing to know the connection's own channel identity
// or matching against the viewing user.
//
// ContactRef is embedded (json inline): external_id/name/email/platform_user
// sit at the top level of each jsonb array element alongside is_self, so the
// roster stays the queryable person dimension — a filter is a containment match
// (participants @> '[{"external_id":"…"}]'::jsonb, GIN-indexed).
type ConversationParticipant struct {
	ContactRef
	// IsSelf marks the connected account's own entry (Gmail: the mailbox address;
	// Slack: the bot user). Set server-side at ingest/send where the self identity
	// is known; the wire and UI never recompute it.
	IsSelf bool `json:"is_self"`
}
