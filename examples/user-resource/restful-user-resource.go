package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	restful "github.com/emicklei/go-restful/v3"
	"github.com/go-openapi/spec"
)

// UserResource is the REST layer to the User domain
type UserResource struct {
	// normally one would use DAO (data access object)
	users map[string]User
}

// WebService creates a new service that can handle REST requests for User resources.
func (u UserResource) WebService() *restful.WebService {
	ws := new(restful.WebService)
	ws.
		Path("/users").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML) // you can specify this per route as well

	tags := []string{"users"}

	ws.Route(ws.GET("/").To(u.findAllUsers).
		// docs
		Doc("get all users").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Writes([]User{}).
		Returns(200, "OK", []User{}))

	ws.Route(ws.GET("/{user-id}").To(u.findUser).
		// docs
		Doc("get a user").
		Param(ws.PathParameter("user-id", "identifier of the user").DataType("integer").DefaultValue("1")).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Writes(User{}). // on the response
		Returns(200, "OK", User{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.PUT("/{user-id}").To(u.upsertUser).
		// docs
		Doc("update a user").
		Param(ws.PathParameter("user-id", "identifier of the user").DataType("string")).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Reads(User{})) // from the request

	ws.Route(ws.POST("").To(u.createUser).
		// docs
		Doc("create a user").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Reads(User{})) // from the request

	ws.Route(ws.DELETE("/{user-id}").To(u.removeUser).
		// docs
		Doc("delete a user").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.PathParameter("user-id", "identifier of the user").DataType("string")))

	return ws
}

// GET http://localhost:8080/users
func (u UserResource) findAllUsers(request *restful.Request, response *restful.Response) {
	log.Println("findAllUsers")
	list := []User{}
	for _, each := range u.users {
		list = append(list, each)
	}
	response.WriteEntity(list)
}

// GET http://localhost:8080/users/1
func (u UserResource) findUser(request *restful.Request, response *restful.Response) {
	log.Println("findUser")
	id := request.PathParameter("user-id")
	usr := u.users[id]
	if len(usr.ID) == 0 {
		response.WriteErrorString(http.StatusNotFound, "User could not be found.")
	} else {
		response.WriteEntity(usr)
	}
}

// PUT http://localhost:8080/users/1
// <User><Id>1</Id><Name>Melissa Raspberry</Name></User>
func (u *UserResource) upsertUser(request *restful.Request, response *restful.Response) {
	log.Println("upsertUser")
	usr := User{ID: request.PathParameter("user-id")}
	err := request.ReadEntity(&usr)
	if err == nil {
		u.users[usr.ID] = usr
		response.WriteEntity(usr)
	} else {
		response.WriteError(http.StatusInternalServerError, err)
	}
}

// POST http://localhost:8080/users
// <User><Id>1</Id><Name>Melissa</Name></User>
func (u *UserResource) createUser(request *restful.Request, response *restful.Response) {
	log.Println("createUser")
	usr := User{ID: fmt.Sprintf("%d", time.Now().Unix())}
	err := request.ReadEntity(&usr)
	if err == nil {
		u.users[usr.ID] = usr
		response.WriteHeaderAndEntity(http.StatusCreated, usr)
	} else {
		response.WriteError(http.StatusInternalServerError, err)
	}
}

// DELETE http://localhost:8080/users/1
func (u *UserResource) removeUser(request *restful.Request, response *restful.Response) {
	log.Println("removeUser")
	id := request.PathParameter("user-id")
	delete(u.users, id)
}

func main() {
	u := UserResource{map[string]User{}}
	restful.DefaultContainer.Add(u.WebService())

	config := restfulspec.Config{
		WebServices:                   restful.RegisteredWebServices(), // you control what services are visible
		APIPath:                       "/apidocs.json",
		PostBuildSwaggerObjectHandler: enrichSwaggerObject}
	restful.DefaultContainer.Add(restfulspec.NewOpenAPIService(config))

	// Optionally, you can install the Swagger Service which provides a nice Web UI on your REST API
	// You need to download the Swagger HTML5 assets and change the FilePath location in the config below.
	// Open http://localhost:8080/apidocs/?url=http://localhost:8080/apidocs.json
	http.Handle("/apidocs/", http.StripPrefix("/apidocs/", http.FileServer(http.Dir("/Users/emicklei/Projects/swagger-ui/dist"))))

	log.Printf("start listening on localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func enrichSwaggerObject(swo *spec.Swagger) {
	swo.Info = &spec.Info{
		InfoProps: spec.InfoProps{
			Title:       "UserService",
			Description: "Resource for managing Users",
			Contact: &spec.ContactInfo{
				ContactInfoProps: spec.ContactInfoProps{
					Name:  "john",
					Email: "john@doe.rp",
					URL:   "http://johndoe.org",
				},
			},
			License: &spec.License{
				LicenseProps: spec.LicenseProps{
					Name: "MIT",
					URL:  "http://mit.org",
				},
			},
			Version: "1.0.0",
		},
	}
	swo.Tags = []spec.Tag{spec.Tag{TagProps: spec.TagProps{
		Name:        "users",
		Description: "Managing users"}}}
}

// User is just a sample type
type User struct {
	ID   string `xml:"id" json:"id" description:"identifier of the user"`
	Name string `xml:"name" json:"name" description:"name of the user" default:"john"`
	Age  int    `xml:"age" json:"age" description:"age of the user" default:"21"`
}
