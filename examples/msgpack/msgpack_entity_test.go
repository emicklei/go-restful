package restPack

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	restful "github.com/emicklei/go-restful"
)

func TestMsgPack(t *testing.T) {

	// register msg pack entity
	restful.RegisterEntityAccessor(MIME_MSGPACK, NewEntityAccessorMsgPack(MIME_MSGPACK))
	type Tool struct {
		Name   string
		Vendor string
	}

	// Write
	httpWriter := httptest.NewRecorder()
	mpack := &Tool{Name: "json", Vendor: "apple"}
	resp := restful.NewResponse(httpWriter)
	resp.SetRequestAccepts("application/x-msgpack,*/*;q=0.8")

	err := resp.WriteEntity(mpack)
	if err != nil {
		t.Errorf("err %v", err)
	}

	// Read
	bodyReader := bytes.NewReader(httpWriter.Body.Bytes())
	httpRequest, _ := http.NewRequest("GET", "/test", bodyReader)
	httpRequest.Header.Set("Content-Type", "application/x-msgpack; charset=UTF-8")
	request := restful.NewRequest(httpRequest)
	readMsgPack := new(Tool)
	err = request.ReadEntity(&readMsgPack)
	if err != nil {
		t.Errorf("err %v", err)
	}
	if equal := reflect.DeepEqual(mpack, readMsgPack); !equal {
		t.Fatalf("should not be error")
	}
}
