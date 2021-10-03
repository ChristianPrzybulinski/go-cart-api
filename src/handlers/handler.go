package handlers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func StartServer(port string) {
	router := mux.NewRouter().PathPrefix("/api/v1/").Subrouter()
	SetupHandlers(router)
	http.ListenAndServe(port, router)
}

func SetupHandlers(router *mux.Router) {

	router.HandleFunc("/cart", homePage).Methods(http.MethodGet)
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world!")
}
