package restful

// Copyright 2013 Ernest Micklei. All rights reserved.
// Use of this source code is governed by a license
// that can be found in the LICENSE file.

import "strings"

// CrossOriginResourceSharing is used to create a Container Filter that implements CORS.
// Cross-origin resource sharing (CORS) is a mechanism that allows JavaScript on a web page
// to make XMLHttpRequests to another domain, not the domain the JavaScript originated from.
//
// http://en.wikipedia.org/wiki/Cross-origin_resource_sharing
// http://enable-cors.org/server.html
// http://www.html5rocks.com/en/tutorials/cors/#toc-handling-a-not-so-simple-request
type CrossOriginResourceSharing struct {
	ExposeHeaders  []string // list of Header names
	AllowedHeaders []string // list of Header names
	AllowedDomains []string // list of allowed values for Http Origin. If empty all are allowed.
	CookiesAllowed bool
	Container      *Container
}

// Filter is a filter function that implements the CORS flow as documented on http://enable-cors.org/server.html
// and http://www.html5rocks.com/static/images/cors_server_flowchart.png
func (c CrossOriginResourceSharing) Filter(req *Request, resp *Response, chain *FilterChain) {
	origin := req.Request.Header.Get(HEADER_Origin)
	if len(origin) == 0 {
		if trace {
			traceLogger.Println("no Http header Origin set")
		}
		chain.ProcessFilter(req, resp)
		return
	}
	if len(c.AllowedDomains) > 0 { // if provided then origin must be included
		included := false
		for _, each := range c.AllowedDomains {
			if each == origin {
				included = true
				break
			}
		}
		if !included {
			if trace {
				traceLogger.Println("HTTP Origin:%s is not part of %v", origin, c.AllowedDomains)
			}
			chain.ProcessFilter(req, resp)
			return
		}
	}
	if req.Request.Method != "OPTIONS" {
		c.doActualRequest(req, resp, chain)
		return
	}
	if acrm := req.Request.Header.Get(HEADER_AccessControlRequestMethod); acrm != "" {
		c.doPreflightRequest(req, resp, chain)
	} else {
		c.doActualRequest(req, resp, chain)
	}
}

func (c CrossOriginResourceSharing) doActualRequest(req *Request, resp *Response, chain *FilterChain) {
	resp.AddHeader(HEADER_AccessControlExposeHeaders, strings.Join(c.AllowedHeaders, ","))
	c.checkAndSetExposeHeaders(resp)
	c.setAllowOriginHeader(req, resp)
	c.checkAndSetAllowCredentials(resp)
	// continue processing the response
	chain.ProcessFilter(req, resp)
}

func (c CrossOriginResourceSharing) doPreflightRequest(req *Request, resp *Response, chain *FilterChain) {
	allowedMethods := c.Container.computeAllowedMethods(req)
	acrm := req.Request.Header.Get(HEADER_AccessControlRequestMethod)
	if !c.isValidAccessControlRequestMethod(acrm, allowedMethods) {
		if trace {
			traceLogger.Printf("Http header %s:%s is not in %v",
				HEADER_AccessControlRequestMethod,
				acrm,
				allowedMethods)
		}
		chain.ProcessFilter(req, resp)
		return
	}
	acrhs := req.Request.Header.Get(HEADER_AccessControlRequestHeaders)
	if len(acrhs) > 0 {
		for _, each := range strings.Split(acrhs, ",") {
			if !c.isValidAccessControlRequestHeader(strings.Trim(each, " ")) {
				if trace {
					traceLogger.Printf("Http header %s:%s is not in %v",
						HEADER_AccessControlRequestHeaders,
						acrhs,
						c.AllowedHeaders)
				}
				chain.ProcessFilter(req, resp)
				return
			}
		}
	}
	resp.AddHeader(HEADER_AccessControlAllowMethods, strings.Join(allowedMethods, ","))
	resp.AddHeader(HEADER_AccessControlAllowHeaders, acrhs)
	c.setAllowOriginHeader(req, resp)
	c.checkAndSetAllowCredentials(resp)
	// return http 200 response, no body
}

func (c CrossOriginResourceSharing) isOriginAllowed(origin string) bool {
	if len(origin) == 0 {
		return false
	}
	if len(c.AllowedDomains) == 0 {
		return true
	}
	allowed := false
	for _, each := range c.AllowedDomains {
		if each == origin {
			allowed = true
			break
		}
	}
	return allowed
}

func (c CrossOriginResourceSharing) setAllowOriginHeader(req *Request, resp *Response) {
	origin := req.Request.Header.Get(HEADER_Origin)
	if c.isOriginAllowed(origin) {
		resp.AddHeader(HEADER_AccessControlAllowOrigin, origin)
	}
}

func (c CrossOriginResourceSharing) checkAndSetExposeHeaders(resp *Response) {
	if len(c.ExposeHeaders) > 0 {
		resp.AddHeader(HEADER_AccessControlExposeHeaders, strings.Join(c.ExposeHeaders, ","))
	}
}

func (c CrossOriginResourceSharing) checkAndSetAllowCredentials(resp *Response) {
	if c.CookiesAllowed {
		resp.AddHeader(HEADER_AccessControlAllowCredentials, "true")
	}
}

func (c CrossOriginResourceSharing) isValidAccessControlRequestMethod(method string, allowedMethods []string) bool {
	for _, each := range allowedMethods {
		if each == method {
			return true
		}
	}
	return false
}

func (c CrossOriginResourceSharing) isValidAccessControlRequestHeader(header string) bool {
	for _, each := range c.AllowedHeaders {
		if strings.ToLower(each) == strings.ToLower(header) {
			return true
		}
	}
	return false
}
