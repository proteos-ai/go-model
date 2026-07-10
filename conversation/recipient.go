package conversationmodel

// Recipient is the send-time addressing value object: who/where an originated
// message goes. It is deliberately NOT a stored entity — the directory rows
// are Participant (a person) and Room (a venue); a Recipient merely points at
// one of them by kind + external id. Senders and reactors are always people
// (ParticipantRef); only a send target can be either.
type Recipient struct {
	Kind RecipientKind `json:"kind"`
	// ExternalId is the connector-side identity of the target — a Room's
	// external id (Slack channel id) or a Participant's (Slack user id, email
	// To address).
	ExternalId string `json:"external_id"`
}
