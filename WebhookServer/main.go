package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func main() {
	TlsCertPath := "/etc/secrets/tls/tls.crt"
	TlsKeyPath := "/etc/secrets/tls/tls.key"

	go func() {
		handler := http.NewServeMux()

		handler.HandleFunc("/test", func(rw http.ResponseWriter, r *http.Request) {
			// kubectl exec -it -n webhookserver-ns $(kubectl get pods --no-headers -o custom-columns=":metadata.name" -n webhookserver-ns) -- wget -q -O- "localhost:8080/test"
			fmt.Fprintf(rw, "test\n")
		})
		fmt.Printf("Starting localhost http server on :8080 with test endpoint\n")
		if err := http.ListenAndServe("localhost:8080", handler); err != nil {
			log.Fatal(err)
		}
	}()

	fmt.Printf("Starting TLS server on: 8443\n")
	//mux := http.NewServeMux()
	//mux.HandleFunc("/validate/pods", postWebhook)
	http.HandleFunc("/validate/pods", postWebhook)
	server := &http.Server{
		Addr:      ":8443",
		TLSConfig: certSetup(TlsCertPath, TlsKeyPath),
		//Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	if err := server.ListenAndServeTLS("", ""); err != nil {
		log.Fatal(err)
	}
}

func certSetup(certPath string, privKeyPath string) (serverTLSConf *tls.Config) {

	cert, err := ioutil.ReadFile(certPath)
	if err != nil {
		log.Fatal(err)
	}
	privKey, err := ioutil.ReadFile(privKeyPath)
	if err != nil {
		log.Fatal(err)
	}

	serverCert, err := tls.X509KeyPair(cert, privKey)
	if err != nil {
		log.Fatal(err)
	}

	/*serverCert, err := tls.LoadX509KeyPair(certPath, privKeyPath)
	if err != nil {
		log.Fatal(err)
	}*/

	//fmt.Println("--- Certificate ---")
	//fmt.Println(string(serverCert.Certificate[0]))
	//fmt.Println(string(cert))
	//fmt.Println("--- Private Key ---")
	//fmt.Println(string(privKey))
	//fmt.Println(serverCert.PrivateKey)

	return &tls.Config{
		Certificates: []tls.Certificate{serverCert},
	}
}
