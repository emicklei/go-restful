package restful

import (
	"fmt"
	"regexp"
	"sync"
)

var (
	customVerbReg   = regexp.MustCompile(":([A-Za-z]+)$")
	customVerbCache sync.Map // Cache for compiled custom verb regexes
)

func hasCustomVerb(routeToken string) bool {
	return customVerbReg.MatchString(routeToken)
}

func isMatchCustomVerb(routeToken string, pathToken string) bool {
	rs := customVerbReg.FindStringSubmatch(routeToken)
	if len(rs) < 2 {
		return false
	}

	customVerb := rs[1]
	regexPattern := fmt.Sprintf(":%s$", customVerb)

	// Check cache first
	if cached, found := customVerbCache.Load(regexPattern); found {
		specificVerbReg := cached.(*regexp.Regexp)
		return specificVerbReg.MatchString(pathToken)
	}

	// Compile and cache the regex
	specificVerbReg := regexp.MustCompile(regexPattern)
	customVerbCache.Store(regexPattern, specificVerbReg)
	return specificVerbReg.MatchString(pathToken)
}

func removeCustomVerb(str string) string {
	return customVerbReg.ReplaceAllString(str, "")
}
