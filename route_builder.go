package restful

// RouteBuilder is a helper to construct Routes.
// httpMethod and function are required fields 
// Produces,Consumes, rootPath and currentPath fields are optional
type RouteBuilder struct {
	rootPath                string
	currentPath             string
	produces                []string
	consumes                []string
	httpMethod              string
	function                RouteFunction
	doc                     string
	readSample, writeSample interface{}
}

func (self *RouteBuilder) To(function RouteFunction) *RouteBuilder {
	self.function = function
	return self
}

func (self *RouteBuilder) Method(method string) *RouteBuilder {
	self.httpMethod = method
	return self
}
func (self *RouteBuilder) Produces(mimeTypes ...string) *RouteBuilder {
	self.produces = mimeTypes
	return self
}
func (self *RouteBuilder) Consumes(mimeTypes ...string) *RouteBuilder {
	self.consumes = mimeTypes
	return self
}
func (self *RouteBuilder) Path(subPath string) *RouteBuilder {
	self.currentPath = subPath
	return self
}
func (self *RouteBuilder) Doc(documentation string) *RouteBuilder {
	self.doc = documentation
	return self
}

func (self *RouteBuilder) Reads(sample interface{}) *RouteBuilder {
	self.readSample = sample
	return self
}

func (self *RouteBuilder) Writes(sample interface{}) *RouteBuilder {
	self.writeSample = sample
	return self
}

func (self *RouteBuilder) QueryParam(name, comment string) *RouteBuilder {
	return self
}

func (self *RouteBuilder) PathParam(name, comment string) *RouteBuilder {
	return self
}

func (self *RouteBuilder) servicePath(path string) *RouteBuilder {
	self.rootPath = path
	return self
}

// If no specific Route path then set to rootPath
// If no specific Produces then set to rootProduces
// If no specific Consumes then set to rootConsumes
func (self *RouteBuilder) copyDefaults(rootProduces, rootConsumes []string) {
	if len(self.produces) == 0 {
		self.produces = rootProduces
	}
	if len(self.consumes) == 0 {
		self.consumes = rootConsumes
	}
}

func (self *RouteBuilder) Build() Route {
	route := Route{
		Method:       self.httpMethod,
		Path:         self.rootPath + self.currentPath,
		Produces:     self.produces,
		Consumes:     self.consumes,
		Function:     self.function,
		relativePath: self.currentPath,
		Doc:          self.doc}
	route.postBuild()
	return route
}
