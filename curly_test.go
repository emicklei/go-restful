package restful

import (
	"io"
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
}

// clear && go test -v -test.run Test_matchesRouteByPathTokens ...restful
func Test_matchesRouteByPathTokens(t *testing.T) {
	router := CurlyRouter{}
	for _, each := range routeMatchers {
		routeToks := tokenizePath(each.route)
		reqToks := tokenizePath(each.path)
		matches, pCount, sCount := router.matchesRouteByPathTokens(routeToks, reqToks)
		if matches != each.matches {
			t.Fatalf("unexpected matches outcome route:%s, path:%s, matches:%v", each.route, each, each.path, each.matches)
		}
		if pCount != each.paramCount {
			t.Fatalf("unexpected paramCount got:%d want:%d ", pCount, each.paramCount)
		}
		if sCount != each.staticCount {
			t.Fatalf("unexpected staticCount got:%d want:%d ", sCount, each.staticCount)
		}
	}
}

// clear && go test -v -test.run TestCurly_ISSUE_34 ...restful
func TestCurly_ISSUE_34(t *testing.T) {
	ws1 := new(WebService).Path("/")
	ws1.Route(ws1.GET("/{type}/{id}").To(curlyDummy))
	ws1.Route(ws1.GET("/network/{id}").To(curlyDummy))
	routes := CurlyRouter{}.selectRoutes(ws1, tokenizePath("/network/12"))
	if len(routes) != 2 {
		t.Fatal("expected 2 routes")
	}
	if routes[0].Path != "/network/{id}" {
		t.Error("first is", routes[0].Path)
	}
}

// clear && go test -v -test.run TestCurly_ISSUE_34_2 ...restful
func TestCurly_ISSUE_34_2(t *testing.T) {
	ws1 := new(WebService).Path("/")
	ws1.Route(ws1.GET("/network/{id}").To(curlyDummy))
	ws1.Route(ws1.GET("/{type}/{id}").To(curlyDummy))
	routes := CurlyRouter{}.selectRoutes(ws1, tokenizePath("/network/12"))
	if len(routes) != 2 {
		t.Fatal("expected 2 routes")
	}
	if routes[0].Path != "/network/{id}" {
		t.Error("first is", routes[0].Path)
	}
}

func curlyDummy(req *Request, resp *Response) { io.WriteString(resp.ResponseWriter, "curlyDummy") }
