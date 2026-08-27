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

	// Scoped twins. The grant a role holds decides HOW MUCH it sees: `read`
	// means every record of the entity in the org — exactly what it means today
	// — while `read_scoped` means only those whose access grants intersect the
	// caller's principals (owned by them, owned by a team they are in, shared
	// with them, or carrying an org-wide grant).
	//
	// Because `read` keeps its current meaning, every existing role grant works
	// unchanged and nothing migrates.
	//
	// The verbs stay independent, so combinations are meaningful on their own:
	// `deals:read` + `deals:write_scoped` reads as "see every deal, edit only
	// your own" without any new concept.
	PermissionReadScoped   Permission = "read_scoped"
	PermissionWriteScoped  Permission = "write_scoped"
	PermissionDeleteScoped Permission = "delete_scoped"
)

var Permissions = []Permission{
	PermissionRead,
	PermissionWrite,
	PermissionDelete,
	PermissionReadScoped,
	PermissionWriteScoped,
	PermissionDeleteScoped,
}

// ScopedTwin returns the scoped counterpart of an unscoped permission, and
// whether one exists. It is the single place the pairing is expressed, so the
// stage-1 gate cannot drift from the enum.
func ScopedTwin(permission Permission) (Permission, bool) {
	switch permission {
	case PermissionRead:
		return PermissionReadScoped, true
	case PermissionWrite:
		return PermissionWriteScoped, true
	case PermissionDelete:
		return PermissionDeleteScoped, true
	default:
		return "", false
	}
}

// IsScoped reports whether a permission is a scoped twin.
func IsScoped(permission Permission) bool {
	switch permission {
	case PermissionReadScoped, PermissionWriteScoped, PermissionDeleteScoped:
		return true
	default:
		return false
	}
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
