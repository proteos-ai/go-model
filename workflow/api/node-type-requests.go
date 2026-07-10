package workflowapi

import (
	workflowmodel "go.proteos.ai/model/workflow"
)

// SyncNodeTypesRequest is the node host's boot-time self-registration payload:
// its complete descriptor set. workflow-service upserts every descriptor and
// marks platform-global types absent from the set inactive.
type SyncNodeTypesRequest struct {
	Descriptors []workflowmodel.NodeDescriptor `json:"descriptors" validate:"required"`
}

// SyncNodeTypesResponse reports the sync outcome.
type SyncNodeTypesResponse struct {
	Upserted    int `json:"upserted"`
	Deactivated int `json:"deactivated"`
}

// GetNodeTypesResponse is the editor catalog: every node type visible to the
// caller's org (platform-global + org-installed), latest and prior versions.
type GetNodeTypesResponse struct {
	Data []workflowmodel.WorkflowNodeType `json:"data"`
}

// InvokeNodeMethodRequest proxies an editor dynamic-method call
// (load_options / list_search / resource_mapper_fields / credential_test) to
// the node host. Parameters carries the node's current (partial) parameter
// values from the editor.
type InvokeNodeMethodRequest struct {
	TypeVersion     int               `json:"type_version,omitempty"`
	Parameters      map[string]any    `json:"parameters,omitempty"`
	CredentialRefs  map[string]string `json:"credential_refs,omitempty"`
	Filter          string            `json:"filter,omitempty"`
	PaginationToken string            `json:"pagination_token,omitempty"`
}
