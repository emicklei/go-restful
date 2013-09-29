package restful

// Copyright 2013 Ernest Micklei. All rights reserved.
// Use of this source code is governed by a license
// that can be found in the LICENSE file.

import (
	//	"log"
	"net/http"
	"sort"
	"strings"
)

// CurlyRouter expects Routes with paths that contain zero or more parameters in curly brackets.
type CurlyRouter struct{}

// SelectRoute finds a Route given the input HTTP Request and report if found (ok).
// The HTTP writer is be used to directly communicate non-200 HTTP stati.
func (c CurlyRouter) SelectRoute(
	webServices []*WebService,
	httpWriter http.ResponseWriter,
	httpRequest *http.Request) (selectedService *WebService, selected *Route, ok bool) {

	requestTokens := tokenizePath(httpRequest.URL.Path)

	detectedService := c.detectWebService(requestTokens, webServices)
	if detectedService == nil {
		return nil, nil, false
	}
	candidateRoutes := c.selectRoutes(detectedService, httpWriter, requestTokens)
	if len(candidateRoutes) == 0 {
		return detectedService, nil, false
	}
	selectedRoute := c.detectRoute(candidateRoutes, httpWriter, httpRequest)
	if selectedRoute == nil {
		return detectedService, nil, false
	}
	return detectedService, selectedRoute, true
}

func (c CurlyRouter) selectRoutes(ws *WebService, httpWriter http.ResponseWriter, requestTokens []string) []Route {
	candidates := &sortableCurlyRoutes{[]*curlyRoute{}}
	for _, each := range ws.routes {
		matches, paramCount, staticCount := c.matchesRouteByPathTokens(each.pathParts, requestTokens)
		if matches {
			candidates.add(&curlyRoute{each, paramCount, staticCount}) // TODO make sure Routes() return pointers?
		}
	}
	sort.Sort(sort.Reverse(candidates))
	return candidates.routes()
}

func (c CurlyRouter) matchesRouteByPathTokens(routeTokens, requestTokens []string) (matches bool, paramCount int, staticCount int) {
	if len(routeTokens) != len(requestTokens) {
		return false, 0, 0
	}
	for i, routeToken := range routeTokens {
		requestToken := requestTokens[i]
		if !strings.HasPrefix(routeToken, "{") {
			if requestToken != routeToken {
				return false, 0, 0
			}
			staticCount++
		} else {
			paramCount++
		}
	}
	return true, paramCount, staticCount
}

func (c CurlyRouter) detectRoute(candidateRoutes []Route, httpWriter http.ResponseWriter, httpRequest *http.Request) *Route {
	route, found := RouterJSR311{}.detectRoute(candidateRoutes, httpWriter, httpRequest) // TODO change signature
	if found {
		return route
	} else {
		return nil
	}
}

func (c CurlyRouter) detectWebService(requestTokens []string, webServices []*WebService) *WebService {
	if len(webServices) == 0 {
		return nil
	}
	var best *WebService
	score := -1
	for _, each := range webServices {
		matches, eachScore := c.computeWebserviceScore(requestTokens, each.pathExpr.tokens)
		if matches && (eachScore > score) {
			best = each
			score = eachScore
		}
	}
	return best
}

func (c CurlyRouter) computeWebserviceScore(requestTokens []string, tokens []string) (bool, int) {
	// return whether tokens match and the weighted score of the longest matching consecutive tokens from the beginning

	if len(tokens) > len(requestTokens) {
		return false, 0
	}
	score := 0
	for i := 0; i < len(tokens); i++ {
		each := requestTokens[i]
		other := tokens[i]
		if len(each) == 0 && len(other) == 0 {
			score++
			continue
		}
		if len(other) > 0 && strings.HasPrefix(other, "{") {
			// no empty match
			if len(each) == 0 {
				return false, score
			}
			score += 1
		} else {
			// not a parameter
			if each != other {
				return false, score
			}
			score += (len(tokens) - i) * 10 //fuzzy
		}
	}
	return true, score
}
