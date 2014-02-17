package swagger

import (
	"github.com/emicklei/go-restful"
	// "github.com/emicklei/hopwatch"
	"log"
	"net/http"
	"reflect"
	"strings"
)

type SwaggerService struct {
	config            Config
	apiDeclarationMap map[string]ApiDeclaration
}

func newSwaggerService(config Config) *SwaggerService {
	return &SwaggerService{
		config:            config,
		apiDeclarationMap: map[string]ApiDeclaration{}}
}

// LogInfo is the function that is called when this package needs to log. It defaults to log.Printf
var LogInfo = log.Printf

// InstallSwaggerService add the WebService that provides the API documentation of all services
// conform the Swagger documentation specifcation. (https://github.com/wordnik/swagger-core/wiki).
func InstallSwaggerService(aSwaggerConfig Config) {
	RegisterSwaggerService(aSwaggerConfig, restful.DefaultContainer)
}

// RegisterSwaggerService add the WebService that provides the API documentation of all services
// conform the Swagger documentation specifcation. (https://github.com/wordnik/swagger-core/wiki).
func RegisterSwaggerService(config Config, wsContainer *restful.Container) {
	sws := newSwaggerService(config)
	ws := new(restful.WebService)
	ws.Path(config.ApiPath)
	ws.Produces(restful.MIME_JSON)
	ws.Filter(enableCORS)
	ws.Route(ws.GET("/").To(sws.getListing))
	ws.Route(ws.GET("/{a}").To(sws.getDeclarations))
	ws.Route(ws.GET("/{a}/{b}").To(sws.getDeclarations))
	ws.Route(ws.GET("/{a}/{b}/{c}").To(sws.getDeclarations))
	ws.Route(ws.GET("/{a}/{b}/{c}/{d}").To(sws.getDeclarations))
	ws.Route(ws.GET("/{a}/{b}/{c}/{d}/{e}").To(sws.getDeclarations))
	ws.Route(ws.GET("/{a}/{b}/{c}/{d}/{e}/{f}").To(sws.getDeclarations))
	ws.Route(ws.GET("/{a}/{b}/{c}/{d}/{e}/{f}/{g}").To(sws.getDeclarations))
	LogInfo("[restful/swagger] listing is available at %v%v", config.WebServicesUrl, config.ApiPath)
	wsContainer.Add(ws)

	// Build all ApiDeclarations
	for _, each := range config.WebServices {
		rootPath := each.RootPath()
		// skip the api service itself
		if rootPath != config.ApiPath {
			if rootPath == "" || rootPath == "/" {
				// use routes
				for _, route := range each.Routes() {
					entry := staticPathFromRoute(route)
					_, exists := sws.apiDeclarationMap[entry]
					if !exists {
						sws.apiDeclarationMap[entry] = sws.composeDeclaration(each, entry)
					}
				}
			} else { // use root path
				sws.apiDeclarationMap[each.RootPath()] = sws.composeDeclaration(each, each.RootPath())
			}
		}
	}

	// Check paths for UI serving
	if config.SwaggerPath != "" && config.SwaggerFilePath != "" {
		swaggerPathSlash := config.SwaggerPath
		// path must end with slash /
		if "/" != config.SwaggerPath[len(config.SwaggerPath)-1:] {
			LogInfo("[restful/swagger] use corrected SwaggerFilePath ; must end with slash (/)")
			swaggerPathSlash += "/"
		}
		LogInfo("[restful/swagger] %v%v is mapped to folder %v", config.WebServicesUrl, swaggerPathSlash, config.SwaggerFilePath)
		wsContainer.Handle(swaggerPathSlash, http.StripPrefix(swaggerPathSlash, http.FileServer(http.Dir(config.SwaggerFilePath))))
	} else {
		LogInfo("[restful/swagger] Swagger(File)Path is empty ; no UI is served")
	}
}

func staticPathFromRoute(r restful.Route) string {
	static := r.Path
	bracket := strings.Index(static, "{")
	if bracket <= 1 { // result cannot be empty
		return static
	}
	if bracket != -1 {
		static = r.Path[:bracket]
	}
	if strings.HasSuffix(static, "/") {
		return static[:len(static)-1]
	} else {
		return static
	}
}

func enableCORS(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	if origin := req.HeaderParameter(restful.HEADER_Origin); origin != "" {
		resp.AddHeader(restful.HEADER_AccessControlAllowOrigin, origin)
	}
	chain.ProcessFilter(req, resp)
}

func (sws SwaggerService) getListing(req *restful.Request, resp *restful.Response) {
	listing := ResourceListing{SwaggerVersion: swaggerVersion}
	for k, _ := range sws.apiDeclarationMap {
		ref := ApiRef{Path: k}
		listing.Apis = append(listing.Apis, ref)
	}
	resp.WriteAsJson(listing)
}

func (sws SwaggerService) getDeclarations(req *restful.Request, resp *restful.Response) {
	resp.WriteAsJson(sws.apiDeclarationMap[composeRootPath(req)])
}

func (sws SwaggerService) composeDeclaration(ws *restful.WebService, pathPrefix string) ApiDeclaration {
	decl := ApiDeclaration{
		SwaggerVersion: swaggerVersion,
		BasePath:       sws.config.WebServicesUrl,
		ResourcePath:   ws.RootPath(),
		Models:         map[string]Model{}}

	// collect any path parameters
	rootParams := []Parameter{}
	for _, param := range ws.PathParameters() {
		rootParams = append(rootParams, asSwaggerParameter(param.Data()))
	}
	// aggregate by path
	pathToRoutes := map[string][]restful.Route{}
	for _, other := range ws.Routes() {
		if strings.HasPrefix(other.Path, pathPrefix) {
			routes := pathToRoutes[other.Path]
			pathToRoutes[other.Path] = append(routes, other)
		}
	}
	for path, routes := range pathToRoutes {
		api := Api{Path: path}
		for _, route := range routes {
			operation := Operation{HttpMethod: route.Method,
				Summary:  route.Doc,
				Type:     asDataType(route.WriteSample),
				Nickname: route.Operation}

			operation.Consumes = route.Consumes
			operation.Produces = route.Produces

			// share root params if any
			for _, swparam := range rootParams {
				operation.Parameters = append(operation.Parameters, swparam)
			}
			// route specific params
			for _, param := range route.ParameterDocs {
				operation.Parameters = append(operation.Parameters, asSwaggerParameter(param.Data()))
			}
			sws.addModelsFromRouteTo(&operation, route, &decl)
			api.Operations = append(api.Operations, operation)
		}
		decl.Apis = append(decl.Apis, api)
	}
	return decl
}

// addModelsFromRoute takes any read or write sample from the Route and creates a Swagger model from it.
func (sws SwaggerService) addModelsFromRouteTo(operation *Operation, route restful.Route, decl *ApiDeclaration) {
	if route.ReadSample != nil {
		sws.addModelFromSampleTo(operation, false, route.ReadSample, decl.Models)
	}
	if route.WriteSample != nil {
		sws.addModelFromSampleTo(operation, true, route.WriteSample, decl.Models)
	}
}

// addModelFromSample creates and adds (or overwrites) a Model from a sample resource
func (sws SwaggerService) addModelFromSampleTo(operation *Operation, isResponse bool, sample interface{}, models map[string]Model) {
	st := reflect.TypeOf(sample)
	isCollection := false
	if st.Kind() == reflect.Slice || st.Kind() == reflect.Array {
		st = st.Elem()
		isCollection = true
	} else {
		if st.Kind() == reflect.Ptr {
			if st.Elem().Kind() == reflect.Slice || st.Elem().Kind() == reflect.Array {
				st = st.Elem().Elem()
				isCollection = true
			}
		}
	}
	modelName := modelBuilder{}.keyFrom(st)
	if isResponse {
		if isCollection {
			modelName = "array[" + modelName + "]"
		}
		operation.Type = modelName
	}
	modelBuilder{models}.addModel(reflect.TypeOf(sample), "")
}

func asSwaggerParameter(param restful.ParameterData) Parameter {
	return Parameter{
		Name:        param.Name,
		Description: param.Description,
		ParamType:   asParamType(param.Kind),
		Type:        param.DataType,
		DataType:    param.DataType,
		Format:      asFormat(param.DataType),
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

func asFormat(name string) string {
	return "" // TODO
}

func asParamType(kind int) string {
	switch {
	case kind == restful.PathParameterKind:
		return "path"
	case kind == restful.QueryParameterKind:
		return "query"
	case kind == restful.BodyParameterKind:
		return "body"
	case kind == restful.HeaderParameterKind:
		return "header"
	case kind == restful.FormParameterKind:
		return "form"
	}
	return ""
}

func asDataType(any interface{}) string {
	if any == nil {
		return "void"
	}
	return reflect.TypeOf(any).Name()
}
