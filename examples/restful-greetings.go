package main 

import (
	"github.com/emicklei/go-restful"
	"net/http"
)

type GreetingService struct {
	restful.WebService
	answer string
}

func main() {
	ws := new(GreetingService)
	ws.answer = "world"
	ws.Route(ws.GET("/hello").To(func(req *restful.Request, resp *restful.Response) { ws.hello(req,resp) ; return } ))
	restful.Add(ws)
	http.ListenAndServe(":8080", nil)
}

func (self *GreetingService) hello(req *restful.Request, resp *restful.Response) {
	resp.Write([]byte(self.answer))
}
