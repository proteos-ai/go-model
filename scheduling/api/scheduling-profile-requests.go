package schedulingapi

import conversationmodel "go.proteos.ai/model/conversation"

// PutSchedulingProfileRequest replaces a user's scheduling profile (PUT
// semantics: every field is the new value). WeeklyHours are {day, from,
// until} windows ("HH:MM", until > from, "24:00" allowed as an end) read in
// Timezone; an empty list means bookable any time.
type PutSchedulingProfileRequest struct {
	Timezone    string                        `json:"timezone" validate:"required,max=64"`
	WeeklyHours []conversationmodel.WindowDay `json:"weekly_hours" validate:"max=64,dive"`
	DisplayName string                        `json:"display_name" validate:"max=120"`
	Headline    string                        `json:"headline" validate:"max=280"`
	IsBookable  bool                          `json:"is_bookable"`
}
