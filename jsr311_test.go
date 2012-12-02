package restful

import (
	"testing"
)

func TestSelectRoutesSlash(t *testing.T) {
	ws1 := WebService{rootPath: "/"}
	ws1.Route(ws1.GET(""))
	ws1.Route(ws1.GET("/"))
	ws1.Route(ws1.GET("/u"))
	ws1.Route(ws1.POST("/u"))
	ws1.Route(ws1.POST("/u/v"))
	ws1.Route(ws1.POST("/u/{w}"))
	ws1.Route(ws1.POST("/u/{w}/z"))
	routes := selectRoutes(ws1, "/u")
	if len(routes) != 2 {
		t.Fatalf("expected:2, actual:%v", len(routes))
	}
}
func TestSelectRoutesU(t *testing.T) {
	ws1 := WebService{rootPath: "/u"}
	ws1.Route(ws1.GET(""))
	ws1.Route(ws1.GET("/"))
	ws1.Route(ws1.GET("/v"))
	ws1.Route(ws1.POST("/{w}"))
	ws1.Route(ws1.POST("/{w}/z"))
	routes := selectRoutes(ws1, "/v")
	if len(routes) != 1 {
		t.Fatalf("expected:1, actual:%v", len(routes))
	}
}

var tempregexs = []struct {
	template, regex        string
	literalCount, varCount int
}{
	{"", "(/.*)?", 0, 0},
	{"/a/{b}/c/", "/a/([^/]+?)/c(/.*)?", 2, 1},
	{"/{a}/{b}/{c-d-e}/", "/([^/]+?)/([^/]+?)/([^/]+?)(/.*)?", 0, 3},
}

func TestTemplateToRegularExpression(t *testing.T) {
	ok := true
	for i, fixture := range tempregexs {
		actual, lCount, vCount := templateToRegularExpression(fixture.template)
		if actual != fixture.regex {
			t.Logf("regex mismatch, expected:%v , actual:%v, line:%v\n", fixture.regex, actual, i+11) // 11 = where the data starts
			ok = false
		}
		if lCount != fixture.literalCount {
			t.Logf("literal count mismatch, expected:%v , actual:%v, line:%v\n", fixture.literalCount, lCount, i+11)
			ok = false
		}
		if vCount != fixture.varCount {
			t.Logf("variable count mismatch, expected:%v , actual:%v, line:%v\n", fixture.varCount, vCount, i+11)
			ok = false
		}
	}
	if !ok {
		t.Fatal("one or more expression did not match")
	}
}

var paths = []struct {
	// url with path is handled by service with root and last capturing group has value final
	path, root, final string
}{
	{"/", "/", "/"},
	{"/p", "/p", ""},
	{"/p/x", "/p/{q}", ""},
	{"/q/x", "/q", "/x"},
	{"/p/x/", "/p/{q}", "/"},
	{"/p/x/y", "/p/{q}", "/y"},
	{"/q/x/y", "/q", "/x/y"},
}

func TestDetectDispatcher(t *testing.T) {
	ws1 := WebService{rootPath: "/"}
	ws2 := WebService{rootPath: "/p"}
	ws3 := WebService{rootPath: "/q"}
	ws4 := WebService{rootPath: "/p/q"}
	ws5 := WebService{rootPath: "/p/{q}"}
	ws6 := WebService{rootPath: "/p/{q}/"}
	var dispatchers = []Dispatcher{ws1, ws2, ws3, ws4, ws5, ws6}

	ok := true
	for i, fixture := range paths {
		who, final, err := detectDispatcher(fixture.path, dispatchers)
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
