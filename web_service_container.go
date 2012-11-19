package restful

import (
	"net/http"
)

type Dispatcher interface {
	Dispatch(http.ResponseWriter,*http.Request) bool
}

// WebServiceContainer hold a collection of Dispatcher implementations
// Dispatching is the process of delegating a Http request to a Route function 
type WebServiceContainer struct {
	services []Dispatcher
}

func (self WebServiceContainer) Add(service Dispatcher) {
	self.services = append(self.services, service)
}

func (self WebServiceContainer) Dispatch(httpWriter http.ResponseWriter, httpRequest *http.Request) {
	for _, each := range self.services {
		if each.Dispatch(httpWriter, httpRequest) {
			break
		}
	}
}
