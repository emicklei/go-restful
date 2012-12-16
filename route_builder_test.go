package restful

import (
	"testing"
)

func TestRouteBuilder_PathParameter(t *testing.T) {
	p := &Parameter{name: "name", description: "desc"}
	p.AllowMultiple(true)
	p.DataType("int")
	p.Required(true)
	values := map[string]string{"a": "b"}
	p.AllowableValues(values)
	p.bePath()

	b := new(RouteBuilder)
	b.Param(p)
	r := b.Build()
	if !r.parameterDocs[0].allowMultiple {
		t.Error("AllowMultiple invalid")
	}
	if r.parameterDocs[0].dataType != "int" {
		t.Error("dataType invalid")
	}
	if !r.parameterDocs[0].required {
		t.Error("required invalid")
	}
	if r.parameterDocs[0].kind != PATH_PARAMETER {
		t.Error("kind invalid")
	}
	if r.parameterDocs[0].allowableValues["a"] != "b" {
		t.Error("allowableValues invalid")
	}
}

func TestRouteBuilder(t *testing.T) {
	json := "application/json"
	b := new(RouteBuilder)
	b.Path("/routes").Method("HEAD").Consumes(json).Produces(json)
	r := b.Build()
	if r.Path != "/routes" {
		t.Error("path invalid")
	}
	if r.Produces[0] != json {
		t.Error("produces invalid")
	}
	if r.Consumes[0] != json {
		t.Error("consumes invalid")
	}
}
