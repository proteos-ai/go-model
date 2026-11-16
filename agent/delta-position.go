package agentmodel

import (
	"fmt"
	"strconv"
	"strings"
)

// DeltaPosition is a place on a session's live delta channel. Seq is the log seq
// of the latest whole event at that place; Index counts the chunk frames sent
// since that event. A whole event sits at Index 0, so its position is its seq, and
// the chunk frames written after it take Index 1, 2, … until the next whole event.
//
// Written "<seq>-<index>", it is the cursor of the delta channel: the id of each
// SSE frame, and what a client sends back to resume.
type DeltaPosition struct {
	Seq   int64
	Index int64
}

func (position DeltaPosition) String() string {
	return strconv.FormatInt(position.Seq, 10) + "-" + strconv.FormatInt(position.Index, 10)
}

// IsWholeEvent reports whether the position holds a whole event of the log, as
// opposed to a chunk frame.
func (position DeltaPosition) IsWholeEvent() bool {
	return position.Index == 0
}

// ParseDeltaPosition reads a cursor: "<seq>-<index>", or a bare "<seq>", which
// is the position of that whole event.
func ParseDeltaPosition(raw string) (DeltaPosition, error) {
	seqPart, indexPart, hasIndex := strings.Cut(raw, "-")
	seq, err := strconv.ParseInt(seqPart, 10, 64)
	if err != nil || seq < 0 {
		return DeltaPosition{}, fmt.Errorf("delta position %q: seq is not a non-negative integer", raw)
	}
	if !hasIndex {
		return DeltaPosition{Seq: seq}, nil
	}
	index, err := strconv.ParseInt(indexPart, 10, 64)
	if err != nil || index < 0 {
		return DeltaPosition{}, fmt.Errorf("delta position %q: index is not a non-negative integer", raw)
	}
	return DeltaPosition{Seq: seq, Index: index}, nil
}

// DeltaEntry is one entry of the live delta channel: a whole event of the log or
// a chunk frame (see DeltaEventTypes), with its Position on the channel.
type DeltaEntry struct {
	Position DeltaPosition
	Event    SessionEvent
}
