package restful

import (
	"net/http"
	"log"
)

type Dispatcher interface {
	Dispatch(http.ResponseWriter,*http.Request) bool
}

var webServices = []Dispatcher{} 

func Add(service Dispatcher) {
	log.Printf("Adding service: %#v\n", service)
	webServices = append(webServices, service)
}

func Dispatch(httpWriter http.ResponseWriter, httpRequest *http.Request) {
	for _, each := range webServices {
		if each.Dispatch(httpWriter, httpRequest) {
			break
		}
	}
}

func init() {
	log.Printf("Initializing go-restful\n")
	http.HandleFunc("/", Dispatch)
}