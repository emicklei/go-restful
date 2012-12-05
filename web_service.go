package restful

import ()

type WebService struct {
	rootPath string
	routes   []Route
	produces []string
	consumes []string
}

// Specify the root URL template path of the WebService.
// All Routes will be relative to this path.
func (self *WebService) Path(root string) *WebService {
	self.rootPath = root
	return self
}

// Document the Path Parameter used in my Root
func (self *WebService) PathParam(name, documentation string) *WebService {
	// TODO
	return self
}

// Create a new Route using the RouteBuilder and add to the ordered list of Routes.
func (self *WebService) Route(builder *RouteBuilder) *WebService {
	builder.copyDefaults(self.produces, self.consumes)
	self.routes = append(self.routes, builder.Build())
	return self
}

// Create a new RouteBuilder and initialize its http method
func (self *WebService) Method(httpMethod string) *RouteBuilder {
	return new(RouteBuilder).servicePath(self.rootPath).Method(httpMethod)
}

// Specify that this WebService can produce one or more MIME types.
func (self *WebService) Produces(contentTypes ...string) *WebService {
	self.produces = contentTypes
	return self
}

// Specify that this WebService can consume one or more MIME types.
func (self *WebService) Consumes(accepts ...string) *WebService {
	self.consumes = accepts
	return self
}

func (self WebService) Routes() []Route {
	return self.routes
}
func (self WebService) RootPath() string {
	return self.rootPath
}

/*
	Convenience methods
*/

// Shortcut for .Method("GET").Path(subPath)
func (self *WebService) GET(subPath string) *RouteBuilder {
	return new(RouteBuilder).servicePath(self.rootPath).Method("GET").Path(subPath)
}

// Shortcut for .Method("POST").Path(subPath)
func (self *WebService) POST(subPath string) *RouteBuilder {
	return new(RouteBuilder).servicePath(self.rootPath).Method("POST").Path(subPath)
}

// Shortcut for .Method("PUT").Path(subPath)
func (self *WebService) PUT(subPath string) *RouteBuilder {
	return new(RouteBuilder).servicePath(self.rootPath).Method("PUT").Path(subPath)
}

// Shortcut for .Method("DELETE").Path(subPath)
func (self *WebService) DELETE(subPath string) *RouteBuilder {
	return new(RouteBuilder).servicePath(self.rootPath).Method("DELETE").Path(subPath)
}
