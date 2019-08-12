package service

import (
	"context"
	"fmt"
	"net/http"
)

type Client struct {
	HttpClient *http.Client
	Host       string
}

func NewClient(httpClient *http.Client, host string, options ...ClientOption) *Client {
	client := &Client{
		HttpClient: httpClient,
		Host:       host,
	}
	for _, o := range options {
		o(client)
	}
	return client
}

type ClientOption func(c *Client)

type GetEmployeeRequest struct {
	Phone string
}

// method name given from operationId parameter
func (c *Client) GetEmployee(ctx context.Context, ger GetEmployeeRequest) (*Employee, error) {
	return &Employee{}, nil
}

type GetEmployeeResp500 struct {}

func (e *GetEmployeeResp500) Error() string {
	return "Internal Server error"
}

type GetEmployeeResp204 struct {}

func (e *GetEmployeeResp204) Error() string {
	return "Employee not found"
}

type GetEmployeeResp404 struct {
	ErrCode int `json:"err_code"`
	Msg string `json:"msg"`
}

func (e *GetEmployeeResp404) Error() string {
	return fmt.Sprintf("Bad request: error %d: message: %s", e.ErrCode, e.Msg);  // ? как выводить кастомные ошибки
}
