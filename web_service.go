package restful

import (
	"net/http"
)

type Dispatcher interface {
	// Dispath the request to a matching Route and call its Function.
	// Return whether the request was handled.
	Dispatch(request *http.Request) bool
}

type WebService struct {
	Root     string
	methods  []Route
	Produces string
	Consumes string
}

func (self *WebService) Path(root string) *WebService {
	self.Root = root
	return self
}
func (self *WebService) Route(builder *RouteBuilder) *WebService {
	self.methods = append(self.methods, builder.Build())
	return self
}
func (self WebService) Dispatch(request *http.Request) bool {
	//	for _, each := range self.methods {
	//		if each.canHandle(request) {
	//			return true
	//		}
	//	}
	return false
}

func (self *WebService) Method(httpMethod string) *RouteBuilder {
	return new(RouteBuilder).Method(httpMethod)
}

func (self *WebService) ContentType(contentType string) *WebService {
	self.Produces = contentType
	return self
}
func (self *WebService) Accept(accept string) *WebService {
	self.Consumes = accept
	return self
}
