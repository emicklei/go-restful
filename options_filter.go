package restful

//import "log"
import "bytes"

func GetOPTIONSFilter(req *Request, resp *Response, chain *FilterChain) {
	if "OPTIONS" != req.Request.Method {
		chain.ProcessFilter(req, resp)
		return
	}
	// Go through all RegisteredWebServices() and all its Routes to collect the options
	methods := []string{}
	requestPath := req.Request.URL.Path
	for _, ws := range RegisteredWebServices() {
		matches := ws.pathExpr.Matcher.FindStringSubmatch(requestPath)
		if matches != nil {
			finalMatch := matches[len(matches)-1]
			for _, rt := range ws.Routes() {
				matches := rt.pathExpr.Matcher.FindStringSubmatch(finalMatch)
				//log.Printf("reg:%v, route:%s,matches:%v", rt.pathExpr.Source, rt.String(), matches)
				if matches != nil {
					lastMatch := matches[len(matches)-1]
					if lastMatch == "" || lastMatch == "/" { // do not include if value is neither empty nor ‘/’.
						//log.Printf("route:%s\n", rt.String())
						methods = append(methods, rt.Method)
					}
				}
			}
		}
	}
	// compose Allow
	buf := new(bytes.Buffer)
	for _, m := range methods {
		if buf.Len() > 0 {
			buf.WriteString(",")
		}
		buf.WriteString(m)
	}
	resp.AddHeader("Allow", buf.String())
}
