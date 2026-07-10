package common

type SortDirection string

const (
	SortDirectionAsc  SortDirection = "asc"
	SortDirectionDesc SortDirection = "desc"
)

var SortDirections = []SortDirection{
	SortDirectionAsc,
	SortDirectionDesc,
}

func (SortDirection) Enum() []interface{} {
	enums := []interface{}{}
	for _, element := range SortDirections {
		enums = append(enums, element)
	}
	return enums
}

// Sorting carries the sort key + direction of a list request. Tag rationale
// mirrors Pagination: `json` for BindQueryAsJson + responses, `form` for gin
// ShouldBindQuery binds.
type Sorting struct {
	SortBy *string        `json:"sort_by" form:"sort_by"`
	Sort   *SortDirection `json:"sort_direction" form:"sort_direction" validate:"omitempty,oneof=asc desc"`
}
