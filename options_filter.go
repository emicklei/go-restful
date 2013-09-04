package restful

// OPTIONSFilter is a filter function that inspects the Http Request for the OPTIONS method
// and provides the response with a set of allowed methods for the request URL Path.
// To install this filter on the Default Container use:
//
// 		restful.Filter(restful.OPTIONSFilter)
func OPTIONSFilter(req *Request, resp *Response, chain *FilterChain) {
	if "OPTIONS" != req.Request.Method {
		chain.ProcessFilter(req, resp)
		return
	}
	resp.AddHeader("Allow", computeAllowedMethods(req))
}
