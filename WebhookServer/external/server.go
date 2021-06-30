package external

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"k8s.io/klog/v2"
)

func HandleRequests() *mux.Router {
	dsl := os.Getenv("DSL")
	def := os.Getenv("DEFAULT_DSL")
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/namespace/{name:[A-Za-z0-9-]+}", func(rw http.ResponseWriter, r *http.Request) {
		patchConfig(rw, r, dsl)
		klog.Infof("Patch: DSL Config")
	}).Methods("PUT")

	router.HandleFunc("/reset", func(rw http.ResponseWriter, r *http.Request) {
		resetConfig(rw, r, dsl, def)
		klog.Infof("Reset: DSL Config")
	}).Methods("GET")

	router.HandleFunc("/clear/{name:[A-Za-z0-9-]+}", func(rw http.ResponseWriter, r *http.Request) {
		clearConfig(rw, r, dsl)
		klog.Infof("Clear: DSL Config")
	}).Methods("GET")

	return router
}
