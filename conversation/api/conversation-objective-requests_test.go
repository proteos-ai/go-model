package conversationapi

import (
	"encoding/json"
	"strings"
	"testing"

	"go.proteos.ai/model/common"
	conversationmodel "go.proteos.ai/model/conversation"
)

// An unset Optional measure must stay out of an encoded patch (an SDK
// update of `name` alone must not clear the stored measure); an explicit
// null must still encode as null.
func TestUpdateConversationObjectiveDefinitionRequestOmitsUnsetMeasure(t *testing.T) {
	name := "Confirm budget"
	encoded, err := json.Marshal(UpdateConversationObjectiveDefinitionRequest{Name: &name})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(encoded), "measure") {
		t.Fatalf("unset measure must be omitted: %s", encoded)
	}
	cleared, err := json.Marshal(UpdateConversationObjectiveDefinitionRequest{
		Measure: common.Optional[conversationmodel.ConversationObjectiveMeasure]{Present: true, Value: nil}})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(cleared), `"measure":null`) {
		t.Fatalf("explicit null must encode as null: %s", cleared)
	}
	var decoded UpdateConversationObjectiveDefinitionRequest
	if err := json.Unmarshal([]byte(`{"name":"x"}`), &decoded); err != nil || decoded.Measure.Present {
		t.Fatalf("absent key must decode as not present: %+v %v", decoded.Measure, err)
	}
}
