package functionsmodel

// ActionScope distinguishes record-scoped actions (invoked against a
// specific record) from global actions (invoked without a record
// context). Backs the `actions.scope` column.
type ActionScope string

const (
	ActionScopeEntity ActionScope = "entity"
	ActionScopeGlobal ActionScope = "global"
	// ActionScopeConnectorMethod marks the wasm behind one custom-connector
	// method. It is the ONLY connector-specific fact stored on the action:
	// the (connector, method) → action binding lives in the connector
	// manifest's methods (action_slug), not here. This scope excludes the row
	// from the action catalog and rejects it on the user action-invoke routes;
	// it is dispatched only by connector-service via
	// `POST /functions/v1/connector-methods/:slug/invoke`.
	ActionScopeConnectorMethod ActionScope = "connector_method"
)

var ActionScopes = []ActionScope{
	ActionScopeEntity,
	ActionScopeGlobal,
	ActionScopeConnectorMethod,
}

func (ActionScope) Enum() []interface{} {
	out := make([]interface{}, len(ActionScopes))
	for i, s := range ActionScopes {
		out[i] = s
	}
	return out
}
