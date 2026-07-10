package common

type UserType string

const (
	// UserTypePerson is a human user.
	UserTypePerson UserType = "person"
	// UserTypeAgent is an AI agent acting on its own.
	UserTypeAgent UserType = "agent"
	// UserTypeApi is a non-interactive API / service client. Reserved for
	// future use — not yet emitted by any writer.
	UserTypeApi UserType = "api"
)

// UserRef is a reference to a user — the party responsible for a write. Every
// actor is modelled as a user with a kind (Type): a person (a human), an agent
// (an AI actor), or an api client. System/bootstrap writes use the sentinel
// {Type: person, Id: PlatformUserId} (see SystemUserRef).
//
// It is stored as JSONB on the audit columns (created_by/updated_by) and as the
// value of any `user` attribute, so agent (and later API) authorship is
// first-class. Filtering keys on Id; the chip resolves the user by Id.
type UserRef struct {
	Type UserType `json:"type"` // person | agent | api
	Id   string   `json:"id"`
}
