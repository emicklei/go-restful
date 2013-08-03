// Copyright (c) 2012 Ernest Micklei. All rights reserved.

package restful

import (
	"net/http"
)

// DefaultContainer is a restful.Container that uses http.DefaultServeMux
var DefaultContainer *Container

func init() {
	DefaultContainer = NewContainer()
	DefaultContainer.serveMux = http.DefaultServeMux
}

// If set the true then panics will not be caught to return HTTP 500.
// In that case, Route functions are responsible for handling any error situation.
// Default value is false = recover from panics. This has performance implications.
// OBSOLETE ; use restful.DefaultContainer.DoNotRecover(true)
var DoNotRecover = false

// Add registers a new WebService add it to the http listeners.
func Add(service *WebService) {
	DefaultContainer.Add(service)
}

// Filter appends a global FilterFunction. These are called before dispatch a http.Request to a WebService.
func Filter(filter FilterFunction) {
	DefaultContainer.Filter(filter)
}

// RegisteredWebServices returns the collections of WebServices
func RegisteredWebServices() []*WebService {
	return DefaultContainer.RegisteredWebServices()
}
