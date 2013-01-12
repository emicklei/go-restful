package main

import (
	"github.com/emicklei/go-restful"
	"log"
	"net/http"
)

type User struct {
	Id, Name string
}

var users = map[string]User{}

func NewUserService() *restful.WebService {
	ws := new(restful.WebService)
	ws.
		Path("/users").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML)

	ws.Route(ws.GET("/{user-id}").To(findUser).
		// docs	
		Doc("get a user").
		Param(ws.PathParameter("user-id", "identifier of the user")))

	ws.Route(ws.POST("").To(updateUser).
		// docs	
		Doc("update a user"))

	ws.Route(ws.PUT("/{user-id}").To(createUser).
		// docs	
		Doc("create a user").
		Param(ws.PathParameter("user-id", "identifier of the user")))

	ws.Route(ws.DELETE("/{user-id}").To(removeUser).
		// docs	
		Doc("delete a user").
		Param(ws.PathParameter("user-id", "identifier of the user")))
	
	return ws
}

func findUser(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("user-id")
	usr := users[id]
	if len(usr.Id) == 0 {
		response.WriteError(http.StatusNotFound, nil)
	} else {
		response.WriteEntity(usr)
	}
}

func updateUser(request *restful.Request, response *restful.Response) {
	usr := new(User)
	err := request.ReadEntity(&usr)
	if err == nil {
		users[usr.Id] = *usr
		response.WriteEntity(usr)
	} else {
		response.WriteError(http.StatusInternalServerError, err)
	}
}

func createUser(request *restful.Request, response *restful.Response) {
	usr := User{Id: request.PathParameter("user-id")}
	err := request.ReadEntity(&usr)
	if err == nil {
		users[usr.Id] = usr
		response.WriteHeader(http.StatusCreated)
		response.WriteEntity(usr)
	} else {
		response.WriteError(http.StatusInternalServerError, err)
	}
}

func removeUser(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("user-id")
	delete(users, id)
}

func main() {
	restful.Add(NewUserService())
	
	// Open http://localhost:8080/apidocs and enter http://localhost:8080/apidocs.json in the api input field.
	config := restful.SwaggerConfig{ 
		WebServicesUrl: "http://localhost:8080",
		ApiPath: "/apidocs.json",
		SwaggerPath: "/apidocs/",
		SwaggerFilePath: "/Users/emicklei/Downloads/swagger-ui-1.1.7" }	
	restful.InstallSwaggerService(config)
	
	log.Printf("start listening on localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
