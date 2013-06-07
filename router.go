package restful

import "net/http"

// A RouteSelector finds the best matching Route given the input HTTP Request
type RouteSelector interface {

	// SelectWebService finds the WebService (typically by inspecting its Path) that matches the input URL path.
	SelectWebService(path string, webServices []*WebService) *WebService

	// SelectRoute finds a Route given the input HTTP Request and report if found (ok).
	// The HTTP writer is be used to directly communicate non-200 HTTP stati.
	SelectRoute(routes []Route, httpWriter http.ResponseWriter, httpRequest *http.Request) (selected Route, ok bool)
}
