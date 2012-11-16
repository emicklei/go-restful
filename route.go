package restful

// Signature of function that can be bound to a Route
type RouteFunction func(*Request, *Response)

// Route binds a HTTP Method,Path,Consumes combination to a RouteFunction
type Route struct {
	Method   string
	Produces string // TODO make this a slice
	Consumes string // TODO make this a slice
	Path     string
	Function RouteFunction
}
