// Copyright 2012 Ernest Micklei. All rights reserved.
// Use of this source code is governed by a license
// that can be found in the LICENSE file.

package restful

import (
	"net/http"
	"strings"
)

// Signature of a function that can be bound to a Route.
type RouteFunction func(*Request, *Response)

// Route binds a HTTP Method,Path,Consumes combination to a RouteFunction.
type Route struct {
	Method   string
	Produces []string
	Consumes []string
	Path     string // webservice root path + described path
	Function RouteFunction
	Filters  []FilterFunction

	// cached values for dispatching
	relativePath string
	pathParts    []string
	pathExpr     *pathExpression // cached compilation of relativePath as RegExp

	// documentation
	Doc                     string
	Operation               string
	ParameterDocs           []*Parameter
	ReadSample, WriteSample interface{} // structs that model an example request or response payload
}

// Initialize for Route
func (self *Route) postBuild() {
	self.pathParts = tokenizePath(self.Path)
}

// Create Request and Response from their http versions
func (self *Route) wrapRequestResponse(httpWriter http.ResponseWriter, httpRequest *http.Request) (*Request, *Response) {
	params := self.extractParameters(httpRequest.URL.Path)
	accept := httpRequest.Header.Get(HEADER_Accept)
	wrappedRequest := &Request{httpRequest, params}
	wrappedResponse := &Response{httpWriter, accept, self.Produces}
	return wrappedRequest, wrappedResponse
}

// Extract any path parameters from the the request URL path and call the function
func (self *Route) dispatch(wrappedRequest *Request, wrappedResponse *Response) {
	if len(self.Filters) > 0 {
		chain := FilterChain{Filters: self.Filters, Target: self.Function}
		chain.ProcessFilter(wrappedRequest, wrappedResponse)
	} else {
		// unfiltered
		self.Function(wrappedRequest, wrappedResponse)
	}
}

// Return whether the mimeType matches to what this Route can produce.
func (self Route) matchesAccept(mimeTypesWithQuality string) bool {
	parts := strings.Split(mimeTypesWithQuality, ",")
	for _, each := range parts {
		withoutQuality := strings.Split(each, ";")[0]
		if withoutQuality == "*/*" {
			return true
		}
		for _, other := range self.Produces {
			if other == withoutQuality {
				return true
			}
		}
	}
	return false
}

// Return whether the mimeType matches to what this Route can consume.
func (self Route) matchesContentType(mimeTypes string) bool {
	parts := strings.Split(mimeTypes, ",")
	for _, each := range parts {
		var contentType string
		if strings.Contains(each, ";") {
			contentType = strings.Split(each, ";")[0]
		} else {
			contentType = each
		}
		for _, other := range self.Consumes {
			if other == "*/*" || other == contentType {
				return true
			}
		}
	}
	return false
}

// Extract the parameters from the request url path
func (self Route) extractParameters(urlPath string) map[string]string {
	urlParts := tokenizePath(urlPath)
	pathParameters := map[string]string{}
	for i, key := range self.pathParts {
		var value string
		if i >= len(urlParts) {
			value = ""
		} else {
			value = urlParts[i]
		}
		if strings.HasPrefix(key, "{") { // path-parameter
			pathParameters[strings.Trim(key, "{}")] = value
		}
	}
	return pathParameters
}

// Tokenize an URL path using the slash separator ; the result does not have empty tokens
func tokenizePath(path string) []string {
	if "/" == path {
		return []string{}
	}
	return strings.Split(strings.Trim(path, "/"), "/")
}

// for debugging
func (r Route) String() string {
	return r.Method + " " + r.Path
}
