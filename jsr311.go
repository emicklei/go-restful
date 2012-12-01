package restful

// This file implements the flow for matching Requests to Routes (and consequently Resource Functions)
// as specified by the JSR311 http://jsr311.java.net/nonav/releases/1.1/spec/spec.html
//
import (
	"bytes"
	"errors"
	"log"
	"regexp"
	"strings"
)

func detectDispatcher(requestPath string, dispatchers []Dispatcher) (Dispatcher, error) {
	filtered := []Dispatcher{}
	for _, each := range dispatchers {
		expression := templateToRegularExpression(each.RootPath())
		compiled, err := regexp.Compile(expression)
		if err != nil {
			log.Printf("Invalid template %v because: %v. Ignore dispatcher\n", each.RootPath(), err)
		} else {
			matches := compiled.MatchString(requestPath)
			if matches {
				filtered = append(filtered, each)
			}
		}
	}
	if len(filtered) == 0 {
		return nil, errors.New("not found")
	}
	return filtered[0], nil
}

//Convert URI Template to Regular Expression.
//URI encode the template, ignoring URI template variable specifications.
//Escape any regular expression characters in the URI template, again ignoring URI template variable specifications.
//Replace each URI template variable with a capturing group containing the specified regular expression or ‘([/]+?)’ if no regular expression is specified.
//If the resulting string ends with ‘/’ then remove the final character.
//Append ‘(/.*)?’ to the result.
func templateToRegularExpression(template string) string {
	var buffer bytes.Buffer
	tokens := strings.Split(template, "/")
	for _, each := range tokens {
		if each == "" {
			continue
		}
		buffer.WriteString("/")
		if strings.HasPrefix(each, "{") {
			// ignore var spec
			buffer.WriteString("([^/]+?)")
		} else {
			encoded := each // TODO URI encode
			buffer.WriteString(regexp.QuoteMeta(encoded))
		}
	}
	return strings.TrimRight(buffer.String(), "/") + "(/.*)?"
}
