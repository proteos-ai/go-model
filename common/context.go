package common

import (
	"context"
)

type Context interface {
	context.Context
	GetUserRef() UserRef
	GetUserId() string
	GetToken() string
	GetOrgId() string
}

type AppContext struct {
	context.Context
}

// GetUserRef returns the UserRef responsible for the current write — a real
// user (Type person/agent/api) for authenticated requests, or the system
// sentinel (SystemUserRef) for system/bootstrap writes. Returns the zero
// UserRef when unset.
func (ctx AppContext) GetUserRef() UserRef {
	v := ctx.Value(ContextTagSource)
	if v == nil {
		return UserRef{}
	}
	ref, _ := v.(UserRef)
	return ref
}

// GetUserId returns just the id of the user responsible for the current write
// (see GetUserRef) — a user id (UUID string) for authenticated requests, or
// PlatformUserId ("platform") for system/bootstrap writes. Returns "" when
// unset.
func (ctx AppContext) GetUserId() string {
	return ctx.GetUserRef().Id
}

// GetToken returns the full Authorization header value the
// authentication middleware stashed under ContextTagToken — e.g.
// `Bearer eyJ...`. The canonical internal token format across the
// codebase is the full header value; callers forward it verbatim to
// downstream services rather than stripping and re-prepending.
func (ctx AppContext) GetToken() string {
	// Nil-safe: a detached/background context may carry no token — return "" rather
	// than panicking on the type assertion (matches GetOrgId's behaviour).
	token, _ := ctx.Value(ContextTagToken).(string)
	return token
}

func (ctx AppContext) GetOrgId() string {
	orgId := ctx.Value(ContextTagOrgId)
	if orgId == nil {
		return ""
	}
	return orgId.(string)
}
