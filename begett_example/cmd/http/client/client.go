package main

import (
	"context"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-gad/gad/begett_example/service"
)

func main() {
	// Some logic
	httpClient := http.DefaultClient

	client := service.NewClient(httpClient, "http://example.com")

	employee, err := client.GetEmployee(
		context.Background(),
		service.GetEmployeeRequest{
			Phone: "+79165177922",
		},
	)

	if err != nil {
		internalServerError, ok := err.(*service.GetEmployeeResp500)
		if ok {
			panic(internalServerError)
		}
		badRequest, ok := err.(*service.GetEmployeeResp400)
		if ok {
			panic(badRequest.Msg)
		}
		employeeNotFound, ok := err.(*service.GetEmployeeResp204)
		if ok {
			panic(employeeNotFound.Error())
		}
		panic("Unknown error")
	}

	spew.Dump(employee)
}
