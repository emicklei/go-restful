/*
go-restful, a petit package for creating REST-style WebServices without magic. (Work-in-Progress)

Design discussed on http://ernestmicklei.com/2012/11/11/go-restful-api-design/

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
		ws.Route(ws.Method("GET").Path("/{id}").To(GetApplication))
		return ws
	}
	func GetApplication(request *Request, writer http.ResponseWriter) {
		id := request.PathParameter("id")
		...
	}

Example main:

	func main() {
		restful.Add(landscapeservice.New())
		log.Fatal(http.ListenAndServe(":8080", nil))	
	}

*/
package restful
