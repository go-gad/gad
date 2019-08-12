// May be we should move it to separate client lib

package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
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

type RequestOption func(req *http.Request)

// method name given from operationId parameter
func (c *Client) GetEmployee(ctx context.Context, req GetEmployeeRequest, options ...RequestOption) (*Employee, error) {
	httpReq, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%s/v1/employee_by_phone/%s", c.Host, req.Phone),
		nil,
	)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't create new request")
	}

	httpReq.Header.Add("User-Agent", "begett/1.0")
	httpReq.Header.Add("Accept", "application/json")

	for _, r := range options {
		r(httpReq)
	}
	
	resp, err := c.HttpClient.Do(httpReq)

	if err != nil {
		return nil, errors.Wrap(err, "Transport layer error")
	}

	if resp.StatusCode == 200 {
		//marshall
		return &Employee{}, nil // Should we return &GetEmployeeResp200{} ?
	}

	if resp.StatusCode == 500 {
		return nil, &GetEmployeeResp500{}
	}

	if resp.StatusCode == 404 {
		return nil, &GetEmployeeResp400{
			// fill params
		}
	}

	if resp.StatusCode == 204 {
		return nil, &GetEmployeeResp204{}
	}

	return nil, errors.Wrap(err, "Unknown error")
}
