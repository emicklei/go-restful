package restful

import (
	"net/http"
)

// WebServiceContainer encapsulates WebService objects
// for Handling an incoming Http Request by Dispatching it to
// the first WebService,Route combination that matches.
type WebServiceContainer struct {
	services []WebService
}

func (self WebServiceContainer) Add(service WebService) {
	self.services = append(self.services, service)
}

func (self WebServiceContainer) ServeHTTP(httpWriter http.ResponseWriter, httpRequest *http.Request) {
	for _,each := range self.services {
		if each.dispatch(httpWriter,httpRequest) {
			break
		}	
	}
}