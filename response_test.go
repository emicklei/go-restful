package restful

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWriteHeader(t *testing.T) {
	httpWriter := httptest.NewRecorder()
	resp := Response{httpWriter, "*/*", []string{"*/*"}, 0, 0}
	resp.WriteHeader(123)
	if resp.StatusCode() != 123 {
		t.Errorf("Unexpected status code:%d", resp.StatusCode())
	}
}

func TestNoWriteHeader(t *testing.T) {
	httpWriter := httptest.NewRecorder()
	resp := Response{httpWriter, "*/*", []string{"*/*"}, 0, 0}
	if resp.StatusCode() != http.StatusOK {
		t.Errorf("Unexpected status code:%d", resp.StatusCode())
	}
}

type food struct {
	Kind string
}

// go test -v -test.run TestMeasureContentLengthXml ...restful
func TestMeasureContentLengthXml(t *testing.T) {
	httpWriter := httptest.NewRecorder()
	resp := Response{httpWriter, "*/*", []string{"*/*"}, 0, 0}
	resp.WriteAsXml(food{"apple"})
	if resp.ContentLength() != 76 {
		t.Errorf("Incorrect measured length:%d", resp.ContentLength())
	}
}

// go test -v -test.run TestMeasureContentLengthJson ...restful
func TestMeasureContentLengthJson(t *testing.T) {
	httpWriter := httptest.NewRecorder()
	resp := Response{httpWriter, "*/*", []string{"*/*"}, 0, 0}
	resp.WriteAsJson(food{"apple"})
	if resp.ContentLength() != 22 {
		t.Errorf("Incorrect measured length:%d", resp.ContentLength())
	}
	//println(httpWriter.Body.String())
}
