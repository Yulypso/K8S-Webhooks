package http

import (
	"K8S-Webhooks/WebhookServer/pods"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"time"
)

func NewServer(port string, tlsCertPath string, tlsKeyPath string) *http.Server {
	mux := http.NewServeMux()
	admissionHandler := newAdmissionHandler()

	/*
	 * Webhooks
	 */
	podMutation := pods.NewMutationWebhook()
	podMutation2 := pods.NewMutationWebhook2()
	podValidation := pods.NewValidationWebhook()

	/*
	 * Routers
	 */
	//mux.HandleFunc("/mutate/podsa", admissionHandler.serve(podMutationA))
	mux.HandleFunc("/mutate/pods/a", admissionHandler.serve(podMutation))
	mux.HandleFunc("/mutate/pods/b", admissionHandler.serve(podMutation2))
	mux.HandleFunc("/validate/pods", admissionHandler.serve(podValidation))

	return &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		TLSConfig:    certSetup(tlsCertPath, tlsKeyPath),
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}

func certSetup(certPath string, privKeyPath string) (serverTLSConf *tls.Config) {
	serverCert, err := tls.LoadX509KeyPair(certPath, privKeyPath)
	if err != nil {
		log.Fatal(err)
	}

	return &tls.Config{
		Certificates: []tls.Certificate{serverCert},
	}
}
