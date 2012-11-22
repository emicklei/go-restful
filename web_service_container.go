package restful

import (
	"log"
	"net/http"
)

type Dispatcher interface {
	Dispatch(http.ResponseWriter, *http.Request) bool
}

// Collection of registered WebServices that can handle Http requests
var webServices = []Dispatcher{}

// Register a new WebService (a Dispatcher)
func Add(service Dispatcher) {
	log.Printf("Adding service: %#v\n", service)
	webServices = append(webServices, service)
}

// Dispatch the incoming Http Request to the first registered WebServices that handled it
func Dispatch(httpWriter http.ResponseWriter, httpRequest *http.Request) {
	for _, each := range webServices {
		if each.Dispatch(httpWriter, httpRequest) {
			break
		}
	}
}

// Hook my Dispatch function as the standard Http handler
func init() {
	log.Printf("Initializing go-restful\n")
	http.HandleFunc("/", Dispatch)
}