package workflowapi

import (
	workflowmodel "go.proteos.ai/model/workflow"

	"go.proteos.ai/model/common"
)

type GetManyExecutionsQuery struct {
	WorkflowKey *string `json:"workflow_key" form:"workflow_key" db:"workflow_key"`
	Status      *string `json:"status" form:"status" db:"status"`
	common.Pagination
	common.Sorting
}

type GetManyExecutionsResponse struct {
	Meta common.ResponseMeta               `json:"meta"`
	Data []workflowmodel.WorkflowExecution `json:"data"`
}

// GetExecutionResponse is the execution detail: the header row plus its
// append-only node executions (ordered by started_at, then run_index).
type GetExecutionResponse struct {
	Execution      workflowmodel.WorkflowExecution `json:"execution"`
	NodeExecutions []workflowmodel.NodeExecution   `json:"node_executions"`
}

// NodeItemsSide selects which side of a node run an item window reads:
// the node's own output ports (default) or the mirrored "in:<port>" copy of
// what fed it.
const (
	NodeItemsSideOutput = "output"
	NodeItemsSideInput  = "input"
)

// GetNodeExecutionItemsQuery windows into one node run's stored items on one
// port — the NDV input/output panels page through large item sets with it.
type GetNodeExecutionItemsQuery struct {
	Port   string `json:"port" form:"port"`
	Side   string `json:"side" form:"side" validate:"omitempty,oneof=input output"`
	Offset int    `json:"offset" form:"offset"`
	Limit  int    `json:"limit" form:"limit"`
}

// GetNodeExecutionItemsResponse is one window of items plus the true total.
type GetNodeExecutionItemsResponse struct {
	Items      []workflowmodel.Item `json:"items"`
	ItemsTotal int                  `json:"items_total"`
	Offset     int                  `json:"offset"`
	Limit      int                  `json:"limit"`
}
