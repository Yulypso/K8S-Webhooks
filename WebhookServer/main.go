package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	//config := Config{}
	//TlsCertPath := "../Certificates/webhookservertls.cert"
	//TlsKeyPath := "../Certificates/webhookservertls.key"

	TlsCertPath := "/etc/secrets/tls/tls.crt"
	TlsKeyPath := "/etc/secrets/tls/tls.key"
	serverTLSConf, err := certSetup(TlsCertPath, TlsKeyPath)
	if err != nil {
		fmt.Println("Error:", err)
	}

	//time.Sleep(9999999 * time.Second)

	go func() {
		handler := http.NewServeMux()

		handler.HandleFunc("/test", func(rw http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(rw, "test\n")
		})
		fmt.Printf("Starting localhost http server on :8080 with test endpoint\n")
		err = http.ListenAndServe("localhost:8080", handler)

		if err != nil {
			log.Fatal(err)
		}
	}()

	fmt.Printf("Starting TLS server on: 8443\n")
	mux := http.NewServeMux()
	mux.HandleFunc("/validate/pods", postWebhook)
	server := &http.Server{
		Addr:      ":8443",
		TLSConfig: serverTLSConf,
		Handler:   mux,
	}

	if err := server.ListenAndServeTLS("", ""); err != nil {
		fmt.Printf("Failed to listen and serve: %v", err)
	}
}

func certSetup(certPath string, privKeyPath string) (serverTLSConf *tls.Config, err error) {

	cert, err := ioutil.ReadFile(certPath)
	if err != nil {
		fmt.Println("Error:", err)
	}
	PrivKey, err := ioutil.ReadFile(privKeyPath)
	if err != nil {
		fmt.Println("Error:", err)
	}

	serverCert, err := tls.X509KeyPair(cert, PrivKey)
	if err != nil {
		return nil, err
	}

	serverTLSConf = &tls.Config{
		Certificates: []tls.Certificate{serverCert},
	}
	return
}
