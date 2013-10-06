package main

import (
	"github.com/emicklei/go-restful"
	"io"
	"log"
	"net/http"
)

var indentLevel int

func container_filter_A(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	trace("container_filter_A: before", 1)
	chain.ProcessFilter(req, resp)
	trace("container_filter_A: after", -1)
}

func container_filter_B(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	trace("container_filter_B: before", 1)
	chain.ProcessFilter(req, resp)
	trace("container_filter_B: after", -1)
}

func service_filter_A(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	trace("service_filter_A: before", 1)
	chain.ProcessFilter(req, resp)
	trace("service_filter_A: after", -1)
}

func service_filter_B(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	trace("service_filter_B: before", 1)
	chain.ProcessFilter(req, resp)
	trace("service_filter_B: after", -1)
}

func route_filter_A(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	trace("route_filter_A: before", 1)
	chain.ProcessFilter(req, resp)
	trace("route_filter_A: after", -1)
}

func route_filter_B(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	trace("route_filter_B: before", 1)
	chain.ProcessFilter(req, resp)
	trace("route_filter_B: after", -1)
}

func trace(what string, delta int) {
	indented := what
	for t := 0; t < indentLevel; t++ {
		indented = "." + indented
	}
	log.Printf("%s", indented)
	indentLevel += delta
}

func main() {
	restful.Filter(container_filter_A)
	restful.Filter(container_filter_B)

	ws1 := new(restful.WebService)
	ws1.Filter(service_filter_A)
	ws1.Filter(service_filter_B)
	ws1.Route(ws1.GET("/1").To(doitA).Filter(route_filter_A).Filter(route_filter_B))

	ws2 := new(restful.WebService)
	ws2.Filter(service_filter_A)
	ws2.Filter(service_filter_B)
	ws2.Route(ws2.GET("/2").To(doitB).Filter(route_filter_A).Filter(route_filter_B))

	restful.Add(ws1)
	restful.Add(ws2)

	log.Print("go-restful example listing on http://localhost:8080/1 and http://localhost:8080/2")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func doitA(req *restful.Request, resp *restful.Response) {
	io.WriteString(resp, "nothing to see in 1")
}

func doitB(req *restful.Request, resp *restful.Response) {
	io.WriteString(resp, "nothing to see in B")
}
