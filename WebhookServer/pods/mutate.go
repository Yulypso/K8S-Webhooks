package pods

import (
	"K8S-Webhooks/WebhookServer/admissioncontroller"
	"fmt"

	admission "k8s.io/api/admission/v1"
)

/** MUTATE CREATE **/
func mutateCreate() admissioncontroller.AdmitFunc {
	return func(r *admission.AdmissionRequest) (*admissioncontroller.Result, error) {
		fmt.Println("Log: POD MUTATING...")

		/* Parse requested pod */
		/*pod, err := unmarshalPod(r.Object.Raw)
		if err != nil {
			return &admissioncontroller.Result{Msg: err.Error()}, nil
		}*/

		/* Mutate Operation list */
		var operations []admissioncontroller.PatchOperation

		/* Apply pod mutating operation conditions */
		operations = patches(operations)

		return &admissioncontroller.Result{
			Allowed:  true,
			PatchOps: operations,
		}, nil
	}
}

/** MUTATE UPDATE **/

/** MUTATE DELETE **/

/** MUTATE CONNECT **/
