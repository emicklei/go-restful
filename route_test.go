package restful

import (
	"testing"
)

// accept should match produces
func TestMatchesAcceptConsumesEmpty(t *testing.T) {
	r := Route{Consumes: []string{}}
	if r.matchesAccept("application/xml") {
		t.Errorf("accept should not match xml")
	}
	if !r.matchesAccept("*/*") {
		t.Errorf("accept should match star")
	}
}

// accept should match produces
func TestMatchesAcceptConsumeStar(t *testing.T) {
	//	r := Route{Consumes: []string{"*/*"}}
	//	if !r.matchesAccept("application/xml") {
	//		t.Errorf("accept should match xml")
	//	}
	//	if !r.matchesAccept("*/*") {
	//		t.Errorf("accept should match star")
	//	}
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
	params := doMatchesPath("/from/{source}", 3, "/from/here", true, t)
	if params["source"] != "here" {
		t.Errorf("parameter mismatch here")
	}

	doMatchesPath("/", 2, "/", true, t)

	params = doMatchesPath("/from/{source}/to/{destination}", 5, "/from/AMS/to/NY", true, t)
	if params["source"] != "AMS" {
		t.Errorf("parameter mismatch AMS")
	}

	params = doMatchesPath("{}/from/{source}/", 4, "/from/SOS/", true, t)
	if params["source"] != "SOS" {
		t.Errorf("parameter mismatch SOS")
	}
}

func doMatchesPath(routePath string, size int, urlPath string, shouldMatch bool, t *testing.T) map[string]string {
	r := Route{Path: routePath}
	r.postBuild()
	if len(r.pathParts) != size {
		t.Fatalf("len not %v %v, but %v", size, r.pathParts, len(r.pathParts))
	}
	matches, params := r.matchesPath(urlPath)
	if matches != shouldMatch {
		t.Errorf("disagree about matches: %v", routePath)
		return params
	}
	return params
}
