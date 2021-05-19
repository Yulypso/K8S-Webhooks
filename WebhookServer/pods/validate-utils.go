package pods

import (
	"fmt"

	v1 "k8s.io/api/core/v1"
)

/* Unvalidate if Pod is run as root, forbidden: true*/
func unvalidateRunAsRoot(pod *v1.Pod, message *string) bool {
	if pod.Spec.SecurityContext.RunAsUser == nil { // Root user or uninitialized int64 type
		fmt.Println("Log: Run as root detected, Unvalidating...")
		*message = "Run as root is forbidden"
		return true
	}
	return false
}
