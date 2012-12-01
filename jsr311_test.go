package restful

import (
	"testing"
)

var tempregexs = []struct {
	template, regex string
}{
	{"", "(/.*)?"},
	{"/a/{b}/c/", "/a/([^/]+?)/c(/.*)?"},
}

func TestTemplateToRegularExpression(t *testing.T) {
	ok := true
	for i, fixture := range tempregexs {
		actual := templateToRegularExpression(fixture.template)
		if actual != fixture.regex {
			t.Logf("regex mismatch, expected:%v , actual:%v, line:%v\n", fixture.regex, actual, i+39)
			ok = false
		}
	}
	if !ok {
		t.Fatal("one or more expression did not match")
	}
}
