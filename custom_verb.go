package restful

import (
	"fmt"
	"regexp"
	"sync"
)

var (
	customVerbReg     = regexp.MustCompile(":([A-Za-z]+)$")
	customVerbCache   sync.Map // Cache for compiled custom verb regexes
	customVerbCacheEnabled = true // Enable/disable custom verb regex caching
)

// SetCustomVerbCacheEnabled enables or disables custom verb regex caching.
// When disabled, custom verb regex patterns will be compiled on every request.
// When enabled (default), compiled custom verb regex patterns are cached for better performance.
func SetCustomVerbCacheEnabled(enabled bool) {
	customVerbCacheEnabled = enabled
}

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

	// Check cache first (if enabled)
	if customVerbCacheEnabled {
		if cached, found := customVerbCache.Load(regexPattern); found {
			if specificVerbReg, ok := cached.(*regexp.Regexp); ok {
				return specificVerbReg.MatchString(pathToken)
			}
		}
	}

	// Compile the regex
	specificVerbReg := regexp.MustCompile(regexPattern)
	
	// Cache the regex (if enabled)
	if customVerbCacheEnabled {
		customVerbCache.Store(regexPattern, specificVerbReg)
	}
	
	return specificVerbReg.MatchString(pathToken)
}

func removeCustomVerb(str string) string {
	return customVerbReg.ReplaceAllString(str, "")
}
