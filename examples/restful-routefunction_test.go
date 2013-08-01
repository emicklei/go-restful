package main

import (
	"github.com/emicklei/go-restful"
	"net/http"
	"net/http/httptest"
	"testing"
)

func getIt(req *restful.Request, resp *restful.Response) {
	resp.WriteHeader(404)
}

func TestCallFunction(t *testing.T) {
	req := new(restful.Request)
	req.Request = new(http.Request)

	resp := new(restful.Response)
	recorder := new(httptest.ResponseRecorder)
	resp.ResponseWriter = recorder

	getIt(req, resp)
	if recorder.Code != 404 {
		t.Logf("Missing or wrong status code:%d", recorder.Code)
	}
}
