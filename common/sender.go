package common

import (
	"encoding/json"
	"fmt"

	"golang.org/x/exp/slices"
)

type Sender string

const (
	SenderUser   Sender = "user"
	SenderAva    Sender = "ava"
	SenderSystem Sender = "system"
)

var Senders = []Sender{
	SenderUser,
	SenderSystem,
	SenderAva,
}

func (Sender) Enum() []interface{} {
	enums := []interface{}{}
	for _, element := range Senders {
		enums = append(enums, element)
	}
	return enums
}

func (sender *Sender) UnmarshalJSON(byteArray []byte) error {
	str := string(byteArray)
	if str == "null" {
		*sender = ""
		return nil
	}

	type _Sender Sender
	var stringValue *_Sender = (*_Sender)(sender)
	err := json.Unmarshal(byteArray, &stringValue)

	if err != nil {
		return err
	}

	if slices.Contains(Senders, *sender) {
		return nil
	}

	return fmt.Errorf("invalid sender: %s", *stringValue)
}
