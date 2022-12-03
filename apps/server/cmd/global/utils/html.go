package utils

import (
	"strings"
	"teniditter-server/cmd/services/html"

	htmlpkg "golang.org/x/net/html"
)

// It'll parse, tokenize the html and goes through every element of the document to check whether there is a script tag or not.
//
// If the html is invalid it'll return "true" by default
func ContainsScript(rawHtml string) (found bool) {
	err := html.FindElements(rawHtml, "script", func(elem *htmlpkg.Node) (stop bool) {
		found = true
		return true
	})
	if err != nil {
		return true
	}
	return
}

// A faster but less accurate "ContainsScript()" function: it will only look for the keyword "script" in the html, so if there is this word in a text or an attribute ect... this will return true even though there is no actual "script" tag
func ContainsScriptFast(rawHtml string) bool {
	return strings.Contains(rawHtml, "script")
}
