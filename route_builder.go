package restful

import ()

// RouteBuilder is a helper to construct Routes
// httpMethod and function are required 
// Produces,Consumes and currentPath are optional
type RouteBuilder struct {
	currentPath string
	Produces    string
	Consumes    string

	httpMethod string
	function   RouteFunction
}

func (self *RouteBuilder) To(function RouteFunction) *RouteBuilder {
	self.function = function
	return self
}
func (self *RouteBuilder) Method(method string) *RouteBuilder {
	self.httpMethod = method
	return self
}
func (self *RouteBuilder) ContentType(contentType string) *RouteBuilder {
	self.Produces = contentType
	return self
}
func (self *RouteBuilder) Accept(accept string) *RouteBuilder {
	self.Consumes = accept
	return self
}
func (self *RouteBuilder) Path(otherPath string) *RouteBuilder {
	self.currentPath = otherPath
	return self
}
func (self *RouteBuilder) Build() Route {
	route := Route{
		Method:   self.httpMethod,
		Path:     self.currentPath,
		Produces: self.Produces,
		Consumes: self.Consumes,
		Function: self.function}
	route.postBuild()
	return route
}
