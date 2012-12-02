package restful

// RouteBuilder is a helper to construct Routes.
// httpMethod and function are required fields 
// Produces,Consumes, rootPath and currentPath fields are optional
type RouteBuilder struct {
	rootPath    string
	currentPath string
	Produces    []string
	Consumes    []string

	httpMethod string
	function   RouteFunction
	doc        string
}

func (self *RouteBuilder) To(function RouteFunction) *RouteBuilder {
	self.function = function
	return self
}
func (self *RouteBuilder) Method(method string) *RouteBuilder {
	self.httpMethod = method
	return self
}
func (self *RouteBuilder) ContentType(contentTypes ...string) *RouteBuilder {
	self.Produces = contentTypes
	return self
}
func (self *RouteBuilder) Accept(accepts ...string) *RouteBuilder {
	self.Consumes = accepts
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
func (self *RouteBuilder) Doc(documentation string) *RouteBuilder {
	self.doc = documentation
	return self
}

// If no specific Route path then set to rootPath
// If no specific Produces then set to rootProduces
// If no specific Consumes then set to rootConsumes
func (self *RouteBuilder) copyDefaults(rootProduces, rootConsumes []string) {
	if len(self.Produces) == 0 {
		self.Produces = rootProduces
	}
	if len(self.Consumes) == 0 {
		self.Consumes = rootConsumes
	}
}

func (self *RouteBuilder) Build() Route {
	route := Route{
		Method:       self.httpMethod,
		Path:         self.rootPath + self.currentPath,
		Produces:     self.Produces,
		Consumes:     self.Consumes,
		Function:     self.function,
		relativePath: self.currentPath,
		Doc:          self.doc}
	route.postBuild()
	return route
}
