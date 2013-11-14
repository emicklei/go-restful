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
	History     []File
	HistoryPtrs []*File
}

// go test -v -test.run TestApi ...swagger
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

// go test -v -test.run TestIssue78 ...swagger
type Response struct {
	Code  int
	Users *[]User
	Items *[]Item
}
type User struct {
	Id, Name string
}
type Item struct {
	Id, Name string
}

func TestIssue78(t *testing.T) {
	sws := newSwaggerService(Config{})
	decl := ApiDeclaration{Models: map[string]Model{}}
	sws.addModelFromSampleTo(&Operation{}, true, Response{Items: &[]Item{}}, &decl)
	model, ok := decl.Models["swagger.Response"]
	if !ok {
		t.Fatal("missing response model")
	}
	if "swagger.Response" != model.Id {
		t.Fatal("wrong model id:" + model.Id)
	}
	code, ok := model.Properties["Code"]
	if !ok {
		t.Fatal("missing code")
	}
	if "int" != code.Type {
		t.Fatal("wrong code type:" + code.Type)
	}
	items, ok := model.Properties["Items"]
	if !ok {
		t.Fatal("missing items")
	}
	if "array" != items.Type {
		t.Fatal("wrong items type:" + items.Type)
	}
	items_items := items.Items
	if items_items == nil {
		t.Fatal("missing items->items")
	}
	ref := items_items["$ref"]
	if ref == "" {
		t.Fatal("missing $ref")
	}
	if ref != "swagger.Item" {
		t.Fatal("wrong $ref:" + ref)
	}
	output, _ := json.MarshalIndent(decl, " ", " ")
	os.Stdout.Write(output)
}
