package restful

import "strings"

type mime struct {
	media   string
	quality string
}

func insertMime(l []mime, e mime) []mime {
	for i, each := range l {
		// if current mime has lower quality then insert before
		if e.quality > each.quality {
			left := append([]mime{}, l[0:i]...)
			return append(append(left, e), l[i:]...)
		}
	}
	return append(l, e)
}

func sortedMimes(accept string) (sorted []mime) {
	for _, each := range strings.Split(accept, ",") {
		q := strings.Split(strings.Trim(each, " "), ";")
		if len(q) == 1 {
			sorted = insertMime(sorted, mime{q[0], "q=1.0"})
		} else {
			sorted = insertMime(sorted, mime{q[0], q[1]})
		}
	}
	return
}
