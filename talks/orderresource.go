package main

import (
	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful/swagger"
	"log"
	"net/http"
)

type Order struct{}

type OrderResource struct {
	// typically reference a DAO (data-access-object)
}

func (o OrderResource) Register() {
	ws := new(restful.WebService).Path("/customers/{customer_id}/orders")
	ws.Consumes(restful.MIME_XML)
	ws.Produces(restful.MIME_XML)

	ws.Route(ws.GET("/").To(o.getOrdersForCustomer).
		Doc("return the orders of a customer").
		Param(ws.PathParameter("customer_id", "identifier of the customer").DataType("string")).
		Writes([]Order{})) // on the response

	restful.Add(ws)
}

func (o OrderResource) getOrdersForCustomer(req *restful.Request, resp *restful.Response) {
	log.Print("enter getOrdersForCustomer:" + req.Request.URL.Path)
}

func main() {
	OrderResource{}.Register()

	config := swagger.Config{
		WebServicesUrl:  "http://localhost:8080",
		ApiPath:         "/apidocs.json",
		SwaggerPath:     "/apidocs/",
		SwaggerFilePath: "/Users/emicklei/Downloads/swagger-ui-1.1.7",
		WebServices:     restful.RegisteredWebServices()} // you control what services are visible
	swagger.InstallSwaggerService(config)

	http.ListenAndServe(":8080", nil)
}
