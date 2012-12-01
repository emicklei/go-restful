package restful

import (
	"testing"
)

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
			t.Logf("regex mismatch, expected:%v , actual:%v, line:%v\n", fixture.regex, actual, i+39)
			ok = false
		}
		if lCount != fixture.literalCount {
			t.Logf("literal count mismatch, expected:%v , actual:%v, line:%v\n", fixture.literalCount, lCount, i+39)
			ok = false
		}
		if vCount != fixture.varCount {
			t.Logf("variable count mismatch, expected:%v , actual:%v, line:%v\n", fixture.varCount, vCount, i+39)
			ok = false
		}
	}
	if !ok {
		t.Fatal("one or more expression did not match")
	}
}

var paths = []struct {
	// url with path is handled by service with root
	path, root string
}{
	{"/", "/"},
	{"/p/x", "/p/{q}"},
	{"/q/x", "/q"},
	{"/p/x/", "/p/{q}"},
	{"/p/x/y", "/p/{q}"},
	{"/q/x/y", "/q"},
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
	for _, fixture := range paths {
		who, err := detectDispatcher(fixture.path, dispatchers)
		if err != nil {
			t.Logf("error in detection:%v", err)
			ok = false
		}
		if who.RootPath() != fixture.root {
			t.Logf("Unexpected dispatcher, expected:%v, actual:%v", fixture.root, who.RootPath())
			ok = false
		}
	}
	if !ok {
		t.Fail()
	}
}
