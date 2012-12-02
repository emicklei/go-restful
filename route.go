// Copyright 2012 Ernest Micklei. All rights reserved.
// Use of this source code is governed by a license 
// that can be found in the LICENSE file.

package restful

import (
	"net/http"
	"strings"
)

const RouteFunctionCalled = 0

// Signature of function that can be bound to a Route.
type RouteFunction func(*Request, *Response)

// Route binds a HTTP Method,Path,Consumes combination to a RouteFunction.
type Route struct {
	Method   string
	Produces []string
	Consumes []string
	Path     string
	Function RouteFunction

	relativePath string
	pathParts    []string
}

func (self *Route) postBuild() {
	self.pathParts = strings.Split(self.Path, "/")
}

// Extract any path parameters from the the request URL path and call the function
func (self *Route) dispatch(httpWriter http.ResponseWriter, httpRequest *http.Request) {
	params := self.extractParameters(httpRequest.URL.Path)
	accept := httpRequest.Header.Get(HEADER_Accept)
	self.Function(&Request{httpRequest, params}, &Response{httpWriter, accept})
}

// Return whether the mimeType matches what this Route can produce.
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

// Return whether the mimeType matches what this Route can consume.
func (self Route) matchesContentType(mimeTypes string) bool {
	parts := strings.Split(mimeTypes, ",")
	for _, each := range parts {
		for _, other := range self.Consumes {
			if other == "*/*" || other == each {
				return true
			}
		}
	}
	return false
}

// Extract the parameters from the urlPath
func (self Route) extractParameters(urlPath string) map[string]string {
	urlParts := strings.Split(urlPath, "/")
	pathParameters := map[string]string{}
	for i, key := range self.pathParts {
		value := urlParts[i]
		if strings.HasPrefix(key, "{") { // path-parameter
			pathParameters[strings.Trim(key, "{}")] = value
		}
	}
	return pathParameters
}
