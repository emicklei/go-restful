package restful

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWriteHeader(t *testing.T) {
	httpWriter := httptest.NewRecorder()
	resp := Response{httpWriter, "*/*", []string{"*/*"}, 0}
	resp.WriteHeader(123)
	if resp.StatusCode() != 123 {
		t.Errorf("Unexpected status code:%d", resp.StatusCode())
	}
}

func TestNoWriteHeader(t *testing.T) {
	httpWriter := httptest.NewRecorder()
	resp := Response{httpWriter, "*/*", []string{"*/*"}, 0}
	if resp.StatusCode() != http.StatusOK {
		t.Errorf("Unexpected status code:%d", resp.StatusCode())
	}
}
