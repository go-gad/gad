package main

import (
	"context"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-gad/gad/begett_example/service"
	"github.com/go-kit/kit/endpoint"
	"github.com/sirupsen/logrus"
)

func main() {
	client, err := service.NewHTTPClient("http://127.0.0.1")
	if err != nil {
		panic(spew.Sprintf("Client error: %#v", err))
	}

	client.GetEmployeeEndpoint = logWrapper(client.GetEmployeeEndpoint)

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

func logWrapper(endpoint endpoint.Endpoint) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		logger := logrus.StandardLogger()
		logger.SetFormatter(&logrus.JSONFormatter{})
		logger.WithField("time", time.Now)
		logger.Info("start request")
		resp, err := endpoint(ctx, request)
		time.Sleep(time.Second)
		logger.Info("end request")
		return resp, err
	}
}