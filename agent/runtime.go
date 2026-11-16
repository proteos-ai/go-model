package agentmodel

import (
	"encoding/json"
	"fmt"

	"golang.org/x/exp/slices"
)

// Runtime says how an agent's sessions run: where the agent loop executes and
// what a turn survives. It names the kind of runtime, not the vendor behind it;
// which provider serves each kind is the agent-service adapters' choice. It is
// configuration on the Agent, and each Session keeps the runtime it was created
// with, so changing an agent's runtime never moves its existing sessions.
type Runtime string

const (
	// RuntimeDurable runs the loop outside agent-service's request path, on a
	// provider that keeps the turn going through restarts and deploys and gives the
	// session a sandbox. Built for long-running work. Anthropic Managed Agents is
	// the provider today.
	RuntimeDurable Runtime = "durable"
	// RuntimeEmbedded runs the loop inside the agent-service process against a
	// model provider's streaming API. The session's event log is the conversation,
	// so the session itself is stored like any other; replies stream while the model
	// writes them, there is no sandbox, and a turn in flight ends if the process does.
	// "Embedded" is about the loop living in agent-service, with nothing keeping the
	// turn outside it; it says nothing about where the chat is shown.
	RuntimeEmbedded Runtime = "embedded"
)

// DefaultRuntime is the runtime of an agent that names none, and of every agent
// and session created before runtimes existed.
const DefaultRuntime = RuntimeDurable

var Runtimes = []Runtime{RuntimeDurable, RuntimeEmbedded}

func (Runtime) Enum() []interface{} {
	enums := []interface{}{}
	for _, element := range Runtimes {
		enums = append(enums, element)
	}
	return enums
}

// UnmarshalJSON accepts a known runtime, and null or "" for "not given".
func (runtime *Runtime) UnmarshalJSON(byteArray []byte) error {
	if string(byteArray) == "null" {
		*runtime = ""
		return nil
	}

	type _Runtime Runtime
	value := (*_Runtime)(runtime)
	if err := json.Unmarshal(byteArray, value); err != nil {
		return err
	}

	if *runtime == "" || slices.Contains(Runtimes, *runtime) {
		return nil
	}

	return fmt.Errorf("invalid runtime: %s", *runtime)
}
