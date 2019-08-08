package main

import (
	"context"
	"net/http"
)

type UpdateOrderAssignedReq struct {
	ID          string
	RequestBody struct {
		DriverID int `json:driver_id`
	}
	Request *http.Request
}

func (r *UpdateOrderAssignedReq) Decode(ctx context.Context, req *http.Request) error {
	r.Request = req
	// set r.ID
	// decode body to r.RequestBody

	return nil
}

type UpdateOrderAssignedResp struct {
}

func (r *UpdateOrderAssignedResp) Encode(ctx context.Context, rw http.ResponseWriter) error {

	// encode json
	return nil
}

func handlerUpdateOrderAssigned(w http.ResponseWriter, r *http.Request) {
	//err := json.NewDecoder(r.Body).Decode(&e)
	// mux.Vars(r)
	w.WriteHeader(http.StatusNoContent)
}
