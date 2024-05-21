package main

import (
	"io"
	"log"
	"net/http"

	restful "github.com/emicklei/go-restful/v3"
)

// This example shows the minimal code needed to get a restful.WebService working.
//
// GET http://localhost:8080/hello

func main() {
	ws := new(restful.WebService)
	ws.Route(ws.GET("/hello").To(hello))
	restful.Add(ws)

	// DO NOT wrap http.ListenAndServe with log.Fatal in production
	// or you won't be able to drain in-flight request gracefully, even you handle sigterm
	log.Fatal(http.ListenAndServe(":8080", nil)) 
}

func hello(req *restful.Request, resp *restful.Response) {
	io.WriteString(resp, "world")
}
