package restful

import (
	"log"
	"net/http"
)

type Dispatcher interface {
	Dispatch(http.ResponseWriter, *http.Request) int
}

// Collection of registered Dispatchers that can handle Http requests
var webServices = []Dispatcher{}

// Register a new Dispatcher
func Add(service Dispatcher) {
	log.Printf("Adding service: %#v\n", service)
	webServices = append(webServices, service)
}

// Dispatch the incoming Http Request to the first registered Dispatcher that handled it
func Dispatch(httpWriter http.ResponseWriter, httpRequest *http.Request) {
	lastStatus := http.StatusNotFound
	for _, each := range webServices {
		lastStatus = each.Dispatch(httpWriter, httpRequest)
		if http.StatusOK == lastStatus {
			// response has been written
			return
		}
	}
	httpWriter.WriteHeader(lastStatus)
}

// Hook my Dispatch function as the standard Http handler
func init() {
	log.Printf("Initializing go-restful\n")
	http.HandleFunc("/", Dispatch)
}
