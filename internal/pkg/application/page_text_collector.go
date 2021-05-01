package application

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"

	"golang.org/x/net/html"
)

type BookmarkContentCollector struct {
	
}

func (pbc BookmarkContentCollector) CollectText(readCloser io.ReadCloser) (io.ReadCloser, error) {
	doc, err := html.Parse(readCloser)
	if err != nil {
		return nil, err
	}

	text := &bytes.Buffer{}
	title, err := getTagByName(doc, "title")
	if err == nil {
		collectText(title, text)
	}

	body, err := getTagByName(doc, "body")
	if err == nil {
		collectText(body, text)
	}

	return ioutil.NopCloser(text), nil
}

func collectText(n *html.Node, buf *bytes.Buffer) {
	if n.Type == html.TextNode {
		buf.WriteString(n.Data)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		collectText(c, buf)
	}
}

func getTagByName(doc *html.Node, tagName string) (*html.Node, error) {
	var body *html.Node
	var crawler func(*html.Node)
	crawler = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == tagName {
			body = node
			return
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			crawler(child)
		}
	}
	crawler(doc)
	if body != nil {
		return body, nil
	}

	return nil, fmt.Errorf("missing <%s> in the node tree", tagName)
}
