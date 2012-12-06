package restful

import (
	"github.com/emicklei/go-restful/swagger"
)

type swaggerService struct {
	WebService
	basePath string
}

var webServicesBasePath string
var swaggerServiceApiPath string

// Return the WebService that provides the API documentation of all services
// conform the Swagger documentation specifcation. (https://github.com/wordnik/swagger-core/wiki).
// The services are defined relative to @basePath, e.g. http://myservice:8989 .
// The JSON documentation is available on @apiPath, e.g. /api-docs.json
func NewSwaggerService(basePath, apiPath string) *WebService {
	webServicesBasePath = basePath
	swaggerServiceApiPath = apiPath

	ws := new(WebService)
	ws.Path(apiPath)
	ws.Produces(MIME_JSON)
	ws.Route(ws.GET("/").To(getListing))
	ws.Route(ws.GET("/{rootPath}").To(getDeclarations))
	return ws
}

func getListing(req *Request, resp *Response) {
	listing := swagger.ResourceListing{SwaggerVersion: "1.1", BasePath: webServicesBasePath}
	for _, each := range webServices {
		// skip the api service itself
		if each.RootPath() != swaggerServiceApiPath {
			api := swagger.Api{Path: each.RootPath()} // url encode , Description: each.Doc}
			listing.Apis = append(listing.Apis, api)
		}
	}
	resp.WriteAsJson(listing)
}

func getDeclarations(req *Request, resp *Response) {
	rootPath := req.PathParameter("rootPath")
	decl := swagger.ApiDeclaration{SwaggerVersion: "1.1", BasePath: rootPath}
	resp.WriteAsJson(decl)
}
