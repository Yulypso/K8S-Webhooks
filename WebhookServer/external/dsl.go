package external

import (
	"K8S-Webhooks/WebhookServer/pods"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

/* Patch dsl config within persistant volume from external request */
func patchConfig(rw http.ResponseWriter, r *http.Request, dsl string) {
	var opType pods.OperationType

	buf := new(strings.Builder)
	if _, err := io.Copy(buf, r.Body); err != nil {
		log.Println("error:", err)
	}

	fmt.Println(buf.String())
	fmt.Println("---bufferr--")
	fmt.Println(buf.String())
	body := strings.ReplaceAll(buf.String(), "\\", "\\\\")
	fmt.Println("---body--")
	fmt.Println(body)
	reqBody := ioutil.NopCloser(strings.NewReader(body))

	decoder := json.NewDecoder(reqBody)
	if err := decoder.Decode(&opType); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Req Body must be type of pods.OperationType"))
		return
	}

	namespace := mux.Vars(r)["name"]

	configBytes, err := SyncReadFile(dsl)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
		return
	}
	config := pods.Byte2Config(configBytes)
	config[pods.Namespace(namespace)] = opType
	fmt.Println("---opType--")
	fmt.Println(opType)
	configBytes = pods.Config2Byte(config)

	if err = SyncWriteFile(dsl, configBytes); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
		return
	}
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("Patch: DSL Config"))
}

/* Reset config to default.json */
func resetConfig(rw http.ResponseWriter, r *http.Request, dsl string, def string) {
	input, err := SyncReadFile(def)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
		return
	}

	err = SyncWriteFile(dsl, input)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
		return
	}
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("Reset: DSL Config"))
}

/* Clear config to empty file */
func clearConfig(rw http.ResponseWriter, r *http.Request, dsl string) {
	namespace := mux.Vars(r)["name"]

	configBytes, err := SyncReadFile(dsl)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
		return
	}
	config := pods.Byte2Config(configBytes)
	delete(config, pods.Namespace(namespace))
	configBytes = pods.Config2Byte(config)

	if err = SyncWriteFile(dsl, configBytes); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
		return
	}
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("Clear: DSL Config"))
}
