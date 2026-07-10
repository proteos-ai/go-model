package workflowmodel

import (
	"encoding/json"
	"time"
)

// NodeRuntime selects which runtime executes a node type. `go` nodes are
// first-party, compiled into the node host; `wasm` nodes are user-defined
// sandboxed bundles (Phase 4); `intrinsic` types are executed by the
// interpreter itself (wait, execute-workflow) or by workflow-service's trigger
// machinery — they appear in the catalog for the editor but are never
// dispatched to a node host.
type NodeRuntime string

const (
	NodeRuntimeGo        NodeRuntime = "go"
	NodeRuntimeWasm      NodeRuntime = "wasm"
	NodeRuntimeIntrinsic NodeRuntime = "intrinsic"
)

// NodeGroup is the coarse catalog grouping of a node type.
type NodeGroup string

const (
	NodeGroupTrigger   NodeGroup = "trigger"
	NodeGroupAction    NodeGroup = "action"
	NodeGroupTransform NodeGroup = "transform"
	NodeGroupFlow      NodeGroup = "flow"
)

// PropertyType is the control type of one descriptor property (n8n
// INodeProperties paradigm, snake_case enum values — node-spec.md §parameter
// schema).
type PropertyType string

const (
	PropertyTypeString            PropertyType = "string"
	PropertyTypeNumber            PropertyType = "number"
	PropertyTypeBoolean           PropertyType = "boolean"
	PropertyTypeOptions           PropertyType = "options"
	PropertyTypeMultiOptions      PropertyType = "multi_options"
	PropertyTypeJson              PropertyType = "json"
	PropertyTypeDateTime          PropertyType = "date_time"
	PropertyTypeColor             PropertyType = "color"
	PropertyTypeNotice            PropertyType = "notice"
	PropertyTypeHidden            PropertyType = "hidden"
	PropertyTypeCollection        PropertyType = "collection"
	PropertyTypeFixedCollection   PropertyType = "fixed_collection"
	PropertyTypeCredentialsSelect PropertyType = "credentials_select"
	PropertyTypeButton            PropertyType = "button"
	// Wave 2 (editor support lands with Phase 3):
	PropertyTypeResourceLocator      PropertyType = "resource_locator"
	PropertyTypeResourceMapper       PropertyType = "resource_mapper"
	PropertyTypeFilter               PropertyType = "filter"
	PropertyTypeAssignmentCollection PropertyType = "assignment_collection"
)

// NodeDescriptor is the full, self-describing contract of one node type: it
// drives the editor's picker and parameter form, backend validation, and
// dispatch (runtime + task queue). JSON envelope keys are snake_case; property
// `name` values are opaque verbatim strings (never case-converted).
type NodeDescriptor struct {
	Type             NodeType         `json:"type"`
	DisplayName      string           `json:"display_name"`
	Description      string           `json:"description,omitempty"`
	Icon             string           `json:"icon,omitempty"` // Lucide icon name, PascalCase
	Group            NodeGroup        `json:"group"`
	Version          int              `json:"version"`
	Runtime          NodeRuntime      `json:"runtime"`
	TaskQueue        string           `json:"task_queue,omitempty"`
	Inputs           []PortSpec       `json:"inputs"`
	Outputs          []PortSpec       `json:"outputs"`
	Properties       []NodeProperty   `json:"properties"`
	Credentials      []CredentialSpec `json:"credentials,omitempty"`
	Methods          []string         `json:"methods,omitempty"` // "<kind>:<method_name>"
	Routing          json.RawMessage  `json:"routing,omitempty"` // declarative spec (Phase 4)
	DocumentationUrl string           `json:"documentation_url,omitempty"`
	Subtitle         string           `json:"subtitle,omitempty"` // Liquid over parameters, shown on canvas
}

// PortSpec declares one named input or output port. Port keys are stable
// identifiers referenced by connections (never indexes).
type PortSpec struct {
	Key            string `json:"key"`
	DisplayName    string `json:"display_name,omitempty"`
	IsRequired     bool   `json:"is_required,omitempty"`
	MaxConnections int    `json:"max_connections,omitempty"`
}

// NodeProperty is one parameter in a descriptor's schema. Name is the key under
// node.parameters. Every value-bearing property accepts a Liquid expression
// unless NoDataExpression is true.
type NodeProperty struct {
	DisplayName      string           `json:"display_name"`
	Name             string           `json:"name"`
	Type             PropertyType     `json:"type"`
	Default          any              `json:"default,omitempty"`
	Description      string           `json:"description,omitempty"`
	Placeholder      string           `json:"placeholder,omitempty"`
	Hint             string           `json:"hint,omitempty"`
	IsRequired       bool             `json:"is_required,omitempty"`
	NoDataExpression bool             `json:"no_data_expression,omitempty"`
	Options          []PropertyOption `json:"options,omitempty"`
	DisplayOptions   *DisplayOptions  `json:"display_options,omitempty"`
	TypeOptions      map[string]any   `json:"type_options,omitempty"`
}

// PropertyOption is one entry of an options/multi_options list, one repeatable
// sub-property group of a collection, or one named group of a fixed_collection
// (in which case Values carries the group's fixed sub-fields).
type PropertyOption struct {
	Name        string         `json:"name"`
	Value       any            `json:"value,omitempty"`
	DisplayName string         `json:"display_name,omitempty"`
	Description string         `json:"description,omitempty"`
	Values      []NodeProperty `json:"values,omitempty"`
}

// DisplayOptions gates a property's visibility on sibling parameter values. The
// editor re-evaluates reactively; the backend applies the same rule during
// validation (hidden properties are ignored).
type DisplayOptions struct {
	Show map[string][]any `json:"show,omitempty"`
	Hide map[string][]any `json:"hide,omitempty"`
}

// CredentialSpec declares that a node type uses a credential of the given
// credential-type slug.
type CredentialSpec struct {
	Type           string          `json:"type"`
	IsRequired     bool            `json:"is_required,omitempty"`
	DisplayOptions *DisplayOptions `json:"display_options,omitempty"`
}

// NodeTypeStatus is the registry lifecycle of one (type, version) row. Hosts
// self-register on boot; platform types absent from a sync are marked inactive
// (hidden from the picker; pinned instances keep resolving).
type NodeTypeStatus string

const (
	NodeTypeStatusActive   NodeTypeStatus = "active"
	NodeTypeStatusInactive NodeTypeStatus = "inactive"
)

// WorkflowNodeType is one registry row: a registered (package, name, version)
// with its full descriptor. OrgId is empty for platform-global types (internal
// Go nodes, engine natives) and set for org-installed (wasm) types, which are
// visible only to that org.
type WorkflowNodeType struct {
	Package    string         `json:"package"`
	Name       string         `json:"name"`
	Version    int            `json:"version"`
	Runtime    NodeRuntime    `json:"runtime"`
	TaskQueue  string         `json:"task_queue,omitempty"`
	Descriptor NodeDescriptor `json:"descriptor"`
	Status     NodeTypeStatus `json:"status"`
	OrgId      string         `json:"org_id,omitempty"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
}

// TypeKey reconstructs the `<package>.<node>` catalog key of a registry row.
func (nodeType WorkflowNodeType) TypeKey() NodeType {
	if nodeType.Package == "" {
		return NodeType(nodeType.Name)
	}
	return NodeType(nodeType.Package + "." + nodeType.Name)
}
