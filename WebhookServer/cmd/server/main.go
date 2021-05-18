package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"time"

	"github.com/yulypso/K8S-Webhooks/http"
	log "k8s.io/klog/v2"
)

func main() {
	tlsCertPath := "/etc/secrets/tls/tls.crt"
	tlsKeyPath := "/etc/secrets/tls/tls.key"
	port := "8443"

	log.Infof("Starting TLS server on port: %s", port)
	mux := http.NewServeMux()
	mux.HandleFunc("/validate/pods", postWebhook)
	server := &http.Server{
		Addr:         ":8443",
		TLSConfig:    certSetup(tlsCertPath, tlsKeyPath),
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	if err := server.ListenAndServeTLS("", ""); err != nil {
		log.Fatal(err)
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
