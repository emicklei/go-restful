/*
Package go-restful, a lean package for creating REST-style WebServices without magic.

Example WebService:

	package landscapeservice

	import (
	    "github.com/emicklei/go-restful"
	)

	func New() *restful.WebService {
		ws := new(restful.WebService)
	   	ws.Path("/applications").
			Consumes(restful.MIME_XML, restful.MIME_JSON).
			Produces(restful.MIME_XML, restful.MIME_JSON)

		ws.Route(ws.GET("/{id}").To(getApplication).
			// for documentation
			Doc("Get the Application node by its id").
			Operation("getApplication").
			Param(ws.PathParameter("id" , "the identifier for an application node")).
			Param(ws.QueryParameter("environment" , "the scope in which the application node lives")).
			Writes(Application{})) // to the response

		ws.Route(ws.POST("/").To(saveApplication).
			// for documentation
			Doc("Create or update the Application node").
			Operation("saveApplication").
			Reads(Application{})) // from the request
		return ws
	}
	func getApplication(request *restful.Request, response *restful.Response) {
			id := request.PathParameter("id")
			env := request.QueryParameter("environment")
			...
	}
	func saveApplication(request *restful.Request, response *restful.Response) {
		   response.AddHeader("X-Something","other")
		   ...
		   response.WriteEntity(anApp) // uses Accept header to detect XML/JSON

	}

Example main:

	func main() {
		restful.Add(landscapeservice.New())
		log.Fatal(http.ListenAndServe(":8080", nil))
	}


WebServices

A WebService has a collection of Route objects that dispatch incoming Http Requests to a function calls.

	type RouteFunction func(*restful.Request, *restful.Response)

A Route is defined by a HTTP method, an URL path and (optionally) the MIME types it consumes (Content-Type) and produces (Accept).
This package has the logic to find the best matching Route and if found, call its Function.
The (*Request, *Response) arguments provide functions for reading information from the request and writing information back to the response.


Filters

A filter dynamically intercepts requests and responses to transform or use the information contained in the requests or responses.
You can use filters to perform generic logging, measurement, authentication, redirect, set response headers etc.
In the restful package there are three hooks into the request,response flow where filters can be added.
Each filter must define a FilterFunction:

	func (req *restful.Request, resp *restful.Response, chain *restful.FilterChain)

Use the following statement to pass the request,response pair to the next filter or RouteFunction

	chain.ProcessFilter(req, resp)

Global Filters

These are processed before any registered WebService.

	// install a global filter (processed before any webservice)
	restful.Filter(globalLogging)


WebService Filters

These are processed before any Route of a WebService.

	// install a webservice filter (processed before any route)
	ws.Filter(webserviceLogging).Filter(measureTime)


Route Filters

These are processed before calling the function associated with the Route.

	// install 2 chained route filters (processed before calling findUser)
	ws.Route(ws.GET("/{user-id}").Filter(routeLogging).Filter(NewCountFilter().routeCounter).To(findUser))


See the example https://github.com/emicklei/go-restful/blob/master/examples/restful-filters.go with full implementations.

Serving files

Use the Go standard http.ServeFile function to serve file system assets.

	ws.Route(ws.GET("/static/{resource}").To(staticFromPathParam))
	...
	// http://localhost:8080/static/test.xml
	// http://localhost:8080/static/
	func staticFromPathParam(req *restful.Request, resp *restful.Response) {
		http.ServeFile(
			resp.ResponseWriter,
			req.Request,
			path.Join(rootdir, req.PathParameter("resource")))
	}

See the example https://github.com/emicklei/go-restful/blob/master/examples/restful-serve-static.go with full implementations.

Response Encoding

Two encodings are supported: gzip and deflate. To enable this for all responses:

	restful.EnableContentEncoding = true

If a Http request includes the Accept-Encoding header then the response content will be compressed using the specified encoding.

Alternatively, you can create a Filter that performs the encoding and install it per WebService or Route.
See the example https://github.com/emicklei/go-restful/blob/master/examples/restful-encoding-filter.go


Error Handling

Unexpected things happen. If a request cannot be processed because of a failure, your service needs to tell the response what happened and why.
For this reason HTTP status codes exist and it is important to use the correct code in every exceptional situation.

400: Bad Request

If path or query parameters are not valid (content or type) then use http.StatusBadRequest.

	id, err := strconv.Atoi(req.PathParameter("id"))
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

404: Not Found

Despite a valid URI, the resource requested may not be available

	resp.WriteHeader(http.StatusNotFound)

500: Internal Server Error

If the application logic could not process the request (or write the response) then use http.StatusInternalServerError.

	question, err := application.SharedLogic.GetQuestionById(id)
	if err != nil {
		log.Printf("GetQuestionById failed:", err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

ServiceError

In addition to setting the correct (error) Http status code, you can choose to write a ServiceError message on the response:

	resp.WriteEntity(restful.NewError(http.StatusBadRequest, "Non-integer {id} path parameter"))

	resp.WriteEntity(restful.NewError(http.StatusInternalServerError, err.Error()))



Resources

[project]: https://github.com/emicklei/go-restful

[example]: https://github.com/emicklei/go-restful/blob/master/examples/restful-user-service.go

[design]:  http://ernestmicklei.com/2012/11/11/go-restful-api-design/

[showcase]: https://github.com/emicklei/landskape

(c) 2013, http://ernestmicklei.com. MIT License
*/
package restful
