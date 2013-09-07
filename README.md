go-restful
==========

package for building REST-style Web Services using Google Go

REST asks developers to use HTTP methods explicitly and in a way that's consistent with the protocol definition. This basic REST design principle establishes a one-to-one mapping between create, read, update, and delete (CRUD) operations and HTTP methods. According to this mapping:

- GET = Retrieve
- POST = Create if you are sending a command to the server to create a subordinate of the specified resource, using some server-side algorithm.
- POST = Update if you are requesting the server to update one or more subordinates of the specified resource.
- PUT = Create iff you are sending the full content of the specified resource (URI).
- PUT = Update iff you are updating the full content of the specified resource.
- DELETE = Delete if you are requesting the server to delete the resource
- PATCH = Update partial content of a resource
- OPTIONS = Get information about the communication options for the Request-URI
    
### Example

	ws := new(restful.WebService)
	ws.
		Path("/users").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML)

	ws.Route(ws.GET("/{user-id}").To(u.findUser).
		Doc("get a user").
		Param(ws.PathParameter("user-id", "identifier of the user").DataType("string")).
		Writes(User{}))
	...
	
	func (u UserResource) findUser(request *restful.Request, response *restful.Response) {
		id := request.PathParameter("user-id")
		...
	
[Full API of a UserResource](https://github.com/emicklei/go-restful/tree/master/examples/restful-user-resource.go) 
		
### Features

- Routes for request -> function mapping with path parameter (e.g. {id}) support
- Routing algorithm after [JSR311](http://jsr311.java.net/nonav/releases/1.1/spec/spec.html); Router can be configured
- Request API for reading structs from JSON/XML and accesing parameters (path,query,header)
- Response API for writing structs to JSON/XML and setting headers
- Filters for intercepting the request, response flow	 on Service or Route level
- Containers for WebServices on different HTTP endpoints
- Content encoding (gzip,deflate) of responses
- Automatic responses on OPTIONS (using a filter)
- Automatic CORS request handling (using a filter)
- API declaration for Swagger UI
- Panic recovery to produce HTTP 500
	
### Resources

- [Documentation go-restful (godoc.org)](http://godoc.org/github.com/emicklei/go-restful)
- [Hello world, plain and simple](https://github.com/emicklei/go-restful/tree/master/examples/restful-hello-world.go)  
- [Full API of a UserResource](https://github.com/emicklei/go-restful/tree/master/examples/restful-user-resource.go) 
- [Other examples](https://github.com/emicklei/go-restful/tree/master/examples)
- [Example posted on blog](http://ernestmicklei.com/2012/11/24/go-restful-first-working-example/)
- [Design explained on blog](http://ernestmicklei.com/2012/11/11/go-restful-api-design/)
- [Showcase: Mora - MongoDB REST Api server](https://github.com/emicklei/mora)
- [Showcase: Landskape tool](https://github.com/emicklei/landskape)

[![Build Status](https://drone.io/github.com/emicklei/go-restful/status.png)](https://drone.io/github.com/emicklei/go-restful/latest)

(c) 2013, http://ernestmicklei.com. MIT License