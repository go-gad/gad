// May be we should move it to separate client lib

package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/pkg/errors"
)

type Client struct {
	GetEmployeeEndpoint endpoint.Endpoint
	Host                string
}

func NewHTTPClient(baseURL string, options ...ClientOption) (*Client, error) {
	getEmployeeURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	getEmployyEndpoint := httptransport.NewClient(
		http.MethodGet,
		getEmployeeURL,
		encodeHTTPGetEmployeeRequest,
		decodeHTTPGetEmployeeResponse,
	).Endpoint()

	client := &Client{
		Host:                baseURL,
		GetEmployeeEndpoint: getEmployyEndpoint,
	}
	for _, o := range options {
		o(client)
	}
	return client, nil
}

type ClientOption func(c *Client)

type RequestOption func(req *http.Request)

func encodeHTTPGetEmployeeRequest(_ context.Context, r *http.Request, request interface{}) error {
	getEmployeeReq := request.(GetEmployeeRequest)
	r.URL.Path = fmt.Sprintf("/v1/employee_by_phone/%s", getEmployeeReq.Phone)
	return nil
}

func decodeHTTPGetEmployeeResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode == 204 {
		return nil, &GetEmployeeResp204{}
	}
	if r.StatusCode == 404 {
		var respErr *GetEmployeeResp400
		err := json.NewDecoder(r.Body).Decode(respErr)
		if err != nil {
			return nil, errors.Wrap(err, "couldn't decode 404 error body")
		}
		return nil, respErr
	}
	if r.StatusCode == 500 {
		return nil, &GetEmployeeResp500{}
	}

	if r.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("Unknown error: code: %s", r.Status))
	}

	var resp GetEmployeeResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}

// method name given from operationId parameter
func (c *Client) GetEmployee(ctx context.Context, req GetEmployeeRequest, options ...RequestOption) (*Employee, error) {
	resp, err := c.GetEmployeeEndpoint(ctx, req)

	if err != nil {
		return nil, err
	}

	response, ok := resp.(GetEmployeeResponse)
	if !ok {
		return nil, errors.New("couldn't cast GetEmployeeResponse")
	}
	return &response.Employee, nil // Should we return &GetEmployeeResp200{} ?
}
