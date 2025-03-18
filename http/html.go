package http

import "golang.org/x/net/html"

// extractMetaImage finds the og:image or similar meta tag from an HTML document
func extractMeta(doc *html.Node) meta {
	var m meta
	// Try to find Open Graph image first
	ogImage := findMetaContent(doc, "property", "og:image")
	if ogImage != "" {
		m.imageURL = ogImage
	}

	// Try Twitter image
	twitterImage := findMetaContent(doc, "name", "twitter:image")
	if twitterImage != "" {
		m.imageURL = twitterImage
	}

	// Try regular meta image
	metaImage := findMetaContent(doc, "name", "image")
	if metaImage != "" {
		m.imageURL = metaImage
	}

	m.title = findMetaContent(doc, "property", "og:title")
	if m.title == "" {
		// TODO: Look for non meta title
		m.title = "recipe"
	}
	m.description = findMetaContent(doc, "name", "description")
	m.siteName = findMetaContent(doc, "property", "og:site_name")
	return m
}

// findMetaContent looks for a meta tag with the specified attribute and value
// and returns its content attribute
func findMetaContent(doc *html.Node, attrName, attrValue string) string {
	var content string

	var crawler func(*html.Node)
	crawler = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "meta" {
			// Check if this meta tag has the attribute we're looking for
			var hasAttr bool
			var contentAttr string

			for _, attr := range n.Attr {
				if attr.Key == attrName && attr.Val == attrValue {
					hasAttr = true
				}
				if attr.Key == "content" {
					contentAttr = attr.Val
				}
			}

			if hasAttr && contentAttr != "" {
				content = contentAttr
				return
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if content == "" { // Only continue if we haven't found it yet
				crawler(c)
			}
		}
	}

	crawler(doc)
	return content
}
