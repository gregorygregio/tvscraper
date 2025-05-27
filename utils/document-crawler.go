package utils

import (
	"fmt"
	"io"
	"strings"

	"golang.org/x/net/html"
)

type DocumentCrawler struct {
	DocumentRoot *html.Node
}

type onElementAction func(*html.Node)

func NewDocumentCrawler(reader io.Reader) (*DocumentCrawler, error) {
	doc, err := html.Parse(reader)
	if err != nil {
		fmt.Printf("Error parsing HTML: %s\n", err.Error())
		return nil, fmt.Errorf("error parsing HTML")
	}

	return &DocumentCrawler{
		DocumentRoot: doc,
	}, nil
}

func (dc *DocumentCrawler) ForEachElement(action onElementAction) {
	parseNode(dc.DocumentRoot, action)
}

func parseNode(n *html.Node, action onElementAction) {
	if n.Data == "script" || n.Data == "style" || n.Data == "head" {
		return
	}
	action(n)

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		parseNode(c, action)
	}
}

func HasClass(n *html.Node, className string) bool {
	return HasAttr(n, "class", className)
}

func HasAttr(n *html.Node, attrName string, val string) bool {
	if n.Type == html.ElementNode {
		for _, attr := range n.Attr {
			if attr.Key == attrName && strings.Contains(attr.Val, val) {
				return true
			}
		}
	}
	return false
}
