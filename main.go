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

var url *string
var articleClass *string
var textTag *string

// write a better error handler
func errHandle(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}

// get the data from the URL
func getData(mainString *string) {
	// get date from the URL
	res, err := http.Get(*url)
	errHandle(err)
	//	fmt.Printf("status code: %s\n", res.Status)
	if res.StatusCode != 200 {
		errHandle(errors.New("error code from the server"))
	}

	defer res.Body.Close()

	// read the document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	errHandle(err)
	doc.Find(*articleClass).Each(func(i int, s *goquery.Selection) {
		// sss store temp value
		var sss string

		s.Find(*textTag).Each(func(j int, ss *goquery.Selection) {
			c := ss.Children()
			n := len(c.Nodes)
			if n == 0 {
				sss = strings.TrimSpace(ss.Text())
				if sss != "" {
					if j != 0 {
						*mainString += (" " + sss)
					} else {
						*mainString += sss
					}
				}
			}

		})

	})
}

func numberOfChild(n *html.Node) int {
	if n == nil {
		return -1
	}
	count := 0
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		count++
	}
	return count
}

func wordCount(s *string) (int, map[string]int) {
	words := strings.Fields(*s)
	m := make(map[string]int)
	for i := range *s {
		m[string((*s)[i])]++
	}
	return len(words), m
}

// display stats
func stats(s *string) {
	// get the no. of characters
	c := len(*s)
	w, m := wordCount(s)

	fmt.Printf("no. of \ncharacters: %d\nwords: %d\nmap: %v", c, w, m)
	for i := range m {
		fmt.Println(i, ":", m[i])
	}
}

func main() {
	var statsB bool
	url = flag.String("u", "", "URL for the web article")
	articleClass = flag.String("a", "", "URL for the web article")
	textTag = flag.String("t", "p", "tag of text")
	flag.BoolVar(&statsB, "s", false, "tag of text")

	flag.Parse()

	var s string

	getData(&s)
	fmt.Println(s)

	if statsB {
		stats(&s)
	}
}
