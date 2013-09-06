package restful

// Cross-origin resource sharing (CORS) is a mechanism that allows JavaScript on a web page
// to make XMLHttpRequests to another domain, not the domain the JavaScript originated from.
//
// http://en.wikipedia.org/wiki/Cross-origin_resource_sharing
// http://enable-cors.org/server.html
// http://www.html5rocks.com/en/tutorials/cors/#toc-handling-a-not-so-simple-request
//
// CrossOriginResourceSharing is used to create a Container Filter that implements CORS
type CrossOriginResourceSharing struct {
	ExposeHeaders  []string // list of Header names
	AllowdHeaders  []string // list of Header names
	CookiesAllowed bool
	Container      *Container
}

// Filter is a filter function that implements the CORS flow as documented on http://enable-cors.org/server.html
// and http://www.html5rocks.com/static/images/cors_server_flowchart.png
// To install this filter on the Default Container use:
//
//		cors := restful.CrossOriginResourceSharing{ExposeHeaders:"X-My-Header", CookiesAllowed:false, restful.DefaultContainer}
// 		restful.Filter(cors.Filter)
func (c CrossOriginResourceSharing) Filter(req *Request, resp *Response, chain *FilterChain) {
	if origin := req.Request.Header.Get("Origin"); len(origin) == 0 {
		chain.ProcessFilter(req, resp)
		return
	}
	if req.Request.Method != "OPTIONS" {
		c.doActualRequest(req, resp, chain)
		return
	}
	if acrm := req.Request.Header.Get("Access-Control-Request-Method"); acrm != "" {
		c.doPreflightRequest(req, resp, chain)
	} else {
		c.doActualRequest(req, resp, chain)
	}

}

func (c CrossOriginResourceSharing) doActualRequest(req *Request, resp *Response, chain *FilterChain) {
	resp.AddHeader("Access-Control-Expose-Headers", toCommaSeparated(c.AllowdHeaders))
	c.checkAndSetExposeHeaders(resp)
	c.setAllowOriginHeader(req, resp)
	c.checkAndSetAllowCredentials(resp)
	// continue processing the response
	chain.ProcessFilter(req, resp)
}

func (c CrossOriginResourceSharing) doPreflightRequest(req *Request, resp *Response, chain *FilterChain) {
	allowedMethods := c.Container.computeAllowedMethods(req)
	if !c.isValidAccessControlRequestMethod(req.Request.Method, allowedMethods) {
		chain.ProcessFilter(req, resp)
		return
	}
	acrhs := req.Request.Header.Get("Access-Control-Request-Headers")
	if acrhs != "" {
		if !c.isValidAccessControlRequestHeader(acrhs) {
			chain.ProcessFilter(req, resp)
			return
		}
	}
	resp.AddHeader("Access-Control-Allow-Methods", toCommaSeparated(allowedMethods))
	resp.AddHeader("Access-Control-Allow-Headers", acrhs)
	c.setAllowOriginHeader(req, resp)
	c.checkAndSetAllowCredentials(resp)
	// return http 200 response, no body
}

func (c CrossOriginResourceSharing) setAllowOriginHeader(req *Request, resp *Response) {
	origin := req.Request.Header.Get("Origin")
	resp.AddHeader("Access-Control-Allow-Origin", origin)
}

func (c CrossOriginResourceSharing) checkAndSetExposeHeaders(resp *Response) {
	if len(c.ExposeHeaders) > 0 {
		resp.AddHeader("Access-Control-Expose-Headers", toCommaSeparated(c.ExposeHeaders))
	}
}

func (c CrossOriginResourceSharing) checkAndSetAllowCredentials(resp *Response) {
	if c.CookiesAllowed {
		resp.AddHeader("Access-Control-Allow-Credentials", "true")
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
	for _, each := range c.AllowdHeaders {
		if each == header {
			return true
		}
	}
	return false
}
