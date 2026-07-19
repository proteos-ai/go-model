package workflowmodel

import (
	"encoding/json"
	"fmt"
	"strings"
)

// NodeType identifies a node's registered type in the node catalog. Catalog
// types follow `<package>.<node>` (e.g. "proteos-nodes-core.http-request");
// engine-native trigger types keep their legacy `trigger.<kind>` keys. The
// pairing of NodeType and WorkflowNode.TypeVersion selects one registered
// NodeDescriptor.
type NodeType string

const (
	NodeTypeTriggerCron      NodeType = "trigger.cron"
	NodeTypeTriggerManual    NodeType = "trigger.manual"
	NodeTypeTriggerWebhook   NodeType = "trigger.webhook"
	NodeTypeTriggerEvent     NodeType = "trigger.event"
	NodeTypeTriggerMessage   NodeType = "trigger.message"
	NodeTypeTriggerConnector NodeType = "trigger.connector"
	NodeTypeActionAgent      NodeType = "action.agent"
)

// Interpreter intrinsics — catalog-registered node types the interpreter
// executes itself (never dispatched to a node host).
const (
	NodeTypeWait            NodeType = "proteos-nodes-core.wait"
	NodeTypeExecuteWorkflow NodeType = "proteos-nodes-core.execute-workflow"
)

// IsTriggerType reports whether a node type is a trigger (starts an execution)
// rather than an action.
func IsTriggerType(nodeType NodeType) bool {
	return strings.HasPrefix(string(nodeType), "trigger.")
}

// IsIntrinsicType reports whether the interpreter executes this type itself
// (Temporal timer / child workflow) instead of dispatching an ExecuteNode
// activity.
func IsIntrinsicType(nodeType NodeType) bool {
	return nodeType == NodeTypeWait || nodeType == NodeTypeExecuteWorkflow
}

// SplitTypeKey splits a `<package>.<node>` type key into package and node name.
// The node name is the segment after the LAST dot (package names may not
// contain dots; node names may not either, so this is exact for well-formed
// keys and degrades safely for legacy `trigger.<kind>` keys).
func SplitTypeKey(nodeType NodeType) (pkg string, name string) {
	key := string(nodeType)
	index := strings.LastIndex(key, ".")
	if index < 0 {
		return "", key
	}
	return key[:index], key[index+1:]
}

// OnErrorPolicy selects how the interpreter routes a failed node run.
type OnErrorPolicy string

const (
	// OnErrorStop fails the execution (default; the zero value means stop).
	OnErrorStop OnErrorPolicy = "stop"
	// OnErrorContinue passes the node's input items through unchanged to its
	// main output and continues.
	OnErrorContinue OnErrorPolicy = "continue"
	// OnErrorContinueErrorOutput emits error items on the node's `error` output
	// port if it declares one; otherwise behaves like stop.
	OnErrorContinueErrorOutput OnErrorPolicy = "continue_error_output"
)

// NodeRetryPolicy is the per-node activity retry policy. MaxAttempts 1 disables
// retries (for non-idempotent nodes); zero values fall back to the descriptor's
// defaults, then the engine defaults.
type NodeRetryPolicy struct {
	MaxAttempts    int `json:"max_attempts,omitempty"`
	BackoffSeconds int `json:"backoff_seconds,omitempty"`
}

// WorkflowGraph is the v2 node graph: nodes plus port-addressed directed
// connections, stored verbatim as JSONB on the workflow_versions row. PinData
// carries editor-only per-node pinned output items (keyed by node id) used for
// partial executions; the production interpreter ignores it unless a run
// explicitly opts in.
type WorkflowGraph struct {
	Nodes       []WorkflowNode       `json:"nodes"`
	Connections []WorkflowConnection `json:"connections"`
	PinData     map[string][]Item    `json:"pin_data,omitempty"`
}

// WorkflowNode is one step in the graph. Parameters holds the RAW parameter
// values for the node's descriptor properties (Liquid expressions unresolved —
// resolution happens at ExecuteNode time, never in the interpreter).
type WorkflowNode struct {
	Id             string            `json:"id"`
	Type           NodeType          `json:"type"`
	TypeVersion    int               `json:"type_version,omitempty"`
	Name           string            `json:"name"`
	Parameters     json.RawMessage   `json:"parameters,omitempty"`
	CredentialRefs map[string]string `json:"credential_refs,omitempty"`
	Position       *NodePosition     `json:"position,omitempty"`
	IsDisabled     bool              `json:"is_disabled,omitempty"`
	Notes          string            `json:"notes,omitempty"`
	OnError        OnErrorPolicy     `json:"on_error,omitempty"`
	Retry          *NodeRetryPolicy  `json:"retry,omitempty"`
}

// NodePosition is the optional editor coordinate of a node (cosmetic; ignored by
// the executor).
type NodePosition struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// DefaultPort is the single-port key used when a connection endpoint omits the
// port (and by the v1→v2 graph upgrade).
const DefaultPort = "main"

// ErrorPort is the reserved output-port key the on_error=continue_error_output
// policy routes failed items to, when the node's descriptor declares it.
const ErrorPort = "error"

// ConnectionEndpoint addresses one side of a connection: a node by id and one of
// its ports by key. Connections are id-keyed (renaming a node never breaks the
// graph) and port-keyed (never index-keyed).
//
// Wire compat: v1 stored connections as bare node-id strings
// (`{"from": "n1", "to": "n2"}`); UnmarshalJSON upgrades that shorthand to
// `{node_id: "n1", port: "main"}` so pre-v2 graphs load without a data
// migration and normalize to the object form on next save.
type ConnectionEndpoint struct {
	NodeId string `json:"node_id"`
	Port   string `json:"port"`
}

// UnmarshalJSON accepts both the v2 object form and the v1 node-id string form.
func (endpoint *ConnectionEndpoint) UnmarshalJSON(data []byte) error {
	var nodeId string
	if err := json.Unmarshal(data, &nodeId); err == nil {
		endpoint.NodeId = nodeId
		endpoint.Port = DefaultPort
		return nil
	}
	type endpointAlias ConnectionEndpoint
	var alias endpointAlias
	if err := json.Unmarshal(data, &alias); err != nil {
		return fmt.Errorf("connection endpoint must be a node id string or {node_id, port}: %w", err)
	}
	if alias.Port == "" {
		alias.Port = DefaultPort
	}
	*endpoint = ConnectionEndpoint(alias)
	return nil
}

// WorkflowConnection is a directed edge from one node's output port to another
// node's input port.
type WorkflowConnection struct {
	From ConnectionEndpoint `json:"from"`
	To   ConnectionEndpoint `json:"to"`
}

// ── Engine-native trigger parameters ────────────────────────────────────────
//
// Trigger nodes are executed by workflow-service itself (cron schedules, webhook
// endpoint, event/message consumers, manual runs); their parameters stay typed
// Go structs. Action/transform/flow node parameters are descriptor-driven and
// remain raw JSON on the node.

// CronTriggerParams fires the workflow on a recurring cron schedule. Temporal
// owns the actual schedule; CronExpression + Timezone are handed to it verbatim.
type CronTriggerParams struct {
	CronExpression string `json:"cron_expression"`
	Timezone       string `json:"timezone"`
}

// ManualTriggerParams carries no configuration — the workflow is fired only by an
// explicit "run now" API call.
type ManualTriggerParams struct{}

// WebhookTriggerParams fires the workflow when its opaque token is POSTed to
// /workflows/v1/webhooks/:token. The token is generated server-side.
type WebhookTriggerParams struct {
	Token string `json:"token"`
}

// EventTriggerParams fires the workflow on a record-change event for EntitySlug
// matching one of Verbs (created|updated|deleted), consumed off the org event bus.
type EventTriggerParams struct {
	EntitySlug string   `json:"entity_slug"`
	Verbs      []string `json:"verbs"`
}

// MessageTriggerParams fires the workflow on a conversation-service message
// event. Direction defaults to inbound (outbound events are usually the
// workflow's own sends). Channels/ConnectionId narrow the match; EventTypes
// defaults to ["message.created"].
type MessageTriggerParams struct {
	Direction    string   `json:"direction,omitempty"`
	Channels     []string `json:"channels,omitempty"`
	ConnectionId string   `json:"connection_id,omitempty"`
	EventTypes   []string `json:"event_types,omitempty"`
}

// ConnectorTriggerParams fires the workflow on a connector-service
// sync.item_changed event — a pre-built connector's sync loop observed a
// changed item on the remote system (a Google Calendar event edited, …).
// ConnectorKey/ConnectionId narrow the match (all when empty). The changed
// item becomes the trigger item.
type ConnectorTriggerParams struct {
	ConnectorKey string `json:"connector_key,omitempty"`
	ConnectionId string `json:"connection_id,omitempty"`
}

// AgentActionParams runs an agent (in agent-service) and waits for it to finish.
// Kickoff selects how the agent is started — a plain message turn, or an
// Anthropic-graded outcome loop.
type AgentActionParams struct {
	AgentKey string       `json:"agent_key"`
	Kickoff  AgentKickoff `json:"kickoff"`
}

// KickoffType discriminates how an agent node starts the agent.
type KickoffType string

const (
	KickoffTypeMessage KickoffType = "message"
	KickoffTypeOutcome KickoffType = "outcome"
)

// KickoffSource selects where a kickoff draws its instruction text from: typed
// inline ("manual"), or resolved from a reusable agent-service Prompt by key
// ("prompt") at execution time. Empty defaults to manual (backward compatible).
type KickoffSource string

const (
	KickoffSourceManual KickoffSource = "manual"
	KickoffSourcePrompt KickoffSource = "prompt"
)

// AgentKickoff is the tagged union of agent-node kickoffs. Exactly one of Message
// / Outcome is populated, selected by Type.
type AgentKickoff struct {
	Type    KickoffType     `json:"type"`
	Message *MessageKickoff `json:"message,omitempty"`
	Outcome *OutcomeKickoff `json:"outcome,omitempty"`
}

// MessageKickoff starts the agent with a user.message. When Source is "prompt"
// the message text is resolved from the Prompt named by PromptKey at execution
// time; otherwise Content is the raw content blocks array forwarded verbatim as
// the agent-service user.message payload.
type MessageKickoff struct {
	Source    KickoffSource   `json:"source,omitempty"`
	Content   json.RawMessage `json:"content,omitempty"`
	PromptKey string          `json:"prompt_key,omitempty"`
}

// OutcomeKickoff starts the agent with a user.define_outcome — Anthropic Managed
// Agents runs the grade/iterate loop server-side against Rubric. The Description
// and Rubric are independently sourced: each is typed inline ("manual") or
// resolved from a reusable Prompt by its *PromptKey ("prompt") at execution time.
type OutcomeKickoff struct {
	DescriptionSource    KickoffSource `json:"description_source,omitempty"`
	Description          string        `json:"description,omitempty"`
	DescriptionPromptKey string        `json:"description_prompt_key,omitempty"`
	RubricSource         KickoffSource `json:"rubric_source,omitempty"`
	Rubric               OutcomeRubric `json:"rubric"`
	RubricPromptKey      string        `json:"rubric_prompt_key,omitempty"`
	MaxIterations        *int          `json:"max_iterations,omitempty"`
}

// OutcomeRubric is the grading rubric for an outcome kickoff: inline text or a
// reference to an uploaded file.
type OutcomeRubric struct {
	Type    string `json:"type"` // "text" | "file"
	Content string `json:"content,omitempty"`
	FileId  string `json:"file_id,omitempty"`
}

// DecodeAgentActionParams decodes an action.agent node's raw parameters.
func DecodeAgentActionParams(raw json.RawMessage) (AgentActionParams, error) {
	var params AgentActionParams
	if err := json.Unmarshal(raw, &params); err != nil {
		return AgentActionParams{}, err
	}
	return params, nil
}

// DecodeCronTriggerParams decodes a trigger.cron node's raw parameters.
func DecodeCronTriggerParams(raw json.RawMessage) (CronTriggerParams, error) {
	var params CronTriggerParams
	if err := json.Unmarshal(raw, &params); err != nil {
		return CronTriggerParams{}, err
	}
	return params, nil
}

// DecodeWebhookTriggerParams decodes a trigger.webhook node's raw parameters.
func DecodeWebhookTriggerParams(raw json.RawMessage) (WebhookTriggerParams, error) {
	var params WebhookTriggerParams
	if err := json.Unmarshal(raw, &params); err != nil {
		return WebhookTriggerParams{}, err
	}
	return params, nil
}

// DecodeEventTriggerParams decodes a trigger.event node's raw parameters.
func DecodeEventTriggerParams(raw json.RawMessage) (EventTriggerParams, error) {
	var params EventTriggerParams
	if err := json.Unmarshal(raw, &params); err != nil {
		return EventTriggerParams{}, err
	}
	return params, nil
}

// DecodeMessageTriggerParams decodes a trigger.message node's raw parameters.
func DecodeMessageTriggerParams(raw json.RawMessage) (MessageTriggerParams, error) {
	var params MessageTriggerParams
	if err := json.Unmarshal(raw, &params); err != nil {
		return MessageTriggerParams{}, err
	}
	return params, nil
}

// DecodeConnectorTriggerParams decodes a trigger.connector node's raw parameters.
func DecodeConnectorTriggerParams(raw json.RawMessage) (ConnectorTriggerParams, error) {
	var params ConnectorTriggerParams
	if err := json.Unmarshal(raw, &params); err != nil {
		return ConnectorTriggerParams{}, err
	}
	return params, nil
}
