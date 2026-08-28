package functionsmodel

// ActionScope distinguishes record-scoped actions (invoked against a
// specific record) from global actions (invoked without a record
// context). Backs the `actions.scope` column.
type ActionScope string

const (
	ActionScopeEntity ActionScope = "entity"
	// ActionScopeEntityBatch marks a BATCH action: invoked against a SET of
	// records of one entity via
	// `POST /api/v1/entities/:entity/actions/:slug/invoke` with a
	// `{record_ids, params}` body. The guest receives the whole id array in
	// ONE invocation (never one dispatch per record) so it can amortise the
	// work — one query, one external call, one aggregate result. Requires an
	// entity slug, exactly like ActionScopeEntity.
	ActionScopeEntityBatch ActionScope = "entity_batch"
	ActionScopeGlobal      ActionScope = "global"
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
	ActionScopeEntityBatch,
	ActionScopeGlobal,
	ActionScopeConnectorMethod,
}

// RequiresEntity reports whether a scope must carry an entity slug.
func (scope ActionScope) RequiresEntity() bool {
	return scope == ActionScopeEntity || scope == ActionScopeEntityBatch
}

func (ActionScope) Enum() []interface{} {
	out := make([]interface{}, len(ActionScopes))
	for i, s := range ActionScopes {
		out[i] = s
	}
	return out
}
