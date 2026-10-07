package restful

import (
	"reflect"
	"sync"
	"testing"
)

var tempregexs = []struct {
	template, regex        string
	names                  []string
	literalCount, varCount int
}{
	{"", "^(/.*)?$", nil, 0, 0},
	{"/a/{b}/c/", "^/a/([^/]+?)/c(/.*)?$", []string{"b"}, 2, 1},
	{"/{a}/{b}/{c-d-e}/", "^/([^/]+?)/([^/]+?)/([^/]+?)(/.*)?$", []string{"a", "b", "c-d-e"}, 0, 3},
	{"/{p}/abcde", "^/([^/]+?)/abcde(/.*)?$", []string{"p"}, 5, 1},
	{"/a/{b:*}", "^/a/(.*)(/.*)?$", []string{"b"}, 1, 1},
	{"/a/{b:[a-z]+}", "^/a/([a-z]+)(/.*)?$", []string{"b"}, 1, 1},
}

func TestTemplateToRegularExpression(t *testing.T) {
	ok := true
	for i, fixture := range tempregexs {
		actual, lCount, varNames, vCount, _ := templateToRegularExpression(fixture.template)
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
		if !reflect.DeepEqual(fixture.names, varNames) {
			t.Logf("variable name mismatch, expected:%v , actual:%v, line:%v\n", fixture.names, varNames, i)
			ok = false
		}
	}
	if !ok {
		t.Fatal("one or more expression did not match")
	}
}

func TestNewPathExpression_CachesCompiledRegexp(t *testing.T) {
	const (
		tmpl1 = "/namespaces/{namespace}/{resource}/{name}"
		tmpl2 = "/namespaces/{ns}/{res}/{id}"
	)
	pe1, err := newPathExpression(tmpl1)
	if err != nil {
		t.Fatalf("newPathExpression(%q) returned unexpected error: %v", tmpl1, err)
	}
	pe2, err := newPathExpression(tmpl2)
	if err != nil {
		t.Fatalf("newPathExpression(%q) returned unexpected error: %v", tmpl2, err)
	}
	if pe1.Matcher != pe2.Matcher {
		t.Errorf("newPathExpression(%q) and newPathExpression(%q) Matcher = %p, %p; want identical *regexp.Regexp pointer for expression %q", tmpl1, tmpl2, pe1.Matcher, pe2.Matcher, pe1.Source)
	}
	if want := []string{"namespace", "resource", "name"}; !reflect.DeepEqual(pe1.VarNames, want) {
		t.Errorf("newPathExpression(%q).VarNames = %v, want %v", tmpl1, pe1.VarNames, want)
	}
	if want := []string{"ns", "res", "id"}; !reflect.DeepEqual(pe2.VarNames, want) {
		t.Errorf("newPathExpression(%q).VarNames = %v, want %v", tmpl2, pe2.VarNames, want)
	}
}

var benchmarkKubernetesPathTemplates = []string{
	"/{resource}",
	"/{resource}/{name}",
	"/{resource}/{name}/status",
	"/{resource}/{name}/scale",
	"/namespaces/{namespace}/{resource}",
	"/namespaces/{namespace}/{resource}/{name}",
	"/namespaces/{namespace}/{resource}/{name}/status",
	"/namespaces/{namespace}/{resource}/{name}/scale",
	"/watch/{resource}",
	"/watch/namespaces/{namespace}/{resource}",
	"/watch/namespaces/{namespace}/{resource}/{name}",
}

func BenchmarkNewPathExpression_KubernetesRoutes(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		regexpCache = sync.Map{}
		for gv := 0; gv < 60; gv++ {
			for _, tmpl := range benchmarkKubernetesPathTemplates {
				if _, err := newPathExpression(tmpl); err != nil {
					b.Fatal(err)
				}
			}
		}
	}
}
