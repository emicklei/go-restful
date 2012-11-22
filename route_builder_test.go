package restful

import (
	"testing"
)

func TestRouteBuilder(t *testing.T) {
	json := "application/json"
	b := new(RouteBuilder)
	b.Path("/routes").Method("HEAD").Accept(json).ContentType(json)
	r := b.Build()
	if r.Path != "/routes" {
		t.Error("path invalid")
	}
	if r.Produces != json {
		t.Error("produces invalid")
	}
	if r.Consumes != json {
		t.Error("consumes invalid")
	}
}
