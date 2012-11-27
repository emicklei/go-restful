package restful

import (
	"encoding/xml"
	"github.com/emicklei/go-restful/wadl"
	"log"
	"net/http"
	"strings"
)

type Dispatcher interface {
	Dispatch(http.ResponseWriter, *http.Request)
	Routes() []Route
	RootPath() string
}

// Collection of registered Dispatchers that can handle Http requests
var webServices = map[string]Dispatcher{}

// Register a new Dispatcher
func Add(service Dispatcher) {
	routedService, present := webServices[service.RootPath()]
	if present {
		log.Panicf("restful: conflict with registered service :%v", routedService)
	}
	webServices[service.RootPath()] = service
}

// Dispatch the incoming Http Request to a matching Dispatcher.
// A Dispatcher is matched when the request URL path starts with the Dispatcher's root path.
func Dispatch(httpWriter http.ResponseWriter, httpRequest *http.Request) {
	requestPath := httpRequest.URL.Path
	for rootPath, each := range webServices {
		if requestPath == rootPath || strings.HasPrefix(requestPath, rootPath+"/") {
			each.Dispatch(httpWriter, httpRequest)
			return
		}
	}
	httpWriter.WriteHeader(http.StatusNotFound)
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
