package main

import (
	"io"
	"log"
	"net/http"

	restful "github.com/emicklei/go-restful/v3"
)

// This example shows the minimal code needed to get a restful.WebService working.
//
// curl http://localhost:8080/hello -> world
// curl http://localhost:8080/hello/ -> to you

func main() {
	restful.MergePathStrategy = restful.TrimSlashStrategy
	ws := new(restful.WebService)
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
