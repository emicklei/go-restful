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
		id := request.PathParameter("id")
		...
	}
	func SaveApplication(request *Request, response *Response) {
		...
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
