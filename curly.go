package restful

// Copyright 2013 Ernest Micklei. All rights reserved.
// Use of this source code is governed by a license
// that can be found in the LICENSE file.

import (
	"errors"
	"net/http"
	"sort"
	"strings"
)

// CurlyRouter expects Routes with paths that contain zero or more parameters in curly brackets.
type CurlyRouter struct{}

// SelectRoute is part of the Router interface and returns the best match
// for the WebService and its Route for the given Request.
func (c CurlyRouter) SelectRoute(
	webServices []*WebService,
	httpRequest *http.Request) (selectedService *WebService, selected *Route, err error) {

	requestTokens := tokenizePath(httpRequest.URL.Path)

	detectedService := c.detectWebService(requestTokens, webServices)
	if detectedService == nil {
		return nil, nil, errors.New("[restful] no detected service")
	}
	candidateRoutes := c.selectRoutes(detectedService, requestTokens)
	if len(candidateRoutes) == 0 {
		return detectedService, nil, errors.New("[restful] no candidate routes")
	}
	selectedRoute, err := c.detectRoute(candidateRoutes, httpRequest)
	if selectedRoute == nil {
		return detectedService, nil, err
	}
	return detectedService, selectedRoute, nil
}

// selectRoutes return a collection of Route from a WebService that matches the path tokens from the request.
func (c CurlyRouter) selectRoutes(ws *WebService, requestTokens []string) []Route {
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

// matchesRouteByPathTokens computes whether it matches, howmany parameters do match and what the number of static path elements are.
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

// detectRoute selectes from a list of Route the first match by inspecting both the Accept and Content-Type
// headers of the Request. See also RouterJSR311 in jsr311.go
func (c CurlyRouter) detectRoute(candidateRoutes []Route, httpRequest *http.Request) (*Route, error) {
	return RouterJSR311{}.detectRoute(candidateRoutes, httpRequest)
}

// detectWebService returns the best matching webService given the list of path tokens.
// see also computeWebserviceScore
func (c CurlyRouter) detectWebService(requestTokens []string, webServices []*WebService) *WebService {
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

// computeWebserviceScore returns whether tokens match and
// the weighted score of the longest matching consecutive tokens from the beginning.
func (c CurlyRouter) computeWebserviceScore(requestTokens []string, tokens []string) (bool, int) {
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
