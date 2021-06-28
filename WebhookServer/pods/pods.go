package pods

import (
	"K8S-Webhooks/WebhookServer/admissioncontroller"
)

/* Mutating Webhooks for Pods */
func NewMutationWebhook() admissioncontroller.Hook {
	return admissioncontroller.Hook{
		Create: mutateCreate(),
	}
}
