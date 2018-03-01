package main

import (
	"context"
	"encoding/json"
	"net/http"
)

func GadDecodeCreateUserRequest(_ context.Context, _ *http.Request) (*CreateUserData, error) {
	return userData(), nil
}

func GadEncodeCreateUserResponse(_ context.Context, w http.ResponseWriter, response *CreateUserData) error {
	return json.NewEncoder(w).Encode(response)
}

// >>> START GENERATED CODE
type DecodeCreateUserRequestFunc func(context.Context, *http.Request) (*CreateUserData, error)

type EncodeCreateUserResponseFunc func(context.Context, http.ResponseWriter, *CreateUserData) error

type CreateUserServiceWrapperFunc func(context.Context, *CreateUserData) (*CreateUserData, error)

type createUserHandler struct {
	dec         DecodeCreateUserRequestFunc
	enc         EncodeCreateUserResponseFunc
	serviceWrap CreateUserServiceWrapperFunc
}

func (h createUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	request, err := h.dec(ctx, r)
	if err != nil {
		return
	}

	response, err := h.serviceWrap(ctx, request)
	if err != nil {
		return
	}

	err = h.enc(ctx, w, response)
	if err != nil {
		return
	}
}

func NewCreateUserHandler(
	serviceWrap CreateUserServiceWrapperFunc,
	dec DecodeCreateUserRequestFunc,
	enc EncodeCreateUserResponseFunc,
) *createUserHandler {
	h := &createUserHandler{
		serviceWrap: serviceWrap,
		dec:         dec,
		enc:         enc,
	}
	return h
}

// <<< END GENERATED CODE

// step 4.1
func MakeCreateUserServiceWrapper(s UserService) CreateUserServiceWrapperFunc {
	return func(_ context.Context, r *CreateUserData) (*CreateUserData, error) {
		id, err := s.CreateUser(r.Name)
		if err != nil {
			return nil, err
		}
		r.Id = id
		return r, nil
	}
}
