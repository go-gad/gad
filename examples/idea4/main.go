package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	log.Println("Start program")

	router := mux.NewRouter()
	router.Path("/orders/{id}/assigned").Methods("PATCH").HandlerFunc(handlerUpdateOrderAssigned)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	log.Println("Start server on `:8080`")
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func handlerUpdateOrderAssigned(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}
