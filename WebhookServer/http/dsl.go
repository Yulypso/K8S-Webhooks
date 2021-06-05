package http

import (
	"K8S-Webhooks/WebhookServer/pods"
	"encoding/json"
	"fmt"
	"net/http"
)

type DslConfig string

/* TODO
 *
 */
func Dsl(rw http.ResponseWriter, r *http.Request) {
	if r.Method == "PUT" {
		fmt.Println("Log: Add DSL Request received")

		var config pods.Config // TODO: Define AddDslConfig JSON : namespace => on lui rajoute des champs
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&config)
		if err != nil {
			panic(err)
		}
		defer r.Body.Close()

		rw.WriteHeader(http.StatusOK)
		json.NewEncoder(rw).Encode(config)
		return
	} else if r.Method == "DELETE" {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		rw.Write([]byte("HTTP Method not allowed."))
	}
}
