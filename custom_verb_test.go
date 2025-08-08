package restful

import (
	"sync"
	"testing"
)

func TestHasCustomVerb(t *testing.T) {
	testCase := []struct {
		path string
		has  bool
	}{
		{"/{userId}:init", true},
		{"/{userId:init}", false},
		{"/users/{id:init}:init", true},
		{"/users/{id}", false},
	}

	for _, v := range testCase {
		rs := hasCustomVerb(v.path)
		if rs != v.has {
			t.Errorf("path: %v should has no custom verb", v.path)
		}
	}
}

func TestRemoveCustomVerb(t *testing.T) {
	testCase := []struct {
		path         string
		expectedPath string
	}{
		{"/{userId}:init", "/{userId}"},
		{"/{userId:init}", "/{userId:init}"},
		{"/users/{id:init}:init", "/users/{id:init}"},
		{"/users/{id}", "/users/{id}"},
		{"/init/users/{id:init}:init", "/init/users/{id:init}"},
	}

	for _, v := range testCase {
		rs := removeCustomVerb(v.path)
		if rs != v.expectedPath {
			t.Errorf("expected value: %v, actual: %v", v.expectedPath, rs)
		}
	}
}

func TestMatchCustomVerb(t *testing.T) {
	testCase := []struct {
		routeToken string
		pathToken  string
		expected   bool
	}{
		{"{userId:regex}:init", "creator-123456789gWe:init", true},
		{"{userId:regex}:init", "creator-123456789gWe", false},
		{"{userId:regex}", "creator-123456789gWe:init", false},
		{"users:init", "users:init", true},
		{"users:init", "tokens:init", true},
	}

	for idx, v := range testCase {
		rs := isMatchCustomVerb(v.routeToken, v.pathToken)
		if rs != v.expected {
			t.Errorf("expected value: %v, actual: %v, index: [%v]", v.expected, rs, idx)
		}
	}
}

func TestCustomVerbCaching(t *testing.T) {
	// Clear cache before test
	customVerbCache = sync.Map{}
	
	routeToken := "{userId:regex}:POST"
	pathToken := "user123:POST"
	
	// First call should cache the regex
	result1 := isMatchCustomVerb(routeToken, pathToken)
	if !result1 {
		t.Error("Expected first call to match")
	}
	
	// Verify cache contains the pattern
	pattern := ":POST$"
	_, found := customVerbCache.Load(pattern)
	if !found {
		t.Error("Expected pattern to be cached")
	}
	
	// Second call should use cached regex
	result2 := isMatchCustomVerb(routeToken, pathToken)
	if !result2 {
		t.Error("Expected second call to match using cache")
	}
	
	// Test with different verb to ensure separate caching
	routeToken2 := "{userId:regex}:GET"
	pathToken2 := "user123:GET"
	
	result3 := isMatchCustomVerb(routeToken2, pathToken2)
	if !result3 {
		t.Error("Expected GET verb to match")
	}
	
	// Verify both patterns are cached
	pattern2 := ":GET$"
	_, found2 := customVerbCache.Load(pattern2)
	if !found2 {
		t.Error("Expected GET pattern to be cached")
	}
}
