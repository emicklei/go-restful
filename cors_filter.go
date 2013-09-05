package restful

import (
	"strings"
)

type CrossOriginResourceSharing struct {
	ExposeHeaders  string // comma separated list of Header names
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
	resp.AddHeader("Access-Control-Expose-Headers", "Content-Type, Accept")
	c.checkAndSetExposeHeaders(resp)
	c.setAllowOriginHeader(req, resp)
	c.checkAndSetAllowCredentials(resp)
	// continue processing the response
	chain.ProcessFilter(req, resp)
}

func (c CrossOriginResourceSharing) doPreflightRequest(req *Request, resp *Response, chain *FilterChain) {
	if !c.isValidAccessControlRequestMethod(req.Request.Method) {
		chain.ProcessFilter(req, resp)
		return
	}
	if acrh := req.Request.Header.Get("Access-Control-Request-Header"); acrh != "" {
		if !c.isValidAccessControlRequestHeader(acrh) {
			chain.ProcessFilter(req, resp)
			return
		}
	}
	resp.AddHeader("Access-Control-Allow-Methods", c.Container.computeAllowedMethods(req))
	resp.AddHeader("Access-Control-Allow-Headers", "Content-Type, Accept")
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
		resp.AddHeader("Access-Control-Expose-Headers", c.ExposeHeaders)
	}
}

func (c CrossOriginResourceSharing) checkAndSetAllowCredentials(resp *Response) {
	if c.CookiesAllowed {
		resp.AddHeader("Access-Control-Allow-Credentials", "true")
	}
}

func (c CrossOriginResourceSharing) isValidAccessControlRequestMethod(method string) bool {
	return strings.Contains("GET PUT POST DELETE HEAD OPTIONS PATCH", method) // I know there are more but my guess is that this is enough
}

func (c CrossOriginResourceSharing) isValidAccessControlRequestHeader(header string) bool {
	return strings.Contains("accept content-type", strings.ToLower(header))
}
