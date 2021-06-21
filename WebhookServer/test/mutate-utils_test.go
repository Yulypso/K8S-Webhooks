package test

import (
	"K8S-Webhooks/WebhookServer/admissioncontroller"
	"K8S-Webhooks/WebhookServer/pods"
	"fmt"
	"io/ioutil"
	"testing"

	admission "k8s.io/api/admission/v1"
	authenticationv1 "k8s.io/api/authentication/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
)

func TestVerifyDeployment(t *testing.T) {
	admissionRequestBytes, err := ioutil.ReadFile("files/admissionRequest.json")
	if err != nil {
		fmt.Println(err)
	}

	optionBytes, err := ioutil.ReadFile("files/options.json")
	if err != nil {
		fmt.Println(err)
	}

	r := admission.AdmissionRequest{
		UID: types.UID("da61dffe-3205-404c-bddf-980e36c95914"),
		Kind: metav1.GroupVersionKind{
			Kind: "Pod",
		},
		Resource: metav1.GroupVersionResource{
			Group:    "",
			Version:  "v1",
			Resource: "pods",
		},
		SubResource: "",
		Operation:   admission.Create,
		Name:        "node-app",
		Namespace:   "admissionwebhook-ns",
		UserInfo: authenticationv1.UserInfo{
			Username: "docker-for-desktop",
			UID:      "",
			Groups:   []string{"system:masters", "system:authenticated"},
			Extra:    make(map[string]authenticationv1.ExtraValue),
		},
		Object: runtime.RawExtension{
			Raw:    admissionRequestBytes,
			Object: nil,
		},
		OldObject: runtime.RawExtension{
			Raw:    []byte{},
			Object: nil,
		},
		DryRun: &[]bool{false}[0],
		Options: runtime.RawExtension{
			Raw:    optionBytes,
			Object: nil,
		},
		RequestKind: &metav1.GroupVersionKind{
			Group:   "",
			Version: "v1",
			Kind:    "Pod",
		},
		RequestResource: &metav1.GroupVersionResource{
			Group:    "",
			Version:  "v1",
			Resource: "pods",
		},
		RequestSubResource: "",
	}

	config := pods.Byte2Config(pods.ReadFile("files/config.json"))

	var jpOperations []admissioncontroller.PatchOperation
	var jpVerifications []admissioncontroller.PatchOperation

	jpOperations = pods.GetJsonPathOperations(config, pods.Namespace((&r).Namespace), jpOperations)
	jpVerifications = pods.GetJsonPathVerifications(config, pods.Namespace((&r).Namespace), jpVerifications)

	operations, _ := pods.VerifyDeployment(jpOperations, (&r))

	_, _ = pods.VerifyDeployment(jpVerifications, (&r))

	/* operations contains only Patch operations */
	if len(operations) != len(jpOperations) {
		t.Errorf("Expected deck length of %v but got %v", len(jpOperations), len(operations))
	}

	/* operations contains only add/remove/replace operations */
	for _, o := range operations {
		if o.Op != "add" && o.Op != "remove" && o.Op != "replace" {
			t.Errorf("Expected operation add or remove or replace but got %v operation", o.Op)
		}
	}

	/* multiple: $.spec.containers */
	if operations[0].Path != "/spec/containers/3" {
		t.Errorf("Expected operations[0].Path to be /spec/containers/3 but got %v", operations[0].Path)
	}

	/* string: $.spec.securityContext.runAsUser */
	if operations[1].Path != "/spec/securityContext/runAsUser" {
		t.Errorf("Expected operations[1].Path to be /spec/securityContext/runAsUser but got %v", operations[1].Path)
	}

	/* filter: $.spec.containers[?(@[name] == 'node-app0')] */
	if operations[2].Path != "/spec/containers/0" {
		t.Errorf("Expected operations[2].Path to be /spec/containers/3 but got %v", operations[2].Path)
	}
}
