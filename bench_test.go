package restful

import (
	"fmt"
	"io"
	"testing"
)

var (
	wsCount = 26
	rtCount = 26
	uris    = []string{}
)

func setup() {
	for i := 0; i < wsCount; i++ {
		root := fmt.Sprintf("/%s/{%s}/", string(i+97), string(i+97))
		ws := new(WebService).Path(root)
		for j := 0; j < rtCount; j++ {
			sub := fmt.Sprintf("/%s2/{%s2}", string(j+97), string(j+97))
			ws.Route(ws.GET(sub).To(echo))
		}
		Add(ws)
		for _, each := range ws.Routes() {
			uris = append(uris, "http://bench.com"+each.Path)
		}
	}
}

func echo(req *Request, resp *Response) {
	io.WriteString(resp.ResponseWriter, "echo")
}

func BenchmarkMany(b *testing.B) {
	setup()
	b.ResetTimer()
	for t := 0; t < b.N; t++ {
		for _, each := range uris {
			// println(each)
			sendIt(each)
		}
	}
}
