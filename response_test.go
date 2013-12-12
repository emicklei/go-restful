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
}

// go test -v -test.run TestMeasureContentLengthWriteErrorString ...restful
func TestMeasureContentLengthWriteErrorString(t *testing.T) {
	httpWriter := httptest.NewRecorder()
	resp := Response{httpWriter, "*/*", []string{"*/*"}, 0, 0}
	resp.WriteErrorString(404, "Invalid")
	if resp.ContentLength() != len("Invalid") {
		t.Errorf("Incorrect measured length:%d", resp.ContentLength())
	}
}

// go test -v -test.run TestStatusCreatedAndContentTypeJson_Issue54 ...restful
func TestStatusCreatedAndContentTypeJson_Issue54(t *testing.T) {
	httpWriter := httptest.NewRecorder()
	resp := Response{httpWriter, "application/json", []string{"application/json"}, 0, 0}
	resp.WriteHeader(201)
	resp.WriteAsJson(food{"Juicy"})
	if httpWriter.HeaderMap.Get("Content-Type") != "application/json" {
		t.Errorf("Expected content type json but got:%d", httpWriter.HeaderMap.Get("Content-Type"))
	}
	if httpWriter.Code != 201 {
		t.Errorf("Expected status 201 but got:%d", httpWriter.Code)
	}
}
