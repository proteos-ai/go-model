package workflowmodel

// ErrorItems turns one input port's items of a failed node run into the items
// the on_error=continue_error_output policy emits on the node's `error`
// output: one per input item, carrying the input's json with the failure
// attached under `error` (code, message, details, is_user_error — the
// ExecutionError shape 1:1) and paired to the input index on inputPort (so a
// multi-input node like merge keeps every item's lineage). Binary refs travel
// along. An input `error` key, if any, is overwritten.
func ErrorItems(inputPort string, items []Item, failure ExecutionError) []Item {
	errorItems := make([]Item, 0, len(items))
	for index, item := range items {
		payload := make(map[string]any, len(item.Json)+1)
		for key, value := range item.Json {
			payload[key] = value
		}
		payload["error"] = failure.asItemValue()
		errorItems = append(errorItems, Item{
			Json:       payload,
			Binary:     item.Binary,
			PairedItem: &PairedItem{Item: index, Input: inputPort},
		})
	}
	return errorItems
}

// asItemValue is the ExecutionError as a snake_case item value (the same keys
// its JSON encoding uses, so Liquid reads `json.error.message` either way).
func (failure ExecutionError) asItemValue() map[string]any {
	value := map[string]any{
		"code":          failure.Code,
		"message":       failure.Message,
		"is_user_error": failure.IsUserError,
	}
	if len(failure.Details) > 0 {
		value["details"] = failure.Details
	}
	return value
}
