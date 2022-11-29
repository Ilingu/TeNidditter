package utils

import "net/url"

func IsValidURL(urlToCheck string) bool {
	_, err := url.ParseRequestURI(urlToCheck)
	return err == nil
}

func IsUrlEncoded(str string) bool {
	dec, err := url.QueryUnescape(str)
	return err != nil && dec != str
}
