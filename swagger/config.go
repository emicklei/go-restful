package swagger

import (
	"net/http"

	"github.com/emicklei/go-restful"
)

type Config struct {
	WebServicesUrl  string                // url where the services are available, e.g. http://localhost:8080
	ApiPath         string                // path where the JSON api is avaiable , e.g. /apidocs
	SwaggerPath     string                // [optional] path where the swagger UI will be served, e.g. /swagger
	SwaggerFilePath string                // [optional] location of folder containing Swagger HTML5 application index.html
	WebServices     []*restful.WebService // api listing is constructed from this list of restful WebServices.
	StaticHandler   http.Handler          // will serve all static content (scripts,pages,images)
	DisableCORS     bool                  // [optional] on default CORS (Cross-Origin-Resource-Sharing) is enabled.
}
