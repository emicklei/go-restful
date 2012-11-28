/*
go-restful, a lean package for creating REST-style WebServices without magic.

Example WebService:

	package landscapeservice

	import (
	    "github.com/emicklei/go-restful"
	)

	type LandscapeService struct {
		restful.WebService
	}
	func New() *LandscapeService {
		ws := new(LandscapeService)
		ws.Path("/applications").Accept("application/xml").ContentType("application/xml")

		ws.Route(ws.GET("/{id}").To(GetApplication))
		ws.Route(ws.POST("/").To(SaveApplication))
		return ws
	}
	func GetApplication(request *Request, response *Response) {
		// id := request.PathParameter("id")
		// env := request.QueryParameter("environment")
	}
	func SaveApplication(request *Request, response *Response) {
		// response.AddHeader("X-Something","other")
		// response.WriteEntity(anApp) , use Accept header to detect XML/JSON
		// response.WriterError(http.StatusInternalServerError,err)
	}	

Example main:

	func main() {
		restful.Add(landscapeservice.New())
		// Show me the WADL spec
		log.Print(restful.Wadl("http://localhost:8080"))
		log.Fatal(http.ListenAndServe(":8080", nil))	
	}

[project]: https://github.com/emicklei/go-restful

[example]: http://ernestmicklei.com/2012/11/24/go-restful-first-working-example/

[design]:  http://ernestmicklei.com/2012/11/11/go-restful-api-design/

[1st use]: https://github.com/emicklei/landskape

(c) 2012, http://ernestmicklei.com. MIT License
*/
package restful
