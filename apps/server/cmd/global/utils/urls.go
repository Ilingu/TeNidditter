package utils

import "net/url"

// Return whether the input url is a valid HTTP url or not.
func IsValidURL(urlToCheck string) bool {
	_, err := url.ParseRequestURI(urlToCheck)
	return err == nil
}

// Return the url is already encoded/QueryEscape()
func IsUrlEncoded(str string) bool {
	dec, err := url.QueryUnescape(str)
	return err != nil && dec != str
}
