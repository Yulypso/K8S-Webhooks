package main

import "fmt"

/***********************************************************************************/

// PatchOperation is an operation of a JSON patch
type PatchOperation struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value,omitempty"`
}

func PrintPatchOperations(operations []PatchOperation) {
	fmt.Println("\nOperations applied:")
	for _, op := range operations {
		fmt.Println(op)
	}
}

const (
	addOperation     = "add"
	removeOperation  = "remove"
	replaceOperation = "replace"
	mandatoryData    = "mandatorydata"
	forbiddenData    = "forbiddendata"
)

// AddPatchOperation returns an add JSON patch operation.
func AddPatchOperation(path string, value interface{}) PatchOperation {
	return PatchOperation{
		Op:    addOperation,
		Path:  path,
		Value: value,
	}
}

// RemovePatchOperation returns a remove JSON patch operation.
func RemovePatchOperation(path string) PatchOperation {
	return PatchOperation{
		Op:   removeOperation,
		Path: path,
	}
}

// ReplacePatchOperation returns an replace JSON patch operation.
func ReplacePatchOperation(path string, value interface{}) PatchOperation {
	return PatchOperation{
		Op:    replaceOperation,
		Path:  path,
		Value: value,
	}
}

// MandatoryDataCheckOperation.
func MandatoryDataCheckOperation(path string, value interface{}) PatchOperation {
	return PatchOperation{
		Op:    mandatoryData,
		Path:  path,
		Value: value,
	}
}

// ForbiddenDataCheckOperation.
func ForbiddenDataCheckOperation(path string, value interface{}) PatchOperation {
	return PatchOperation{
		Op:    forbiddenData,
		Path:  path,
		Value: value,
	}
}
