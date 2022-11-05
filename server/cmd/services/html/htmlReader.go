package html

import (
	"bytes"

	"golang.org/x/net/html"
)

func FindElements(rawHtml, tag string, callback func(elem *html.Node) (stop bool)) error {
	reader := bytes.NewReader([]byte(rawHtml))
	doc, err := html.Parse(reader)
	if err != nil {
		return err
	}

	findInHtml(doc, tag, &callback)
	return nil
}
func findInHtml(doc *html.Node, tag string, cb *func(elem *html.Node) (stop bool)) {
	for child := doc.FirstChild; child != nil; child = child.NextSibling {
		if child.Type == html.ElementNode && child.Data == tag {
			if stop := (*cb)(child); stop {
				return
			}
		}
		findInHtml(child, tag, cb)
	}
}
