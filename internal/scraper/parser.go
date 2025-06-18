package scraper

import (
	"io"

	"golang.org/x/net/html"
)


type Link struct {
	Text string
	Href string
}


func ParseLinks(r io.Reader, url string)([]Link, error){
	doc, err := html.Parse(r);

	if err != nil {
         return nil, &ErrParseFailed{ // <-- NEW: Return custom error
			URL:        url,
			Reason:     "failed to parse HTML document",
			WrappedErr: err,
		}
	}
    
	var links []Link
	var f func(*html.Node)

	f = func (n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			var href string

			for _, a:= range n.Attr {
				if a.Key == "href" {
					href = a.Val
					break
				}
			}

			var text string
			if n.FirstChild != nil && n.FirstChild.Type == html.TextNode {
				text = n.FirstChild.Data
			}

			if href != "" {
				links = append(links, Link{Text: text, Href: href})
			}
		}

		for c := n.FirstChild; c!= nil; c = c.NextSibling {
			f(c)
		}
	}

	f(doc)
	return links, nil
}