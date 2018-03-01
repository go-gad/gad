package main

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	"net/http"
)

func KitDecodeDetermineRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return userData(), nil
}

func KitEncodeCreateUserResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func MakeCreateUserEndpoint(s UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		data := request.(*CreateUserData)
		id, err := s.CreateUser(data.Name)
		if err != nil {
			return nil, err
		}
		data.Id = id
		return data, nil
	}
}
