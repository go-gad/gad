package main

import (
	kithttp "github.com/go-kit/kit/transport/http"
	"log"
	"net/http"
	"time"
)

func main() {
	log.Println("Start program")

	// step 5.1
	s := &service{}

	// step 5.2
	gh := NewCreateUserHandler(
		MakeCreateUserServiceWrapper(s),
		GadDecodeCreateUserRequest,
		GadEncodeCreateUserResponse,
	)

	kh := kithttp.NewServer(
		MakeCreateUserEndpoint(s),
		KitDecodeDetermineRequest,
		KitEncodeCreateUserResponse,
	)

	mux := http.NewServeMux()
	mux.Handle("/gad", gh)
	mux.Handle("/kit", kh)
	// step 5.3 - make server and start listen
	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Println("Start server on `:8080`")
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

type UserService interface {
	CreateUser(string) (int, error)
}

type service struct{}

func (s *service) CreateUser(name string) (int, error) {
	return 100500, nil
}

type CreateUserData struct {
	Id       int
	Name     string
	CreateAt time.Time
	Items    []Item
}

type Item struct {
	Id    int
	Title string
}

func userData() *CreateUserData {
	return &CreateUserData{
		Name:     "foo-bar",
		CreateAt: time.Now(),
		Items: []Item{
			Item{Id: 1, Title: "ssss111"},
			Item{Id: 2, Title: "ssss222"},
			Item{Id: 3, Title: "ssss333"},
		},
	}
}
