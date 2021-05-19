package http

import (
	"K8S-Webhooks/WebhookServer/admissioncontroller"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	admission "k8s.io/api/admission/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/klog/v2"
)

type admissionHandler struct {
	decoder runtime.Decoder
}

// To deserialize provided data into Kubernetes objects
func newAdmissionHandler() *admissionHandler {
	return &admissionHandler{
		decoder: serializer.NewCodecFactory(runtime.NewScheme()).UniversalDeserializer(),
	}
}

func checkIsPostMethod(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Error: Invalid method, Only POST requests are allowed", http.StatusMethodNotAllowed)
		return
	}
}

func checkContentType(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Error: Invalid content type, Only content type 'application/json' is supported", http.StatusBadRequest)
		return
	}
}

func getBodyRequest(w http.ResponseWriter, r *http.Request) []byte {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not read request body: %v", err), http.StatusBadRequest)
	}
	return body
}

func getReview(h *admissionHandler, body []byte, w http.ResponseWriter) admission.AdmissionReview {
	var review admission.AdmissionReview
	if _, _, err := h.decoder.Decode(body, nil, &review); err != nil {
		http.Error(w, fmt.Sprintf("Error: Could not deserialize request: %v", err), http.StatusBadRequest)
	}
	if review.Request == nil {
		http.Error(w, "Error: Malformed admission review: request is nil", http.StatusBadRequest)
	}
	return review
}

func (h *admissionHandler) serve(hook admissioncontroller.Hook) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Log: Post request received")

		checkIsPostMethod(w, r)
		fmt.Println("Log: 1")
		checkContentType(w, r)
		fmt.Println("Log: 2")
		body := getBodyRequest(w, r)
		fmt.Println("Log: 3")
		review := getReview(h, body, w)
		fmt.Println("Log: 4")

		/* Mutating/Validating Webhook execution */
		result, err := hook.Execute(review.Request)
		if err != nil {
			klog.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		fmt.Println("Log: 5")
		admissionResponse := admission.AdmissionReview{
			Response: &admission.AdmissionResponse{
				UID:     review.Request.UID,
				Allowed: result.Allowed,
				Result:  &meta.Status{Message: result.Msg},
			},
		}
		fmt.Println("Log: 6")
		if len(result.PatchOps) > 0 {
			patchBytes, err := json.Marshal(result.PatchOps)
			if err != nil {
				klog.Error(err)
				http.Error(w, fmt.Sprintf("Error: Could not marshal JSON patch: %v", err), http.StatusInternalServerError)
			}
			admissionResponse.Response.Patch = patchBytes
		}
		fmt.Println("Log: 7")
		res, err := json.Marshal(admissionResponse)
		if err != nil {
			klog.Error(err)
			http.Error(w, fmt.Sprintf("Error: Could not marshal response: %v", err), http.StatusInternalServerError)
			return
		}
		fmt.Println("Log: 8")
		klog.Infof("Webhook [%s - %s] - Allowed: %t", r.URL.Path, review.Request.Operation, result.Allowed)
		w.WriteHeader(http.StatusOK)
		w.Write(res)
	}
}
