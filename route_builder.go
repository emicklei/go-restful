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
	queryParametersDoc      map[string]string
	pathParametersDoc       map[string]string
}

// If this route is matched with the incoming Http Request then call this function with the *Request,*Response pair. Required.
func (self *RouteBuilder) To(function RouteFunction) *RouteBuilder {
	self.function = function
	return self
}

// Specify what HTTP method to match. Required.
func (self *RouteBuilder) Method(method string) *RouteBuilder {
	self.httpMethod = method
	return self
}

// Specify what MIME types can be produced ; the matched one will appear in the Content-Type Http header.
func (self *RouteBuilder) Produces(mimeTypes ...string) *RouteBuilder {
	self.produces = mimeTypes
	return self
}

// Specify what MIME types can be consumes ; the Accept Http header must matched any of these
func (self *RouteBuilder) Consumes(mimeTypes ...string) *RouteBuilder {
	self.consumes = mimeTypes
	return self
}

// Specify the relative (w.r.t WebService root path) URL path to match. Default is "/".
func (self *RouteBuilder) Path(subPath string) *RouteBuilder {
	self.currentPath = subPath
	return self
}

// Tell what this route is all about. Optional.
func (self *RouteBuilder) Doc(documentation string) *RouteBuilder {
	self.doc = documentation
	return self
}

// Tell what resource type will be read from the request payload. Optional.
func (self *RouteBuilder) Reads(sample interface{}) *RouteBuilder {
	//self.readSample = sample
	return self
}

// Tell what resource type will be written as the response payload. Optional.
func (self *RouteBuilder) Writes(sample interface{}) *RouteBuilder {
	//self.writeSample = sample
	return self
}

// Tell what the query param means. Optional.
func (self *RouteBuilder) QueryParam(name, comment string) *RouteBuilder {
	if self.queryParametersDoc == nil {
		self.queryParametersDoc = map[string]string{}
	}
	self.queryParametersDoc[name] = comment
	return self
}

// Tell what the path param means. Optional.
func (self *RouteBuilder) PathParam(name, comment string) *RouteBuilder {
	if self.pathParametersDoc == nil {
		self.pathParametersDoc = map[string]string{}
	}
	self.pathParametersDoc[name] = comment
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

// Create a new Route using the specification details collected by the RouteBuilder
func (self *RouteBuilder) Build() Route {
	route := Route{
		Method:       self.httpMethod,
		Path:         concatPath(self.rootPath, self.currentPath),
		Produces:     self.produces,
		Consumes:     self.consumes,
		Function:     self.function,
		relativePath: self.currentPath,
		Doc:          self.doc}
	route.postBuild()
	return route
}

func concatPath(path1, path2 string) string {
	return strings.TrimRight(path1, "/") + "/" + strings.TrimLeft(path2, "/")
}
