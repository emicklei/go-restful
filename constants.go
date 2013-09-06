package restful

import "bytes"

const (
	MIME_XML  = "application/xml"
	MIME_JSON = "application/json"

	HEADER_Allow           = "Allow"
	HEADER_Accept          = "Accept"
	HEADER_ContentType     = "Content-Type"
	HEADER_LastModified    = "Last-Modified"
	HEADER_AcceptEncoding  = "Accept-Encoding"
	HEADER_ContentEncoding = "Content-Encoding"

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
