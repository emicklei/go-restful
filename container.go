package restful

// Copyright 2013 Ernest Micklei. All rights reserved.
// Use of this source code is governed by a license
// that can be found in the LICENSE file.

import (
	"bytes"
	"fmt"
	//"github.com/emicklei/hopwatch"
	"log"
	"net/http"
	"runtime"
	"strings"
)

// Container holds a collection of WebServices and a http.ServeMux to dispatch http requests.
// The requests are further dispatched to routes of WebServices using a RouteSelector
type Container struct {
	webServices            []*WebService
	serveMux               *http.ServeMux
	isRegisteredOnRoot     bool
	containerFilters       []FilterFunction
	doNotRecover           bool          // default is false
	router                 RouteSelector // default is a RouterJSR311
	contentEncodingEnabled bool          // default is false
}

// NewContainer creates a new Container using a new ServeMux and default router (RouterJSR311)
func NewContainer() *Container {
	return &Container{
		webServices:            []*WebService{},
		serveMux:               http.NewServeMux(),
		isRegisteredOnRoot:     false,
		containerFilters:       []FilterFunction{},
		doNotRecover:           false,
		router:                 RouterJSR311{},
		contentEncodingEnabled: false}
}

// If DoNotRecover then panics will not be caught to return HTTP 500.
// In that case, Route functions are responsible for handling any error situation.
// Default value is false = recover from panics. This has performance implications.
func (c *Container) DoNotRecover(doNot bool) {
	c.doNotRecover = doNot
}

// Router changes the default Router (currently RouterJSR311)
func (c *Container) Router(aRouter RouteSelector) {
	c.router = aRouter
}

// EnableContentEncoding (default=false) allows for GZIP or DEFLATE encoding of responses.
func (c *Container) EnableContentEncoding(enabled bool) {
	c.contentEncodingEnabled = enabled
}

func (c *Container) Add(service *WebService) *Container {
	if service.pathExpr == nil {
		service.Path("") // lazy initialize path
	}
	// If registered on root then no additional specific mapping is needed
	if !c.isRegisteredOnRoot {
		pattern := c.fixedPrefixPath(service.RootPath())
		// check if root path registration is needed
		if "/" == pattern || "" == pattern {
			c.serveMux.HandleFunc("/", c.dispatch)
			c.isRegisteredOnRoot = true
		} else {
			// detect if registration already exists
			alreadyMapped := false
			for _, each := range c.webServices {
				if each.RootPath() == service.RootPath() {
					alreadyMapped = true
					break
				}
			}
			if !alreadyMapped {
				c.serveMux.HandleFunc(pattern, c.dispatch)
				if !strings.HasSuffix(pattern, "/") {
					c.serveMux.HandleFunc(pattern+"/", c.dispatch)
				}
			}
		}
	}
	c.webServices = append(c.webServices, service)
	//hopwatch.Dump(c)
	return c
}

// Dispatch the incoming Http Request to a matching WebService.
func (c *Container) dispatch(httpWriter http.ResponseWriter, httpRequest *http.Request) {
	// Instal panic recovery unless told otherwise
	if !c.doNotRecover { // catch all for 500 response
		defer func() {
			if r := recover(); r != nil {
				var buffer bytes.Buffer
				buffer.WriteString(fmt.Sprintf("[restful] recover from panic situation: - %v\r\n", r))
				for i := 1; ; i += 1 {
					_, file, line, ok := runtime.Caller(i)
					if !ok {
						break
					}
					buffer.WriteString(fmt.Sprintf("    %s:%d\r\n", file, line))
				}

				log.Println(buffer.String())
				httpWriter.WriteHeader(http.StatusInternalServerError)
				httpWriter.Write(buffer.Bytes())
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
	// assume without compression, test for override
	writer := httpWriter
	if c.contentEncodingEnabled {
		doCompress, encoding := wantsCompressedResponse(httpRequest)
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
		}
	}

	// Process any container filters
	if len(c.containerFilters) > 0 {
		wrappedRequest, wrappedResponse := newBasicRequestResponse(writer, httpRequest)
		proceed := false
		chain := FilterChain{Filters: c.containerFilters, Target: func(req *Request, resp *Response) {
			// we passed all filters
			proceed = true
		}}
		chain.ProcessFilter(wrappedRequest, wrappedResponse)
		if !proceed {
			return
		}
	}
	// Find best match Route ; detected is false if no match was found
	webService, route, detected := c.router.SelectRoute(
		c.webServices,
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

// fixedPrefixPath returns the fixed part of the partspec ; it may include template vars {}
func (c Container) fixedPrefixPath(pathspec string) string {
	varBegin := strings.Index(pathspec, "{")
	if -1 == varBegin {
		return pathspec
	}
	return pathspec[:varBegin]
}

// ServeHTTP implements net/http.Handler therefore a Container can be a Handler in a http.Server
func (c Container) ServeHTTP(httpwriter http.ResponseWriter, httpRequest *http.Request) {
	c.serveMux.ServeHTTP(httpwriter, httpRequest)
}

// Handle registers the handler for the given pattern. If a handler already exists for pattern, Handle panics.
func (c Container) Handle(pattern string, handler http.Handler) {
	c.serveMux.Handle(pattern, handler)
}

// Filter appends a container FilterFunction. These are called before dispatching
// a http.Request to a WebService from the container
func (c *Container) Filter(filter FilterFunction) {
	c.containerFilters = append(c.containerFilters, filter)
}

// RegisteredWebServices returns the collections of added WebServices
func (c Container) RegisteredWebServices() []*WebService {
	return c.webServices
}

// computeAllowedMethods returns a list of HTTP methods that are valid for a Request
func (c Container) computeAllowedMethods(req *Request) []string {
	// Go through all RegisteredWebServices() and all its Routes to collect the options
	methods := []string{}
	requestPath := req.Request.URL.Path
	for _, ws := range c.RegisteredWebServices() {
		matches := ws.pathExpr.Matcher.FindStringSubmatch(requestPath)
		if matches != nil {
			finalMatch := matches[len(matches)-1]
			for _, rt := range ws.Routes() {
				matches := rt.pathExpr.Matcher.FindStringSubmatch(finalMatch)
				if matches != nil {
					lastMatch := matches[len(matches)-1]
					if lastMatch == "" || lastMatch == "/" { // do not include if value is neither empty nor ‘/’.
						methods = append(methods, rt.Method)
					}
				}
			}
		}
	}
	// methods = append(methods, "OPTIONS")  not sure about this
	return methods
}

// newBasicRequestResponse creates a pair of Request,Response from its http versions.
// It is basic because no parameter or (produces) content-type information is given.
func newBasicRequestResponse(httpWriter http.ResponseWriter, httpRequest *http.Request) (*Request, *Response) {
	accept := httpRequest.Header.Get(HEADER_Accept)
	return &Request{httpRequest, map[string]string{}}, // empty parameters
		&Response{httpWriter, accept, []string{}} // empty content-types
}
