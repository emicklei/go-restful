package restful

// This file implements the flow for matching Requests to Routes (and consequently Resource Functions)
// as specified by the JSR311 http://jsr311.java.net/nonav/releases/1.1/spec/spec.html.
// Concept of locators is not implemented.
import (
	"bytes"
	"errors"
	//	"github.com/emicklei/hopwatch"
	"log"
	"net/http"
	"regexp"
	"sort"
	"strings"
)

// http://jsr311.java.net/nonav/releases/1.1/spec/spec3.html#x3-360003.7.2
func detectRoute(routes []Route, httpWriter http.ResponseWriter, httpRequest *http.Request) (Route, bool) {
	// http method
	methodOk := []Route{}
	for _, each := range routes {
		if httpRequest.Method == each.Method {
			methodOk = append(methodOk, each)
		}
	}
	if len(methodOk) == 0 {
		httpWriter.WriteHeader(http.StatusMethodNotAllowed)
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
			return Route{}, false
		}
	}
	// accept
	outputMediaOk := []Route{}
	accept := httpRequest.Header.Get(HEADER_Accept)
	for _, each := range inputMediaOk {
		if each.matchesAccept(accept) {
			outputMediaOk = append(outputMediaOk, each)
		}
	}
	if len(outputMediaOk) == 0 {
		httpWriter.WriteHeader(http.StatusNotAcceptable)
		return Route{}, false
	}
	return bestMatchByMedia(outputMediaOk, contentType, accept), true
}

// http://jsr311.java.net/nonav/releases/1.1/spec/spec3.html#x3-360003.7.2
func bestMatchByMedia(routes []Route, contentType string, accept string) Route {
	// TODO
	return routes[0]
}

// http://jsr311.java.net/nonav/releases/1.1/spec/spec3.html#x3-360003.7.2  (step 2)
func selectRoutes(dispatcher Dispatcher, pathRemainder string) []Route {
	if pathRemainder == "" || pathRemainder == "/" {
		return dispatcher.Routes()
	}
	filtered := sortableRouteCandidates{}
	for _, each := range dispatcher.Routes() {
		expression, literalCount, varCount := templateToRegularExpression(each.relativePath)
		compiled, err := regexp.Compile(expression)
		if err != nil {
			log.Printf("Invalid template %v because: %v. Ignore route\n", each.Path, err)
		} else {
			matches := compiled.FindStringSubmatch(pathRemainder)
			if matches != nil {
				lastMatch := matches[len(matches)-1]
				if lastMatch == "" || lastMatch == "/" { // do not include if value is neither empty nor ‘/’.
					filtered.candidates = append(filtered.candidates,
						routeCandidate{each, expression, len(matches), literalCount, varCount})
				}
			}
		}
	}
	if len(filtered.candidates) == 0 {
		return []Route{}
	}
	sort.Sort(filtered)
	rmatch := filtered.candidates[0].regexpression
	// select routes from candidates whoes expression matches rmatch
	matchingRoutes := []Route{}
	for _, each := range filtered.candidates {
		if each.regexpression == rmatch {
			matchingRoutes = append(matchingRoutes, each.route)
		}
	}
	return matchingRoutes
}

// http://jsr311.java.net/nonav/releases/1.1/spec/spec3.html#x3-360003.7.2
func detectDispatcher(requestPath string, dispatchers []Dispatcher) (Dispatcher, string, error) {
	filtered := sortableDispatcherCandidates{}
	for _, each := range dispatchers {
		expression, literalCount, varCount := templateToRegularExpression(each.RootPath())
		compiled, err := regexp.Compile(expression)
		if err != nil {
			log.Printf("Invalid template %v because: %v. Ignore dispatcher\n", each.RootPath(), err)
		} else {
			matches := compiled.FindStringSubmatch(requestPath)
			if matches != nil {
				filtered.candidates = append(filtered.candidates,
					dispatcherCandidate{each, matches[len(matches)-1], len(matches), literalCount, varCount})
			}
		}
	}
	if len(filtered.candidates) == 0 {
		return nil, "", errors.New("not found")
	}
	sort.Sort(filtered)
	//hopwatch.Break("detectDispatcher", "filtered", filtered)
	return filtered.candidates[0].dispatcher, filtered.candidates[0].finalMatch, nil
}

// http://jsr311.java.net/nonav/releases/1.1/spec/spec3.html#x3-370003.7.3
func templateToRegularExpression(template string) (expression string, literalCount int, varCount int) {
	var buffer bytes.Buffer
	buffer.WriteString("^")
	tokens := strings.Split(template, "/")
	for _, each := range tokens {
		if each == "" {
			continue
		}
		buffer.WriteString("/")
		if strings.HasPrefix(each, "{") {
			// ignore var spec
			varCount += 1
			buffer.WriteString("([^/]+?)")
		} else {
			literalCount += len(each)
			encoded := each // TODO URI encode
			buffer.WriteString(regexp.QuoteMeta(encoded))
		}
	}
	return strings.TrimRight(buffer.String(), "/") + "(/.*)?$", literalCount, varCount
}

// Types and functions to support the sorting of Routes

type routeCandidate struct {
	route           Route
	regexpression   string
	matchesCount    int
	literalCount    int
	nonDefaultCount int
}
type sortableRouteCandidates struct {
	candidates []routeCandidate
}

func (self sortableRouteCandidates) Len() int {
	return len(self.candidates)
}
func (self sortableRouteCandidates) Swap(i, j int) {
	self.candidates[i], self.candidates[j] = self.candidates[j], self.candidates[i]
}
func (self sortableRouteCandidates) Less(j, i int) bool { // Do reverse so the i and j are in this order
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

// Types and functions to support the sorting of Dispatchers

type dispatcherCandidate struct {
	dispatcher      Dispatcher
	finalMatch      string
	matchesCount    int
	literalCount    int
	nonDefaultCount int
}
type sortableDispatcherCandidates struct {
	candidates []dispatcherCandidate
}

func (self sortableDispatcherCandidates) Len() int {
	return len(self.candidates)
}
func (self sortableDispatcherCandidates) Swap(i, j int) {
	self.candidates[i], self.candidates[j] = self.candidates[j], self.candidates[i]
}
func (self sortableDispatcherCandidates) Less(j, i int) bool { // Do reverse so the i and j are in this order
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
