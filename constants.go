package restful

import "bytes"

const (
	MIME_XML  = "application/xml"
	MIME_JSON = "application/json"

	HEADER_Allow                         = "Allow"
	HEADER_Accept                        = "Accept"
	HEADER_Origin                        = "Origin"
	HEADER_ContentType                   = "Content-Type"
	HEADER_LastModified                  = "Last-Modified"
	HEADER_AcceptEncoding                = "Accept-Encoding"
	HEADER_ContentEncoding               = "Content-Encoding"
	HEADER_AccessControlExposeHeaders    = "Access-Control-Expose-Headers"
	HEADER_AccessControlRequestMethod    = "Access-Control-Request-Method"
	HEADER_AccessControlRequestHeaders   = "Access-Control-Request-Headers"
	HEADER_AccessControlAllowMethods     = "Access-Control-Allow-Methods"
	HEADER_AccessControlAllowOrigin      = "Access-Control-Allow-Origin"
	HEADER_AccessControlAllowCredentials = "Access-Control-Allow-Credentials"
	HEADER_AccessControlAllowHeaders     = "Access-Control-Allow-Headers"

	ENCODING_GZIP    = "gzip"
	ENCODING_DEFLATE = "deflate"
)

func toCommaSeparated(names []string) string {
	buf := new(bytes.Buffer)
	for _, each := range names {
		if buf.Len() > 0 {
			buf.WriteString(",")
		}
		buf.WriteString(each)
	}
	return buf.String()
}
