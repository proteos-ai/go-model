package conversationmodel

import "go.proteos.ai/model/common"

// EmailSenderDomain is what a sending-platform connection may send AS: every
// local part on Domain is a valid From (the provider signs the whole domain),
// and replies route through ReplyHostname (the Inbound Parse hostname). A
// connection in single-sender mode has none — it sends only as its one
// verified address.
type EmailSenderDomain struct {
	Domain        string `json:"domain"`
	ReplyHostname string `json:"reply_hostname"`
	// ConfiguredAddress is the from address the connection was installed
	// with — the fallback sending identity while no default sender is pinned
	// and the acting user has no address on the domain.
	ConfiguredAddress string `json:"configured_address"`
}

// EmailSender is ONE sending identity of a connection: an email contact
// address on the connection's sender domain. A personal sender belongs to a
// contact bound to a platform user (the colleague sends as themselves); a
// shared mailbox (sales@) belongs to a role contact without one. Computed on
// read from the contact directory ⊕ the connection's default — never stored
// as its own row.
type EmailSender struct {
	ContactAddressId string `json:"contact_address_id"`
	ContactId        string `json:"contact_id"`
	FromEmail        string `json:"from_email"`
	FromName         string `json:"from_name"`
	// PlatformUser is the colleague this sender belongs to; nil for a shared
	// mailbox.
	PlatformUser *common.UserRef `json:"platform_user,omitempty"`
	// IsDefault marks the connection's fallback sender — used when neither the
	// conversation nor the acting user pins one.
	IsDefault bool `json:"is_default"`
}
