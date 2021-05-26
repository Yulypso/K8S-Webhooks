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
	 * Load Config file
	 */
	config := pods.ParseConfig("patches.conf")

	/*
	 * Webhooks
	 */
	podMutation := pods.NewMutationWebhook(config)
	podValidation := pods.NewValidationWebhook()

	/*
	 * Routers
	 */
	mux.HandleFunc("/mutate", admissionHandler.serve(podMutation))
	mux.HandleFunc("/validate", admissionHandler.serve(podValidation))

	/* TODO Update Config Endpoint*/
	//add namespace
	//remove namespace

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
