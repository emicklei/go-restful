package restful

import (
	"github.com/emicklei/go-restful/swagger"
)

type swaggerService struct {
	WebService
	basePath string
}

var apiService = new(swaggerService)

// Return the WebService that provides the API documentation of all services
// using the Swagger documentation specifcation. (https://github.com/wordnik/swagger-core/wiki)
// The JSON documentation is available on /api/api-docs.json
func SwaggerService() *swaggerService {
	return apiService
}

// Set the Http base path of all restful WebServices (http://some.domain)
func SwaggerBasePath(basePath string) {
	apiService.basePath = basePath
}

func init() {
	apiService.Path("/api").Produces(MIME_JSON)
	apiService.Route(apiService.GET("/api-docs.json").To(getListing))
	apiService.Route(apiService.GET("/api-docs.json/{rootPath}").To(getDeclarations))
}

func getListing(req *Request, resp *Response) {
	listing := swagger.ResourceListing{SwaggerVersion: "1.1", BasePath: apiService.basePath}
	for _, each := range webServices {
		if each != apiService {
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
