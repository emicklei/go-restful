package restful

import (
	"net/http"
)

type Request struct {
	http.Request
}

func (self *Request) PathParameter(name string) string {
	return name
}
func (self *Request) QueryParameter(name string) string {
	return name
}
