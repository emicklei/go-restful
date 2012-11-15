package restful

type Route struct {
	Method   string
	Produces string
	Consumes string
	Path     string
	Function func(*Request, *Response)
}
