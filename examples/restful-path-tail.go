package main

import (
	. "github.com/emicklei/go-restful"
	"io"
	"log"
	"net/http"
)

// This example shows how to a Route that matches the "tail" of a path.
//
// GET http://localhost:8080/basepath/some/other/location/test.xml

func main() {
	DefaultContainer.Router(LoggingRouterJSR311{})
	ws := new(WebService)
	ws.Route(ws.GET("/basepath/{resource:.*}").To(staticFromPathParam))
	Add(ws)

	println("[go-restful] serve path tails from http://localhost:8080/basepath")
	http.ListenAndServe(":8080", nil)
}

func staticFromPathParam(req *Request, resp *Response) {
	io.WriteString(resp, req.PathParameter("resource"))
}

type LoggingRouterJSR311 struct {
	router RouterJSR311
}

func (l LoggingRouterJSR311) SelectRoute(
	webServices []*WebService,
	httpRequest *http.Request) (selectedService *WebService, selectedRoute *Route, err error) {
	log.Printf("SelectRoute\nwebServices=\n%v\nhttpRequest=%v\n)\nselectedService=%v\nselectedRoute=%v\nerror=%v\n", webServices, httpRequest, selectedService, selectedRoute, err)
	return l.router.SelectRoute(webServices, httpRequest)
}
