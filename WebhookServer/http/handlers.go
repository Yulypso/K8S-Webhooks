package http

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"k8s.io/klog/v2"
)

func podMutation(w http.ResponseWriter, r *http.Request) {
	fmt.Println("log: post request received")

	var request AdmissionReviewRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, fmt.Sprintf("JSON body in invalid format: %s\n", err.Error()), http.StatusBadRequest)
		return
	}

	if request.APIVersion != "admission.k8s.io/v1" || request.Kind != "AdmissionReview" {
		http.Error(w, fmt.Sprintf("wrong APIVersion or kind: %s - %s", request.APIVersion, request.Kind), http.StatusBadRequest)
		return
	}

	fmt.Printf("debug: %+v\n", request.Request)

	response := AdmissionReviewResponse{
		APIVersion: "admission.k8s.io/v1",
		Kind:       "AdmissionReview",
		Response: Response{
			UID:     request.Request.UID,
			Allowed: true,
		},
	}

	// add patches if we're creating a pod
	if request.Request.Kind.Group == "" && request.Request.Kind.Version == "v1" && request.Request.Kind.Kind == "Pod" && request.Request.Operation == "CREATE" {
		patch, err := ioutil.ReadFile("../../http/Patches/securityContext.conf")
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
			http.Error(w, fmt.Sprintf("Read config file error: %s\n", err.Error()), http.StatusBadRequest)
		}
		patchEnc := base64.StdEncoding.EncodeToString(patch)
		response.Response.PatchType = "JSONPatch"
		response.Response.Patch = patchEnc
	}

	out, err := json.Marshal(response)
	if err != nil {
		http.Error(w, fmt.Sprintf("JSON output marshal error: %s\n", err.Error()), http.StatusBadRequest)
		return
	}

	klog.Infof("Webhook [%s - %s] - Allowed: %s - Response: %s\n", r.URL.Path, "Mutate pod", "true", string(out))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}
