package business_logic

import (
	"context"

	"github.com/go-gad/gad/begett_example/service"
)

type Service struct{}

func (s *Service) GetEmployee(ctx context.Context, phone string) (service.Employee, error) {
	return service.Employee{
		UUID:  "TestUUID",
		Email: "TestEmail",
		Phone: "TestPhone",
		Company: &service.Company{
			UUID: "TestCompanyUUID",
			Name: "TestCompanyName",
		},
	}, nil
}
