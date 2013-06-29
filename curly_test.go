package restful

import (
	"testing"
)

var request_paths = []struct {
	// url with path (1) is handled by service with root (2) and remainder has value final (3)
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

func _TestDetectWebService(t *testing.T) {
	ws1 := new(WebService).Path("/")
	ws2 := new(WebService).Path("/p")
	ws3 := new(WebService).Path("/q")
	ws4 := new(WebService).Path("/p/q")
	ws5 := new(WebService).Path("/p/{q}")
	ws6 := new(WebService).Path("/p/{q}/")
	ws7 := new(WebService).Path("/{p}/q")
	var wss = []*WebService{ws1, ws2, ws3, ws4, ws5, ws6, ws7}

	router := CurlyRouter{}

	ok := true
	for i, fixture := range request_paths {
		who, final, err := router.detectWebService(fixture.path, wss)
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
