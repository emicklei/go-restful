package restful

import (

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