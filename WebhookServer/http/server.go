package http

import (
	"K8S-Webhooks/WebhookServer/pods"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func NewServer(port string, tlsCertPath string, tlsKeyPath string) *http.Server {
	/* Verify if dsl config exist */
	dsl := os.Getenv("DSL")
	def := os.Getenv("DEFAULT_DSL")
	if _, err := os.Stat(dsl); os.IsNotExist(err) {
		fmt.Printf(dsl + " does not exist, creating ...\n")
		initConfig(def, dsl)
	}

	/* Used in dev */
	initConfig(def, dsl)
	/***************/

	/* Load Config file */
	config := pods.Byte2Config(pods.ReadFile(dsl))

	/* Webhooks */
	podMutation := pods.NewMutationWebhook(config)
	podValidation := pods.NewValidationWebhook()

	/* Routers */
	mux := http.NewServeMux()
	admissionHandler := newAdmissionHandler()
	mux.HandleFunc("/mutate", admissionHandler.serve(podMutation))
	mux.HandleFunc("/validate", admissionHandler.serve(podValidation))

	/* TODO
	 * Update Config Endpoints
	 * - Add
	 * - Remove (Namespace)
	 */
	mux.HandleFunc("/namespace/{id}", Dsl)

	return &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		TLSConfig:    certSetup(tlsCertPath, tlsKeyPath),
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}

/* Copy default dsl config into persistant volume */
func initConfig(def string, dsl string) {
	input, err := ioutil.ReadFile(def)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = ioutil.WriteFile(dsl, input, 0644)
	if err != nil {
		fmt.Println("Error creating", dsl)
		fmt.Println(err)
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
