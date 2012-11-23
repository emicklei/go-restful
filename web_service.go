package restful

import (
	"net/http"
	"strings"
)

type WebService struct {
	RootPath string
	routes   []Route
	Produces string
	Consumes string
}

func (self *WebService) Path(root string) *WebService {
	self.RootPath = root
	return self
}
func (self *WebService) Route(builder *RouteBuilder) *WebService {
	builder.copyDefaults(self.Produces, self.Consumes)
	self.routes = append(self.routes, builder.Build())
	return self
}
func (self WebService) Dispatch(httpWriter http.ResponseWriter, httpRequest *http.Request) int {
	// cheap test before iterating the routes
	if !strings.HasPrefix(httpRequest.URL.Path, self.RootPath) {
		return http.StatusNotFound
	}
	var lastStatus int
	for _, each := range self.routes {
		lastStatus = each.dispatch(httpWriter, httpRequest)
		if http.StatusOK == lastStatus {
			return lastStatus
		}
	}
	return lastStatus
}

func (self *WebService) Method(httpMethod string) *RouteBuilder {
	return new(RouteBuilder).RootPath(self.RootPath).Method(httpMethod)
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
	return new(RouteBuilder).RootPath(self.RootPath).Method("GET").Path(subPath)
}

// Shortcut for .Method("POST").Path(subPath)
func (self *WebService) POST(subPath string) *RouteBuilder {
	return new(RouteBuilder).RootPath(self.RootPath).Method("POST").Path(subPath)
}

// Shortcut for .Method("PUT").Path(subPath)
func (self *WebService) PUT(subPath string) *RouteBuilder {
	return new(RouteBuilder).RootPath(self.RootPath).Method("PUT").Path(subPath)
}

// Shortcut for .Method("DELETE").Path(subPath)
func (self *WebService) DELETE(subPath string) *RouteBuilder {
	return new(RouteBuilder).RootPath(self.RootPath).Method("DELETE").Path(subPath)
}
