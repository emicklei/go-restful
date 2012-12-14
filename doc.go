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

		ws.Route(ws.GET("/{id}").
			Doc("Get the Application node by its id").
			Param(ws.PathParameter("id" , the unique string identifier for an application node")).
			To(getApplication).
			Writes(Application{}))
		ws.Route(ws.POST("/").
			Doc("Create or update the Application node").
			To(saveApplication).
			Reads(Application{}))
		return ws
	}
	func getApplication(request *Request, response *Response) {
		// id := request.PathParameter("id")
		// env := request.QueryParameter("environment")
	}
	func saveApplication(request *Request, response *Response) {
		// response.AddHeader("X-Something","other")
		// response.WriteEntity(anApp) , uses Accept header to detect XML/JSON
		// response.WriterError(http.StatusInternalServerError,err)
	}	

Example main:

	func main() {
		restful.Add(landscapeservice.New())	
		log.Fatal(http.ListenAndServe(":8080", nil))	
	}

[project]: https://github.com/emicklei/go-restful

[example]: http://ernestmicklei.com/2012/11/24/go-restful-first-working-example/

[design]:  http://ernestmicklei.com/2012/11/11/go-restful-api-design/

[1st use]: https://github.com/emicklei/landskape

(c) 2012, http://ernestmicklei.com. MIT License
*/
package restful
