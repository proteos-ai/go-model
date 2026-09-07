package conversationmodel

import (
	"encoding/json"
	"fmt"
)

// ChannelActionParams is the typed, per-type input of a ChannelAction — a
// tagged union discriminated by the sibling ActionType (the SendingRuleConfig
// / ConversationFilterConfig pattern). Types without parameters carry the
// empty variant.
type ChannelActionParams interface {
	isChannelActionParams()
	ActionType() ChannelActionType
}

// InvitationParams — a connection request. Note is the optional free text the
// invitee sees (LinkedIn: ≤300 characters); Email is the invitee's email some
// providers require for out-of-network invitations.
type InvitationParams struct {
	Note  string `json:"note,omitempty"`
	Email string `json:"email,omitempty"`
}

func (InvitationParams) isChannelActionParams()        {}
func (InvitationParams) ActionType() ChannelActionType { return ChannelActionTypeInvitation }

// ProfileVisitParams — a notifying profile view. No parameters.
type ProfileVisitParams struct{}

func (ProfileVisitParams) isChannelActionParams()        {}
func (ProfileVisitParams) ActionType() ChannelActionType { return ChannelActionTypeProfileVisit }

// InmailParams — a paid / open-profile direct message to someone outside the
// network. Content is the message body (the same blocks a Send carries);
// Subject is the provider's message subject.
type InmailParams struct {
	Subject string         `json:"subject,omitempty"`
	Content []ContentBlock `json:"content"`
}

func (InmailParams) isChannelActionParams()        {}
func (InmailParams) ActionType() ChannelActionType { return ChannelActionTypeInmail }

// MarshalChannelActionParams encodes a variant to its stored (JSONB) bare
// shape. The discriminator lives in the sibling action_type column. nil → '{}'.
func MarshalChannelActionParams(params ChannelActionParams) (json.RawMessage, error) {
	if params == nil {
		return json.RawMessage("{}"), nil
	}
	return json.Marshal(params)
}

// DecodeChannelActionParams rebuilds the typed variant from the stored bare
// JSON plus the discriminating action type.
func DecodeChannelActionParams(actionType ChannelActionType, raw []byte) (ChannelActionParams, error) {
	switch actionType {
	case ChannelActionTypeInvitation:
		params := InvitationParams{}
		if err := unmarshalConfig(raw, &params); err != nil {
			return nil, err
		}
		return params, nil
	case ChannelActionTypeProfileVisit:
		return ProfileVisitParams{}, nil
	case ChannelActionTypeInmail:
		params := InmailParams{}
		if err := unmarshalConfig(raw, &params); err != nil {
			return nil, err
		}
		return params, nil
	default:
		return nil, fmt.Errorf("unknown channel action type %q", actionType)
	}
}
