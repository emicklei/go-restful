package restful

import (
	"encoding/xml"
	"github.com/emicklei/go-restful/wadl"
	"log"
	"net/http"
	"strings"
)

type Dispatcher interface {
	Dispatch(http.ResponseWriter, *http.Request) int
	Routes() []Route
}

// Collection of registered Dispatchers that can handle Http requests
var webServices = []Dispatcher{}

// Register a new Dispatcher
func Add(service Dispatcher) {
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
	log.Printf("restful: initializing\n")
	http.HandleFunc("/", Dispatch)
}

// Return the api in XML
func Wadl(base string) string {
	resources := wadl.Resources{Base: base}
	for _, eachWebService := range webServices {

		for _, eachRoute := range eachWebService.Routes() {
			response := wadl.Response{}
			for _, mimeType := range strings.Split(eachRoute.Produces, ",") {
				response.AddRepresentation(wadl.Representation{MediaType: mimeType})
			}
			method := wadl.Method{Name: eachRoute.Method, Response: response}
			resource := wadl.Resource{Path: eachRoute.Path, Method: method}
			resources.AddResource(resource)
		}
	}
	app := wadl.Application{Resources: resources}
	bytes, _ := xml.MarshalIndent(app, "", " ")
	return string(bytes)
}
