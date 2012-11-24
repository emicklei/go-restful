package restful

import (
	"net/http"
	"strings"
)

type WebService struct {
	rootPath string
	routes   []Route
	produces string
	consumes string
}

// Specify the root URL path of the WebService.
// All Routes will be relative to this path.
func (self *WebService) Path(root string) *WebService {
	self.rootPath = root
	return self
}

// Create a new Route using the RouteBuilder and add to the ordered list of Routes.
func (self *WebService) Route(builder *RouteBuilder) *WebService {
	builder.copyDefaults(self.produces, self.consumes)
	self.routes = append(self.routes, builder.Build())
	return self
}

// Dispatch the incoming Http Request to a matching Route.
// The first matching Route will process the request and write any response.
// Return the Http Status as the result.
func (self WebService) Dispatch(httpWriter http.ResponseWriter, httpRequest *http.Request) int {
	// cheap test before iterating the routes
	if !strings.HasPrefix(httpRequest.URL.Path, self.rootPath) {
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
	return new(RouteBuilder).RootPath(self.rootPath).Method(httpMethod)
}

// Specify that this WebService can produce one or more MIME types.
func (self *WebService) Produces(contentTypes ...string) *WebService {
	self.produces = strings.Join(contentTypes, ",")
	return self
}

// Specify that this WebService can consume one or more MIME types.
func (self *WebService) Consumes(accepts ...string) *WebService {
	self.consumes = strings.Join(accepts, ",")
	return self
}

// TODO make routes public?
func (self *WebService) Routes() []Route {
	return self.routes
}

/*
	Convenience methods
*/

// Shortcut for .Method("GET").Path(subPath)
func (self *WebService) GET(subPath string) *RouteBuilder {
	return new(RouteBuilder).RootPath(self.rootPath).Method("GET").Path(subPath)
}

// Shortcut for .Method("POST").Path(subPath)
func (self *WebService) POST(subPath string) *RouteBuilder {
	return new(RouteBuilder).RootPath(self.rootPath).Method("POST").Path(subPath)
}

// Shortcut for .Method("PUT").Path(subPath)
func (self *WebService) PUT(subPath string) *RouteBuilder {
	return new(RouteBuilder).RootPath(self.rootPath).Method("PUT").Path(subPath)
}

// Shortcut for .Method("DELETE").Path(subPath)
func (self *WebService) DELETE(subPath string) *RouteBuilder {
	return new(RouteBuilder).RootPath(self.rootPath).Method("DELETE").Path(subPath)
}
