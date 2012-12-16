package restful

import (
	"github.com/emicklei/go-restful/swagger"
	//	"github.com/emicklei/hopwatch"
	//	"fmt"
	"log"
	"net/http"

//	"net/url"
)

type SwaggerConfig struct {
	WebServicesUrl  string // url where the services are available, e.g. http://localhost:8080
	ApiPath         string // path where the JSON api is avaiable , e.g. /apidocs
	SwaggerPath     string // path where the swagger UI will be served, e.g. /swagger
	SwaggerFilePath string // location of folder containing Swagger HTML5 application index.html
}

var webServicesBasePath string
var swaggerServiceApiPath string

// InstallSwaggerService add the WebService that provides the API documentation of all services
// conform the Swagger documentation specifcation. (https://github.com/wordnik/swagger-core/wiki).
func InstallSwaggerService(config SwaggerConfig) {
	// TODO just keep config around
	webServicesBasePath = config.WebServicesUrl
	swaggerServiceApiPath = config.ApiPath

	ws := new(WebService)
	ws.Path(config.ApiPath)
	ws.Produces(MIME_JSON)
	ws.Route(ws.GET("/").To(getListing))
	ws.Route(ws.GET("/{rootPath}").To(getDeclarations))
	Add(ws)

	// Install FileServer
	log.Printf("[restful] %v%v is mapped to folder %v", config.WebServicesUrl, config.SwaggerPath, config.SwaggerFilePath)
	http.Handle(config.SwaggerPath, http.StripPrefix(config.SwaggerPath, http.FileServer(http.Dir(config.SwaggerFilePath))))
}

func getListing(req *Request, resp *Response) {
	listing := swagger.ResourceListing{SwaggerVersion: "1.1", BasePath: webServicesBasePath}
	for _, each := range webServices {
		// skip the api service itself
		if each.RootPath() != swaggerServiceApiPath {
			api := swagger.Api{
				Path: swaggerServiceApiPath + "/" + each.RootPath()}
			//Description: each.Doc}
			listing.Apis = append(listing.Apis, api)
		}
	}
	resp.WriteAsJson(listing)
}

func getDeclarations(req *Request, resp *Response) {
	rootPath := "/" + req.PathParameter("rootPath")
	// log.Printf("rootPath:%V", rootPath)
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
					for _, param := range route.parameterDocs {
						swparam := swagger.Parameter{
							Name:        param.name,
							Description: param.description,
							ParamType:   asParamType(param.kind),
							DataType:    "String",
							Required:    param.required}
						operation.Parameters = append(operation.Parameters, swparam)
					}
					api.Operations = append(api.Operations, operation)
				}
				decl.Apis = append(decl.Apis, api)
			}
		}
	}
	resp.WriteAsJson(decl)
}

func asParamType(kind int) string {
	switch {
	case kind == PATH_PARAMETER:
		return "path"
	case kind == QUERY_PARAMETER:
		return "query"
	case kind == BODY_PARAMETER:
		return "body"
	}
	return ""
}
