package pods

import (
	"K8S-Webhooks/WebhookServer/admissioncontroller"
	"fmt"
	"log"
	"os"

	admission "k8s.io/api/admission/v1"
)

/** MUTATE CREATE **/
func mutateCreate() admissioncontroller.AdmitFunc {
	return func(r *admission.AdmissionRequest) (*admissioncontroller.Result, error) {

		/* Load Config file */
		dsl := os.Getenv("DSL")
		config := Byte2Config(ReadFile(dsl))

		/* Mutate Operation list */
		var jpOperations []admissioncontroller.PatchOperation
		var jpVerifications []admissioncontroller.PatchOperation

		if r.Kind.Kind == "Pod" && r.Kind.Version == "v1" { /* Pod Mutations */
			fmt.Println("Log: POD MUTATING")

			/* get pod patches */
			jpOperations = getJsonPathOperations(config, Namespace(r.Namespace), jpOperations)
			jpVerifications = getJsonPathVerifications(config, Namespace(r.Namespace), jpVerifications)

		} else if r.Kind.Kind == "Deployment" && r.Kind.Version == "v1" { /* Deployment Mutations */
			fmt.Println("Log: DEPLOYMENT MUTATING")
			// TODO
		}

		operations, err := verifyDeployment(jpOperations, r)
		if err != nil {
			log.Println(err)
			return &admissioncontroller.Result{
				Allowed:  false,
				PatchOps: operations,
				Msg:      err.Error(),
			}, nil
		}

		_, err = verifyDeployment(jpVerifications, r)
		if err != nil {
			log.Println(err)
			return &admissioncontroller.Result{
				Allowed:  false,
				PatchOps: operations,
				Msg:      err.Error(),
			}, nil
		}

		fmt.Println("Config: Patch Operations")
		admissioncontroller.PrintPatchOperations(jpOperations)
		fmt.Println("Config: Verification Operations")
		admissioncontroller.PrintPatchOperations(jpVerifications)
		fmt.Println("Applied Operations")
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
