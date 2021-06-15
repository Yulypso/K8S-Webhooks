package admissioncontroller

import (
	"fmt"

	admission "k8s.io/api/admission/v1"
)

// Result contains the result of an admission request
type Result struct {
	Allowed  bool
	Msg      string
	PatchOps []PatchOperation
}

// AdmitFunc defines how to process an admission request
type AdmitFunc func(request *admission.AdmissionRequest) (*Result, error)

// Hook represents the set of functions for each operation in an admission webhook.
type Hook struct {
	Create  AdmitFunc
	Delete  AdmitFunc
	Update  AdmitFunc
	Connect AdmitFunc
}

func (h *Hook) Execute(r *admission.AdmissionRequest) (*Result, error) {
	switch r.Operation {
	case admission.Create:
		return h.Create(r)
	case admission.Update:
		return h.Update(r)
	case admission.Delete:
		return h.Delete(r)
	case admission.Connect:
		return h.Connect(r)
	}
	return &Result{Msg: fmt.Sprintf("Error: Invalid operation: %s", r.Operation)}, nil
}

/***********************************************************************************/

// PatchOperation is an operation of a JSON patch
type PatchOperation struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value,omitempty"`
}

func PrintPatchOperations(operations []PatchOperation) {
	fmt.Println("\nOperations:", len(operations))
	for _, op := range operations {
		fmt.Println("["+op.Op+"]\n- Path:", op.Path, "\n- Value:", op.Value)
		fmt.Println()
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
