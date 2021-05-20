package pods

import (
	"K8S-Webhooks/WebhookServer/admissioncontroller"
	"fmt"
)

/** UTILS **/

/* Annotate mutations */
func annotate(key string, value string, operations []admissioncontroller.PatchOperation) []admissioncontroller.PatchOperation {
	fmt.Println("Log: Add /metadata/annotations...")
	metadata := map[string]string{key: value}
	operations = append(operations, admissioncontroller.AddPatchOperation("/metadata/annotations", metadata))
	return operations
}

/* Apply operations */
func patches(operations []admissioncontroller.PatchOperation) []admissioncontroller.PatchOperation {
	fmt.Println("Log: Operating...")

	// Operate
	patches := admissioncontroller.ParseConfig("patches.conf")
	for path, value := range patches { // for each patch within the config file
		operations = append(operations, admissioncontroller.AddPatchOperation(path, value))
	}

	// Add annotation
	operations = annotate("mutate-CREATE", "mutateRunAsRoot", operations)
	return operations
}
