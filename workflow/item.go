package workflowmodel

import (
	"encoding/json"
	"fmt"
)

// Item is the unit of data flowing between workflow nodes (n8n 1:1 — see
// docs/workflow-engine/node-spec.md §Items model). Json is the payload; Binary
// holds storage-service file references keyed by property name (binary content
// never travels inline); PairedItem records which input item produced this item.
type Item struct {
	Json       map[string]any       `json:"json"`
	Binary     map[string]BinaryRef `json:"binary,omitempty"`
	PairedItem *PairedItem          `json:"paired_item,omitempty"`
}

// BinaryRef is a storage-service file reference carried by an item.
type BinaryRef struct {
	FileId    string `json:"file_id"`
	MimeType  string `json:"mime_type,omitempty"`
	FileName  string `json:"file_name,omitempty"`
	SizeBytes int64  `json:"size_bytes,omitempty"`
}

// PairedItem is item lineage: the index of the input item (and the input port it
// arrived on) that produced an output item. It drives the editor's hover-pairing
// and cross-node expression resolution. The wire shape allows an int shorthand
// (`"paired_item": 0`) equivalent to `{"item": 0}`.
type PairedItem struct {
	Item  int    `json:"item"`
	Input string `json:"input,omitempty"`
}

// UnmarshalJSON accepts both the object form and the int shorthand.
func (paired *PairedItem) UnmarshalJSON(data []byte) error {
	var index int
	if err := json.Unmarshal(data, &index); err == nil {
		paired.Item = index
		paired.Input = ""
		return nil
	}
	type pairedItemAlias PairedItem
	var alias pairedItemAlias
	if err := json.Unmarshal(data, &alias); err != nil {
		return fmt.Errorf("paired_item must be an int or an object: %w", err)
	}
	*paired = PairedItem(alias)
	return nil
}

// PayloadRef addresses one stored item array in the claim-check payload store
// (table workflow_node_execution_payloads), keyed by the producing node run and
// output port.
type PayloadRef struct {
	ExecutionId string `json:"execution_id"`
	NodeId      string `json:"node_id"`
	RunIndex    int    `json:"run_index"`
	Port        string `json:"port"`
}
