// Copyright (c) 2012 Ernest Micklei. All rights reserved.

package restful

import (
	"log"
	"net/http"
	"strings"
)

// The Dispatch function is responsible to delegating to the appropriate WebService that has been registered via Add.
// The default implementation is DefaultDispatch which also does some basic panic handling and closes the request body.
//
// Example of overriding it to add basic request logging:
//	restful.Dispatch = func(w http.ResponseWriter, r *http.Request) {
//		fmt.Println(r.Method, r.URL)
//		restful.DefaultDispatch(w, r)
//	}
//
//  Deprecated: Use filters instead.
//
var Dispatch http.HandlerFunc

// If set the true then panics will not be caught to return HTTP 500.
// In that case, Route functions are responsible for handling any error situation.
// Default value is false = recover from panics. This has performance implications.
var DoNotRecover = false

// The Router is responsible for selecting the best matching Route given the input (request,response)
// See jsr311.go
var Router = RouterJSR311{}

// Collection of registered WebServices that can handle Http requests
var webServices = []*WebService{}

// Remember if any WebService is mapped on root /
var isRegisteredOnRoot = false

// Collection of Filter functions that apply to all requests.
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
	// Instal panic recovery unless told otherwise
	if !DoNotRecover { // catch all for 500 response
		defer func() {
			if r := recover(); r != nil {
				log.Println("[restful] recover from panic situation:", r)
				httpWriter.WriteHeader(http.StatusInternalServerError)
				return
			}
		}()
	}
	// Install closing the request body (if any)
	defer func() {
		if nil != httpRequest.Body {
			httpRequest.Body.Close()
		}
	}()

	// Detect if compression is needed
	doCompress, encoding := WantsCompressedResponse(httpRequest)
	var writer http.ResponseWriter
	if doCompress {
		var err error
		writer, err = NewCompressingResponseWriter(httpWriter, encoding)
		if err != nil {
			log.Println("[restful] unable to install compressor:", err)
			httpWriter.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer func() {
			writer.(*CompressingResponseWriter).Close()
		}()
	} else {
		// without compression
		writer = httpWriter
	}

	// Process any global filters
	if len(globalFilters) > 0 {
		wrappedRequest, wrappedResponse := newBasicRequestResponse(writer, httpRequest)
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
	// Find best match Route ; detected is false if no match was found
	webService, route, detected := Router.SelectRoute(
		httpRequest.URL.Path,
		webServices,
		httpWriter,
		httpRequest)
	if detected {
		// pass through filters (if any)
		filters := webService.filters
		wrappedRequest, wrappedResponse := route.wrapRequestResponse(writer, httpRequest)
		if len(filters) > 0 {
			chain := FilterChain{Filters: filters, Target: func(req *Request, resp *Response) {
				// handle request by route after passed all filters
				route.dispatch(wrappedRequest, wrappedResponse)
			}}
			chain.ProcessFilter(wrappedRequest, wrappedResponse)
		} else {
			// handle request by route
			route.dispatch(wrappedRequest, wrappedResponse)
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
