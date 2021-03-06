package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func main() {
	log.Println("Start program")

	s := &service{}

	h := NewCreateUserHandler(
		MakeCreateUserServiceWrapper(s),
		decodeCreateUserRequest,
		encodeCreateUserResponse,
		HandlerBefore(Logger), // step 3.1
		ServiceBefore(Auth),   // step 3.2
	)

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

// step 2.1
func Logger(ctx context.Context, r *http.Request) (context.Context, HandlerResponseFunc) {
	start := time.Now()
	log.Printf("start process request %s %s", r.Method, r.RequestURI)

	f := func(ctx context.Context, w http.ResponseWriter) context.Context {
		log.Printf("request %s %s tooks %v", r.Method, r.RequestURI, time.Since(start))
		return ctx
	}

	return ctx, f
}

// step 2.2
func Auth(ctx context.Context, request interface{}) (context.Context, ServiceResponseFunc) {
	// do smth with request
	log.Println("authorize request")
	user := "Ivan"

	f := func(ctx context.Context, response interface{}) context.Context {
		// do smth with response
		// use variables from after hook scope
		log.Printf("User `%s` finished work", user)
		return ctx
	}

	return ctx, f
}

// >>> step 1.1
type HanlderRequestFunc func(context.Context, *http.Request) (context.Context, HandlerResponseFunc)

type ServiceRequestFunc func(_ context.Context, request interface{}) (context.Context, ServiceResponseFunc)

type ServiceResponseFunc func(_ context.Context, response interface{}) context.Context

type HandlerResponseFunc func(context.Context, http.ResponseWriter) context.Context

type HookOption func(*Hooks)

type Hooks struct {
	beforeH []HanlderRequestFunc
	beforeS []ServiceRequestFunc
	afterS  []ServiceResponseFunc
	afterH  []HandlerResponseFunc
}

func HandlerBefore(hooks ...HanlderRequestFunc) HookOption {
	return func(h *Hooks) { h.beforeH = append(h.beforeH, hooks...) }
}

func ServiceBefore(hooks ...ServiceRequestFunc) HookOption {
	return func(h *Hooks) { h.beforeS = append(h.beforeS, hooks...) }
}

func ServiceAfter(hooks ...ServiceResponseFunc) HookOption {
	return func(h *Hooks) { h.afterS = append(h.afterS, hooks...) }
}

func HandlerAfter(hooks ...HandlerResponseFunc) HookOption {
	return func(h *Hooks) { h.afterH = append(h.afterH, hooks...) }
}

// <<< step 1.1

type UserService interface {
	CreateUser(string) (int, error)
}

type service struct{}

func (s *service) CreateUser(name string) (int, error) {

	log.Println("called CreateUser method of service")
	return 100500, nil
}

type CreateUserRequest struct {
	Name string
}

type CreateUserResponse struct {
	Id int
}

func decodeCreateUserRequest(_ context.Context, _ *http.Request) (*CreateUserRequest, error) {
	return &CreateUserRequest{Name: "foobar"}, nil
}

func encodeCreateUserResponse(_ context.Context, w http.ResponseWriter, response *CreateUserResponse) error {
	return json.NewEncoder(w).Encode(response)
}

// go generate CreateUserRequest CreateUserResponse CreateUserHanlder

// >>> START GENERATED CODE
type DecodeCreateUserRequestFunc func(context.Context, *http.Request) (*CreateUserRequest, error)

type EncodeCreateUserResponseFunc func(context.Context, http.ResponseWriter, *CreateUserResponse) error

type CreateUserServiceWrapperFunc func(context.Context, *CreateUserRequest) (*CreateUserResponse, error)

type createUserHandler struct {
	dec         DecodeCreateUserRequestFunc
	enc         EncodeCreateUserResponseFunc
	serviceWrap CreateUserServiceWrapperFunc
	hooks       *Hooks
}

func (h createUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// step 1.2
	for _, f := range h.hooks.beforeH {
		var after HandlerResponseFunc
		ctx, after = f(ctx, r)
		if after != nil {
			HandlerAfter(after)(h.hooks)
		}
	}

	request, err := h.dec(ctx, r)
	if err != nil {
		// write to w
		return
	}

	// step 1.4
	for _, f := range h.hooks.beforeS {
		var after ServiceResponseFunc
		ctx, after = f(ctx, r)
		if after != nil {
			ServiceAfter(after)(h.hooks)
		}
	}

	response, err := h.serviceWrap(ctx, request)
	if err != nil {
		// write to w
		return
	}

	// step 1.5
	for _, f := range h.hooks.afterS {
		ctx = f(ctx, response)
	}

	// step 1.3
	for _, f := range h.hooks.afterH {
		ctx = f(ctx, w)
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
	options ...HookOption,
) *createUserHandler {
	h := &createUserHandler{
		serviceWrap: serviceWrap,
		dec:         dec,
		enc:         enc,
		hooks:       &Hooks{},
	}
	for _, opt := range options {
		opt(h.hooks)
	}
	return h
}

// <<< END GENERATED CODE

func MakeCreateUserServiceWrapper(s UserService) CreateUserServiceWrapperFunc {
	return func(_ context.Context, r *CreateUserRequest) (*CreateUserResponse, error) {
		id, err := s.CreateUser(r.Name)
		if err != nil {
			return nil, err
		}
		return &CreateUserResponse{Id: id}, nil
	}
}
