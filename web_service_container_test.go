package restful

import ()

var pathtests = []struct {
	root, path     string
	matchingTokens int
}{
	{"", "", 1},
	{"/", "/", 1},
	{"/p", "/p/", 2},
	{"/p/", "/p/", 1},
	{"/p/", "/p/q", 2},
	{"/{p}", "/p/q", 2},
	{"/p/{q}", "/p/q", 3},
	{"/p", "/a", 1},
}
