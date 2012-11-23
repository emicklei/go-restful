package restful

import (
	"log"
	"net/http"
	"strings"
)

// Signature of function that can be bound to a Route.
type RouteFunction func(*Request, *Response)

// Route binds a HTTP Method,Path,Consumes combination to a RouteFunction.
type Route struct {
	Method   string
	Produces string // TODO make this a slice
	Consumes string // TODO make this a slice
	Path     string
	Function RouteFunction

	pathParts []string
}

func (self *Route) postBuild() {
	self.pathParts = strings.Split(self.Path, "/")
}

// If the Route matches the request then handle it and return http.StatusOK.
// Return other appropriate http status values otherwise.
func (self *Route) dispatch(httpWriter http.ResponseWriter, httpRequest *http.Request) int {
	log.Printf("restful: does %v matches Path: %v", httpRequest.URL.Path, self.Path)
	// the order of matching types are relevant
	matches, params := self.matchesPath(httpRequest.URL.Path)
	if !matches {
		return http.StatusNotFound
	}
	if self.Method != httpRequest.Method {
		return http.StatusMethodNotAllowed
	}
	accept := httpRequest.Header.Get("Accept")
	if !self.matchesAccept(accept) {
		return http.StatusUnsupportedMediaType
	}
	self.Function(&Request{httpRequest, params}, &Response{httpWriter, accept})
	return http.StatusOK
}

// Return whether the mimeType matches what this Route can consume.
func (self Route) matchesAccept(mimeType string) bool {
	log.Printf("restful: does %v matches Accept: %v", mimeType, self.Consumes)
	// cheap test first
	if len(self.Consumes) == 0 || strings.HasPrefix(self.Consumes, "*/*") {
		return true
	}
	parts := strings.Split(mimeType, ",")
	for _, each := range parts {
		if strings.HasPrefix(each, self.Consumes) {
			return true
		}
	}
	return false
}

// Check if the URL path matches the parameterized path of the Route.
// If it does then return a map(s->s) with the values for each path parameter.
func (self Route) matchesPath(urlPath string) (bool, map[string]string) {
	urlParts := strings.Split(urlPath, "/")
	if len(self.pathParts) != len(urlParts) {
		return false, nil
	}
	pathParameters := map[string]string{}
	for i, key := range self.pathParts {
		value := urlParts[i]
		if strings.HasPrefix(key, "{") { // path-parameter
			pathParameters[strings.Trim(key, "{}")] = value
		} else { // fixed
			if key != value {
				return false, nil
			}
		}
	}
	return true, pathParameters
}
