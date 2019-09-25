package handlers

import (
	"encoding/json"
	"fmt"
)

type Response struct {
	Status  string
	Message string
	Err     error
	Type    string
	Service interface{}
	Obj     interface{}
}

func (r *Response) AddService(serv interface{}) {
	r.Service = serv
}

func NewResponse(s, msg string, err error) *Response {
	var newmsg string
	if err != nil {
		newmsg = fmt.Sprintf("%s > %s", msg, err.Error())
	} else {
		newmsg = fmt.Sprintf("%s", msg)
	}
	return &Response{
		Status:  s,
		Message: newmsg,
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
