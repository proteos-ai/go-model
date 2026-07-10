package accountmodel

import (
	"encoding/json"
	"fmt"

	"golang.org/x/exp/slices"
)

type Permission string

const (
	PermissionRead   Permission = "read"
	PermissionWrite  Permission = "write"
	PermissionDelete Permission = "delete"
)

var Permissions = []Permission{
	PermissionRead,
	PermissionWrite,
	PermissionDelete,
}

func (Permission) Enum() []interface{} {
	enums := []interface{}{}
	for _, element := range Permissions {
		enums = append(enums, element)
	}
	return enums
}

func (permission *Permission) UnmarshalJSON(byteArray []byte) error {
	str := string(byteArray)
	if str == "null" {
		*permission = ""
		return nil
	}

	type _Permission Permission
	var stringValue *_Permission = (*_Permission)(permission)
	err := json.Unmarshal(byteArray, &stringValue)

	if err != nil {
		return err
	}

	if slices.Contains(Permissions, *permission) {
		return nil
	}

	return fmt.Errorf("invalid permission: %s", *stringValue)
}
