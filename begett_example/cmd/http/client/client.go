package main

import (
	"context"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-gad/gad/begett_example/service"
)

func main() {
	client, err := service.NewHTTPClient("http://127.0.0.1")
	if err != nil {
		panic(spew.Sprintf("Client error: %#v", err))
	}

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
