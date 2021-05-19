package pods

import (
	"K8S-Webhooks/WebhookServer/admissioncontroller"
	"fmt"

	admission "k8s.io/api/admission/v1"
)

func mutateCreate() admissioncontroller.AdmitFunc {
	return func(r *admission.AdmissionRequest) (*admissioncontroller.Result, error) {
		fmt.Println("Log: POD MUTATING...")
		// Mutate Operation list
		var operations []admissioncontroller.PatchOperation

		// Parse pod
		pod, err := unmarshalPod(r.Object.Raw)
		fmt.Println("Log: DONE UNMARSHAL POD")
		if err != nil {
			fmt.Println("Error: FAILED UNMARSHAL POD")
			return &admissioncontroller.Result{Msg: err.Error()}, nil
		}
		fmt.Println("Log: CONDITIONS")
		fmt.Println("IS WORKING ?")
		fmt.Println(pod.GetLabels())
		fmt.Println("IS BUGGING ?")
		fmt.Println(pod.Spec.SecurityContext.RunAsUser)
		fmt.Println("HEHEHEEH")

		// Mutate if Pod is run as root
		if pod.Spec.SecurityContext.RunAsUser == nil { // Root user or uninitialized int64 type
			fmt.Println("Log: Run as root detected, Mutating...")
			patches := admissioncontroller.ParseConfig("securityContext.conf")
			for path, value := range patches { // for each patch within the config file
				operations = append(operations, admissioncontroller.AddPatchOperation(path, value))
			}

			// Add a simple annotation using `AddPatchOperation`
			fmt.Println("Log: Add annotations...")
			metadata := map[string]string{"origin": "Mutation"}
			operations = append(operations, admissioncontroller.AddPatchOperation("/metadata/annotations", metadata))
		}

		fmt.Println("Log: DONE MUTATING...")
		return &admissioncontroller.Result{
			Allowed:  true,
			PatchOps: operations,
		}, nil
	}
}
