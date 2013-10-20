package restful

// Copyright 2013 Ernest Micklei. All rights reserved.
// Use of this source code is governed by a license
// that can be found in the LICENSE file.

import (
	"log"
	"reflect"
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
	filters     []FilterFunction
	// documentation
	doc                     string
	operation               string
	readSample, writeSample interface{}
	parameters              []*Parameter
}

// To bind the route to a function.
// If this route is matched with the incoming Http Request then call this function with the *Request,*Response pair. Required.
func (b *RouteBuilder) To(function RouteFunction) *RouteBuilder {
	b.function = function
	return b
}

// Method specifies what HTTP method to match. Required.
func (b *RouteBuilder) Method(method string) *RouteBuilder {
	b.httpMethod = method
	return b
}

// Produces specifies what MIME types can be produced ; the matched one will appear in the Content-Type Http header.
func (b *RouteBuilder) Produces(mimeTypes ...string) *RouteBuilder {
	b.produces = mimeTypes
	return b
}

// Consumes specifies what MIME types can be consumes ; the Accept Http header must matched any of these
func (b *RouteBuilder) Consumes(mimeTypes ...string) *RouteBuilder {
	b.consumes = mimeTypes
	return b
}

// Path specifies the relative (w.r.t WebService root path) URL path to match. Default is "/".
func (b *RouteBuilder) Path(subPath string) *RouteBuilder {
	b.currentPath = subPath
	return b
}

// Doc tells what this route is all about. Optional.
func (b *RouteBuilder) Doc(documentation string) *RouteBuilder {
	b.doc = documentation
	return b
}

// Reads tells what resource type will be read from the request payload. Optional.
// A parameter of type "body" is added ,required is set to true and the dataType is set to the qualified name of the sample's type.
func (b *RouteBuilder) Reads(sample interface{}) *RouteBuilder {
	b.readSample = sample
	typeAsName := reflect.TypeOf(sample).String()
	bodyParameter := &Parameter{&ParameterData{Name: typeAsName}}
	bodyParameter.beBody()
	bodyParameter.Required(true)
	bodyParameter.DataType(typeAsName)
	b.Param(bodyParameter)
	return b
}

// Writes tells what resource type will be written as the response payload. Optional.
func (b *RouteBuilder) Writes(sample interface{}) *RouteBuilder {
	b.writeSample = sample
	return b
}

// Param allows you to document the parameters of the Route.
func (b *RouteBuilder) Param(parameter *Parameter) *RouteBuilder {
	if b.parameters == nil {
		b.parameters = []*Parameter{}
	}
	b.parameters = append(b.parameters, parameter)
	return b
}

// Operation allows you to document what the acutal method/function call is of the Route.
func (b *RouteBuilder) Operation(name string) *RouteBuilder {
	b.operation = name
	return b
}

func (b *RouteBuilder) servicePath(path string) *RouteBuilder {
	b.rootPath = path
	return b
}

// Filter appends a FilterFunction to the end of filters for this Route to build.
func (b *RouteBuilder) Filter(filter FilterFunction) *RouteBuilder {
	b.filters = append(b.filters, filter)
	return b
}

// If no specific Route path then set to rootPath
// If no specific Produces then set to rootProduces
// If no specific Consumes then set to rootConsumes
func (b *RouteBuilder) copyDefaults(rootProduces, rootConsumes []string) {
	if len(b.produces) == 0 {
		b.produces = rootProduces
	}
	if len(b.consumes) == 0 {
		b.consumes = rootConsumes
	}
}

// Build creates a new Route using the specification details collected by the RouteBuilder
func (b *RouteBuilder) Build() Route {
	pathExpr, err := newPathExpression(b.currentPath)
	if err != nil {
		log.Fatalf("[restful] Invalid path:%s because:%v", b.currentPath, err)
	}
	if b.function == nil {
		log.Fatalf("[restful] No function specified for route:" + b.currentPath)
	}
	route := Route{
		Method:        b.httpMethod,
		Path:          concatPath(b.rootPath, b.currentPath),
		Produces:      b.produces,
		Consumes:      b.consumes,
		Function:      b.function,
		Filters:       b.filters,
		relativePath:  b.currentPath,
		pathExpr:      pathExpr,
		Doc:           b.doc,
		Operation:     b.operation,
		ParameterDocs: b.parameters,
		ReadSample:    b.readSample,
		WriteSample:   b.writeSample}
	route.postBuild()
	return route
}

func concatPath(path1, path2 string) string {
	return strings.TrimRight(path1, "/") + "/" + strings.TrimLeft(path2, "/")
}
