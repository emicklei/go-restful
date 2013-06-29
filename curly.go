package restful

import "net/http"

// CurlyRouter expects Routes with paths that contain zero or more parameters in curly brackets.
type CurlyRouter struct{}

// SelectRoute finds a Route given the input HTTP Request and report if found (ok).
// The HTTP writer is be used to directly communicate non-200 HTTP stati.
func (c CurlyRouter) SelectRoute(
	path string,
	webServices []*WebService,
	httpWriter http.ResponseWriter,
	httpRequest *http.Request) (selectedService *WebService, selected Route, ok bool) {

	// TODO
	return webServices[0], webServices[0].Routes()[0], true
}

func (c CurlyRouter) detectWebService(requestPath string, webServices []*WebService) (*WebService, string, error) {
	// TODO
	return webServices[0], "", nil
}
