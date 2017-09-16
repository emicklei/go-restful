package restful

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// go test -v -test.run TestOptionsFilter ...restful
func TestOptionsFilter(t *testing.T) {
	tearDown()
	ws := new(WebService)
	optionsFilter := OptionsFilter{Container: DefaultContainer}
	ws.Route(ws.GET("/candy/{kind}").To(dummy))
	ws.Route(ws.DELETE("/candy/{kind}").To(dummy))
	ws.Route(ws.POST("/candies").To(dummy))
	ws.Route(ws.PUT("/candies").To(dummy))
	Add(ws)
	Filter(optionsFilter.Filter)

	httpRequest, _ := http.NewRequest("OPTIONS", "http://here.io/candy/gum", nil)
	httpWriter := httptest.NewRecorder()
	DefaultContainer.dispatch(httpWriter, httpRequest)
	actual := httpWriter.Header().Get(HEADER_Allow)
	if "GET,DELETE" != actual {
		t.Fatal("expected: GET,DELETE but got:" + actual)
	}

	httpRequest, _ = http.NewRequest("OPTIONS", "http://here.io/candies", nil)
	httpWriter = httptest.NewRecorder()
	DefaultContainer.dispatch(httpWriter, httpRequest)
	actual = httpWriter.Header().Get(HEADER_Allow)
	if "POST,PUT" != actual {
		t.Fatal("expected: POST,PUT but got:" + actual)
	}

}
