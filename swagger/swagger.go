// Package swagger implements the structures of the Swagger (https://github.com/wordnik/swagger-core/wiki) specification
package swagger

type ResourceListing struct {
	ApiVersion     string `json:"apiVersion"`
	SwaggerVersion string `json:"swaggerVersion"`
	BasePath       string `json:"basePath"`
	Apis           []Api  `json:"apis"`
}

type Api struct {
	Path        string      `json:"path"`
	Description string      `json:"description"`
	Operations  []Operation `json:"operations"`
	Models      []Model     `json:"models"`
}

type ApiDeclaration struct {
	ApiVersion     string `json:"apiVersion"`
	SwaggerVersion string `json:"swaggerVersion"`
	BasePath       string `json:"basePath"`
	ResourcePath   string `json:"resourcePath"`
	Apis           []Api  `json:"apis"`
}

type Operation struct {
	HttpMethod     string          `json:"httpMethod"`
	Nickname       string          `json:"nickname"`
	ResponseClass  string          `json:"responseClass"`
	Summary        string          `json:"summary"`
	Notes          string          `json:"notes"`
	Parameters     []Parameter     `json:"parameters"`
	ErrorResponses []ErrorResponse `json:"errorResponses"`
}

type Parameter struct {
	ParamType       string            `json:"paramType"`
	Name            string            `json:"name"`
	Description     string            `json:"description"`
	DataType        string            `json:"dataType"`
	Required        bool              `json:"required"`
	AllowableValues map[string]string `json:"allowableValues"`
	AllowMultiple   bool              `json:"allowMultiple"`
}

type ErrorResponse struct {
	Code   int    `json:"code"`
	Reason string `json:"reason"`
}

type Model struct{}
