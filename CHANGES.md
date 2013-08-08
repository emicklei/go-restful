Change history of go-restful
=

2013-08-08
 - (api add) Added implementation Container: a WebServices collection with its own http.ServeMux allowing multiple endpoints per program. Existing uses of go-restful will register their services to the DefaultContainer.
 - (api add) the swagger package has be extended to have a UI per container.
 - if panic is detected then a small stack trace is printed (thanks to runner-mei)
 - (api add) WriteErrorString to Response

Important API changes:

 - (api remove) package variable DoNotRecover no longer works ; use restful.DefaultContainer.DoNotRecover(true) instead.
 - (api remove) package variable EnableContentEncoding no longer works ; use restful.DefaultContainer.EnableContentEncoding(true) instead.
 
 
2013-07-06

 - (api add) Added support for response encoding (gzip and deflate(zlib)). This feature is disabled on default (for backwards compatibility). Use restful.EnableContentEncoding = true in your initialization to enable this feature.

2013-06-19

 - (improve) DoNotRecover option, moved request body closer, improved ReadEntity

2013-06-03

 - (api change) removed Dispatcher interface, hide PathExpression
 - changed receiver names of type functions to be more idiomatic Go

2013-06-02

 - (optimize) Cache the RegExp compilation of Paths.

2013-05-22
	
 - (api add) Added support for request/response filter functions

2013-05-18


 - (api add) Added feature to change the default Http Request Dispatch function (travis cline)
 - (api change) Moved Swagger Webservice to swagger package (see example restful-user)

[2012-11-14 .. 2013-05-18>
 
 - See https://github.com/emicklei/go-restful/commits

2012-11-14

 - Initial commit


