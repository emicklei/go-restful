package main

import (
	. "github.com/emicklei/go-restful"
	"io"
	"log"
	"net/http"
)

// This example shows how to create a Route with google custom method
// Requires the use of a CurlyRouter and path should end with the custom method
//
// GET http://localhost:8080/resource/some-resource-id:init

func main() {
	DefaultContainer.Router(CurlyRouter{})
	ws := new(WebService)
	ws.Route(ws.GET("/basepath/{resourceId}:init").To(fromPathParam))
	Add(ws)

	println("[go-restful] serve path tails from http://localhost:8080/basepath")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func fromPathParam(req *Request, resp *Response) {
	io.WriteString(resp, "resourceId: "+req.PathParameter("resourceId"))
}
