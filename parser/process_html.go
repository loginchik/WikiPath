package parser

import "golang.org/x/net/html"

// FindBodyContent recursively iterates through HTML until it finds
// a div with ID = "bodyContent"
func FindBodyContent(n *html.Node) *html.Node {
	if n.Type == html.ElementNode && n.Data == "div" {
		// Get ID of a div
		for _, attr := range n.Attr {
			if attr.Key == "id" && attr.Val == "bodyContent" {
				return n
			}
		}
	}
	// Delve into current node
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if res := FindBodyContent(c); res != nil {
			return res
		}
	}
	return nil
}

// FindAllLinks extracts all <a> tags from HTML Node
// and gets their href attribute values
func FindAllLinks(n *html.Node, links *[]string) {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, attr := range n.Attr {
			if attr.Key == "href" && attr.Val != "" {
				// Save the link
				*links = append(*links, attr.Val)
			}
		}
	} else {
		// Delve into current Node to get its links
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			FindAllLinks(c, links)
		}
	}
}
