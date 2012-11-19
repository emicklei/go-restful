package restful

import (
	"net/http"
)

type Request struct {
	Request        *http.Request
	pathParameters map[string]string
}

func (self *Request) PathParameter(name string) string {
	return self.pathParameters[name]
}
func (self *Request) QueryParameter(name string) string {
	return self.Request.FormValue(name)
}
