package restful

// Copyright 2013 Ernest Micklei. All rights reserved.
// Use of this source code is governed by a license
// that can be found in the LICENSE file.

import "net/http"

// A RouteSelector finds the best matching Route given the input HTTP Request
type RouteSelector interface {

	// SelectRoute finds a Route given the input HTTP Request and report if found (ok).
	// The HTTP writer is be used to directly communicate non-200 HTTP stati.
	SelectRoute(
		webServices []*WebService,
		httpWriter http.ResponseWriter,
		httpRequest *http.Request) (selectedService *WebService, selected *Route, err error)
}
