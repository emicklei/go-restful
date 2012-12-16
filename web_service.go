package restful

// WebService holds a collection of Route values that bind a Http Method + URL Path to a function.
type WebService struct {
	rootPath       string
	routes         []Route
	produces       []string
	consumes       []string
	pathParameters []*Parameter
}

// Path specifies the root URL template path of the WebService.
// All Routes will be relative to this path.
func (self *WebService) Path(root string) *WebService {
	self.rootPath = root
	return self
}

// AddParameter adds a PathParameter to document parameters used in the root path.
func (self *WebService) Param(parameter *Parameter) *WebService {
	if self.pathParameters == nil {
		self.pathParameters = []*Parameter{}
	}
	self.pathParameters = append(self.pathParameters, parameter)
	return self
}

// PathParameter creates a new Parameter of kind Path for documentation purposes.
func (self *WebService) PathParameter(name, description string) *Parameter {
	p := &Parameter{name: name, description: description, required: true}
	p.bePath()
	return p
}

// QueryParameter creates a new Parameter of kind Query for documentation purposes.
func (self *WebService) QueryParameter(name, description string) *Parameter {
	p := &Parameter{name: name, description: description, required: false}
	p.beQuery()
	return p
}

// BodyParameter creates a new Parameter of kind Body for documentation purposes.
func (self *WebService) BodyParameter(name, description string) *Parameter {
	p := &Parameter{name: name, description: description, required: true}
	p.beBody()
	return p
}

// Route creates a new Route using the RouteBuilder and add to the ordered list of Routes.
func (self *WebService) Route(builder *RouteBuilder) *WebService {
	builder.copyDefaults(self.produces, self.consumes)
	self.routes = append(self.routes, builder.Build())
	return self
}

// Method creates a new RouteBuilder and initialize its http method
func (self *WebService) Method(httpMethod string) *RouteBuilder {
	return new(RouteBuilder).servicePath(self.rootPath).Method(httpMethod)
}

// Produces specifies that this WebService can produce one or more MIME types.
func (self *WebService) Produces(contentTypes ...string) *WebService {
	self.produces = contentTypes
	return self
}

// Produces specifies that this WebService can consume one or more MIME types.
func (self *WebService) Consumes(accepts ...string) *WebService {
	self.consumes = accepts
	return self
}

// Routes returns the Routes associated with this WebService
func (self WebService) Routes() []Route {
	return self.routes
}

// RootPath returns the RootPath associated with this WebService. Default "/"
func (self WebService) RootPath() string {
	return self.rootPath
}

func (self WebService) PathParameters() []*Parameter {
	return self.pathParameters
}

/*
	Convenience methods
*/

// GET is a shortcut for .Method("GET").Path(subPath)
func (self *WebService) GET(subPath string) *RouteBuilder {
	return new(RouteBuilder).servicePath(self.rootPath).Method("GET").Path(subPath)
}

// POST is a shortcut for .Method("POST").Path(subPath)
func (self *WebService) POST(subPath string) *RouteBuilder {
	return new(RouteBuilder).servicePath(self.rootPath).Method("POST").Path(subPath)
}

// PUT is a shortcut for .Method("PUT").Path(subPath)
func (self *WebService) PUT(subPath string) *RouteBuilder {
	return new(RouteBuilder).servicePath(self.rootPath).Method("PUT").Path(subPath)
}

// DELETE is a shortcut for .Method("DELETE").Path(subPath)
func (self *WebService) DELETE(subPath string) *RouteBuilder {
	return new(RouteBuilder).servicePath(self.rootPath).Method("DELETE").Path(subPath)
}
