package pods

import (
	"K8S-Webhooks/WebhookServer/admissioncontroller"
	"fmt"

	v1 "k8s.io/api/core/v1"
)

/** UTILS **/

/* Annotate mutations */
func annotateMutate(key string, value string, operations []admissioncontroller.PatchOperation) []admissioncontroller.PatchOperation {
	fmt.Println("Log: Add /metadata/annotations...")
	metadata := map[string]string{key: value}
	operations = append(operations, admissioncontroller.AddPatchOperation("/metadata/annotations", metadata))
	return operations
}

/* Mutate if Pod is run as root */
func mutateRunAsRoot(pod *v1.Pod, operations []admissioncontroller.PatchOperation) []admissioncontroller.PatchOperation {
	if pod.Spec.SecurityContext.RunAsUser == nil { // Root user or uninitialized int64 type
		fmt.Println("Log: Run as root detected, Mutating...")

		patches := admissioncontroller.ParseConfig("securityContext.conf")
		for path, value := range patches { // for each patch within the config file
			operations = append(operations, admissioncontroller.AddPatchOperation(path, value))
		}

		// Add annotation
		operations = annotateMutate("mutate-CREATE", "mutateRunAsRoot", operations)
	}
	return operations
}
