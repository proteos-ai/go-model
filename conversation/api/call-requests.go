package conversationapi

import "time"

// MintCallTokenRequest is the body of POST /conversations/v1/calls/token —
// the softphone's browser-calling grant for ONE connection. The client
// identity is derived server-side from the authenticated user, never from the
// request.
type MintCallTokenRequest struct {
	ConnectionId string `json:"connection_id" validate:"required"`
}

// PhoneNumberStatus says which system answers a provider phone number:
// connected to this platform (calls become conversations and ring in-app),
// in use by another integration provider-side (connecting takes it over,
// disconnecting restores it), or configured nowhere (available).
type PhoneNumberStatus string

const (
	PhoneNumberStatusConnected PhoneNumberStatus = "connected"
	PhoneNumberStatusAvailable PhoneNumberStatus = "available"
	PhoneNumberStatusInUse     PhoneNumberStatus = "in_use"
)

// RoutingTargetType discriminates a routing target: one platform user, or a
// team (expanded to its direct members at ring time, so team edits apply to
// the very next call).
type RoutingTargetType string

const (
	RoutingTargetUser RoutingTargetType = "user"
	RoutingTargetTeam RoutingTargetType = "team"
)

// RoutingTarget is one dispatch target. Id is a platform user id for
// type=user and a team slug for type=team (teams are slug-keyed config).
type RoutingTarget struct {
	Type RoutingTargetType `json:"type"`
	Id   string            `json:"id"`
}

// PhoneNumberRouting is the user-level dispatch config of a connected number:
// who a call rings. v1 rings every expanded target in parallel (first accept
// wins; only registered softphones actually ring); the object shape leaves
// room for later strategy/ordering options without a wire break. Empty
// targets fall back to the connection-wide routing, then busy.
type PhoneNumberRouting struct {
	Targets []RoutingTarget `json:"targets"`
}

// PhoneNumber is one provider phone number on a phone connection. ExternalId
// is the provider-side number identity (house convention — vendor vocabulary
// like "sid" stays inside the adapter).
type PhoneNumber struct {
	ExternalId   string            `json:"external_id"`
	PhoneNumber  string            `json:"phone_number"`
	FriendlyName string            `json:"friendly_name"`
	Status       PhoneNumberStatus `json:"status"`
	// UsedBy is a DISPLAY-ONLY, provider-worded description of what an in_use
	// number is used by ("https://…", "another Twilio application"). Never
	// parsed, never round-tripped; empty unless status=in_use.
	UsedBy  string             `json:"used_by,omitempty"`
	Routing PhoneNumberRouting `json:"routing"`
}

// UpdatePhoneNumberRequest is the body of
// PUT /conversations/v1/connections/:id/phone-numbers/:externalId — the
// routing update. Connecting/disconnecting the number itself are actions
// (POST …/:externalId/connect, POST …/:externalId/disconnect).
type UpdatePhoneNumberRequest struct {
	// Routing replaces the number's dispatch config; empty targets clear the
	// mapping back to the connection-wide fallback.
	Routing *PhoneNumberRouting `json:"routing,omitempty"`
}

// CallTokenResponse is the minted browser-calling grant the softphone hands
// to the vendor's JS SDK (Twilio Voice: `new Device(token, { edge })`).
// Clients refresh by calling the endpoint again (the SDK's tokenWillExpire
// event is the cue).
type CallTokenResponse struct {
	Token    string `json:"token"`
	Identity string `json:"identity"`
	// Edge is the vendor edge the SDK should connect through ("dublin" for the
	// EU default); empty means the SDK default.
	Edge      string    `json:"edge,omitempty"`
	ExpiresAt time.Time `json:"expires_at"`
}
