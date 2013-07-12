package swagger

import "github.com/emicklei/go-restful"

type Config struct {
	WebServicesUrl  string // url where the services are available, e.g. http://localhost:8080
	ApiPath         string // path where the JSON api is avaiable , e.g. /apidocs
	SwaggerPath     string // [optional] path where the swagger UI will be served, e.g. /swagger
	SwaggerFilePath string // [optional] location of folder containing Swagger HTML5 application index.html
	WebServices     []*restful.WebService
}
