package main

import (
	"fmt"
	"github.com/emicklei/go-restful"
	"log"
	"net/http"
	"os"
	"strings"
)

const MIME_ENV = "application/env"

type ENVEntity struct{}

func (e *ENVEntity) New() restful.EntityEncoder {
	return e
}

func (e *ENVEntity) MIME() string {
	return MIME_ENV
}

func (e *ENVEntity) SetRequest(r *restful.Request) {}

func (e *ENVEntity) SetResponse(r *restful.Response) {}

func (e *ENVEntity) Marshal(v interface{}) ([]byte, error) {
	var r []byte
	m, ok := v.(map[string]string)
	if !ok {
		return []byte{}, fmt.Errorf("failed to parse map from interface")
	}
	for k, v := range m {
		r = append(r, []byte(strings.ToUpper(k)+"="+v+"\n")...)
	}
	return r, nil
}

func (e *ENVEntity) MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	return e.Marshal(v)
}

func (e *ENVEntity) Unmarshal(b []byte, entityPonter interface{}) error {
	return nil
}

func main() {
	restful.RegisterEntityEncoder(&ENVEntity{})
	restful.TraceLogger(log.New(os.Stdout, "[restful] ", log.LstdFlags|log.Lshortfile))
	ws := new(restful.WebService)
	ws.Route(ws.GET("/env").To(getEnv).Produces(MIME_ENV))
	restful.Add(ws)
	http.ListenAndServe(":8080", nil)
}

func getEnv(req *restful.Request, resp *restful.Response) {
	a := map[string]string{
		"mykey": "myvalue",
	}
	resp.WriteEntity(a)
}
