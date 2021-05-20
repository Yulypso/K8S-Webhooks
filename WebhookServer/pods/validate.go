package pods

import (
	"K8S-Webhooks/WebhookServer/admissioncontroller"
	"fmt"

	admission "k8s.io/api/admission/v1"
)

/* VALIDATE CREATE */
func validateCreate() admissioncontroller.AdmitFunc {
	return func(r *admission.AdmissionRequest) (*admissioncontroller.Result, error) {
		fmt.Println("Log: POD VALIDATING...")

		/* Parse pod */
		pod, err := unmarshalPod(r.Object.Raw)
		if err != nil {
			return &admissioncontroller.Result{Msg: err.Error()}, nil
		}

		/* Result Message */
		var message string

		//TODO: verify if configuration is Guaranteed by verifications.conf

		/* Apply pod validating operation conditions */
		if authorized := validate(pod, &message); !authorized {
			return &admissioncontroller.Result{Msg: message}, nil
		}

		return &admissioncontroller.Result{
			Allowed: true,
		}, nil
	}
}

/* VALIDATE UPDATE */

/* VALIDATE DELETE */

/* VALIDATE CONNECT */
