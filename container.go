package restful

import (
	"net/http"
	"strings"
)

var DefaultContainer *Container

func init() {
	DefaultContainer = NewContainer()
	DefaultContainer.serveMux = http.DefaultServeMux
}

type Container struct {
	webServices        []*WebService
	serveMux           *http.ServeMux
	isRegisteredOnRoot bool
	containerFilters   []FilterFunction
	doNotRecover       bool
	router             RouteSelector
}

func NewContainer() *Container {
	return &Container{
		webServices:        []*WebService{},
		serveMux:           http.NewServeMux(),
		isRegisteredOnRoot: false,
		containerFilters:   []FilterFunction{},
		doNotRecover:       false,
		router:             RouterJSR311{}}
}

func (c *Container) Add(service *WebService) *Container {
	if service.pathExpr == nil {
		service.Path("") // lazy initialize path
	}
	// If registered on root then no additional specific mapping is needed
	if !c.isRegisteredOnRoot {
		pattern := c.fixedPrefixPath(service.RootPath())
		// check if root path registration is needed
		if "/" == pattern || "" == pattern {
			c.serveMux.HandleFunc("/", Dispatch)
			c.isRegisteredOnRoot = true
		} else {
			// detect if registration already exists
			alreadyMapped := false
			for _, each := range webServices {
				if each.RootPath() == service.RootPath() {
					alreadyMapped = true
					break
				}
			}
			if !alreadyMapped {
				c.serveMux.HandleFunc(pattern, Dispatch)
				if !strings.HasSuffix(pattern, "/") {
					c.serveMux.HandleFunc(pattern+"/", Dispatch)
				}
			}
		}
	}
	c.webServices = append(c.webServices, service)
	return c
}

// fixedPrefixPath returns the fixed part of the partspec ; it may include template vars {}
func (c Container) fixedPrefixPath(pathspec string) string {
	varBegin := strings.Index(pathspec, "{")
	if -1 == varBegin {
		return pathspec
	}
	return pathspec[:varBegin]
}

// implements net/http.Handler
func (c Container) ServeHTTP(httpwriter http.ResponseWriter, httpRequest *http.Request) {
	c.serveMux.ServeHTTP(httpwriter, httpRequest)
}

// Filter appends a container FilterFunction. These are called before dispatch a http.Request to a WebService.
func (c *Container) Filter(filter FilterFunction) {
	c.containerFilters = append(c.containerFilters, filter)
}
