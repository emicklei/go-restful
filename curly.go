package restful

import (
	"net/http"
	"strings"
)

// CurlyRouter expects Routes with paths that contain zero or more parameters in curly brackets.
type CurlyRouter struct{}

// SelectRoute finds a Route given the input HTTP Request and report if found (ok).
// The HTTP writer is be used to directly communicate non-200 HTTP stati.
func (c CurlyRouter) SelectRoute(
	path string,
	webServices []*WebService,
	httpWriter http.ResponseWriter,
	httpRequest *http.Request) (selectedService *WebService, selected *Route, ok bool) {

	requestTokens := strings.Split(httpRequest.URL.Path, "/")

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
	candidates := []Route{}
	for _, each := range ws.Routes() {
		if c.matchesRouteByPathTokens(each.pathParts, requestTokens) {
			candidates = append(candidates, each)
		}
	}
	return candidates
}

func (c CurlyRouter) matchesRouteByPathTokens(routeTokens, requestTokens []string) bool {
	return true
}

func (c CurlyRouter) detectRoute(candidateRoutes []Route, httpWriter http.ResponseWriter, httpRequest *http.Request) *Route {
	return nil
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
	// return a weighted score of the longest matching consecutive tokens from the beginning
	min := len(requestTokens)
	if len(tokens) < min {
		min = len(tokens)
	}
	score := 0
	for i := 0; i < min; i++ {
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
			score += (min - i) * 10 //fuzzy
		}
	}
	return true, score
}
