// Copyright (c) 2012 Ernest Micklei. All rights reserved.

package restful

import (
	"net/http"
	"strings"
)

type Dispatcher interface {
	Routes() []Route
	RootPath() string
	//	rootRegEx
}

// Collection of registered Dispatchers that can handle Http requests
var webServices = []Dispatcher{}
var isRegisteredOnRoot = false

// Register a new Dispatcher add it to the http listeners.
// Check its root path to see if 
func Add(service Dispatcher) {
	webServices = append(webServices, service)
	path := service.RootPath()
	if len(service.RootPath()) == 0 {
		path = "/"
	} else {
		if varIndex := strings.Index(path, "{"); varIndex != -1 {
			// Use the fixed part of the service rootpath
			path = service.RootPath()[:varIndex]
		}
	}
	if path == "/" {
		// Have to listen to / , but hook only once		
		if !isRegisteredOnRoot {
			http.HandleFunc("/", Dispatch)
			isRegisteredOnRoot = true
		}
	} else {
		http.HandleFunc(path, Dispatch)
	}
}

// Dispatch the incoming Http Request to a matching Dispatcher.
// Matching algorithm is conform http://jsr311.java.net/nonav/releases/1.1/spec/spec.html, see jsr311.go
func Dispatch(httpWriter http.ResponseWriter, httpRequest *http.Request) {
	dispatcher, finalMatch, err := detectDispatcher(httpRequest.URL.Path, webServices)
	if err != nil {
		httpWriter.WriteHeader(http.StatusNotFound)
		return
	}
	routes := selectRoutes(dispatcher, finalMatch)
	route, detected := detectRoute(routes, httpWriter, httpRequest)
	if detected {
		route.dispatch(httpWriter, httpRequest)
	}
	// not detected also means that a response has been written
}
