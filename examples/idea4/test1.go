package main

import (
	"context"
	"log"
)

func main() {
	log.Println("start")

	h := &Handler{
		Req: CreateReq{},
		End: func(ctx context.Context, r CreateReq) error {
			log.Println("endp func called")
			return nil
		},
	}

	h.Serve()
}

type IDec interface {
	Decode()
}

type Handler struct {
	Req IDec
	End EndpCreate
}

type EndpCreate func(ctx context.Context, r CreateReq) error

func (h *Handler) Serve() {
	log.Println("start serve")
	h.End(nil, h.Req)
}

type CreateReq struct{}

func (r CreateReq) Decode() {}
