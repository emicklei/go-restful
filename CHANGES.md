Change history of go-restful
=
2014-03-12
- (api add) Route path parameters can use wildcard or regular expressions. (requires CurlyRouter)

2014-02-26
- (api add) Request now provides information about the matched Route, see method SelectedRoutePath 

2014-02-17
- (api change) renamed parameter constants (go-lint checks)

2014-01-10
 - (api add) support for CloseNotify, see http://golang.org/pkg/net/http/#CloseNotifier

2014-01-07
 - (api change) Write* methods in Response now return the error or nil.
 - added example of serving HTML from a Go template.
 - fixed comparing Allowed headers in CORS (is now case-insensitive)

2013-11-13
 - (api add) Response knows how many bytes are written to the response body.

2013-10-29
 - (api add) RecoverHandler(handler RecoverHandleFunction) to change how panic recovery is handled. Default behavior is to log and return a stacktrace. This may be a security issue as it exposes sourcecode information.

2013-10-04
 - (api add) Response knows what HTTP status has been written
 - (api add) Request can have attributes (map of string->interface, also called request-scoped variables

2013-09-12
 - (api change) Router interface simplified
 - Implemented CurlyRouter, a Router that does not use|allow regular expressions in paths

2013-08-05
 - add OPTIONS support
 - add CORS support

2013-08-27
 - fixed some reported issues (see github)
 - (api change) deprecated use of WriteError; use WriteErrorString instead

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


