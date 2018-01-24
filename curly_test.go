package restful

import (
	"io"
	"net/http"
	"testing"
)

var requestPaths = []struct {
	// url with path (1) is handled by service with root (2) and remainder has value final (3)
	path, root string
}{
	{"/", "/"},
	{"/p", "/p"},
	{"/p/x", "/p/{q}"},
	{"/q/x", "/q"},
	{"/p/x/", "/p/{q}"},
	{"/p/x/y", "/p/{q}"},
	{"/q/x/y", "/q"},
	{"/z/q", "/{p}/q"},
	{"/a/b/c/q", "/"},
}

// go test -v -test.run TestCurlyDetectWebService ...restful
func TestCurlyDetectWebService(t *testing.T) {
	ws1 := new(WebService).Path("/")
	ws2 := new(WebService).Path("/p")
	ws3 := new(WebService).Path("/q")
	ws4 := new(WebService).Path("/p/q")
	ws5 := new(WebService).Path("/p/{q}")
	ws7 := new(WebService).Path("/{p}/q")
	var wss = []*WebService{ws1, ws2, ws3, ws4, ws5, ws7}

	for _, each := range wss {
		t.Logf("path=%s,toks=%v\n", each.pathExpr.Source, each.pathExpr.tokens)
	}

	router := CurlyRouter{}

	ok := true
	for i, fixture := range requestPaths {
		requestTokens := tokenizePath(fixture.path)
		who := router.detectWebService(requestTokens, wss)
		if who != nil && who.RootPath() != fixture.root {
			t.Logf("[line:%v] Unexpected dispatcher, expected:%v, actual:%v", i, fixture.root, who.RootPath())
			ok = false
		}
	}
	if !ok {
		t.Fail()
	}
}

var serviceDetects = []struct {
	path  string
	found bool
	root  string
}{
	{"/a/b", true, "/{p}/{q}/{r}"},
	{"/p/q", true, "/p/q"},
	{"/q/p", true, "/q"},
	{"/", true, "/"},
	{"/p/q/r", true, "/p/q"},
}

// go test -v -test.run Test_detectWebService ...restful
func Test_detectWebService(t *testing.T) {
	router := CurlyRouter{}
	ws1 := new(WebService).Path("/")
	ws2 := new(WebService).Path("/p")
	ws3 := new(WebService).Path("/q")
	ws4 := new(WebService).Path("/p/q")
	ws5 := new(WebService).Path("/p/{q}")
	ws6 := new(WebService).Path("/p/{q}/")
	ws7 := new(WebService).Path("/{p}/q")
	ws8 := new(WebService).Path("/{p}/{q}/{r}")
	var wss = []*WebService{ws8, ws7, ws6, ws5, ws4, ws3, ws2, ws1}
	for _, fix := range serviceDetects {
		requestPath := fix.path
		requestTokens := tokenizePath(requestPath)
		for _, ws := range wss {
			serviceTokens := ws.pathExpr.tokens
			matches, score := router.computeWebserviceScore(requestTokens, serviceTokens)
			t.Logf("req=%s,toks:%v,ws=%s,toks:%v,score=%d,matches=%v", requestPath, requestTokens, ws.RootPath(), serviceTokens, score, matches)
		}
		best := router.detectWebService(requestTokens, wss)
		if best != nil {
			if fix.found {
				t.Logf("best=%s", best.RootPath())
			} else {
				t.Fatalf("should have found:%s", fix.root)
			}
		}
	}
}

var routeMatchers = []struct {
	route       string
	path        string
	matches     bool
	paramCount  int
	staticCount int
}{
	// route, request-path
	{"/a", "/a", true, 0, 1},
	{"/a", "/b", false, 0, 0},
	{"/a", "/b", false, 0, 0},
	{"/a/{b}/c/", "/a/2/c", true, 1, 2},
	{"/{a}/{b}/{c}/", "/a/b", false, 0, 0},
	{"/{x:*}", "/", false, 0, 0},
	{"/{x:*}", "/a", true, 1, 0},
	{"/{x:*}", "/a/b", true, 1, 0},
	{"/a/{x:*}", "/a/b", true, 1, 1},
	{"/a/{x:[A-Z][A-Z]}", "/a/ZX", true, 1, 1},
	{"/basepath/{resource:*}", "/basepath/some/other/location/test.xml", true, 1, 1},
}

// clear && go test -v -test.run Test_matchesRouteByPathTokens ...restful
func Test_matchesRouteByPathTokens(t *testing.T) {
	router := CurlyRouter{}
	for i, each := range routeMatchers {
		routeToks := tokenizePath(each.route)
		reqToks := tokenizePath(each.path)
		matches, pCount, sCount := router.matchesRouteByPathTokens(routeToks, reqToks)
		if matches != each.matches {
			t.Fatalf("[%d] unexpected matches outcome route:%s, path:%s, matches:%v", i, each.route, each.path, matches)
		}
		if pCount != each.paramCount {
			t.Fatalf("[%d] unexpected paramCount got:%d want:%d ", i, pCount, each.paramCount)
		}
		if sCount != each.staticCount {
			t.Fatalf("[%d] unexpected staticCount got:%d want:%d ", i, sCount, each.staticCount)
		}
	}
}

// clear && go test -v -test.run TestExtractParameters_Wildcard1 ...restful
func TestExtractParameters_Wildcard1(t *testing.T) {
	params := doExtractParams("/fixed/{var:*}", 2, "/fixed/remainder", t)
	if params["var"] != "remainder" {
		t.Errorf("parameter mismatch var: %s", params["var"])
	}
}

// clear && go test -v -test.run TestExtractParameters_Wildcard2 ...restful
func TestExtractParameters_Wildcard2(t *testing.T) {
	params := doExtractParams("/fixed/{var:*}", 2, "/fixed/remain/der", t)
	if params["var"] != "remain/der" {
		t.Errorf("parameter mismatch var: %s", params["var"])
	}
}

// clear && go test -v -test.run TestExtractParameters_Wildcard3 ...restful
func TestExtractParameters_Wildcard3(t *testing.T) {
	params := doExtractParams("/static/{var:*}", 2, "/static/test/sub/hi.html", t)
	if params["var"] != "test/sub/hi.html" {
		t.Errorf("parameter mismatch var: %s", params["var"])
	}
}

func TestExtractParameters_Wildcard4(t *testing.T) {
	params := doExtractParams("/static/{var:*}/sub", 3, "/static/test/sub", t)
	if params["var"] != "test/sub" {
		t.Errorf("parameter mismatch var: %s", params["var"])
	}
}

func TestMatchesPath_OneParam(t *testing.T) {
	params := doExtractParams("/from/{source}", 2, "/from/here", t)
	if params["source"] != "here" {
		t.Errorf("parameter mismatch here")
	}
}

func TestMatchesPath_Slash(t *testing.T) {
	params := doExtractParams("/", 0, "/", t)
	if len(params) != 0 {
		t.Errorf("expected empty parameters")
	}
}

func TestMatchesPath_SlashNonVar(t *testing.T) {
	params := doExtractParams("/any", 1, "/any", t)
	if len(params) != 0 {
		t.Errorf("expected empty parameters")
	}
}

func TestMatchesPath_TwoVars(t *testing.T) {
	params := doExtractParams("/from/{source}/to/{destination}", 4, "/from/AMS/to/NY", t)
	if params["source"] != "AMS" {
		t.Errorf("parameter mismatch AMS")
	}
}

func TestMatchesPath_VarOnFront(t *testing.T) {
	params := doExtractParams("{what}/from/{source}/", 3, "who/from/SOS/", t)
	if params["source"] != "SOS" {
		t.Errorf("parameter mismatch SOS")
	}
}

func doExtractParams(routePath string, size int, urlPath string, t *testing.T) map[string]string {
	ws1 := new(WebService).Path("/")
	ws1.Route(ws1.GET(routePath).To(curlyDummy))
	router := CurlyRouter{}
	req, _ := http.NewRequest(http.MethodGet, urlPath, nil)
	_, r, params, err := router.SelectRoute([]*WebService{ws1}, req)
	if err != nil {
		t.Fatalf("Unexpected error selecting route: %v", err.Error())
	}
	if len(r.pathParts) != size {
		t.Fatalf("len not %v %v, but %v", size, r.pathParts, len(r.pathParts))
	}
	return params
}

// clear && go test -v -test.run TestCurly_ISSUE_34 ...restful
func TestCurly_ISSUE_34(t *testing.T) {
	ws1 := new(WebService).Path("/")
	ws1.Route(ws1.GET("/{type}/{id}").To(curlyDummy))
	ws1.Route(ws1.GET("/network/{id}").To(curlyDummy))
	croutes := CurlyRouter{}.selectRoutes(ws1, tokenizePath("/network/12"))
	if len(croutes) != 2 {
		t.Fatal("expected 2 routes")
	}
	if got, want := croutes[0].route.Path, "/network/{id}"; got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

// clear && go test -v -test.run TestCurly_ISSUE_34_2 ...restful
func TestCurly_ISSUE_34_2(t *testing.T) {
	ws1 := new(WebService)
	ws1.Route(ws1.GET("/network/{id}").To(curlyDummy))
	ws1.Route(ws1.GET("/{type}/{id}").To(curlyDummy))
	croutes := CurlyRouter{}.selectRoutes(ws1, tokenizePath("/network/12"))
	if len(croutes) != 2 {
		t.Fatal("expected 2 routes")
	}
	if got, want := croutes[0].route.Path, "/network/{id}"; got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

// clear && go test -v -test.run TestCurly_JsonHtml ...restful
func TestCurly_JsonHtml(t *testing.T) {
	ws1 := new(WebService)
	ws1.Path("/")
	ws1.Route(ws1.GET("/some.html").To(curlyDummy).Consumes("*/*").Produces("text/html"))
	req, _ := http.NewRequest("GET", "/some.html", nil)
	req.Header.Set("Accept", "application/json")
	_, route, _, err := CurlyRouter{}.SelectRoute([]*WebService{ws1}, req)
	if err == nil {
		t.Error("error expected")
	}
	if route != nil {
		t.Error("no route expected")
	}
}

// go test -v -test.run TestCurly_ISSUE_137 ...restful
func TestCurly_ISSUE_137(t *testing.T) {
	ws1 := new(WebService)
	ws1.Route(ws1.GET("/hello").To(curlyDummy))
	ws1.Path("/")
	req, _ := http.NewRequest("GET", "/", nil)
	_, route, _, _ := CurlyRouter{}.SelectRoute([]*WebService{ws1}, req)
	t.Log(route)
	if route != nil {
		t.Error("no route expected")
	}
}

// go test -v -test.run TestCurly_ISSUE_137_2 ...restful
func TestCurly_ISSUE_137_2(t *testing.T) {
	ws1 := new(WebService)
	ws1.Route(ws1.GET("/hello").To(curlyDummy))
	ws1.Path("/")
	req, _ := http.NewRequest("GET", "/hello/bob", nil)
	_, route, _, _ := CurlyRouter{}.SelectRoute([]*WebService{ws1}, req)
	t.Log(route)
	if route != nil {
		t.Errorf("no route expected, got %v", route)
	}
}

func curlyDummy(req *Request, resp *Response) { io.WriteString(resp.ResponseWriter, "curlyDummy") }
