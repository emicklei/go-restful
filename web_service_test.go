package restful

import (
	"net/http"
	"net/http/httptest"
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

func TestCapturePanic(t *testing.T) {
	tearDown()
	Add(newPanicingService())
	httpRequest, _ := http.NewRequest("GET", "http://here.com/fire", nil)
	httpRequest.Header.Set("Accept", "*/*")
	httpWriter := httptest.NewRecorder()
	DefaultContainer.dispatch(httpWriter, httpRequest)
	if 500 != httpWriter.Code {
		t.Error("500 expected on fire")
	}
}

func TestNotFound(t *testing.T) {
	tearDown()
	httpRequest, _ := http.NewRequest("GET", "http://here.com/missing", nil)
	httpRequest.Header.Set("Accept", "*/*")
	httpWriter := httptest.NewRecorder()
	DefaultContainer.dispatch(httpWriter, httpRequest)
	if 404 != httpWriter.Code {
		t.Error("404 expected on missing")
	}
}

func TestMethodNotAllowed(t *testing.T) {
	tearDown()
	Add(newGetOnlyService())
	httpRequest, _ := http.NewRequest("POST", "http://here.com/get", nil)
	httpRequest.Header.Set("Accept", "*/*")
	httpWriter := httptest.NewRecorder()
	DefaultContainer.dispatch(httpWriter, httpRequest)
	if 405 != httpWriter.Code {
		t.Error("405 expected method not allowed")
	}
}

func newPanicingService() *WebService {
	ws := new(WebService).Path("")
	ws.Route(ws.GET("/fire").To(doPanic))
	return ws
}

func newGetOnlyService() *WebService {
	ws := new(WebService).Path("")
	ws.Route(ws.GET("/get").To(doPanic))
	return ws
}

func doPanic(req *Request, resp *Response) {
	println("lightning...")
	panic("fire")
}
