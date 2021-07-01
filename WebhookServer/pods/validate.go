package pods

import (
	"K8S-Webhooks/WebhookServer/admissioncontroller"
	"fmt"
	"log"
	"os"

	admission "k8s.io/api/admission/v1"
)

/* VALIDATE CREATE */
func validateCreate() admissioncontroller.AdmitFunc {
	return func(r *admission.AdmissionRequest) (*admissioncontroller.Result, error) {
		/* Load Config file */
		dsl := os.Getenv("DSL")
		common := os.Getenv("COMMON")
		config := Byte2Config(ReadFile(dsl))
		commonConfig := Byte2OpType(ReadFile(common))

		/* Mutate Operation list */
		var jpVerifications []admissioncontroller.PatchOperation

		if r.Kind.Kind == "Pod" && r.Kind.Version == "v1" { /* Pod Verifications */
			/* get pod verifications */
			jpVerifications = GetJsonPathVerifications(config, Namespace(r.Namespace), jpVerifications)
			jpVerifications = GetJsonPathCommonVerifications(commonConfig, jpVerifications)

		} else if r.Kind.Kind == "Deployment" && r.Kind.Version == "v1" { /* Deployment Mutations */
			fmt.Println("Log: DEPLOYMENT MUTATING")
			// TODO
		}
		err := VerifyValidation(jpVerifications, r)
		if err != nil {
			log.Println(err)
			return &admissioncontroller.Result{
				Allowed: false,
				Msg:     err.Error(),
			}, nil
		}

		fmt.Print("Config: Verification Operations")
		admissioncontroller.PrintPatchOperations(jpVerifications)

		return &admissioncontroller.Result{
			Allowed: true,
		}, nil
	}
}
