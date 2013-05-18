go-restful
==========

package for building REST-style Web Services using Google Go

REST asks developers to use HTTP methods explicitly and in a way that's consistent with the protocol definition. This basic REST design principle establishes a one-to-one mapping between create, read, update, and delete (CRUD) operations and HTTP methods. According to this mapping:

- Create = PUT iff you are sending the full content of the specified resource (URL).
- Create = POST if you are sending a command to the server to create a subordinate of the specified resource, using some server-side algorithm.
- Retrieve = GET.
- Update = PUT iff you are updating the full content of the specified resource.
- Update = POST if you are requesting the server to update one or more subordinates of the specified resource.
- Delete = DELETE if you are requesting the server to delete the resource
    
### Resources

- [Documentation go-restful (godoc.org)](http://godoc.org/github.com/emicklei/go-restful)
- [Hello world, plain and simple](https://github.com/emicklei/go-restful/tree/master/examples/hello/restful-hello-world.go)  
- [Full API of a UserService](https://github.com/emicklei/go-restful/tree/master/examples/restful-user-service.go) 
- [Example posted on blog](http://ernestmicklei.com/2012/11/24/go-restful-first-working-example/)
- [Design explained on blog](http://ernestmicklei.com/2012/11/11/go-restful-api-design/)
- [Showcase: Landskape tool](https://github.com/emicklei/landskape)

[![Build Status](https://travis-ci.org/emicklei/go-restful.png)](https://travis-ci.org/emicklei/go-restful)

(c) 2012+, http://ernestmicklei.com. MIT License