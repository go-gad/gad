package service

import "fmt"

type Company struct {
	UUID string `json:"uuid"`
	Name string	`json:"name"`
}

type Employee struct {
	UUID string `json:"uuid"`
	Email string `json:"email"`
	Phone string `json:"phone"`
	Company *Company `json:"company"`
}

type GetEmployeeRequest struct {
	Phone string `json:"phone"`
}

type GetEmployeeResponse struct {
	Employee Employee `json:"employee,omitempty"`
	Err error `json:"err,omitempty"`
}

func (r *GetEmployeeResponse) Error() error {
	return r.Err
}

type GetEmployeeResp500 struct{}

func (e *GetEmployeeResp500) Error() string {
	return "Internal Server error"
}

type GetEmployeeResp204 struct{}

func (e *GetEmployeeResp204) Error() string {
	return "Employee not found"
}

type GetEmployeeResp400 struct {
	ErrCode int    `json:"err_code"`
	Msg     string `json:"msg"`
}

func (e *GetEmployeeResp400) Error() string {
	return fmt.Sprintf("Bad request: error %d: message: %s", e.ErrCode, e.Msg); // ? как выводить кастомные ошибки
}
