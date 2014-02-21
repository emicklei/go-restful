package restful

import (
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestQueryParameter(t *testing.T) {
	hreq := http.Request{Method: "GET"}
	hreq.URL, _ = url.Parse("http://www.google.com/search?q=foo&q=bar")
	rreq := Request{Request: &hreq}
	if rreq.QueryParameter("q") != "foo" {
		t.Errorf("q!=foo %#v", rreq)
	}
}

type Sample struct {
	Value string
}

func TestReadEntityXml(t *testing.T) {
	bodyReader := strings.NewReader("<Sample><Value>42</Value></Sample>")
	httpRequest, _ := http.NewRequest("GET", "/test", bodyReader)
	httpRequest.Header.Set("Content-Type", "application/xml")
	request := &Request{Request: httpRequest}
	sam := new(Sample)
	request.ReadEntity(sam)
	if sam.Value != "42" {
		t.Fatal("read failed")
	}
}

func TestReadEntityJson(t *testing.T) {
	bodyReader := strings.NewReader(`{"Value" : "42"}`)
	httpRequest, _ := http.NewRequest("GET", "/test", bodyReader)
	httpRequest.Header.Set("Content-Type", "application/json")
	request := &Request{Request: httpRequest}
	sam := new(Sample)
	request.ReadEntity(sam)
	if sam.Value != "42" {
		t.Fatal("read failed")
	}
}

func TestReadEntityJsonCharset(t *testing.T) {
	bodyReader := strings.NewReader(`{"Value" : "42"}`)
	httpRequest, _ := http.NewRequest("GET", "/test", bodyReader)
	httpRequest.Header.Set("Content-Type", "application/json; charset=UTF-8")
	request := newRequest(httpRequest)
	sam := new(Sample)
	request.ReadEntity(sam)
	if sam.Value != "42" {
		t.Fatal("read failed")
	}
}

func TestBodyParameter(t *testing.T) {
	bodyReader := strings.NewReader(`value1=42&value2=43`)
	httpRequest, _ := http.NewRequest("POST", "/test?value1=44", bodyReader) // POST and PUT body parameters take precedence over URL query string
	httpRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	request := newRequest(httpRequest)
	v1, err := request.BodyParameter("value1")
	if err != nil {
		t.Error(err)
	}
	v2, err := request.BodyParameter("value2")
	if err != nil {
		t.Error(err)
	}
	if v1 != "42" || v2 != "43" {
		t.Fatal("read failed")
	}
}

func TestReadEntityUnkown(t *testing.T) {
	bodyReader := strings.NewReader("?")
	httpRequest, _ := http.NewRequest("GET", "/test", bodyReader)
	httpRequest.Header.Set("Content-Type", "application/rubbish")
	request := newRequest(httpRequest)
	sam := new(Sample)
	err := request.ReadEntity(sam)
	if err == nil {
		t.Fatal("read should be in error")
	}
}

func TestSetAttribute(t *testing.T) {
	bodyReader := strings.NewReader("?")
	httpRequest, _ := http.NewRequest("GET", "/test", bodyReader)
	request := newRequest(httpRequest)
	request.SetAttribute("go", "there")
	there := request.Attribute("go")
	if there != "there" {
		t.Fatalf("missing request attribute:%v", there)
	}
}
