package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"golang.org/x/exp/errors"
)

// IService - implement business logic interface
type IService interface {
	GetEmployee(ctx context.Context, phone string) (Employee, error)
}

func NewRouter(svc IService, options ...httptransport.ServerOption) *mux.Router {
	router := mux.NewRouter()

	endpoints := MakeEndoints(svc)

	// GetEmployee
	{
		getEmployeeHandler := MakeGetEmployeeHandler(endpoints.GetEmployeeEndpoint, options...)
		r := router.NewRoute()
		r.Name("GetEmployee")
		r.Methods(http.MethodGet)
		r.Path("/v1/employee_by_phone/{phone}")
		r.Handler(getEmployeeHandler)
	}

	return router
}

func AddGetEmployeeCustomRouter(r *mux.Router, handler http.Handler) *mux.Router {
	route := r.NewRoute()
	route.Methods(http.MethodGet)
	route.Path("/v1/employee_by_phone/{phone}")
	route.Handler(handler)
	return r
}

func MakeGetEmployeeHandler(ep endpoint.Endpoint, options ...httptransport.ServerOption) *httptransport.Server {
	return httptransport.NewServer(
		ep,
		DecodeGetEmployeeRequest,
		EncodeResponse,
		options...,
	)
}

func MakeGetEmployeeCustomHandler(endpoint endpoint.Endpoint, dec httptransport.DecodeRequestFunc, enc httptransport.EncodeResponseFunc, options ...httptransport.ServerOption) *httptransport.Server {
	return httptransport.NewServer(
		endpoint,
		dec,
		enc,
		options...,
	)
}

func DecodeGetEmployeeRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	phone, ok := mux.Vars(r)["phone"]
	if !ok {
		return nil, &GetEmployeeResp400{
			ErrCode: 1,
			Msg:     "incorrect phone value",
		}
	}
	return GetEmployeeRequest{Phone: phone}, nil
}

func EncodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(error); ok && e != nil && e.Error() != "" {
		// Not a Go kit transport error, but a business-logic error.
		// Provide those as HTTP errors.
		encodeError(ctx, e, w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codeFrom(err))
	encodingErr := json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})

	if encodingErr != nil {
		panic(fmt.Sprintf("encoding err: %#v", err))
	}
}

func codeFrom(err error) int {
	if errors.Is(err, &GetEmployeeResp500{}) {
		return http.StatusInternalServerError
	}
	if errors.Is(err, &GetEmployeeResp400{}) {
		return http.StatusBadRequest
	}
	if errors.Is(err, &GetEmployeeResp204{}) {
		return http.StatusNoContent
	}
	return http.StatusInternalServerError
}
