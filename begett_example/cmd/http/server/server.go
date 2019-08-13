package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-gad/gad/begett_example/business_logic"
	"github.com/go-gad/gad/begett_example/service"
	"github.com/gorilla/mux"
)

func main() {
	businessLogicSvc := &business_logic.Service{}

	router := service.NewRouter(businessLogicSvc)

	err := http.ListenAndServe(":80", router)
	if err != nil {
		panic(err)
	}
	os.Exit(0)
}

func StartCustomServer() {
	businessLogicSvc := &business_logic.Service{}

	getEmployeeCustomEndpoint := service.MakeGetEmployeeEndpoint(businessLogicSvc)

	getEmployeeCustomHandler := service.MakeGetEmployeeCustomHandler(
		getEmployeeCustomEndpoint,
		service.DecodeGetEmployeeRequest,
		service.EncodeResponse,
	)

	r := mux.NewRouter()
	service.AddGetEmployeeCustomRouter(r, getEmployeeCustomHandler)
	r.Use(loggingMiddleware)

	err := http.ListenAndServe(":80", r)
	if err != nil {
		panic(err)
	}
	os.Exit(0)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI)
		next.ServeHTTP(w, r)
	})
}