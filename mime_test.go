package restful

import (
	"fmt"
	"testing"
)

// go test -v -test.run TestSortMimes ...restful
func TestSortMimes(t *testing.T) {
	accept := "text/html; q=0.8, text/plain, image/gif,  */*; q=0.01, image/jpeg"
	result := sortedMimes(accept)
	got := fmt.Sprintf("%v", result)
	want := "[{text/plain q=1.0} {image/gif q=1.0} {image/jpeg q=1.0} {text/html  q=0.8} {*/*  q=0.01}]"
	if got != want {
		t.Errorf("bad sort order of mime types:%s", got)
	}
}
