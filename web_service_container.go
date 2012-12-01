package restful

// http://jsr311.java.net/nonav/releases/1.1/spec/spec.html

import (
	"encoding/xml"
	"github.com/emicklei/go-restful/wadl"
	"log"
	"net/http"
)

type Dispatcher interface {
	Dispatch(http.ResponseWriter, *http.Request)
	Routes() []Route
	RootPath() string
	//	rootRegEx
}

// Collection of registered Dispatchers that can handle Http requests
var webServices = []Dispatcher{}

// Register a new Dispatcher
func Add(service Dispatcher) {
	webServices = append(webServices, service)
}

// Dispatch the incoming Http Request to a matching Dispatcher.
// Matching algoritm is conform http://jsr311.java.net/nonav/releases/1.1/spec/spec.html

func Dispatch(httpWriter http.ResponseWriter, httpRequest *http.Request) {
	match, finalMatch, err := detectDispatcher(httpRequest.URL.Path, webServices)
	log.Printf("final match:&v", finalMatch)
	if err == nil {
		match.Dispatch(httpWriter, httpRequest)
	} else {
		httpWriter.WriteHeader(http.StatusNotFound)
	}
}

// Hook my Dispatch function as the standard Http handler
func init() {
	log.Printf("restful: register the Dispatch function to the Default Http handlers.\n")
	http.HandleFunc("/", Dispatch)
}

// Return the api in XML
func Wadl(base string) string {
	resources := wadl.Resources{Base: base}
	for _, eachWebService := range webServices {
		for _, eachRoute := range eachWebService.Routes() {
			response := wadl.Response{}
			for _, mimeType := range eachRoute.Produces {
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
