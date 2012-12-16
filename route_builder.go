package restful

import (
	"strings"
)

// RouteBuilder is a helper to construct Routes.
type RouteBuilder struct {
	rootPath    string
	currentPath string
	produces    []string
	consumes    []string
	httpMethod  string        // required
	function    RouteFunction // required
	// documentation
	doc                     string
	readSample, writeSample string
	parameters              []*Parameter
}

// To bind the route to a function. 
// If this route is matched with the incoming Http Request then call this function with the *Request,*Response pair. Required.
func (self *RouteBuilder) To(function RouteFunction) *RouteBuilder {
	self.function = function
	return self
}

// Method specifies what HTTP method to match. Required.
func (self *RouteBuilder) Method(method string) *RouteBuilder {
	self.httpMethod = method
	return self
}

// Produces specifies what MIME types can be produced ; the matched one will appear in the Content-Type Http header.
func (self *RouteBuilder) Produces(mimeTypes ...string) *RouteBuilder {
	self.produces = mimeTypes
	return self
}

// Specify what MIME types can be consumes ; the Accept Http header must matched any of these
func (self *RouteBuilder) Consumes(mimeTypes ...string) *RouteBuilder {
	self.consumes = mimeTypes
	return self
}

// Path specifies the relative (w.r.t WebService root path) URL path to match. Default is "/".
func (self *RouteBuilder) Path(subPath string) *RouteBuilder {
	self.currentPath = subPath
	return self
}

// Doc tells what this route is all about. Optional.
func (self *RouteBuilder) Doc(documentation string) *RouteBuilder {
	self.doc = documentation
	return self
}

// Reads tells what resource type will be read from the request payload. Optional.
func (self *RouteBuilder) Reads(sample interface{}) *RouteBuilder {
	//self.readSample = sample
	return self
}

// Writes tells what resource type will be written as the response payload. Optional.
func (self *RouteBuilder) Writes(sample interface{}) *RouteBuilder {
	//self.writeSample = sample
	return self
}

// Param allows you to document the parameters of the Route.
func (self *RouteBuilder) Param(parameter *Parameter) *RouteBuilder {
	if self.parameters == nil {
		self.parameters = []*Parameter{}
	}
	self.parameters = append(self.parameters, parameter)
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

// Build creates a new Route using the specification details collected by the RouteBuilder
func (self *RouteBuilder) Build() Route {
	route := Route{
		Method:        self.httpMethod,
		Path:          concatPath(self.rootPath, self.currentPath),
		Produces:      self.produces,
		Consumes:      self.consumes,
		Function:      self.function,
		relativePath:  self.currentPath,
		Doc:           self.doc,
		parameterDocs: self.parameters}
	route.postBuild()
	return route
}

func concatPath(path1, path2 string) string {
	return strings.TrimRight(path1, "/") + "/" + strings.TrimLeft(path2, "/")
}
