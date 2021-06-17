package main

import (
	"K8S-Webhooks/WebhookServer/external"
	"K8S-Webhooks/WebhookServer/internal"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"k8s.io/klog/v2"
)

func main() {
	err := godotenv.Load("../../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	tlsCertPath := "/etc/secrets/tls/tls.crt"
	tlsKeyPath := "/etc/secrets/tls/tls.key"
	internalPort := "8443"
	externalPort := "8080"

	server := internal.NewServer(internalPort, tlsCertPath, tlsKeyPath)

	go func() {
		klog.Infof("Starting TLS internal server on port: %s", internalPort)
		if err := server.ListenAndServeTLS("", ""); err != nil {
			klog.Errorf("Failed to listen and serve: %v", err)
		}
	}()

	go func() {
		klog.Infof("Starting external server on port: %s", externalPort)
		router := external.HandleRequests()
		if err := http.ListenAndServe(":8080", router); err != nil {
			klog.Errorf("Failed to listen and serve: %v", err)
		}
	}()

	// listen to any shutdown signal
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	klog.Infof("Shutdown...")
	if err := server.Shutdown(context.Background()); err != nil {
		klog.Error(err)
	}
}
