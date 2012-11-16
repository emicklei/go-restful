package restful

import (
	"net/http"
)

type Request struct {
	*http.Request
	pathParameters map[string]string
}

func (self *Request) PathParameter(name string) string {
	return self.pathParameters[name]
}
func (self *Request) QueryParameter(name string) string {
	return self.FormValue(name)
}
