package external

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func HandleRequests() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/namespace/{id:[0-9]+}", func(rw http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(rw, mux.Vars(r)["id"])
	}).Methods("PUT")
	return router
}
