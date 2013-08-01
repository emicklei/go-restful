package restful

import (
	"testing"
)

//
// Step 1 tests
//
var paths = []struct {
	// url with path (1) is handled by service with root (2) and last capturing group has value final (3)
	path, root, final string
}{
	{"/", "/", "/"},
	{"/p", "/p", ""},
	{"/p/x", "/p/{q}", ""},
	{"/q/x", "/q", "/x"},
	{"/p/x/", "/p/{q}", "/"},
	{"/p/x/y", "/p/{q}", "/y"},
	{"/q/x/y", "/q", "/x/y"},
	{"/z/q", "/{p}/q", ""},
	{"/a/b/c/q", "/", "/a/b/c/q"},
}

func TestDetectDispatcher(t *testing.T) {
	ws1 := new(WebService).Path("/")
	ws2 := new(WebService).Path("/p")
	ws3 := new(WebService).Path("/q")
	ws4 := new(WebService).Path("/p/q")
	ws5 := new(WebService).Path("/p/{q}")
	ws6 := new(WebService).Path("/p/{q}/")
	ws7 := new(WebService).Path("/{p}/q")
	var dispatchers = []*WebService{ws1, ws2, ws3, ws4, ws5, ws6, ws7}

	router := RouterJSR311{}

	ok := true
	for i, fixture := range paths {
		who, final, err := router.detectDispatcher(fixture.path, dispatchers)
		if err != nil {
			t.Logf("error in detection:%v", err)
			ok = false
		}
		if who.RootPath() != fixture.root {
			t.Logf("[line:%v] Unexpected dispatcher, expected:%v, actual:%v", i, fixture.root, who.RootPath())
			ok = false
		}
		if final != fixture.final {
			t.Logf("[line:%v] Unexpected final, expected:%v, actual:%v", i, fixture.final, final)
			ok = false
		}
	}
	if !ok {
		t.Fail()
	}
}

//
// Step 2 tests
//

// go test -v -test.run TestISSUE_30 ...restful
func TestISSUE_30(t *testing.T) {
	ws1 := new(WebService).Path("/users")
	ws1.Route(ws1.GET("/{id}"))
	ws1.Route(ws1.POST("/login"))
	routes := RouterJSR311{}.selectRoutes(ws1, "/login")
	if len(routes) != 2 {
		t.Fatal("expected 2 routes")
	}
	//t.Logf("routes:%v", routes)
}

// go test -v -test.run TestISSUE_34 ...restful
func TestISSUE_34(t *testing.T) {
	ws1 := new(WebService).Path("/")
	ws1.Route(ws1.GET("/{type}/{id}"))
	ws1.Route(ws1.GET("/network/{id}"))
	routes := RouterJSR311{}.selectRoutes(ws1, "/network/12")
	if len(routes) != 2 {
		t.Fatal("expected 2 routes")
	}
	if routes[0].Path != "/network/{id}" {
		t.Error("first is", routes[0].Path)
	}
}

func TestSelectRoutesSlash(t *testing.T) {
	ws1 := new(WebService).Path("/")
	ws1.Route(ws1.GET(""))
	ws1.Route(ws1.GET("/"))
	ws1.Route(ws1.GET("/u"))
	ws1.Route(ws1.POST("/u"))
	ws1.Route(ws1.POST("/u/v"))
	ws1.Route(ws1.POST("/u/{w}"))
	ws1.Route(ws1.POST("/u/{w}/z"))
	routes := RouterJSR311{}.selectRoutes(ws1, "/u")
	checkRoutesContains(routes, "/u", t)
}
func TestSelectRoutesU(t *testing.T) {
	ws1 := new(WebService).Path("/u")
	ws1.Route(ws1.GET(""))
	ws1.Route(ws1.GET("/"))
	ws1.Route(ws1.GET("/v"))
	ws1.Route(ws1.POST("/{w}"))
	ws1.Route(ws1.POST("/{w}/z"))                    // so full path = /u/{w}/z
	routes := RouterJSR311{}.selectRoutes(ws1, "/v") // test against /u/v
	checkRoutesContains(routes, "/u/{w}", t)
}

func TestSelectRoutesUsers1(t *testing.T) {
	ws1 := new(WebService).Path("/users")
	ws1.Route(ws1.POST(""))
	ws1.Route(ws1.POST("/"))
	ws1.Route(ws1.PUT("/{id}"))
	routes := RouterJSR311{}.selectRoutes(ws1, "/1")
	checkRoutesContains(routes, "/users/{id}", t)
}
func checkRoutesContains(routes []Route, path string, t *testing.T) {
	if !containsRoutePath(routes, path, t) {
		for _, r := range routes {
			t.Logf("route %v %v", r.Method, r.Path)
		}
		t.Fatalf("routes should include [%v]:", path)
	}
}
func containsRoutePath(routes []Route, path string, t *testing.T) bool {
	for _, each := range routes {
		if each.Path == path {
			return true
		}
	}
	return false
}

var tempregexs = []struct {
	template, regex        string
	literalCount, varCount int
}{
	{"", "^(/.*)?$", 0, 0},
	{"/a/{b}/c/", "^/a/([^/]+?)/c(/.*)?$", 2, 1},
	{"/{a}/{b}/{c-d-e}/", "^/([^/]+?)/([^/]+?)/([^/]+?)(/.*)?$", 0, 3},
	{"/{p}/q", "^/([^/]+?)/q(/.*)?$", 1, 1},
}

func TestTemplateToRegularExpression(t *testing.T) {
	ok := true
	for i, fixture := range tempregexs {
		actual, lCount, vCount := templateToRegularExpression(fixture.template)
		if actual != fixture.regex {
			t.Logf("regex mismatch, expected:%v , actual:%v, line:%v\n", fixture.regex, actual, i) // 11 = where the data starts
			ok = false
		}
		if lCount != fixture.literalCount {
			t.Logf("literal count mismatch, expected:%v , actual:%v, line:%v\n", fixture.literalCount, lCount, i)
			ok = false
		}
		if vCount != fixture.varCount {
			t.Logf("variable count mismatch, expected:%v , actual:%v, line:%v\n", fixture.varCount, vCount, i)
			ok = false
		}
	}
	if !ok {
		t.Fatal("one or more expression did not match")
	}
}
