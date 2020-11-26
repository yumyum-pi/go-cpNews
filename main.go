package main

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

// flag variables
var url *string
var articleClass *string
var textTag *string
var ifStats bool

// write a better error handler
func errHandle(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}

// get the data from the URL
func scrapeData(finalText *string) {
	// get date from the URL
	res, err := http.Get(*url)
	errHandle(err)
	if res.StatusCode != 200 {
		errHandle(errors.New("error code from the server"))
	}

	defer res.Body.Close()

	// read the document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	errHandle(err)

	// query the article element
	doc.Find(*articleClass).Each(func(i int, a *goquery.Selection) {
		// temp store text
		var text string

		// query the text element
		a.Find(*textTag).Each(func(j int, t *goquery.Selection) {
			// get the number of children
			c := t.Children()
			n := len(c.Nodes)

			// ignore the elements that have children
			if n == 0 {
				// get the text form the element
				text = strings.TrimSpace(t.Text())
				// add the text to the final text
				// if the text is not blank
				if text != "" {
					// add extra space if the no the 1st element
					if j != 0 {
						*finalText += (" " + text)
					} else {
						*finalText += text
					}
				}
			}

		})

	})
}

// count the no of HTML children
func numberOfChild(n *html.Node) int {
	if n == nil {
		return -1
	}
	count := 0

	// c is the current child
	// loop until c is null and iterate the next sibling
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		count++
	}

	return count
}

// map the character in the given text
func mapChar(s *string) map[string]int {
	m := make(map[string]int)
	for i := range *s {
		m[string((*s)[i])]++
	}
	return m
}

// display stats
func stats(s *string) {
	// get the no. of characters
	c := len(*s)
	m := mapChar(s)

	// divide the string to words
	words := strings.Fields(*s)

	fmt.Printf("no. of \ncharacters: %d\nwords: %d\nmap: %v", c, len(words), m)

	// display in a different format
	// sort the data
	// print the map of the character
	for i := range m {
		fmt.Println(i, ":", m[i])
	}
}

func main() {
	// get the flag data
	url = flag.String("u", "", "URL for the web article")
	articleClass = flag.String("a", "", "URL for the web article")
	textTag = flag.String("t", "p", "tag of text")
	flag.BoolVar(&ifStats, "s", false, "show stats of the text")

	flag.Parse()

	// check if the required flags are not nil

	// store the text
	var text string

	// scrape data from the URL
	scrapeData(&text)

	fmt.Println(text)

	// display stats if the flag is true
	if ifStats {
		stats(&text)
	}
}
