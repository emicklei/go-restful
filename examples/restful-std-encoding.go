package main

import (
	"log"
	"net/http"

	"github.com/emicklei/go-restful"
)

type User struct {
	ID     int
	Active bool
}

func main() {
	restful.Add(NewUserService())
	restful.DefaultContainer.EnableContentEncoding(true)
	log.Print("start listening on localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func NewUserService() *restful.WebService {
	ws := new(restful.WebService)
	ws.
		Path("/users").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)
	ws.Route(ws.GET("/").To(listUsers))
	return ws
}

// curl -vvv -H accept-encoding:gzip http://localhost:8080/users
// curl -vvv http://localhost:8080/users
func listUsers(request *restful.Request, response *restful.Response) {
	users := []User{
		User{
			ID:     1,
			Active: true,
		},
		User{
			ID:     2,
			Active: false,
		},
	}
	response.WriteEntity(users)
}

/**
~/go/src/github.com/emicklei/go-restful (v3 *=)
$ curl -vvv -H accept-encoding:gzip http://localhost:8080/users
*   Trying ::1...
* TCP_NODELAY set
* Connected to localhost (::1) port 8080 (#0)
> GET /users HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.54.0
> Accept: */ /*
> accept-encoding:gzip
>
< HTTP/1.1 200 OK
< Content-Encoding: gzip
< Content-Type: application/json
< Date: Wed, 21 Oct 2020 09:20:49 GMT
< Content-Length: 94
<
��1
1�������i��na��A,$��Ҙ&���=���nE࣍>����S�"���tp/m�ޯ�/���
* Connection #0 to host localhost left intact
���?��I
**/

/**
$ curl -vvv http://localhost:8080/users
*   Trying ::1...
* TCP_NODELAY set
* Connected to localhost (::1) port 8080 (#0)
> GET /users HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.54.0
> Accept: */ /*
>
< HTTP/1.1 200 OK
< Content-Type: application/json
< Date: Wed, 21 Oct 2020 09:23:13 GMT
< Content-Length: 73
<
[
 {
  "ID": 1,
  "Active": true
 },
 {
  "ID": 2,
  "Active": false
 }
* Connection #0 to host localhost left intact
]
**/
