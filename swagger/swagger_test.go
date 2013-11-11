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
	value := Api{Path: "/", Description: "Some Path", Operations: []Operation{}}
	output, _ := json.MarshalIndent(value, " ", " ")
	print(string(output))
}

func TestModelToJsonSchema(t *testing.T) {
	sws := newSwaggerService(Config{})
	decl := ApiDeclaration{Models: map[string]Model{}}
	op := new(Operation)
	op.Nickname = "getSome"
	sws.addModelFromSampleTo(op, true, sample{items: []item{}}, &decl)
	output, _ := json.MarshalIndent(decl.Models, " ", " ")
	os.Stdout.Write(output)
}

// go test -v -test.run TestCreateModelFromRecursiveDataStructure ...swagger
func TestCreateModelFromRecursiveDataStructure(t *testing.T) {
	sws := newSwaggerService(Config{})
	decl := ApiDeclaration{Models: map[string]Model{}}
	op := new(Operation)
	op.Nickname = "getSome"
	sws.addModelFromSampleTo(op, true, File{}, &decl)
	output, _ := json.MarshalIndent(decl.Models, " ", " ")
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
	sws := newSwaggerService(cfg)
	decl := sws.composeDeclaration("/tests")
	output, _ := json.MarshalIndent(decl, " ", " ")
	os.Stdout.Write(output)
}

func dummy(i *restful.Request, o *restful.Response) {}
