package restful

import (
	"net/http"
)

type Response struct {
	http.ResponseWriter
}

func (self *Response) StatusOK() *Response {
	self.WriteHeader(http.StatusOK)
	return self
}
func (self *Response) StatusError() *Response {
	self.WriteHeader(http.StatusInternalServerError)
	return self
}
func (self *Response) AddHeader(header string, value string) *Response {
	return self
}
func (self *Response) Entity(entity interface{}) *Response {
	return self
}
