package main

import (
	"github.com/emicklei/go-restful"
	"net/http"
	"path"
)

// This example shows how to define methods that serve static files
// It uses the standard http.ServeFile method
//
// GET http://localhost:8080/static/test.xml
// GET http://localhost:8080/static/
//
// GET http://localhost:8080/static?resource=subdir/test.xml

var rootdir = "/tmp"

func main() {

	ws := new(restful.WebService)
	ws.Route(ws.GET("/static/{resource}").To(staticFromPathParam))
	ws.Route(ws.GET("/static").To(staticFromQueryParam))
	restful.Add(ws)

	println("[go-restful] serving files on http://localhost:8080/static from local /tmp")
	http.ListenAndServe(":8080", nil)
}

func staticFromPathParam(req *restful.Request, resp *restful.Response) {
	http.ServeFile(
		resp.ResponseWriter,
		req.Request,
		path.Join(rootdir, req.PathParameter("resource")))
}

func staticFromQueryParam(req *restful.Request, resp *restful.Response) {
	http.ServeFile(
		resp.ResponseWriter,
		req.Request,
		path.Join(rootdir, req.QueryParameter("resource")))
}
