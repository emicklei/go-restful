/*
go-restful, a library for creating REST-style WebServices  (Work-in-Progress)

Design discussed on http://ernestmicklei.com/2012/11/11/go-restful-api-design/

Example:

import (
    "github.com/emicklei/go-restful"
)

type LandscapeService struct {
	restful.WebService
}
func New() *LandscapeService {
	ws := new(LandscapeService)
	ws.Path("/")
	ws.Route(ws.Method("GET").Path("/applications").To(GetApplications))
}
func GetApplications(request *Request, writer http.ResponseWriter) {
	...
}


*/
package restful