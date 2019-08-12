package main

import (
	"net/http"
	"os"

	"github.com/go-gad/gad/begett_example/business_logic"
	"github.com/go-gad/gad/begett_example/service"
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
