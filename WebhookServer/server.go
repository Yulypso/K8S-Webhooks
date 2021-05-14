package main

import (
	"fmt"
	"net/http"
)

func postWebhook(w http.ResponseWriter, r *http.Request) {
	fmt.Println("post request received")
	fmt.Fprintf(w, "post request received\n")
	/*var request AdmissionReviewRequest
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

	// add label if we're creating a pod
	if request.Request.Kind.Group == "" && request.Request.Kind.Version == "v1" && request.Request.Kind.Kind == "Pod" && request.Request.Operation == "CREATE" {
		patch := `[{
			"op": "add",
			"path": "/metadata/labels/myExtraLabel",
			"value": "webhook-was-here"
		}]`
		patchEnc := base64.StdEncoding.EncodeToString([]byte(patch))
		response.Response.PatchType = "JSONPatch"
		response.Response.Patch = patchEnc
	}

	out, err := json.Marshal(response)
	if err != nil {
		http.Error(w, fmt.Sprintf("JSON output marshal error: %s\n", err.Error()), http.StatusBadRequest)
		return
	}
	fmt.Printf("Got request, response: %s\n", string(out))
	fmt.Fprintln(w, string(out))*/
}
