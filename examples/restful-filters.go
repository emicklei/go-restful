package main

import (
	"github.com/emicklei/go-restful"
	"log"
	"net/http"
)

type User struct {
	Id, Name string
}

type UserList struct {
	Users []User
}

func NewUserService() *restful.WebService {
	ws := new(restful.WebService)
	ws.
		Path("/users").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML)

	logging := LoggingFilter{findUser}.handleFilter

	counter := new(CountFilter)
	counter.WrappedFunction = logging
	counting := counter.handleFilter
	ws.Route(ws.GET("/{user-id}").To(counting))

	//ws.Filter("/users/", handleLogging)
	return ws
}

// GlobalFilter
func globalHandleLogging(w http.ResponseWriter, r *http.Request) {
	log.Printf("[global-filter] %s,%s\n", r.Method, r.URL)
	restful.DefaultDispatch(w, r)
}

// WebServiceFilter
func webserviceHandleLogging(httpWriter http.ResponseWriter, httpRequest *http.Request) {
	log.Printf("[webservice-filter] %s,%s\n", httpRequest.Method, httpRequest.URL)
}

// RouteFunctionFilter
type LoggingFilter struct {
	WrappedFunction restful.RouteFunction
}

func (l LoggingFilter) handleFilter(request *restful.Request, response *restful.Response) {
	log.Printf("[function-filter (logging)] req:%v resp:%v", request, response)
	l.WrappedFunction(request, response)
}

type CountFilter struct {
	Count           int
	WrappedFunction restful.RouteFunction
}

func (c *CountFilter) handleFilter(request *restful.Request, response *restful.Response) {
	c.Count++
	log.Printf("[function-filter (count)] count:%d, req:%v resp:%v", c.Count, request, response)
	c.WrappedFunction(request, response)
}

// Global Filter > replace Dispatch function , type HandlerFunc func(ResponseWriter, *Request
// WebService Filter >  filter on pattern?
// Route Filter > RouteFunction func(*Request, *Response)

//  A filter dynamically intercepts requests and responses to transform or use the information contained in the requests or responses.
// http://www.oracle.com/technetwork/java/filters-137243.html
//func handleLogging(request *restful.Request, response *restful.Response, chain *FilterChain) {
//	log.Printf("req:%v resp:%v", request, response)

//	chain.handleNextFilter(request, response)
//}

// GET http://localhost:8080/users
//
func getAllUsers(request *restful.Request, response *restful.Response) {
	response.WriteEntity(UserList{[]User{User{"42", "Gandalf"}, User{"3.14", "Pi"}}})
}

// GET http://localhost:8080/users/42
//
func findUser(request *restful.Request, response *restful.Response) {
	response.WriteEntity(User{"42", "Gandalf"})
}

func main() {
	// Install global filter (directly using replacement of Dispatch)
	restful.Dispatch = globalHandleLogging

	restful.Add(NewUserService())
	log.Printf("start listening on localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
