go-restful
==========

package for building REST-style Web Services using Google Go

REST asks developers to use HTTP methods explicitly and in a way that's consistent with the protocol definition. This basic REST design principle establishes a one-to-one mapping between create, read, update, and delete (CRUD) operations and HTTP methods. According to this mapping:

- To create a resource on the server, use POST.
- To retrieve a resource, use GET.
- To change the state of a resource or to update it, use PUT.
- To remove or delete a resource, use DELETE.
    
###Example: [Hello world, plain and simple](https://github.com/emicklei/go-restful/tree/master/examples/restful-hello-world.go)
    
###Example: [Hello world, as a GreetingsService](https://github.com/emicklei/go-restful/tree/master/examples/restful-greetings.go)    
    
##Example: LandscapeService:

	package landscapeservice

	import (
	    "github.com/emicklei/go-restful"
	)

	func New() *restful.WebService {
		ws := new(restful.WebService)
	   	ws.Path("/applications").
			Consumes(restful.MIME_XML, restful.MIME_JSON).
			Produces(restful.MIME_XML, restful.MIME_JSON)

		ws.Route(ws.GET("/{id}").
			Doc("Get the Application node by its id").
			PathParam("id" , the unique string identifier for an application node").
			To(getApplication).
			Writes(Application{}))  // for api doc
			
		ws.Route(ws.POST("/").
			Doc("Create or update the Application node").
			To(saveApplication).
			Reads(Application{}))  // for api doc
		return ws
	}
	func getApplication(request *Request, response *Response) {
		id := request.PathParameter("id")
		env := request.QueryParameter("environment")
		...
	}
	func saveApplication(request *Request, response *Response) {
		// response.AddHeader("X-Something","other")
		// request.ReadEntity(anApp), uses Content-Type header to detect XML/JSON
		// response.WriteEntity(anApp) , uses Accept header to detect XML/JSON
		// response.WriterError(http.StatusInternalServerError,err) , set the response status and write an err
	}

### Resources

- [project on github](https://github.com/emicklei/go-restful)
- [example on blog](http://ernestmicklei.com/2012/11/24/go-restful-first-working-example/)
- [design on blog](http://ernestmicklei.com/2012/11/11/go-restful-api-design/)
- [landskape tool](https://github.com/emicklei/landskape)

(c) 2012, http://ernestmicklei.com. MIT License