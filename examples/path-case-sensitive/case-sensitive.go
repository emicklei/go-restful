package main

import (
	"io"
	"log"
	"net/http"

	restful "github.com/emicklei/go-restful/v3"
)

// This example shows how to handle different casing of path template.
//
// GET http://localhost:8080/hola/Juan
// GET http://localhost:8080/HOLA/Juan
// GET http://localhost:8080/Hola/Juan

func main() {
	ws := new(restful.WebService)

	// hola is path template, to accept different casing of hola, we use regex matching with syntax {name:regex}
	// - {: is nesserary to trigger the regex matching.
	// - (?i) is to make the regex case-insensitive.
	// it seems solve the issue, but there is a issue you might hit: https://github.com/emicklei/go-restful/issues/545
	// to avoid partial matching, another regex pattern
	// - ^$ is needed to match the whole route token.
	ws.Route(ws.GET("/{:(?i)^hola$}/{name:*}").To(hola))
	restful.Add(ws)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func hola(req *restful.Request, resp *restful.Response) {
	io.WriteString(resp, "hola "+req.PathParameter("name"))
}
