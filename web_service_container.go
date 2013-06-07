// Copyright (c) 2012 Ernest Micklei. All rights reserved.

package restful

import (
	"log"
	"net/http"
	"strings"
)

// The Dispatch function is responsible to delegating to the appropriate WebService that has been registered via Add.
// The default implementation is DefaultDispatch which also does some basic panic handling.
//
// Example of overriding it to add basic request logging:
//	restful.Dispatch = func(w http.ResponseWriter, r *http.Request) {
//		fmt.Println(r.Method, r.URL)
//		restful.DefaultDispatch(w, r)
//	}
var Dispatch http.HandlerFunc

// The Router is responsible for selecting the best matching Route given the input (request,response)
var Router RouteSelector

// Collection of registered Dispatchers that can handle Http requests
var webServices = []*WebService{}
var isRegisteredOnRoot = false
var globalFilters = []FilterFunction{}

// Add registers a new WebService add it to the http listeners.
func Add(service *WebService) {
	if service.pathExpr == nil {
		service.Path("") // lazy initialize path
	}
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

// Filter appends a global FilterFunction. These are called before dispatch a http.Request to a WebService.
func Filter(filter FilterFunction) {
	globalFilters = append(globalFilters, filter)
}

// RegisteredWebServices returns the collections of added Dispatchers (WebService is an implementation)
func RegisteredWebServices() []*WebService {
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

// Dispatch the incoming Http Request to a matching WebService.
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
	// step 0. Process any global filters
	if len(globalFilters) > 0 {
		wrappedRequest, wrappedResponse := newBasicRequestResponse(httpWriter, httpRequest)
		proceed := false
		chain := FilterChain{Filters: globalFilters, Target: func(req *Request, resp *Response) {
			// we passed all filters
			proceed = true
		}}
		chain.ProcessFilter(wrappedRequest, wrappedResponse)
		if !proceed {
			return
		}
	}
	// step 1. Identify the root resource class (WebService)
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
		// pass through filters (if any)
		filters := dispatcher.filters
		if len(filters) > 0 {
			wrappedRequest, wrappedResponse := newBasicRequestResponse(httpWriter, httpRequest)
			chain := FilterChain{Filters: filters, Target: func(req *Request, resp *Response) {
				// handle request by route
				route.dispatch(resp, req.Request)
			}}
			chain.ProcessFilter(wrappedRequest, wrappedResponse)
		} else {
			// handle request by route
			route.dispatch(httpWriter, httpRequest)
		}
	}
	// else a non-200 response has already been written
}

// newBasicRequestResponse create a pair of Request,Response from its http versions.
// It is basic because no parameter or (produces) content-type information is given.
func newBasicRequestResponse(httpWriter http.ResponseWriter, httpRequest *http.Request) (*Request, *Response) {
	accept := httpRequest.Header.Get(HEADER_Accept)
	return &Request{httpRequest, map[string]string{}}, // empty parameters
		&Response{httpWriter, accept, []string{}} // empty content-types
}

func init() {
	Dispatch = DefaultDispatch
}
