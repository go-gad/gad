package main

import (
	"net/http"

	"github.com/go-gad/gad/begett_example/business_logic"
	"github.com/go-gad/gad/begett_example/service"
)

func main() {
	businessLogicSvc := &business_logic.Service{}

	router := service.NewRouter(businessLogicSvc)

	http.ListenAndServe(":80", router)
}
