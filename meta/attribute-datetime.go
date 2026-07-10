package metamodel

// DatetimeFormat represents datetime storage/display behavior
// Enterprise pattern: date-time stored in UTC, date/time/duration are naive
type DatetimeFormat string

const (
	// DatetimeFormatDateTime - Full datetime, stored in UTC, displayed in user's timezone
	DatetimeFormatDateTime DatetimeFormat = "date-time"
	// DatetimeFormatDate - Naive date (no timezone), e.g., "2024-01-15"
	DatetimeFormatDate DatetimeFormat = "date"
	// DatetimeFormatTime - Naive time (no timezone), e.g., "14:30:00"
	DatetimeFormatTime DatetimeFormat = "time"
	// DatetimeFormatDuration - ISO 8601 duration, e.g., "P1DT2H30M"
	DatetimeFormatDuration DatetimeFormat = "duration"
)

// DatetimeAttributeMeta holds datetime-specific validation options
type DatetimeAttributeMeta struct {
	Format  DatetimeFormat `json:"format"`            // Required
	Minimum *string        `json:"minimum,omitempty"` // ISO 8601
	Maximum *string        `json:"maximum,omitempty"` // ISO 8601
}
