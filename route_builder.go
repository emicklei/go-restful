package restful

// RouteBuilder is a helper to construct Routes.
// httpMethod and function are required fields 
// Produces,Consumes, rootPath and currentPath fields are optional
type RouteBuilder struct {
	rootPath    string
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
func (self *RouteBuilder) Path(subPath string) *RouteBuilder {
	self.currentPath = subPath
	return self
}
func (self *RouteBuilder) RootPath(path string) *RouteBuilder {
	self.rootPath = path
	return self
}
func (self *RouteBuilder) Build() Route {
	route := Route{
		Method:   self.httpMethod,
		Path:     self.rootPath + self.currentPath,
		Produces: self.Produces,
		Consumes: self.Consumes,
		Function: self.function}
	route.postBuild()
	return route
}
