package main

import (
	"log"
	"net/http"

	restful "github.com/emicklei/go-restful/v3"

	// httpin is a Go package for decoding HTTP requests into structs
	// an alternative to gorilla/mux.
	//
	// See its documentation on integration with go-restful:
	// https://ggicci.github.io/httpin/integrations/go-restful
	//
	// httpin can decode data from:
	// - query parameters
	// - headers
	// - form data
	// - JSON/XML request body
	// - URL path variables
	// - file uploads
	// by defining an input struct and composing fileds struct tags
	"github.com/ggicci/httpin"
)

type ListUsersInput struct {
	Gender   string `in:"query=gender"`
	AgeRange []int  `in:"query=age_range"`
	IsMember bool   `in:"query=is_member"`
	Token    string `in:"header=x-client-token;query=access_token"`
}

func handleListUsers(request *restful.Request, response *restful.Response) {
	// Retrieve you data in one line of code!
	input := request.Request.Context().Value(httpin.Input).(*ListUsersInput)

	response.WriteAsJson(input)
}

func main() {
	ws := new(restful.WebService)

	// Bind input struct with handler.
	ws.Route(ws.GET("/users").Filter(
		restful.HttpMiddlewareHandlerToFilter(httpin.NewInput(ListUsersInput{})),
	).To(handleListUsers))

	restful.Add(ws)

	// Visit http://localhost:8080/users?gender=male&age_range=18&age_range=22&is_member=1&access_token=my-private-secret
	// in your browser to see the result.
	log.Fatal(http.ListenAndServe(":8080", nil))
}
