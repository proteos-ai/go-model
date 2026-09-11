package conversationmodel

import "time"

// EmailSuppressionKind names one of a sending platform's suppression lists —
// the addresses the provider itself refuses to send to. Provider-agnostic
// vocabulary (SendGrid, Mailgun and Postmark all keep these five).
type EmailSuppressionKind string

const (
	EmailSuppressionKindBounces       EmailSuppressionKind = "bounces"
	EmailSuppressionKindBlocks        EmailSuppressionKind = "blocks"
	EmailSuppressionKindSpamReports   EmailSuppressionKind = "spam_reports"
	EmailSuppressionKindInvalidEmails EmailSuppressionKind = "invalid_emails"
	EmailSuppressionKindUnsubscribes  EmailSuppressionKind = "unsubscribes"
)

func IsKnownEmailSuppressionKind(kind EmailSuppressionKind) bool {
	switch kind {
	case EmailSuppressionKindBounces, EmailSuppressionKindBlocks, EmailSuppressionKindSpamReports,
		EmailSuppressionKindInvalidEmails, EmailSuppressionKindUnsubscribes:
		return true
	}
	return false
}

// EmailHealth is the deliverability picture of one sending-platform
// connection: rates computed from OUR channel_event ledger (Windows — the
// authoritative per-connection view, every event we received), plus the
// provider's own account-level snapshot (Provider — engagement quality
// score, its counters, suppression lists, sender authentication, dedicated
// IPs), refreshed daily by the connector's sweep or on demand.
type EmailHealth struct {
	ConnectionId string              `json:"connection_id"`
	Windows      []EmailHealthWindow `json:"windows"`
	Provider     *EmailProviderHealth `json:"provider,omitempty"`
}

// EmailHealthWindow is the ledger rollup over the last Days days. Counts are
// DISTINCT messages per event type (a mail opened three times is one opened
// message); machine opens are excluded. Rates are fractions (0–1):
// delivery / bounce / hard-bounce over sent (processed), open / click /
// spam / unsubscribe over delivered.
type EmailHealthWindow struct {
	Days            int     `json:"days"`
	Sent            int     `json:"sent"`
	Delivered       int     `json:"delivered"`
	Deferred        int     `json:"deferred"`
	Bounced         int     `json:"bounced"`
	HardBounced     int     `json:"hard_bounced"`
	Dropped         int     `json:"dropped"`
	Opened          int     `json:"opened"`
	Clicked         int     `json:"clicked"`
	SpamReported    int     `json:"spam_reported"`
	Unsubscribed    int     `json:"unsubscribed"`
	DeliveryRate    float64 `json:"delivery_rate"`
	BounceRate      float64 `json:"bounce_rate"`
	HardBounceRate  float64 `json:"hard_bounce_rate"`
	OpenRate        float64 `json:"open_rate"`
	ClickRate       float64 `json:"click_rate"`
	SpamRate        float64 `json:"spam_rate"`
	UnsubscribeRate float64 `json:"unsubscribe_rate"`
}

// EmailProviderHealth is the provider's account-level snapshot for one
// connection. Every section is optional: a plan without the feature, or a
// failed pull, leaves it nil and adds a Warning — the rest still renders.
type EmailProviderHealth struct {
	EngagementQuality *EmailEngagementQuality `json:"engagement_quality,omitempty"`
	// Stats are the provider's own counters for THIS connection's traffic
	// (its category) over the last Days days.
	Stats        *EmailProviderStats     `json:"stats,omitempty"`
	Suppressions *EmailSuppressionCounts `json:"suppressions,omitempty"`
	Domain       *EmailDomainStatus      `json:"domain,omitempty"`
	Sender       *EmailSenderStatus      `json:"sender,omitempty"`
	// Ips are the account's dedicated IPs (empty on shared pools).
	Ips         []EmailIp `json:"ips"`
	Warnings    []string  `json:"warnings,omitempty"`
	RefreshedAt time.Time `json:"refreshed_at"`
}

// EmailEngagementQuality is the provider's composite sender-quality score
// (SendGrid Engagement Quality: 1–5 overall, one sub-score per dimension,
// computed daily; needs open tracking and volume).
type EmailEngagementQuality struct {
	Date                 string  `json:"date"`
	Score                float64 `json:"score"`
	BounceClassification float64 `json:"bounce_classification"`
	BounceRate           float64 `json:"bounce_rate"`
	EngagementRecency    float64 `json:"engagement_recency"`
	OpenRate             float64 `json:"open_rate"`
	SpamRate             float64 `json:"spam_rate"`
}

// EmailProviderStats mirrors the provider's statistics metrics verbatim
// (their names ARE the vocabulary senders know from the provider's UI).
type EmailProviderStats struct {
	Days            int `json:"days"`
	Requests        int `json:"requests"`
	Processed       int `json:"processed"`
	Delivered       int `json:"delivered"`
	Deferred        int `json:"deferred"`
	Bounces         int `json:"bounces"`
	Blocks          int `json:"blocks"`
	BounceDrops     int `json:"bounce_drops"`
	InvalidEmails   int `json:"invalid_emails"`
	Opens           int `json:"opens"`
	UniqueOpens     int `json:"unique_opens"`
	Clicks          int `json:"clicks"`
	UniqueClicks    int `json:"unique_clicks"`
	SpamReports     int `json:"spam_reports"`
	SpamReportDrops int `json:"spam_report_drops"`
	Unsubscribes    int `json:"unsubscribes"`
	UnsubscribeDrops int `json:"unsubscribe_drops"`
}

// EmailSuppressionCounts sizes the provider's suppression lists.
// IsApproximate: a list hit the page cap, the count is a floor.
type EmailSuppressionCounts struct {
	Bounces       int  `json:"bounces"`
	Blocks        int  `json:"blocks"`
	SpamReports   int  `json:"spam_reports"`
	InvalidEmails int  `json:"invalid_emails"`
	Unsubscribes  int  `json:"unsubscribes"`
	IsApproximate bool `json:"is_approximate,omitempty"`
}

// EmailDomainStatus is the sending domain's authentication state with the
// DNS records the customer must publish.
type EmailDomainStatus struct {
	Domain  string           `json:"domain"`
	IsValid bool             `json:"is_valid"`
	Dns     []EmailDnsRecord `json:"dns"`
}

type EmailDnsRecord struct {
	Host    string `json:"host"`
	Type    string `json:"type"`
	Data    string `json:"data"`
	IsValid bool   `json:"is_valid"`
	// Reason is the provider's explanation of an invalid record ("" when
	// valid or unknown).
	Reason string `json:"reason,omitempty"`
}

// EmailSenderStatus is the from identity's verification state.
type EmailSenderStatus struct {
	FromEmail  string `json:"from_email"`
	Mode       string `json:"mode"`
	IsVerified bool   `json:"is_verified"`
}

// EmailIp is one dedicated IP of the account with its warmup state.
type EmailIp struct {
	Ip              string     `json:"ip"`
	IsWarmingUp     bool       `json:"is_warming_up"`
	WarmupStartedAt *time.Time `json:"warmup_started_at,omitempty"`
	Pools           []string   `json:"pools"`
}

// EmailSuppression is one address on a provider suppression list.
type EmailSuppression struct {
	Kind      EmailSuppressionKind `json:"kind"`
	Email     string               `json:"email"`
	Reason    string               `json:"reason,omitempty"`
	Status    string               `json:"status,omitempty"`
	CreatedAt time.Time            `json:"created_at"`
}
