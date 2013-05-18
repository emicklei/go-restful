package restful

import (
	"testing"
)

func TestParameter(t *testing.T) {
	p := &Parameter{&ParameterData{Name: "name", Description: "desc"}}
	p.AllowMultiple(true)
	p.DataType("int")
	p.Required(true)
	values := map[string]string{"a": "b"}
	p.AllowableValues(values)
	p.bePath()

	ws := new(WebService)
	ws.Param(p)
	if ws.pathParameters[0].Data().Name != "name" {
		t.Error("path parameter (or name) invalid")
	}
}
func TestWebService_CanCreateParameterKinds(t *testing.T) {
	ws := new(WebService)
	if ws.BodyParameter("b", "b").Kind() != BODY_PARAMETER {
		t.Error("body parameter expected")
	}
	if ws.PathParameter("p", "p").Kind() != PATH_PARAMETER {
		t.Error("path parameter expected")
	}
	if ws.QueryParameter("q", "q").Kind() != QUERY_PARAMETER {
		t.Error("query parameter expected")
	}
}
