package restful

import (
	"github.com/emicklei/go-restful/swagger"
	// "github.com/emicklei/hopwatch"
	"log"
	"net/http"
	"reflect"
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
	ws.Route(ws.GET("/{a}").To(getDeclarations))
	ws.Route(ws.GET("/{a}/{b}").To(getDeclarations))
	ws.Route(ws.GET("/{a}/{b}/{c}").To(getDeclarations))
	ws.Route(ws.GET("/{a}/{b}/{c}/{d}").To(getDeclarations)) // TODO maybe support * in the path spec?
	log.Printf("[restful/swagger] listing is available at %v%v", config.WebServicesUrl, config.ApiPath)
	Add(ws)

	// Install FileServer
	log.Printf("[restful/swagger] %v%v is mapped to folder %v", config.WebServicesUrl, config.SwaggerPath, config.SwaggerFilePath)
	http.Handle(config.SwaggerPath, http.StripPrefix(config.SwaggerPath, http.FileServer(http.Dir(config.SwaggerFilePath))))
}

func getListing(req *Request, resp *Response) {
	listing := swagger.ResourceListing{SwaggerVersion: "1.1", BasePath: webServicesBasePath}
	for _, each := range webServices {
		// skip the api service itself
		if each.RootPath() != swaggerServiceApiPath {
			api := swagger.Api{
				Path: swaggerServiceApiPath + each.RootPath()}
			//Description: each.Doc}
			listing.Apis = append(listing.Apis, api)
		}
	}
	resp.WriteAsJson(listing)
}

func getDeclarations(req *Request, resp *Response) {
	rootPath := composeRootPath(req)
	// log.Printf("rootPath:%V", rootPath)
	decl := swagger.ApiDeclaration{SwaggerVersion: "1.1", BasePath: webServicesBasePath, ResourcePath: rootPath}
	for _, each := range webServices {
		// find the webservice
		if each.RootPath() == rootPath {
			// collect any path parameters
			rootParams := []swagger.Parameter{}
			for _, param := range each.PathParameters() {
				rootParams = append(rootParams, asSwaggerParameter(param))
			}
			// aggregate by path
			pathToRoutes := map[string][]Route{}
			for _, other := range each.Routes() {
				routes := pathToRoutes[other.Path]
				pathToRoutes[other.Path] = append(routes, other)
			}
			for path, routes := range pathToRoutes {
				api := swagger.Api{Path: path, Models: map[string]swagger.Model{}}
				for _, route := range routes {
					operation := swagger.Operation{HttpMethod: route.Method, Summary: route.Doc}
					// share root params if any
					for _, swparam := range rootParams {
						operation.Parameters = append(operation.Parameters, swparam)
					}
					// route specific params					
					for _, param := range route.parameterDocs {
						operation.Parameters = append(operation.Parameters, asSwaggerParameter(param))
					}
					api.Operations = append(api.Operations, operation)
					addModelsFromRoute(&api, route)
				}
				decl.Apis = append(decl.Apis, api)
			}
		}
	}
	resp.WriteAsJson(decl)
}

// addModelsFromRoute takes any read or write sample from the Route and creates a Swagger model from it.
func addModelsFromRoute(api *swagger.Api, route Route) {
	if route.readSample != nil {
		addModelFromSample(api, route.readSample)
	}
	if route.writeSample != nil {
		addModelFromSample(api, route.writeSample)
	}
}

// addModelFromSample creates and adds (or overwrites) a Model from a sample resource
func addModelFromSample(api *swagger.Api, sample interface{}) {
	st := reflect.TypeOf(sample)
	sm := swagger.Model{map[string]swagger.ModelProperty{}}
	for i := 0; i < st.NumField(); i++ {
		sf := st.Field(i)
		sp := swagger.ModelProperty{Type: sf.Type.Name()}
		sm.Properties[sf.Name] = sp
	}
	api.Models[st.String()] = sm
}

func asSwaggerParameter(param *Parameter) swagger.Parameter {
	return swagger.Parameter{
		Name:        param.name,
		Description: param.description,
		ParamType:   asParamType(param.kind),
		DataType:    param.dataType,
		Required:    param.required}
}

// Between 1..4 path parameters supported
func composeRootPath(req *Request) string {
	path := "/" + req.PathParameter("a")
	b := req.PathParameter("b")
	if b == "" {
		return path
	}
	path = path + "/" + b
	c := req.PathParameter("c")
	if c == "" {
		return path
	}
	path = path + "/" + c
	d := req.PathParameter("d")
	if d == "" {
		return path
	}
	return path + "/" + d
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
