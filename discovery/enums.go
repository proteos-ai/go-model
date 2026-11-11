// Package discoverymodel is the shared model of discovery-service: the
// evidence users send Proteos to show how they work — a screenshot with ink,
// a screen recording with ink-annotated frames, or written feedback. Evidence
// lives in the submitter's org; Proteos staff read it across orgs.
package discoverymodel

// DiscoveryEvidenceKind is what the user sent.
type DiscoveryEvidenceKind string

const (
	DiscoveryEvidenceKindScreenshot DiscoveryEvidenceKind = "screenshot"
	DiscoveryEvidenceKindRecording  DiscoveryEvidenceKind = "recording"
	DiscoveryEvidenceKindFeedback   DiscoveryEvidenceKind = "feedback"
)

// DiscoveryEvidenceCategory is the user's own framing of why they sent it.
type DiscoveryEvidenceCategory string

const (
	DiscoveryEvidenceCategoryBug      DiscoveryEvidenceCategory = "bug"
	DiscoveryEvidenceCategoryIdea     DiscoveryEvidenceCategory = "idea"
	DiscoveryEvidenceCategoryQuestion DiscoveryEvidenceCategory = "question"
	DiscoveryEvidenceCategoryPraise   DiscoveryEvidenceCategory = "praise"
	DiscoveryEvidenceCategoryOther    DiscoveryEvidenceCategory = "other"
)

// DiscoveryEvidenceStatus is Proteos' review state. Only platform admins move
// it; the submitter never sees a status other than their own record's.
type DiscoveryEvidenceStatus string

const (
	DiscoveryEvidenceStatusNew      DiscoveryEvidenceStatus = "new"
	DiscoveryEvidenceStatusReviewed DiscoveryEvidenceStatus = "reviewed"
	DiscoveryEvidenceStatusArchived DiscoveryEvidenceStatus = "archived"
)

// DiscoveryEvidenceSource is which client captured it.
type DiscoveryEvidenceSource string

const (
	DiscoveryEvidenceSourceWeb       DiscoveryEvidenceSource = "web"
	DiscoveryEvidenceSourceExtension DiscoveryEvidenceSource = "extension"
)

// DiscoveryEvidenceEnrichmentStatus is the state of the background call that
// writes the title and category the sender no longer types. pending until the
// call ran (or was found unnecessary), then completed or failed.
type DiscoveryEvidenceEnrichmentStatus string

const (
	DiscoveryEvidenceEnrichmentStatusPending   DiscoveryEvidenceEnrichmentStatus = "pending"
	DiscoveryEvidenceEnrichmentStatusCompleted DiscoveryEvidenceEnrichmentStatus = "completed"
	DiscoveryEvidenceEnrichmentStatusFailed    DiscoveryEvidenceEnrichmentStatus = "failed"
)

// DiscoveryEvidenceFieldSource says who wrote a field: "" when it is empty or
// was never set, generated for the title and category the model wrote,
// narration for notes and frame comments drafted from what the sender said,
// user for anything a person typed or edited. A user value is never
// overwritten by a model.
type DiscoveryEvidenceFieldSource string

const (
	DiscoveryEvidenceFieldSourceNone      DiscoveryEvidenceFieldSource = ""
	DiscoveryEvidenceFieldSourceGenerated DiscoveryEvidenceFieldSource = "generated"
	DiscoveryEvidenceFieldSourceNarration DiscoveryEvidenceFieldSource = "narration"
	DiscoveryEvidenceFieldSourceUser      DiscoveryEvidenceFieldSource = "user"
)

// DiscoveryEvidenceNarrationStatus is where a narration job stands.
type DiscoveryEvidenceNarrationStatus string

const (
	DiscoveryEvidenceNarrationStatusTranscribing DiscoveryEvidenceNarrationStatus = "transcribing"
	DiscoveryEvidenceNarrationStatusDrafting     DiscoveryEvidenceNarrationStatus = "drafting"
	DiscoveryEvidenceNarrationStatusCompleted    DiscoveryEvidenceNarrationStatus = "completed"
	DiscoveryEvidenceNarrationStatusFailed       DiscoveryEvidenceNarrationStatus = "failed"
)
