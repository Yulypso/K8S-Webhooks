package pods

import (
	"K8S-Webhooks/WebhookServer/admissioncontroller"
	"fmt"
	"log"
	"os"
	"sort"

	admission "k8s.io/api/admission/v1"
)

/** MUTATE CREATE **/
func mutateCreate() admissioncontroller.AdmitFunc {
	return func(r *admission.AdmissionRequest) (*admissioncontroller.Result, error) {
		/* Load Config file */
		dsl := os.Getenv("DSL")
		common := os.Getenv("COMMON")
		config := Byte2Config(ReadFile(dsl))
		commonConfig := Byte2OpType(ReadFile(common))

		/* Mutate Operation list */
		var jpOperations []admissioncontroller.PatchOperation

		if r.Kind.Kind == "Pod" && r.Kind.Version == "v1" { /* Pod Mutations */
			/* get pod patches */
			jpOperations = GetJsonPathOperations(config, Namespace(r.Namespace), jpOperations)
			jpOperations = GetJSONPathCommonOperations(commonConfig, jpOperations)

		} else if r.Kind.Kind == "Deployment" && r.Kind.Version == "v1" { /* Deployment Mutations */
			fmt.Println("Log: DEPLOYMENT MUTATING")
			// TODO
		}

		operations, err := VerifyMutation(jpOperations, r)
		if err != nil {
			log.Println(err)
			return &admissioncontroller.Result{
				Allowed:  false,
				PatchOps: operations,
				Msg:      err.Error(),
			}, nil
		}
		sort.Slice(operations, func(i, j int) bool { return operations[i].Op < operations[j].Op })

		fmt.Print("Config: Patch Operations")
		admissioncontroller.PrintPatchOperations(jpOperations)
		fmt.Print("Applied Operations")
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
