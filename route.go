package restful

import (
	"strings"
	"net/http"
)

// Signature of function that can be bound to a Route
type RouteFunction func(*Request, http.ResponseWriter)

// Route binds a HTTP Method,Path,Consumes combination to a RouteFunction
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
// If the Route matches the request then handle it and return true ; false otherwise
func (self *Route) dispatch(httpWriter http.ResponseWriter, httpRequest *http.Request) bool {
	if (self.Method != httpRequest.Method) {
		return false
	}
	matches, params := self.MatchesPath(httpRequest.URL.Path)
	if (!matches) {
		return false
	}
	// TODO match accept
	//writerWrapper := responseWriter{httpWriter}
	restRequest := Request{httpRequest,params}
	restResponse := Response{httpWriter}	
	self.Function(&restRequest,restResponse)
	return true
}

func (self Route) MatchesPath(urlPath string) (bool, map[string]string) {
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
