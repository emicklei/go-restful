package restful

import (
    "testing"
)

func TestMatchesPath(t *testing.T) {
	r := Route{Path: "/from/{source}/to/{destination}"}
	matches, params := r.MatchesPath("/from/AMS/to/NY")
	if (!matches) {
		t.Error("should have matched")
		return
	}
	if (params["source"] != "AMS") {
		t.Errorf("parameter mismatch %v",params)	
	}
}