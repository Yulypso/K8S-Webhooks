package admissioncontroller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

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
	From  string      `json:"from"`
	Value interface{} `json:"value,omitempty"`
}

const (
	addOperation    = "add"
	removeOperation = "remove"
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

/***********************************************************************************/

// Parse JSON Patches/*.json files
type Config map[string]interface{}

func ParseConfig(configName string) Config {
	configPath := "../../Patches/"
	jsonFile, err := os.Open(configPath + configName)
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()

	bytes, _ := ioutil.ReadAll(jsonFile)
	var parsed Config
	json.Unmarshal([]byte(bytes), &parsed)
	return parsed
}
