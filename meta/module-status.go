package metamodel

import (
	"encoding/json"
	"fmt"

	"golang.org/x/exp/slices"
)

type ModuleStatus string

const (
	ModuleStatusPending      ModuleStatus = "pending"
	ModuleStatusDeploying    ModuleStatus = "deploying"
	ModuleStatusActive       ModuleStatus = "active"
	ModuleStatusFailed       ModuleStatus = "failed"
	ModuleStatusDeactivating ModuleStatus = "deactivating"
	ModuleStatusInactive     ModuleStatus = "inactive"
)

var ModuleStatuses = []ModuleStatus{
	ModuleStatusDeploying,
	ModuleStatusActive,
	ModuleStatusFailed,
	ModuleStatusInactive,
	ModuleStatusPending,
	ModuleStatusDeactivating,
}

func (ModuleStatus) Enum() []interface{} {
	enums := []interface{}{}
	for _, element := range ModuleStatuses {
		enums = append(enums, element)
	}
	return enums
}

func (moduleStatus *ModuleStatus) UnmarshalJSON(byteArray []byte) error {
	str := string(byteArray)
	if str == "null" {
		*moduleStatus = ""
		return nil
	}

	type _ModuleStatus ModuleStatus
	var stringValue *_ModuleStatus = (*_ModuleStatus)(moduleStatus)
	err := json.Unmarshal(byteArray, &stringValue)

	if err != nil {
		return err
	}

	if slices.Contains(ModuleStatuses, *moduleStatus) {
		return nil
	}

	return fmt.Errorf("invalid module deployment status: %s", *stringValue)
}
