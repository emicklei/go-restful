package restful

import "net/http"

// A RouteSelector finds the best matching Route given the input HTTP Request
type RouteSelector interface {

	// SelectRoute finds a Route given the input HTTP Request and report if found (ok).
	// The HTTP writer is be used to directly communicate non-200 HTTP stati.
	SelectRoute(
		path string,
		webServices []*WebService,
		httpWriter http.ResponseWriter,
		httpRequest *http.Request) (selectedService *WebService, selected Route, ok bool)
}
