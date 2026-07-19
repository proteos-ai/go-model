package conversationmodel

// ContactObservation is ONE sighting of a person as a connector saw them — the
// wire identity plus optional provider enrichment. It is the input to contact
// resolution: connectors return observations from directory sweeps (and ingest
// derives them from message senders/recipients); the domain canonicalizes them
// into contact_address rows and resolves/mints the owning contact. Connectors
// never stamp ids, sources, or audit — that is the domain's job.
type ContactObservation struct {
	// ExternalId is the connector-side identity as observed on the wire (Slack
	// user id, email address, WhatsApp JID) — the raw form, never canonicalized
	// by the connector.
	ExternalId string `json:"external_id"`
	Name       string `json:"name,omitempty"`
	// Email is the person's address on the provider side (Slack profile email,
	// the correspondent address itself for email) when the provider exposes
	// one. Yields the secondary email address key — the cross-channel join.
	Email string `json:"email,omitempty"`
	// Metadata is per-channel provider enrichment (Slack handle, avatar, …)
	// destined for ContactAddress.Metadata.
	Metadata map[string]any `json:"metadata,omitempty"`
}
