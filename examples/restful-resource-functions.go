package main

import (
	"github.com/emicklei/go-restful"
	"log"
	"net/http"
)

type Product struct {
	Id, Title string
}

type ProductResource struct {
	// typically reference a DAO (data-access-object)
}

func (p ProductResource) getOne(req *restful.Request, resp *restful.Response) {
	id := req.PathParameter("id")
	log.Println("getting product with id:" + id)
	resp.WriteEntity(Product{Id: id, Title: "test"})
}

func (p ProductResource) postOne(req *restful.Request, resp *restful.Response) {
	updatedProduct := new(Product)
	err := req.ReadEntity(updatedProduct)
	if err != nil { // bad request
		resp.WriteError(http.StatusBadRequest, err)
	}
	log.Println("updating product with id:" + updatedProduct.Id)
}

func (p ProductResource) Register() {
	ws := new(restful.WebService).Path("/products")
	ws.Consumes(restful.MIME_XML)
	ws.Produces(restful.MIME_XML)
	ws.Route(ws.GET("/{id}").To(p.getOne))
	ws.Route(ws.POST("").To(p.postOne))
	restful.Add(ws)
}

func main() {
	ProductResource{}.Register()
	http.ListenAndServe(":8080", nil)
}
