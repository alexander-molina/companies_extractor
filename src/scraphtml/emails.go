package scraphtml

import (
	"log"
	"net/http"
	"net/url"

	"github.com/alexander-molina/companies_extractor/3rdparty/goquery"
)

// SearcEmail searches all eamils in current page
func SearcEmail(pageURL string) []string {
	// Request the HTML page.
	res, err := http.Get(pageURL)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var emails []string
	// Find the review items
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		link, _ := s.Attr("href")
		parsedLink, error := url.Parse(link)
		if error != nil {
			return
		}
		if parsedLink.Scheme == "mailto" {
			emails = append(emails, parsedLink.Opaque)
		}
	})

	return emails
}
