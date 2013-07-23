package main

import (
	"github.com/emicklei/go-restful"
	"net/http"
)

func main() {
	ws := new(restful.WebService)
	ws.Route(ws.GET("/secret").Filter(basicAuthenticate).To(secret))
	restful.Add(ws)
	http.ListenAndServe(":8080", nil)
}

func basicAuthenticate(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	encoded := req.Request.Header.Get("Authorization")
	// usr/pwd = admin/admin
	// real code does some decoding
	if len(encoded) == 0 || "Basic YWRtaW46YWRtaW4=" != encoded {
		resp.AddHeader("WWW-Authenticate", "Basic realm=Protected Area")
		resp.WriteHeader(401)
		resp.Write([]byte("401: Not Authorized"))
		return
	}
	chain.ProcessFilter(req, resp)
}

func secret(req *restful.Request, resp *restful.Response) {
	resp.Write([]byte("42"))
}
