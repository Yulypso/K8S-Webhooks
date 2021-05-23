package pods

import (
	"K8S-Webhooks/WebhookServer/admissioncontroller"
	"encoding/json"
	"fmt"

	v1 "k8s.io/api/core/v1"
)

/* Parse the requested pod from yaml object */
func unmarshalPod(object []byte) (*v1.Pod, error) {
	var pod v1.Pod
	if err := json.Unmarshal(object, &pod); err != nil {
		fmt.Printf("Error: JSON Unmarshal failed %s", err)
		return nil, err
	}
	return &pod, nil
}

/* Mutating Webhooks for Pods */
func NewMutationWebhook(config Config) admissioncontroller.Hook {
	return admissioncontroller.Hook{
		Create: mutateCreate(config),
	}
}

/* Validating Webhooks for Pods */
func NewValidationWebhook() admissioncontroller.Hook {
	return admissioncontroller.Hook{
		Create: validateCreate(),
	}
}
