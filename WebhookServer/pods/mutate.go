package pods

import (
	"K8S-Webhooks/WebhookServer/admissioncontroller"
	"fmt"

	admission "k8s.io/api/admission/v1"
)

/** MUTATE CREATE **/
func mutateCreate(config Config) admissioncontroller.AdmitFunc {
	return func(r *admission.AdmissionRequest) (*admissioncontroller.Result, error) {
		fmt.Println("Log: POD MUTATING...")

		fmt.Println(r.Kind.Kind, r.Kind.Version)
		fmt.Println(r.Namespace)
		fmt.Println("----")
		fmt.Println(config)

		/* Mutate Operation list */
		var operations []admissioncontroller.PatchOperation

		/* Apply pod mutating operation conditions */
		operations = getPatches(config, operations)

		return &admissioncontroller.Result{
			Allowed:  true,
			PatchOps: operations,
		}, nil
	}
}

/** MUTATE UPDATE **/

/** MUTATE DELETE **/

/** MUTATE CONNECT **/
