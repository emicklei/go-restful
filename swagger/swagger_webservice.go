package swagger

import (
	"github.com/emicklei/go-restful"
	// "github.com/emicklei/hopwatch"
	"log"
	"net/http"
	"reflect"
)

var config Config

// InstallSwaggerService add the WebService that provides the API documentation of all services
// conform the Swagger documentation specifcation. (https://github.com/wordnik/swagger-core/wiki).
func InstallSwaggerService(aSwaggerConfig Config) {
	config = aSwaggerConfig

	ws := new(restful.WebService)
	ws.Path(config.ApiPath)
	ws.Produces(restful.MIME_JSON)
	ws.Route(ws.GET("/").To(getListing))
	ws.Route(ws.GET("/{a}").To(getDeclarations))
	ws.Route(ws.GET("/{a}/{b}").To(getDeclarations))
	ws.Route(ws.GET("/{a}/{b}/{c}").To(getDeclarations))
	ws.Route(ws.GET("/{a}/{b}/{c}/{d}").To(getDeclarations)) // TODO maybe support * in the path spec?
	log.Printf("[restful/swagger] listing is available at %v%v", config.WebServicesUrl, config.ApiPath)
	restful.Add(ws)

	// Install FileServer
	log.Printf("[restful/swagger] %v%v is mapped to folder %v", config.WebServicesUrl, config.SwaggerPath, config.SwaggerFilePath)
	http.Handle(config.SwaggerPath, http.StripPrefix(config.SwaggerPath, http.FileServer(http.Dir(config.SwaggerFilePath))))
}

func getListing(req *restful.Request, resp *restful.Response) {
	listing := ResourceListing{SwaggerVersion: "1.1", BasePath: config.WebServicesUrl}
	for _, each := range config.WebServices {
		// skip the api service itself
		if each.RootPath() != config.ApiPath {
			api := Api{
				Path: config.ApiPath + each.RootPath()}
			//Description: each.Doc}
			listing.Apis = append(listing.Apis, api)
		}
	}
	resp.WriteAsJson(listing)
}

func getDeclarations(req *restful.Request, resp *restful.Response) {
	rootPath := composeRootPath(req)
	// log.Printf("rootPath:%V", rootPath)
	decl := ApiDeclaration{SwaggerVersion: "1.1", BasePath: config.WebServicesUrl, ResourcePath: rootPath}
	for _, each := range config.WebServices {
		// find the webservice
		if each.RootPath() == rootPath {
			// collect any path parameters
			rootParams := []Parameter{}
			for _, param := range each.PathParameters() {
				rootParams = append(rootParams, asSwaggerParameter(param.Data()))
			}
			// aggregate by path
			pathToRoutes := map[string][]restful.Route{}
			for _, other := range each.Routes() {
				routes := pathToRoutes[other.Path]
				pathToRoutes[other.Path] = append(routes, other)
			}
			for path, routes := range pathToRoutes {
				api := Api{Path: path, Models: map[string]Model{}}
				for _, route := range routes {
					operation := Operation{HttpMethod: route.Method, Summary: route.Doc, ResponseClass: "void"}

					// share root params if any
					for _, swparam := range rootParams {
						operation.Parameters = append(operation.Parameters, swparam)
					}
					// route specific params
					for _, param := range route.ParameterDocs {
						operation.Parameters = append(operation.Parameters, asSwaggerParameter(param.Data()))
					}
					addModelsFromRoute(&api, &operation, route)
					api.Operations = append(api.Operations, operation)
				}
				decl.Apis = append(decl.Apis, api)
			}
		}
	}
	resp.WriteAsJson(decl)
}

// addModelsFromRoute takes any read or write sample from the Route and creates a Swagger model from it.
func addModelsFromRoute(api *Api, operation *Operation, route restful.Route) {
	if route.ReadSample != nil {
		addModelFromSample(api, operation, true, route.ReadSample)
	}
	if route.WriteSample != nil {
		addModelFromSample(api, operation, false, route.WriteSample)
	}
}

// addModelFromSample creates and adds (or overwrites) a Model from a sample resource
func addModelFromSample(api *Api, operation *Operation, isResponse bool, sample interface{}) {
	st := reflect.TypeOf(sample)
	if isResponse {
		operation.ResponseClass = st.String()
	}
	sm := Model{st.String(), map[string]ModelProperty{}}
	// TODO handle recursive structures, hidden and array fields
	for i := 0; i < st.NumField(); i++ {
		sf := st.Field(i)
		sp := ModelProperty{Type: sf.Type.Name()}
		sm.Properties[sf.Name] = sp
	}
	api.Models[st.String()] = sm
}

func asSwaggerParameter(param restful.ParameterData) Parameter {
	return Parameter{
		Name:        param.Name,
		Description: param.Description,
		ParamType:   asParamType(param.Kind),
		DataType:    param.DataType,
		Required:    param.Required}
}

// Between 1..4 path parameters supported
func composeRootPath(req *restful.Request) string {
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
	case kind == restful.PATH_PARAMETER:
		return "path"
	case kind == restful.QUERY_PARAMETER:
		return "query"
	case kind == restful.BODY_PARAMETER:
		return "body"
	}
	return ""
}
