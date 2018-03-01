package main

import (
	kithttp "github.com/go-kit/kit/transport/http"
	"net/http"
	"net/http/httptest"
	"testing"
)

// START-GAD OMIT
func BenchmarkGad(b *testing.B) {
	s := &service{}

	h := NewCreateUserHandler(
		MakeCreateUserServiceWrapper(s),
		GadDecodeCreateUserRequest,
		GadEncodeCreateUserResponse,
	)
	for n := 0; n < b.N; n++ {
		req, _ := http.NewRequest("GET", "/gad", nil)
		rw := httptest.NewRecorder()
		h.ServeHTTP(rw, req)
	}
}

// END-GAD OMIT

// START-KIT OMIT
func BenchmarkKit(b *testing.B) {
	s := &service{}

	kh := kithttp.NewServer(
		MakeCreateUserEndpoint(s),
		KitDecodeDetermineRequest,
		KitEncodeCreateUserResponse,
	)
	for n := 0; n < b.N; n++ {
		req, _ := http.NewRequest("GET", "/kit", nil)
		rw := httptest.NewRecorder()
		kh.ServeHTTP(rw, req)
	}
}

// END-KIT OMIT
