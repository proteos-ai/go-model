package conversationmodel

// MessageRecipient is one addressee of a message together with the role it was
// addressed under (to/cc/bcc). It embeds the identity snapshot (ParticipantRef)
// so a recipient carries the same resolved name/email/platform_user as a sender.
// Only email distinguishes cc/bcc; other channels address everyone as `to`. Bcc
// entries are only ever present on messages WE sent — inbound bcc is stripped by
// SMTP and unknowable.
type MessageRecipient struct {
	ParticipantRef
	Role RecipientRole `json:"role"`
}
