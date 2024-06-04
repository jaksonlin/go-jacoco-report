package util

import (
	"bytes"
	"fmt"

	"golang.org/x/net/html"
)

// findTable recursively searches for the first table element in the HTML document
func GetHtmlTable(n *html.Node) *html.Node {
	if n.Type == html.ElementNode && n.Data == "table" {
		return n
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if table := GetHtmlTable(c); table != nil {
			return table
		}
	}
	return nil
}

func GetHtmlNodeFromUrl(url string) (*html.Node, error) {
	// download web page
	body, err := DownloadWebPage(url)
	if err != nil {
		return nil, err
	}
	// parse html
	doc, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	return doc, nil
}

// func to retrieve data from table, first arg is data, second arg is total column count,
// implementation is responsible to maintain the column index inside the method
type DataRetrieverFunc func(string, int) error

func TraverseJacocoHtmlTable(table *html.Node, dataRetrieverFunc DataRetrieverFunc) error {
	// Iterate over rows and columns
	columnCount := 0
	for row := table.FirstChild; row != nil; row = row.NextSibling {
		if row.Type == html.ElementNode && row.Data == "thead" {
			// check how many columns are there
			for col := row.FirstChild; col != nil; col = col.NextSibling {
				if col.Type == html.ElementNode && col.Data == "tr" {
					for trCol := col.FirstChild; trCol != nil; trCol = trCol.NextSibling {
						if trCol.Type == html.ElementNode && (trCol.Data == "td" || trCol.Data == "th") {
							columnCount++
						}
					}
				}
			}
			break
		}
	}
	if columnCount == 0 {
		return fmt.Errorf("no column found in the table")
	}
	for row := table.FirstChild; row != nil; row = row.NextSibling {
		if row.Type == html.ElementNode {
			if row.Data == "tbody" {
				if err := findJacocoTableTr(row, dataRetrieverFunc, columnCount); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func findJacocoTableTr(tbody *html.Node, dataRetrieverFunc DataRetrieverFunc, totalColumnCount int) error {
	for col := tbody.FirstChild; col != nil; col = col.NextSibling {
		if col.Type == html.ElementNode && col.Data == "tr" {
			for trCol := col.FirstChild; trCol != nil; trCol = trCol.NextSibling {
				if trCol.Type == html.ElementNode && (trCol.Data == "td" || trCol.Data == "th") {
					if err := dataRetrieverFunc(renderNode(trCol), totalColumnCount); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

// renderNode extracts the text content of an HTML node and its children
func renderNode(n *html.Node) string {
	var text string
	if n.Type == html.TextNode {
		text = fmt.Sprintf("text=%s;", n.Data)
	} else if n.Type == html.ElementNode {
		if n.Data == "a" {
			text = fmt.Sprintf("href=%s;", n.Attr[0].Val)
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		text += renderNode(c)
	}
	return text
}
