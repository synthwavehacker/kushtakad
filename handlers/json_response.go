package handlers

import (
	"encoding/json"
)

type Response struct {
	Status  string
	Message string
	Err     error
	Service interface{}
}

func (r *Response) AddService(serv interface{}) {
	r.Service = serv
}

func NewResponse(s, msg string, err error) *Response {
	return &Response{
		Status:  s,
		Message: msg,
		Err:     err,
	}
}

func (r *Response) JSON() []byte {
	b, err := json.Marshal(r)
	if err != nil {
		log.Error(err)
		return b
	}
	return b
}
