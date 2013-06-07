go-restful
==========

`go-restful` is a package for building REST-style Web Services using Google Go.

REST asks developers to use HTTP methods explicitly and in a way that's consistent with the protocol definition. This basic REST design principle establishes a one-to-one mapping between create, read, update, and delete (CRUD) operations and HTTP methods. According to this mapping:

- Create = PUT iff you are sending the full content of the specified resource (URL).
- Create = POST if you are sending a command to the server to create a subordinate of the specified resource, using some server-side algorithm.
- Retrieve = GET.
- Update = PUT iff you are updating the full content of the specified resource.
- Update = POST if you are requesting the server to update one or more subordinates of the specified resource.
- Delete = DELETE if you are requesting the server to delete the resource
    

## Documentation

See Godoc for [automatically generated API documentation](http://godoc.org/github.com/emicklei/go-restful).


## Status

[![Build Status](https://drone.io/github.com/emicklei/go-restful/status.png)](https://drone.io/github.com/emicklei/go-restful/latest)
[![Coverage Status](https://coveralls.io/repos/jmcvetta/go-restful/badge.png?branch=master)](https://coveralls.io/r/jmcvetta/go-restful?branch=master)


## Resources

- [Hello world, plain and simple](https://github.com/emicklei/go-restful/tree/master/examples/restful-hello-world.go)  
- [Full API of a UserService](https://github.com/emicklei/go-restful/tree/master/examples/restful-user-service.go) 
- [Example posted on blog](http://ernestmicklei.com/2012/11/24/go-restful-first-working-example/)
- [Design explained on blog](http://ernestmicklei.com/2012/11/11/go-restful-api-design/)
- [Showcase: Landskape tool](https://github.com/emicklei/landskape)


(c) 2012+, http://ernestmicklei.com. MIT License
