package http

import (
	"K8S-Webhooks/WebhookServer/pods"
	"encoding/json"
	"fmt"
	"net/http"
)

/* TODO
 *
 */
func dslAdd(rw http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		fmt.Println("Log: dslAdd Request received")

		var config pods.Config
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&config)
		if err != nil {
			panic(err)
		}
		defer r.Body.Close()

		rw.WriteHeader(http.StatusOK)
		json.NewEncoder(rw).Encode(config)
		return
	} else {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		rw.Write([]byte("HTTP Method not allowed."))
	}
}

func dslRemove(rw http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		fmt.Println("Log: dslRemove Request received")

		var config pods.Config
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&config)
		if err != nil {
			panic(err)
		}
		defer r.Body.Close()

		rw.WriteHeader(http.StatusOK)
		json.NewEncoder(rw).Encode(config)
		return
	} else {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		rw.Write([]byte("HTTP Method not allowed."))
	}
}
