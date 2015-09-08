package restful

import (
	"fmt"
	"io"
	"net/http/httptest"
	"reflect"
	"testing"
)

type keyvalue struct{}

func (kv keyvalue) Read(req *Request, v interface{}) error {
	return nil
}

func (kv keyvalue) Write(resp *Response, v interface{}) error {
	t := reflect.TypeOf(v)
	rv := reflect.ValueOf(v)
	for ix := 0; ix < t.NumField(); ix++ {
		sf := t.Field(ix)
		io.WriteString(resp, sf.Name)
		io.WriteString(resp, "=")
		io.WriteString(resp, fmt.Sprintf("%v\n", rv.Field(ix).Interface()))
	}
	return nil
}

// go test -v -test.run TestKeyValueEncoding ...restful
func TestKeyValueEncoding(t *testing.T) {
	type Book struct {
		Title         string
		Author        string
		PublishedYear int
	}
	RegisterEntityAccessor("application/kv", keyvalue{})
	b := Book{"Singing for Dummies", "john doe", 2015}
	httpWriter := httptest.NewRecorder()
	//								Accept									Produces
	resp := Response{httpWriter, "application/kv,*/*;q=0.8", []string{"application/kv"}, 0, 0, true, nil}
	resp.WriteEntity(b)
	t.Log(string(httpWriter.Body.Bytes()))
}
