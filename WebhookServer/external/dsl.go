package external

import (
	"K8S-Webhooks/WebhookServer/pods"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

/* Patch dsl config within persistant volume from external request */
func patchConfig(rw http.ResponseWriter, r *http.Request, dsl string) {
	var opType pods.OperationType
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&opType); err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Req Body must be type of pods.OperationType"))
		return
	}

	namespace := mux.Vars(r)["name"]

	configBytes, err := ioutil.ReadFile(dsl)
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
		return
	}
	config := pods.Byte2Config(configBytes)
	config[pods.Namespace(namespace)] = opType
	configBytes = pods.Config2Byte(config)

	if err = ioutil.WriteFile(dsl, configBytes, 0644); err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
		return
	}
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("Patch: DSL Config"))
}

/* Reset config to default.json */
func resetConfig(rw http.ResponseWriter, r *http.Request, dsl string, def string) {
	input, err := ioutil.ReadFile(def)
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
		return
	}

	err = ioutil.WriteFile(dsl, input, 0644)
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
		return
	}
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("Reset: DSL Config"))
}

/* Clear config to empty file */
func clearConfig(rw http.ResponseWriter, r *http.Request, dsl string) {
	if err := ioutil.WriteFile(dsl, nil, 0644); err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
		return
	}
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("Clear: DSL Config"))
}
