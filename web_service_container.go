// Copyright (c) 2012 Ernest Micklei. All rights reserved.

package restful

import (
	"log"
	"net/http"
	"strings"
)

// The Dispatch function is responsible to delegating to the appropriate Dispatcher that has been registered via Add.
// The default implementation is DefaultDispatch which also does some basic panic handling.
//
// Example of overriding it to add basic request logging:
//	restful.Dispatch = func(w http.ResponseWriter, r *http.Request) {
//		fmt.Println(r.Method, r.URL)
//		restful.DefaultDispatch(w, r)
//	}
var Dispatch http.HandlerFunc

type Dispatcher interface {
	Routes() []Route
	RootPath() string
	PathParameters() []*Parameter
	//	rootRegEx
}

// Collection of registered Dispatchers that can handle Http requests
var webServices = []Dispatcher{}
var isRegisteredOnRoot = false

// Add registers a new Dispatcher add it to the http listeners.
func Add(service Dispatcher) {
	// If registered on root then no additional specific mapping is needed
	if !isRegisteredOnRoot {
		pattern := fixedPrefixPath(service.RootPath())
		// check if root path registration is needed
		if "/" == pattern || "" == pattern {
			http.HandleFunc("/", Dispatch)
			isRegisteredOnRoot = true
		} else {
			// detect if registration already exists
			alreadyMapped := false
			for _, each := range webServices {
				if each.RootPath() == service.RootPath() {
					alreadyMapped = true
					break
				}
			}
			if !alreadyMapped {
				http.HandleFunc(pattern, Dispatch)
				if !strings.HasSuffix(pattern, "/") {
					http.HandleFunc(pattern+"/", Dispatch)
				}
			}
		}
	}
	webServices = append(webServices, service)
}

// RegisteredWebServices returns the collections of added Dispatchers (WebService is an implementation)
func RegisteredWebServices() []Dispatcher {
	return webServices
}

// fixedPrefixPath returns the fixed part of the partspec ; it may include template vars {}
func fixedPrefixPath(pathspec string) string {
	varBegin := strings.Index(pathspec, "{")
	if -1 == varBegin {
		return pathspec
	}
	return pathspec[:varBegin]
}

// Dispatch the incoming Http Request to a matching Dispatcher.
// Matching algorithm is conform http://jsr311.java.net/nonav/releases/1.1/spec/spec.html, see jsr311.go
func DefaultDispatch(httpWriter http.ResponseWriter, httpRequest *http.Request) {
	// catch all for 500 response
	defer func() {
		if r := recover(); r != nil {
			log.Println("[restful] recover from panic situation:", r)
			httpWriter.WriteHeader(http.StatusInternalServerError)
			return
		}
	}()
	// step 1. Identify the root resource class (Dispatcher)
	dispatcher, finalMatch, err := detectDispatcher(httpRequest.URL.Path, webServices)
	if err != nil {
		httpWriter.WriteHeader(http.StatusNotFound)
		return
	}
	// step 2. Obtain the set of candidate methods (Routes)
	routes := selectRoutes(dispatcher, finalMatch)
	// step 3. Identify the method (Route) that will handle the request
	route, detected := detectRoute(routes, httpWriter, httpRequest)
	if detected {
		route.dispatch(httpWriter, httpRequest)
	}
	// else a non-200 response has already been written
}

func init() {
	Dispatch = DefaultDispatch
}
