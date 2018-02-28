package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	log.Println("Start program")

	// step 5.1
	s := &service{}

	// step 5.2
	h := NewCreateUserHandler(
		MakeCreateUserServiceWrapper(s),
		decodeCreateUserRequest,
		encodeCreateUserResponse,
	)

	// step 5.3 - make server and start listen
	srv := &http.Server{
		Addr:    ":8080",
		Handler: h,
	}

	log.Println("Start server on `:8080`")
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

// step 1.1
type UserService interface {
	CreateUser(string) (int, error)
}

// step 1.2
type service struct{}

// step 1.3
func (s *service) CreateUser(name string) (int, error) {
	log.Println("called CreateUser method of service")
	return 100500, nil
}

// step 2.1
type CreateUserRequest struct {
	Name string
}

// step 2.2
type CreateUserResponse struct {
	Id int
}

// step 2.3
func decodeCreateUserRequest(_ context.Context, _ *http.Request) (*CreateUserRequest, error) {
	return &CreateUserRequest{Name: "foobar"}, nil
}

// step 2.4
func encodeCreateUserResponse(_ context.Context, w http.ResponseWriter, response *CreateUserResponse) error {
	return json.NewEncoder(w).Encode(response)
}

// step 3.1
// go generate CreateUserRequest CreateUserResponse CreateUserHanlder

// >>> START GENERATED CODE
type DecodeCreateUserRequestFunc func(context.Context, *http.Request) (*CreateUserRequest, error)

type EncodeCreateUserResponseFunc func(context.Context, http.ResponseWriter, *CreateUserResponse) error

type CreateUserServiceWrapperFunc func(context.Context, *CreateUserRequest) (*CreateUserResponse, error)

type createUserHandler struct {
	dec         DecodeCreateUserRequestFunc
	enc         EncodeCreateUserResponseFunc
	serviceWrap CreateUserServiceWrapperFunc
}

func (h createUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	request, err := h.dec(ctx, r)
	if err != nil {
		// write to w
		return
	}

	response, err := h.serviceWrap(ctx, request)
	if err != nil {
		// write to w
		return
	}

	err = h.enc(ctx, w, response)
	if err != nil {
		// write to w
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
	return func(_ context.Context, r *CreateUserRequest) (*CreateUserResponse, error) {
		id, err := s.CreateUser(r.Name)
		if err != nil {
			return nil, err
		}
		return &CreateUserResponse{Id: id}, nil
	}
}
