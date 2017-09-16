package restful

import "strings"

// Copyright 2013 Ernest Micklei. All rights reserved.
// Use of this source code is governed by a license
// that can be found in the LICENSE file.

// OptionsFilter is used to create a Filter that implements proper headers for CORS OPTIONS
// requests. Cross-origin resource sharing (CORS) is a mechanism that allows JavaScript on
// a web page to make XMLHttpRequests to another domain, not the domain the JavaScript
// originated from.
//
// http://en.wikipedia.org/wiki/Cross-origin_resource_sharing
// http://enable-cors.org/server.html
// http://www.html5rocks.com/en/tutorials/cors/#toc-handling-a-not-so-simple-request
type OptionsFilter struct {
	Container *Container
}

// Filter is a filter function that inspects the Http Request for the OPTIONS method
// and provides the response with a set of allowed methods for the request URL Path.
// As for any filter, you can also install it for a particular WebService within a Container
func (o *OptionsFilter) Filter(req *Request, resp *Response, chain *FilterChain) {
	if "OPTIONS" != req.Request.Method {
		chain.ProcessFilter(req, resp)
		return
	}

	archs := req.Request.Header.Get(HEADER_AccessControlRequestHeaders)
	methods := strings.Join(o.Container.computeAllowedMethods(req), ",")
	origin := req.Request.Header.Get(HEADER_Origin)

	resp.AddHeader(HEADER_Allow, methods)
	resp.AddHeader(HEADER_AccessControlAllowOrigin, origin)
	resp.AddHeader(HEADER_AccessControlAllowHeaders, archs)
	resp.AddHeader(HEADER_AccessControlAllowMethods, methods)

}
