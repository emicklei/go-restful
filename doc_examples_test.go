package restful

import "net/http"

func ExampleOptionsFilter() {
	// Install the OPTIONS filter on the default Container
	optionsFilter := OptionsFilter{Container: DefaultContainer}
	Filter(optionsFilter.Filter)
}
func ExampleContainer_OptionsFilter() {
	// Install the OPTIONS filter on a Container
	myContainer := new(Container)
	optionsFilter := OptionsFilter{Container: myContainer}
	myContainer.Filter(optionsFilter.Filter)
}

func ExampleContainer() {
	// The Default container of go-restful uses the http.DefaultServeMux.
	// You can create your own Container using restful.NewContainer() and create a new http.Server for that particular container

	ws := new(WebService)
	wsContainer := NewContainer()
	wsContainer.Add(ws)
	server := &http.Server{Addr: ":8080", Handler: wsContainer}
	server.ListenAndServe()
}

func ExampleCrossOriginResourceSharing() {
	// To install this filter on the Default Container use:
	cors := CrossOriginResourceSharing{ExposeHeaders: []string{"X-My-Header"}, CookiesAllowed: false, Container: DefaultContainer}
	Filter(cors.Filter)
}

func ExampleServiceError() {
	resp := new(Response)
	resp.WriteEntity(NewError(http.StatusBadRequest, "Non-integer {id} path parameter"))
}

func ExampleBoundedCachedCompressors() {
	// Register a compressor provider (gzip/deflate read/write) that uses
	// a bounded cache with a maximum of 20 writers and 20 readers.
	SetCompressorProvider(NewBoundedCachedCompressors(20, 20))
}
