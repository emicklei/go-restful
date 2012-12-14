package restful

import (
	"github.com/emicklei/go-restful/swagger"
//	"fmt"
)

type swaggerService struct {
	WebService
	basePath string
}

var webServicesBasePath string
var swaggerServiceApiPath string

// NewSwaggerService returns the WebService that provides the API documentation of all services
// conform the Swagger documentation specifcation. (https://github.com/wordnik/swagger-core/wiki).
// The services are defined relative to @basePath, e.g. http://myservice:8989 .
// The JSON documentation is available on @apiPath, e.g. /api-docs.json
func NewSwaggerService(wsPath, apiPath string) *WebService {
	webServicesBasePath = wsPath
	swaggerServiceApiPath = apiPath

	ws := new(WebService)
	ws.Path(apiPath)
	ws.Produces(MIME_JSON)
	ws.Route(ws.GET("/").To(getListing))
	ws.Route(ws.GET("/{rootPath}").To(getDeclarations))
	return ws
}

func getListing(req *Request, resp *Response) {
	resp.AddHeader("Access-Control-Allow-Origin", "*")
	listing := swagger.ResourceListing{SwaggerVersion: "1.1", BasePath: webServicesBasePath}
	for _, each := range webServices {
		// skip the api service itself
		if each.RootPath() != swaggerServiceApiPath {
			api := swagger.Api{Path: swaggerServiceApiPath + each.RootPath()} // url encode , Description: each.Doc}
			listing.Apis = append(listing.Apis, api)
		}
	}
	resp.WriteAsJson(listing)
}

func getDeclarations(req *Request, resp *Response) {
	resp.AddHeader("Access-Control-Allow-Origin", "*")
	rootPath := "/" + req.PathParameter("rootPath")
	decl := swagger.ApiDeclaration{SwaggerVersion: "1.1", BasePath: webServicesBasePath, ResourcePath: rootPath}
	for _, each := range webServices {
		// find the webservice
		if each.RootPath() == rootPath {
			// aggregate by path
			pathToRoutes := map[string][]Route{}
			for _, other := range each.Routes() {
				routes := pathToRoutes[other.Path]
				pathToRoutes[other.Path] = append(routes, other)
			}
			for path, routes := range pathToRoutes {
				api := swagger.Api{Path: path}
				for _, route := range routes {
					operation := swagger.Operation{HttpMethod: route.Method, Summary: route.Doc}
					api.Operations = append(api.Operations, operation)
				}
				decl.Apis = append(decl.Apis, api)
			}
		}
	}
	resp.WriteAsJson(decl)
}
