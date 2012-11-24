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

// Specify the root URL path of the WebService.
// All Routes will be relative to this path.
func (self *WebService) Path(root string) *WebService {
	self.RootPath = root
	return self
}

// Create a new Route using the RouteBuilder and add to the ordered list of Routes.
func (self *WebService) Route(builder *RouteBuilder) *WebService {
	builder.copyDefaults(self.Produces, self.Consumes)
	self.routes = append(self.routes, builder.Build())
	return self
}

// Dispatch the incoming Http Request to a matching Route.
// The first matching Route will process the request and write any response.
// Return the Http Status as the result.
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

// Create a new RouteBuilder and initialize its http method
func (self *WebService) Method(httpMethod string) *RouteBuilder {
	return new(RouteBuilder).RootPath(self.RootPath).Method(httpMethod)
}

// Specify that this WebService can produce one or more MIME types.
func (self *WebService) ContentType(contentTypes ...string) *WebService {
	self.Produces = strings.Join(contentTypes, ",")
	return self
}

// Specify that this WebService can consume one or more MIME types.
func (self *WebService) Accept(accepts ...string) *WebService {
	self.Consumes = strings.Join(accepts, ",")
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

// func (self *WebService) String() string {}
