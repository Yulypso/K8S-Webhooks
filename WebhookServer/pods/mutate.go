package pods

import (
	"K8S-Webhooks/WebhookServer/admissioncontroller"
	"fmt"

	admission "k8s.io/api/admission/v1"
)

/** MUTATE CREATE **/
func mutateCreate(config Config) admissioncontroller.AdmitFunc {
	return func(r *admission.AdmissionRequest) (*admissioncontroller.Result, error) {

		/* Mutate Operation list */
		var operations []admissioncontroller.PatchOperation

		if r.Kind.Kind == "Pod" && r.Kind.Version == "v1" { /* Pod Mutations */
			fmt.Println("Log: POD MUTATING")

			/* get pod patches */
			operations = getPatches(config, Namespace(r.Namespace), operations)

		} else if r.Kind.Kind == "Deployment" && r.Kind.Version == "v1" { /* Deployment Mutations */
			fmt.Println("Log: DEPLOYMENT MUTATING")
			// TODO
		}

		/* DOING (Error pod deployment)
		 * Add: Check if the field already exist or not, (if YES, remove the operation from operations)
		 * Delete: Check if the field already exist or not, if NOT, remove the operation from operations
		 */
		operations = verifyOperations(operations, r)
		admissioncontroller.PrintPatchOperations(operations)

		return &admissioncontroller.Result{
			Allowed:  true,
			PatchOps: operations,
		}, nil
	}
}

/** MUTATE UPDATE **/

/** MUTATE DELETE **/

/** MUTATE CONNECT **/
