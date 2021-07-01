package internal

import (
	"K8S-Webhooks/WebhookServer/pods"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

func NewServer(port string, tlsCertPath string, tlsKeyPath string) *http.Server {
	/* Verify if dsl config exist */
	def := os.Getenv("DEFAULT_DSL")
	dsl := os.Getenv("DSL")

	if _, err := os.Stat(dsl); os.IsNotExist(err) {
		fmt.Printf(dsl + " does not exist, creating ...\n")
		InitConfig(def, dsl)
	}

	/* Used in dev */
	InitConfig(def, dsl)
	/***************/

	/* Webhooks */
	podMutation := pods.NewMutationWebhook()
	podValidation := pods.NewValidationWebhook()

	/* Routers */
	mux := http.NewServeMux()
	admissionHandler := newAdmissionHandler()
	mux.HandleFunc("/mutate", admissionHandler.serve(podMutation))
	mux.HandleFunc("/validate", admissionHandler.serve(podValidation))

	return &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		TLSConfig:    certSetup(tlsCertPath, tlsKeyPath),
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}

/* Copy default dsl config into persistant volume */
func InitConfig(def string, dsl string) {
	var mutex sync.Mutex
	mutex.Lock()
	input, err := ioutil.ReadFile(def)
	mutex.Unlock()
	if err != nil {
		log.Println(err)
		return
	}

	mutex.Lock()
	err = ioutil.WriteFile(dsl, input, 0644)
	mutex.Unlock()
	if err != nil {
		log.Println("Error creating", dsl)
		log.Println(err)
		return
	}
}

/* Load tls config */
func certSetup(certPath string, privKeyPath string) (serverTLSConf *tls.Config) {
	serverCert, err := tls.LoadX509KeyPair(certPath, privKeyPath)
	if err != nil {
		log.Fatal(err)
	}

	return &tls.Config{
		Certificates: []tls.Certificate{serverCert},
	}
}
