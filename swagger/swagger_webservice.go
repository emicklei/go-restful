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
// DEPRECATED , use RegisterSwaggerService(...)
func InstallSwaggerService(aSwaggerConfig Config) {
	RegisterSwaggerService(aSwaggerConfig, restful.DefaultContainer)
}

// RegisterSwaggerService add the WebService that provides the API documentation of all services
// conform the Swagger documentation specifcation. (https://github.com/wordnik/swagger-core/wiki).
func RegisterSwaggerService(aSwaggerConfig Config, wsContainer *restful.Container) {
	config = aSwaggerConfig

	ws := new(restful.WebService)
	ws.Path(config.ApiPath)
	ws.Produces(restful.MIME_JSON)
	ws.Filter(enableCORS)
	ws.Route(ws.GET("/").To(getListing))
	ws.Route(ws.GET("/{a}").To(getDeclarations))
	ws.Route(ws.GET("/{a}/{b}").To(getDeclarations))
	ws.Route(ws.GET("/{a}/{b}/{c}").To(getDeclarations))
	ws.Route(ws.GET("/{a}/{b}/{c}/{d}").To(getDeclarations))
	ws.Route(ws.GET("/{a}/{b}/{c}/{d}/{e}").To(getDeclarations))
	ws.Route(ws.GET("/{a}/{b}/{c}/{d}/{e}/{f}").To(getDeclarations))
	ws.Route(ws.GET("/{a}/{b}/{c}/{d}/{e}/{f}/{g}").To(getDeclarations))
	log.Printf("[restful/swagger] listing is available at %v%v", config.WebServicesUrl, config.ApiPath)
	wsContainer.Add(ws)

	// Check paths for UI serving
	if config.SwaggerPath != "" && config.SwaggerFilePath != "" {
		log.Printf("[restful/swagger] %v%v is mapped to folder %v", config.WebServicesUrl, config.SwaggerPath, config.SwaggerFilePath)
		wsContainer.Handle(config.SwaggerPath, http.StripPrefix(config.SwaggerPath, http.FileServer(http.Dir(config.SwaggerFilePath))))
	} else {
		log.Printf("[restful/swagger] Swagger(File)Path is empty ; no UI is served")
	}
}

func enableCORS(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
    if origin := req.HeaderParameter("Origin"); origin != "" {
        resp.AddHeader("Access-Control-Allow-Origin", origin)
    }
    chain.ProcessFilter(req, resp)
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
					operation := Operation{HttpMethod: route.Method,
						Summary:       route.Doc,
						ResponseClass: asDataType(route.WriteSample),
						Nickname:      route.Operation}

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
		addModelFromSample(api, operation, false, route.ReadSample)
	}
	if route.WriteSample != nil {
		addModelFromSample(api, operation, true, route.WriteSample)
	}
}

// addModelFromSample creates and adds (or overwrites) a Model from a sample resource
func addModelFromSample(api *Api, operation *Operation, isResponse bool, sample interface{}) {
	st := reflect.TypeOf(sample)
	isCollection := false
	if st.Kind() == reflect.Slice || st.Kind() == reflect.Array {
		st = st.Elem()
		isCollection = true
	}
	if isResponse {
		modelName := st.String()
		if isCollection {
			modelName = "Array[" + modelName + "]"
		}
		operation.ResponseClass = modelName
	}
	addModelToApi(api, reflect.TypeOf(sample))
}

func addModelToApi(api *Api, st reflect.Type) {
	modelName := st.String()
	sm := Model{modelName, map[string]ModelProperty{}}
	// check for structure or primitive type
	if st.Kind() == reflect.Struct {
		for i := 0; i < st.NumField(); i++ {
			sf := st.Field(i)
			jsonName := sf.Name
			// see if a tag overrides this
			if override := st.Field(i).Tag.Get("json"); override != "" {
				jsonName = override
			}
			sm.Properties[jsonName] = asModelProperty(sf, api)
		}
	}
	api.Models[modelName] = sm
}

func asModelProperty(sf reflect.StructField, api *Api) ModelProperty {
	prop := ModelProperty{}
	st := sf.Type
	if st.Kind() == reflect.Slice || st.Kind() == reflect.Array {
		prop.Type = "List"
		prop.Items = map[string]string{"$ref": st.Elem().String()}
		// add|overwrite mode for element type
		addModelToApi(api, st.Elem())
	} else {
		prop.Type = st.String() // include pkg path
	}
	return prop
}

func asSwaggerParameter(param restful.ParameterData) Parameter {
	return Parameter{
		Name:        param.Name,
		Description: param.Description,
		ParamType:   asParamType(param.Kind),
		DataType:    param.DataType,
		Required:    param.Required}
}

// Between 1..7 path parameters is supported
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
	path = path + "/" + d
	e := req.PathParameter("e")
	if e == "" {
		return path
	}
	path = path + "/" + e
	f := req.PathParameter("f")
	if f == "" {
		return path
	}
	path = path + "/" + f
	g := req.PathParameter("g")
	if g == "" {
		return path
	}
	return path + "/" + g
}

func asParamType(kind int) string {
	switch {
	case kind == restful.PATH_PARAMETER:
		return "path"
	case kind == restful.QUERY_PARAMETER:
		return "query"
	case kind == restful.BODY_PARAMETER:
		return "body"
	case kind == restful.HEADER_PARAMETER:
		return "header"
	}
	return ""
}

func asDataType(any interface{}) string {
	if any == nil {
		return "void"
	}
	return reflect.TypeOf(any).Name()
}
