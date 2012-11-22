package restful

import (
	"net/http"
	"strings"
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
	// cheap test before iterating the routes
	if (strings.HasPrefix(self.Root,httpRequest.URL.Path)) {
		return false
	}
	for _, each := range self.routes {
		if each.dispatch(httpWriter, httpRequest) {
			return true
		}
	}
	return false
}

func (self *WebService) Method(httpMethod string) *RouteBuilder {
	return new(RouteBuilder).RootPath(self.Root).Method(httpMethod)
}

func (self *WebService) ContentType(contentType string) *WebService {
	self.Produces = contentType
	return self
}
func (self *WebService) Accept(accept string) *WebService {
	self.Consumes = accept
	return self
}

/*
	Convenience methods
*/

// Shortcut for .Method("GET").Path(subPath)
func (self *WebService) GET(subPath string) *RouteBuilder {
	return new(RouteBuilder).RootPath(self.Root).Method("GET").Path(subPath)
}
// Shortcut for .Method("POST").Path(subPath)
func (self *WebService) POST(subPath string) *RouteBuilder {
	return new(RouteBuilder).RootPath(self.Root).Method("POST").Path(subPath)
}
// Shortcut for .Method("PUT").Path(subPath)
func (self *WebService) PUT(subPath string) *RouteBuilder {
	return new(RouteBuilder).RootPath(self.Root).Method("PUT").Path(subPath)
}
// Shortcut for .Method("DELETE").Path(subPath)
func (self *WebService) DELETE(subPath string) *RouteBuilder {
	return new(RouteBuilder).RootPath(self.Root).Method("DELETE").Path(subPath)
}