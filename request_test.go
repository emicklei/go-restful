package restful

import (
    "testing"    
    "net/http"
    "net/url"
)

func TestPathParameter(t *testing.T) {
	hreq := http.Request{Method: "GET"}
	hreq.URL, _ = url.Parse("http://www.google.com/search?q=foo&q=bar")
	rreq := &Request{http.Request:hreq}
	if (rreq.QueryParameter("q") != "foo") {
		t.Errorf("q!=foo %#v", rreq)
	}
}	