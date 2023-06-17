package main

import (
	"io"
	"log"
	"net/http"

	restful "github.com/emicklei/go-restful/v3"
)

// This example shows the effect of the 3.9.0 behavior.
//
// curl http://localhost:8080/say/hello -> to you
// curl http://localhost:8080/say/hello/ -> to you

func main() {
	ws := new(restful.WebService)
	ws.Path("/say")
	ws.Route(ws.GET("/hello").To(hello1))
	ws.Route(ws.GET("/hello/").To(hello2))
	restful.Add(ws)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func hello1(req *restful.Request, resp *restful.Response) {
	io.WriteString(resp, "world")
}

func hello2(req *restful.Request, resp *restful.Response) {
	io.WriteString(resp, "to you")
}
