package pods

import (
	"K8S-Webhooks/WebhookServer/admissioncontroller"
	"fmt"

	admission "k8s.io/api/admission/v1"
	v1 "k8s.io/api/core/v1"
)

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

		/* Apply pod validating operation conditions */
		if forbidden := unvalidateRunAsRoot(pod, &message); forbidden {
			return &admissioncontroller.Result{Msg: message}, nil
		}

		return &admissioncontroller.Result{
			Allowed: true,
		}, nil
	}
}

/* Unvalidate if Pod is run as root, forbidden: true*/
func unvalidateRunAsRoot(pod *v1.Pod, message *string) bool {
	if pod.Spec.SecurityContext.RunAsUser == nil { // Root user or uninitialized int64 type
		fmt.Println("Log: Run as root detected, Unvalidating...")
		*message = "Run as root is forbidden"
		return true
	}
	return false
}
