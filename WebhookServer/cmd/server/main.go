package main

import (
	"K8S-Webhooks/WebhookServer/http"
	"context"
	"log"
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
	port := "8443"

	klog.Infof("Starting TLS server on port: %s", port)
	server := http.NewServer(port, tlsCertPath, tlsKeyPath)

	go func() {
		if err := server.ListenAndServeTLS("", ""); err != nil {
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
