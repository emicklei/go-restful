package restful

import (
	"testing"
)

func TestParameter(t *testing.T) {
	p := &Parameter{name: "name", description: "desc"}
	p.AllowMultiple(true)
	p.DataType("int")
	p.Required(true)
	values := map[string]string{"a": "b"}
	p.AllowableValues(values)
	p.bePath()

	ws := new(WebService)
	ws.Param(p)
	if ws.pathParameters[0].name != "name" {
		t.Error("path parameter (or name) invalid")
	}
}
func TestWebService_CanCreateParameterKinds(t *testing.T) {
	ws := new(WebService)
	if ws.BodyParameter("b", "b").kind != BODY_PARAMETER {
		t.Error("body parameter expected")
	}
	if ws.PathParameter("p", "p").kind != PATH_PARAMETER {
		t.Error("path parameter expected")
	}
	if ws.QueryParameter("q", "q").kind != QUERY_PARAMETER {
		t.Error("query parameter expected")
	}
}
