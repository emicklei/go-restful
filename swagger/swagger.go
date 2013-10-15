// Package swagger implements the structures of the Swagger (https://github.com/wordnik/swagger-core/wiki) specification
package swagger

const swaggerVersion = "1.2"

type ResourceListing struct {
	ApiVersion     string `json:"apiVersion"`
	SwaggerVersion string `json:"swaggerVersion"` // e.g 1.2
	// BasePath       string `json:"basePath"`  obsolete in 1.1
	Apis []Api `json:"apis"`
}

type Api struct {
	Path        string           `json:"path"` // relative or absolute, must start with /
	Description string           `json:"description"`
	Operations  []Operation      `json:"operations"`
	Models      map[string]Model `json:"models"`
}

type ApiDeclaration struct {
	ApiVersion     string   `json:"apiVersion"`
	SwaggerVersion string   `json:"swaggerVersion"`
	BasePath       string   `json:"basePath"`
	ResourcePath   string   `json:"resourcePath"` // must start with /
	Apis           []Api    `json:"apis,omitempty"`
	Consumes       []string `json:"consumes,omitempty"`
	Produces       []string `json:"produces,omitempty"`
}

type Operation struct {
	HttpMethod string `json:"httpMethod"`
	Nickname   string `json:"nickname"`
	Type       string `json:"type"` // in 1.1 = DataType
	// ResponseClass    string            `json:"responseClass"` obsolete in 1.2
	Summary          string            `json:"summary,omitempty"`
	Notes            string            `json:"notes,omitempty"`
	Parameters       []Parameter       `json:"parameters,omitempty"`
	ResponseMessages []ResponseMessage `json:"responseMessages,omitempty"` // optional
	Consumes         []string          `json:"consumes,omitempty"`
	Produces         []string          `json:"produces,omitempty"`
	Authorizations   []Authorization   `json:"authorizations,omitempty"`
	Protocols        []Protocol        `json:"protocols,omitempty"`
}

type Protocol struct {
}

type ResponseMessage struct {
	Code          int    `json:"code"`
	Message       string `json:"message"`
	ResponseModel string `json:"responseModel"`
}

type Parameter struct {
	ParamType   string `json:"paramType"` // path,query,body,header,form
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`   // integer
	Format      string `json:"format"` // int64
	Required    bool   `json:"required"`
	Minimum     int    `json:"minimum"`
	Maximum     int    `json:"maximum"`
}

type ErrorResponse struct {
	Code   int    `json:"code"`
	Reason string `json:"reason"`
}

type Model struct {
	Id         string                   `json:"id"`
	Required   []string                 `json:"required"`
	Properties map[string]ModelProperty `json:"properties"`
}

type ModelProperty struct {
	Type        string            `json:"type"`
	Description string            `json:"description"`
	Items       map[string]string `json:"items,omitempty"`
}

// https://github.com/wordnik/swagger-core/wiki/authorizations
type Authorization struct {
	LocalOAuth OAuth  `json:"local-oauth"`
	ApiKey     ApiKey `json:"apiKey"`
}

// https://github.com/wordnik/swagger-core/wiki/authorizations
type OAuth struct {
	Type       string               `json:"type"`   // e.g. oauth2
	Scopes     []string             `json:"scopes"` // e.g. PUBLIC
	GrantTypes map[string]GrantType `json:"grantTypes"`
}

// https://github.com/wordnik/swagger-core/wiki/authorizations
type GrantType struct {
	LoginEndpoint        Endpoint `json:"loginEndpoint"`
	TokenName            string   `json:"tokenName"` // e.g. access_code
	TokenRequestEndpoint Endpoint `json:"tokenRequestEndpoint"`
	TokenEndpoint        Endpoint `json:"tokenEndpoint"`
}

// https://github.com/wordnik/swagger-core/wiki/authorizations
type Endpoint struct {
	Url              string `json:"url"`
	ClientIdName     string `json:"clientIdName"`
	ClientSecretName string `json:"clientSecretName"`
	TokenName        string `json:"tokenName"`
}

// https://github.com/wordnik/swagger-core/wiki/authorizations
type ApiKey struct {
	Type   string `json:"type"`   // e.g. apiKey
	PassAs string `json:"passAs"` // e.g. header
}
