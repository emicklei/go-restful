package restful

import (
	"net/http"
	"log"
)
type WebService struct {
	Root     string
	routes   []Route
	Produces string
	Consumes string
}

func (self *WebService) Path(root string) *WebService {
	self.Root = root
	return self
}
func (self *WebService) Route(builder *RouteBuilder) *WebService {
	self.routes = append(self.routes, builder.Build())
	return self
}
func (self WebService) Dispatch(httpWriter http.ResponseWriter, httpRequest *http.Request) bool {
	log.Printf("Webservice %#v",self)
	for _, each := range self.routes {
		if each.dispatch(httpWriter, httpRequest) {
			return true
		}
	}
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
