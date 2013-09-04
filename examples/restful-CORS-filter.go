package main

import (
	"github.com/emicklei/go-restful"
)

// Cross-origin resource sharing (CORS) is a mechanism that allows JavaScript on a web page
// to make XMLHttpRequests to another domain, not the domain the JavaScript originated from.
//
// http://en.wikipedia.org/wiki/Cross-origin_resource_sharing
// http://enable-cors.org/server.html

func applyCORSFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	// If the origin Header was set then return it in the response
	if origin := req.Request.Header.Get("Origin"); origin != "" {
		resp.AddHeader("Access-Control-Allow-Origin", origin)
	} else {
		// Otherwise set to Allow All domains (see docs)
		resp.AddHeader("Access-Control-Allow-Origin", "*")
	}

	resp.AddHeader("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	resp.AddHeader("Access-Control-Allow-Headers", "Content-Type, Accept")
	chain.ProcessFilter(req, resp)
}

func main() {}
