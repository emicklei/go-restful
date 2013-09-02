go-restful
==========

`go-restful` is a package for building REST-style Web Services using Google Go.

REST asks developers to use HTTP methods explicitly and in a way that's consistent with the protocol definition. This basic REST design principle establishes a one-to-one mapping between create, read, update, and delete (CRUD) operations and HTTP methods. According to this mapping:

- GET = Retrieve
- POST = Create if you are sending a command to the server to create a subordinate of the specified resource, using some server-side algorithm.
- POST = Update if you are requesting the server to update one or more subordinates of the specified resource.
- PUT = Create iff you are sending the full content of the specified resource (URL).
- PUT = Update iff you are updating the full content of the specified resource.
- DELETE = Delete if you are requesting the server to delete the resource
- PATCH = Update partial content of a resource
    

## Documentation

See Godoc for [automatically generated API documentation](http://godoc.org/github.com/emicklei/go-restful).


## Status

[![Build Status](https://drone.io/github.com/emicklei/go-restful/status.png)](https://drone.io/github.com/emicklei/go-restful/latest)
[![Coverage Status](https://coveralls.io/repos/emicklei/go-restful/badge.png?branch=master)](https://coveralls.io/r/emicklei/go-restful?branch=master)


## Resources

- [Hello world, plain and simple](https://github.com/emicklei/go-restful/tree/master/examples/restful-hello-world.go)  
- [Full API of a UserResource](https://github.com/emicklei/go-restful/tree/master/examples/restful-user-resource.go) 
- [Example posted on blog](http://ernestmicklei.com/2012/11/24/go-restful-first-working-example/)
- [Design explained on blog](http://ernestmicklei.com/2012/11/11/go-restful-api-design/)
- [Showcase: Mora - MongoDB REST Api server](https://github.com/emicklei/mora)
- [Showcase: Landskape tool](https://github.com/emicklei/landskape)


(c) 2012+, http://ernestmicklei.com. MIT License
