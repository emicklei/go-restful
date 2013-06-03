package main

import (
	"github.com/emicklei/go-restful"
	"net/http"
	"path"
)

var rootdir = "/tmp"

func main() {

	ws := new(restful.WebService)
	ws.Route(ws.GET("/static/{resource}").To(staticFromPathParam))
	ws.Route(ws.GET("/static").To(staticFromQueryParam))
	restful.Add(ws)

	println("[go-restful] serving files on http://localhost:8080/static from local /tmp")
	http.ListenAndServe(":8080", nil)
}

// http://localhost:8080/static/test.xml
// http://localhost:8080/static/
func staticFromPathParam(req *restful.Request, resp *restful.Response) {
	http.ServeFile(
		resp.ResponseWriter,
		req.Request,
		path.Join(rootdir, req.PathParameter("resource")))
}

// http://localhost:8080/static?resource=subdir/test.xml
func staticFromQueryParam(req *restful.Request, resp *restful.Response) {
	http.ServeFile(
		resp.ResponseWriter,
		req.Request,
		path.Join(rootdir, req.QueryParameter("resource")))
}
