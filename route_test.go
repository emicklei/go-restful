package restful

import (
	"testing"
)

// accept should match produces
func TestMatchesAcceptStar(t *testing.T) {
	r := Route{Produces: []string{"application/xml"}}
	if !r.matchesAccept("*/*") {
		t.Errorf("accept should match star")
	}
}

// accept should match produces
func TestMatchesAcceptXml(t *testing.T) {
	r := Route{Produces: []string{"application/xml"}}
	if r.matchesAccept("application/json") {
		t.Errorf("accept should not match json")
	}
	if !r.matchesAccept("application/xml") {
		t.Errorf("accept should match xml")
	}
}

// content type should match consumes
func TestMatchesContentTypeXml(t *testing.T) {
	r := Route{Consumes: []string{"application/xml"}}
	if r.matchesContentType("application/json") {
		t.Errorf("accept should not match json")
	}
	if !r.matchesContentType("application/xml") {
		t.Errorf("accept should match xml")
	}
}

func TestMatchesPath(t *testing.T) {
	params := doExtractParams("/from/{source}", 3, "/from/here", t)
	if params["source"] != "here" {
		t.Errorf("parameter mismatch here")
	}

	params = doExtractParams("/", 2, "/", t)
	if len(params) != 0 {
		t.Errorf("expected empty parameters")
	}

	params = doExtractParams("/from/{source}/to/{destination}", 5, "/from/AMS/to/NY", t)
	if params["source"] != "AMS" {
		t.Errorf("parameter mismatch AMS")
	}

	params = doExtractParams("{}/from/{source}/", 4, "/from/SOS/", t)
	if params["source"] != "SOS" {
		t.Errorf("parameter mismatch SOS")
	}
}

func doExtractParams(routePath string, size int, urlPath string, t *testing.T) map[string]string {
	r := Route{Path: routePath}
	r.postBuild()
	if len(r.pathParts) != size {
		t.Fatalf("len not %v %v, but %v", size, r.pathParts, len(r.pathParts))
	}
	return r.extractParameters(urlPath)
}
