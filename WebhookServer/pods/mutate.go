package pods

import (
	"K8S-Webhooks/WebhookServer/admissioncontroller"
	"fmt"
	"log"

	admission "k8s.io/api/admission/v1"
)

/** MUTATE CREATE **/
func mutateCreate(config Config) admissioncontroller.AdmitFunc {
	return func(r *admission.AdmissionRequest) (*admissioncontroller.Result, error) {

		/* Mutate Operation list */
		var jpOperations []admissioncontroller.PatchOperation

		if r.Kind.Kind == "Pod" && r.Kind.Version == "v1" { /* Pod Mutations */
			fmt.Println("Log: POD MUTATING")

			/* get pod patches */
			jpOperations = getJsonPathOperations(config, Namespace(r.Namespace), jpOperations)

		} else if r.Kind.Kind == "Deployment" && r.Kind.Version == "v1" { /* Deployment Mutations */
			fmt.Println("Log: DEPLOYMENT MUTATING")
			// TODO
		}

		/* DOING (Error pod deployment)
		 * Add: Check if the field already exist or not, (if YES, remove the operation from operations)
		 * Delete: Check if the field already exist or not, if NOT, remove the operation from operations
		 */
		operations, err := verifyDeployment(jpOperations, r)
		if err != nil {
			log.Println(err)
			return &admissioncontroller.Result{
				Allowed:  false,
				PatchOps: operations,
				Msg:      err.Error(),
			}, nil
		}
		admissioncontroller.PrintPatchOperations(jpOperations)
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
