package swagger

import (
	"encoding/json"
	"github.com/emicklei/go-restful"
	"os"
	"testing"
)

type sample struct {
	id       string `swagger:"required"` // TODO
	items    []item
	rootItem item `json:"root"`
}

type item struct {
	itemName string `json:"name"`
}

type File struct {
	History []File
}

func TestApi(t *testing.T) {
	value := Api{Path: "/", Description: "Some Path", Operations: []Operation{}, Models: map[string]Model{}}
	output, _ := json.MarshalIndent(value, " ", " ")
	print(string(output))
}

func TestModelToJsonSchema(t *testing.T) {
	api := new(Api)
	api.Models = map[string]Model{}
	op := new(Operation)
	op.Nickname = "getSome"
	addModelFromSample(api, op, true, sample{items: []item{}})
	output, _ := json.MarshalIndent(api, " ", " ")
	os.Stdout.Write(output)
}

// go test -v -test.run TestCreateModelFromRecursiveDataStructure ...swagger
func TestCreateModelFromRecursiveDataStructure(t *testing.T) {
	api := new(Api)
	api.Models = map[string]Model{}
	op := new(Operation)
	op.Nickname = "getSome"
	addModelFromSample(api, op, true, File{})
	output, _ := json.MarshalIndent(api, " ", " ")
	os.Stdout.Write(output)
}

// go test -v -test.run TestServiceToApi ...swagger
func TestServiceToApi(t *testing.T) {
	ws := new(restful.WebService)
	ws.Path("/tests")
	ws.Consumes(restful.MIME_JSON)
	ws.Produces(restful.MIME_XML)
	ws.Route(ws.GET("/all").To(dummy).Writes(sample{}))
	cfg := Config{
		WebServicesUrl: "http://here.com",
		ApiPath:        "/apipath",
		WebServices:    []*restful.WebService{ws}}
	decl := composeDeclaration("/tests", cfg)
	output, _ := json.MarshalIndent(decl, " ", " ")
	os.Stdout.Write(output)
}

func dummy(i *restful.Request, o *restful.Response) {}
