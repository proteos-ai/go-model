package agentmodel

import (
	"encoding/json"
	"fmt"
	"time"

	"go.proteos.ai/model/common"

	"golang.org/x/exp/slices"
)

// ToolsetKind distinguishes the two toolset origins: platform toolsets are the
// hardcoded groups of the platform MCP server (mcp-service) — read-only, defined
// in code, one per server mount; custom toolsets are org-authored groups of the
// org's own Tool rows.
type ToolsetKind string

const (
	ToolsetKindPlatform ToolsetKind = "platform"
	ToolsetKindCustom   ToolsetKind = "custom"
)

var ToolsetKinds = []ToolsetKind{ToolsetKindPlatform, ToolsetKindCustom}

func (ToolsetKind) Enum() []interface{} {
	enums := []interface{}{}
	for _, element := range ToolsetKinds {
		enums = append(enums, element)
	}
	return enums
}

func (toolsetKind *ToolsetKind) UnmarshalJSON(byteArray []byte) error {
	if string(byteArray) == "null" {
		*toolsetKind = ""
		return nil
	}

	type _ToolsetKind ToolsetKind
	value := (*_ToolsetKind)(toolsetKind)
	if err := json.Unmarshal(byteArray, value); err != nil {
		return err
	}

	if slices.Contains(ToolsetKinds, *toolsetKind) {
		return nil
	}

	return fmt.Errorf("invalid toolset kind: %s", *toolsetKind)
}

// Toolset is a named group of tools an agent can attach as one unit
// (Agent.Toolsets lists toolset keys). Platform and custom toolsets share one
// key namespace: the platform keys are reserved, so a custom toolset can never
// shadow one. For platform toolsets Tools is empty on the wire — the members
// live in mcp-service and are listed via the toolset's tools endpoint; for
// custom toolsets Tools carries the member Tool keys. Keyed by (org_id, key).
type Toolset struct {
	OrgId       string         `json:"org_id"`
	Key         string         `json:"key" sortable:""`
	Name        string         `json:"name" sortable:""`
	ModuleSlug  string         `json:"module_slug" sortable:""`
	Description string         `json:"description"`
	Kind        ToolsetKind    `json:"kind"`
	Tools       []string       `json:"tools"`
	Version     int            `json:"version"`
	CreatedAt   time.Time      `json:"created_at" sortable:""`
	CreatedBy   common.UserRef `json:"created_by"`
	UpdatedAt   time.Time      `json:"updated_at" sortable:""`
	UpdatedBy   common.UserRef `json:"updated_by"`
}

// ToolsetToolSummary describes one tool of a toolset for pickers: the wire
// name the model calls plus display metadata. Used by the toolset tools
// endpoint for both kinds (platform members proxied from mcp-service, custom
// members summarized from the org's Tool rows).
type ToolsetToolSummary struct {
	Name        string `json:"name"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
}
