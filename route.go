package restful

import (
	"strings"
)

// Signature of function that can be bound to a Route
type RouteFunction func(*Request, *Response)

// Route binds a HTTP Method,Path,Consumes combination to a RouteFunction
type Route struct {
	Method   string
	Produces string // TODO make this a slice
	Consumes string // TODO make this a slice
	Path     string
	Function RouteFunction
	
	pathParts []string
}

func (self Route) MatchesPath(urlPath string) (bool, map[string]string) {
	self.pathParts = strings.Split(self.Path,"/")
	urlParts := strings.Split(urlPath,"/")
	if len(self.pathParts) != len(urlParts) {
		return false, nil
	}
	pathParameters := map[string]string{}
	for i,key := range self.pathParts {
		value := urlParts[i]
		if strings.HasPrefix(key,"{") { // path-parameter
			pathParameters[strings.Trim(key,"{}")]=value
		} else { // fixed
			if (key != value) {
				return false, nil
			}	
		}		
	} 
	return true, pathParameters
}