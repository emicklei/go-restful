package restful

// Copyright 2013 Ernest Micklei. All rights reserved.
// Use of this source code is governed by a license
// that can be found in the LICENSE file.

import "net/http"

// A RouteSelector finds the best matching Route given the input HTTP Request
type RouteSelector interface {

	// SelectRoute finds a Route given the input HTTP Request and a list of WebServices.
	// It returns a selected Route, its containing WebService and the path parameters
	// or an error indicating a problem.
	SelectRoute(
		webServices []*WebService,
		httpRequest *http.Request) (selectedService *WebService, selected *Route, pathParameters map[string]string, err error)
}
