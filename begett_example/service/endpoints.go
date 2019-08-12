package service

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	GetEmployeeEndpoint endpoint.Endpoint
}

func MakeEndoints(s IService) Endpoints {
	return Endpoints{
		GetEmployeeEndpoint: MakeGetEmployeeEndpoint(s),
	}
}

func (e Endpoints) GetEmployee(ctx context.Context, phone string) (Employee, error) {
	req := GetEmployeeRequest{Phone: phone}
	resp, err := e.GetEmployeeEndpoint(ctx, req)
	if err != nil {
		return Employee{}, err
	}
	responce := resp.(GetEmployeeResponse)
	return responce.Employee, responce.Err
}

func MakeGetEmployeeEndpoint(s IService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetEmployeeRequest)
		e, err := s.GetEmployee(ctx, req.Phone)
		return GetEmployeeResponse{
			Employee: e,
			Err:      err,
		}, nil
	}
}
