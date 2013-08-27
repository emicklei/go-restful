package restful

// This file implements the flow for matching Requests to Routes (and consequently Resource Functions)
// as specified by the JSR311 http://jsr311.java.net/nonav/releases/1.1/spec/spec.html.
// Concept of locators is not implemented.
import (
	"errors"
	"net/http"
	"sort"
	//"strconv"
	//"github.com/emicklei/hopwatch"
	"fmt"
)

type RouterJSR311 struct{}

func (r RouterJSR311) SelectRoute(
	path string,
	webServices []*WebService,
	httpWriter http.ResponseWriter,
	httpRequest *http.Request) (selectedService *WebService, selectedRoute Route, ok bool) {

	// Identify the root resource class (WebService)
	dispatcher, finalMatch, err := r.detectDispatcher(path, webServices)
	if err != nil {
		httpWriter.WriteHeader(http.StatusNotFound)
		return nil, Route{}, false
	}
	// Obtain the set of candidate methods (Routes)
	routes := r.selectRoutes(dispatcher, finalMatch)

	// Identify the method (Route) that will handle the request
	route, ok := r.detectRoute(routes, httpWriter, httpRequest)
	return dispatcher, route, ok
}

// http://jsr311.java.net/nonav/releases/1.1/spec/spec3.html#x3-360003.7.2
func (r RouterJSR311) detectRoute(routes []Route, httpWriter http.ResponseWriter, httpRequest *http.Request) (Route, bool) {
	// http method
	methodOk := []Route{}
	for _, each := range routes {
		if httpRequest.Method == each.Method {
			methodOk = append(methodOk, each)
		}
	}
	if len(methodOk) == 0 {
		httpWriter.WriteHeader(http.StatusMethodNotAllowed)
		httpWriter.Write([]byte("405: Method Not Allowed"))
		return Route{}, false
	}
	inputMediaOk := methodOk
	// content-type
	contentType := httpRequest.Header.Get(HEADER_ContentType)
	if httpRequest.ContentLength > 0 {
		inputMediaOk = []Route{}
		for _, each := range methodOk {
			if each.matchesContentType(contentType) {
				inputMediaOk = append(inputMediaOk, each)
			}
		}
		if len(inputMediaOk) == 0 {
			httpWriter.WriteHeader(http.StatusUnsupportedMediaType)
			httpWriter.Write([]byte("415: Unsupported Media Type"))
			return Route{}, false
		}
	}
	// accept
	outputMediaOk := []Route{}
	accept := httpRequest.Header.Get(HEADER_Accept)
	if accept == "" {
		accept = "*/*"
	}
	for _, each := range inputMediaOk {
		if each.matchesAccept(accept) {
			outputMediaOk = append(outputMediaOk, each)
		}
	}
	if len(outputMediaOk) == 0 {
		httpWriter.WriteHeader(http.StatusNotAcceptable)
		httpWriter.Write([]byte("406: Not Acceptable"))
		return Route{}, false
	}
	return r.bestMatchByMedia(outputMediaOk, contentType, accept), true
}

// http://jsr311.java.net/nonav/releases/1.1/spec/spec3.html#x3-360003.7.2
// n/m > n/* > */*
func (r RouterJSR311) bestMatchByMedia(routes []Route, contentType string, accept string) Route {
	// TODO
	return routes[0]
}

// http://jsr311.java.net/nonav/releases/1.1/spec/spec3.html#x3-360003.7.2  (step 2)
func (r RouterJSR311) selectRoutes(dispatcher *WebService, pathRemainder string) []Route {
	if pathRemainder == "" || pathRemainder == "/" {
		return dispatcher.Routes()
	}
	filtered := &sortableRouteCandidates{}
	for _, each := range dispatcher.Routes() {
		pathExpr := each.pathExpr
		matches := pathExpr.Matcher.FindStringSubmatch(pathRemainder)
		if matches != nil {
			lastMatch := matches[len(matches)-1]
			if lastMatch == "" || lastMatch == "/" { // do not include if value is neither empty nor ‘/’.
				filtered.candidates = append(filtered.candidates,
					routeCandidate{each, len(matches) - 1, pathExpr.LiteralCount, pathExpr.VarCount})
			}
		}
	}
	if len(filtered.candidates) == 0 {
		return []Route{}
	}
	sort.Sort(sort.Reverse(filtered))
	//hopwatch.Dump(filtered)

	// select other routes from candidates whoes expression matches rmatch
	matchingRoutes := []Route{filtered.candidates[0].route}
	for c := 1; c < len(filtered.candidates); c++ {
		each := filtered.candidates[c]
		if each.route.pathExpr.Matcher.MatchString(pathRemainder) {
			matchingRoutes = append(matchingRoutes, each.route)
		}
	}
	return matchingRoutes
}

// http://jsr311.java.net/nonav/releases/1.1/spec/spec3.html#x3-360003.7.2
func (r RouterJSR311) detectDispatcher(requestPath string, dispatchers []*WebService) (*WebService, string, error) {
	filtered := &sortableDispatcherCandidates{}
	for _, each := range dispatchers {
		pathExpr := each.pathExpr
		matches := pathExpr.Matcher.FindStringSubmatch(requestPath)
		if matches != nil {
			filtered.candidates = append(filtered.candidates,
				dispatcherCandidate{each, matches[len(matches)-1], len(matches), pathExpr.LiteralCount, pathExpr.VarCount})
		}
	}
	if len(filtered.candidates) == 0 {
		return nil, "", errors.New("not found")
	}
	sort.Sort(sort.Reverse(filtered))
	return filtered.candidates[0].dispatcher, filtered.candidates[0].finalMatch, nil
}

// Types and functions to support the sorting of Routes

type routeCandidate struct {
	route           Route
	matchesCount    int // the number of capturing groups
	literalCount    int // the number of literal characters (means those not resulting from template variable substitution)
	nonDefaultCount int // the number of capturing groups with non-default regular expressions (i.e. not ‘([^  /]+?)’)
}

func (r routeCandidate) expressionToMatch() string {
	return r.route.pathExpr.Source
}

func (r routeCandidate) String() string {
	return fmt.Sprintf("(m=%d,l=%d,n=%d)", r.matchesCount, r.literalCount, r.nonDefaultCount)
}

type sortableRouteCandidates struct {
	candidates []routeCandidate
}

func (self *sortableRouteCandidates) Len() int {
	return len(self.candidates)
}
func (self *sortableRouteCandidates) Swap(i, j int) {
	self.candidates[i], self.candidates[j] = self.candidates[j], self.candidates[i]
}
func (self *sortableRouteCandidates) Less(i, j int) bool {
	ci := self.candidates[i]
	cj := self.candidates[j]
	// primary key
	if ci.literalCount < cj.literalCount {
		return true
	}
	if ci.literalCount > cj.literalCount {
		return false
	}
	// secundary key
	if ci.matchesCount < cj.matchesCount {
		return true
	}
	if ci.matchesCount > cj.matchesCount {
		return false
	}
	// tertiary key
	if ci.nonDefaultCount < cj.nonDefaultCount {
		return true
	}
	if ci.nonDefaultCount > cj.nonDefaultCount {
		return false
	}
	// quaternary key ("source" is interpreted as Path)
	return ci.route.Path < cj.route.Path
}

// Types and functions to support the sorting of Dispatchers

type dispatcherCandidate struct {
	dispatcher      *WebService
	finalMatch      string
	matchesCount    int // the number of capturing groups
	literalCount    int // the number of literal characters (means those not resulting from template variable substitution)
	nonDefaultCount int // the number of capturing groups with non-default regular expressions (i.e. not ‘([^  /]+?)’)
}
type sortableDispatcherCandidates struct {
	candidates []dispatcherCandidate
}

func (self *sortableDispatcherCandidates) Len() int {
	return len(self.candidates)
}
func (self *sortableDispatcherCandidates) Swap(i, j int) {
	self.candidates[i], self.candidates[j] = self.candidates[j], self.candidates[i]
}
func (self *sortableDispatcherCandidates) Less(i, j int) bool {
	ci := self.candidates[i]
	cj := self.candidates[j]
	// primary key
	if ci.matchesCount < cj.matchesCount {
		return true
	}
	if ci.matchesCount > cj.matchesCount {
		return false
	}
	// secundary key
	if ci.literalCount < cj.literalCount {
		return true
	}
	if ci.literalCount > cj.literalCount {
		return false
	}
	// tertiary key
	return ci.nonDefaultCount < cj.nonDefaultCount
}
