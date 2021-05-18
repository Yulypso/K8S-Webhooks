package main

import (
	"K8S-Webhooks/WebhookServer/http"
	"context"
	"os"
	"os/signal"
	"syscall"

	log "k8s.io/klog/v2"
)

func main() {
	tlsCertPath := "/etc/secrets/tls/tls.crt"
	tlsKeyPath := "/etc/secrets/tls/tls.key"
	port := "8443"

	log.Infof("Starting TLS server on port: %s", port)
	server := http.NewServer(port, tlsCertPath, tlsKeyPath)

	go func() {
		if err := server.ListenAndServeTLS("", ""); err != nil {
			log.Errorf("Failed to listen and serve: %v", err)
		}
	}()

	// listen to any shutdown signal
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	log.Infof("Shutdown...")
	if err := server.Shutdown(context.Background()); err != nil {
		log.Error(err)
	}
}
