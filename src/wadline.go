package main

import (
	"net/url"
	"strconv"
	"strings"
	"sync"

	"github.com/alexander-molina/companies_extractor/3rdparty/goquery"
)

var foundCompanies sync.Map

// Wadline предоставляет алгоритм для извлечения данных из сайта https://wadline.ru/
func Wadline() {
	siteName := "wadline"
	urls := [...]string{"https://wadline.ru/design", "https://wadline.ru/web"}
	const lowerIdx = 1
	const higherIdx = 55

	for _, URL := range urls {
		for pageIdx := lowerIdx; pageIdx < higherIdx; pageIdx++ {
			page, _ := GetPage(URL + "?page=" + strconv.Itoa(pageIdx))
			go getCompanies(page)
		}
	}

	WriteToExcel(siteName, &foundCompanies)

	// foundCompanies.Range(func(k, v interface{}) bool {
	// 	fmt.Printf("key:%s, val:%+v", k, v)
	// 	return true
	// })
}

func getCompanies(doc *goquery.Document) {
	doc.Find(".info").Each(func(i int, s *goquery.Selection) {
		name := s.Find(".h1 a").Text()
		name = strings.Split(name, "\n")[1]
		name = strings.ToLower(name)

		location := s.Find("p").Text()
		location = strings.Split(location, "\t")[0]
		l := strings.Split(location, "\n")
		location = l[0] + l[1]

		companyLink, _ := s.Find(".h1 a").Attr("href")
		companyLink = "https://wadline.ru" + companyLink
		companyURL := getCompanyURL(companyLink)

		emails := make(map[string]string)
		if companyURL != "" {
			companyURLs := NavigateSite(companyURL)
			for URL := range companyURLs {
				if URL != "" {
					go SearchEmails(URL, &emails)
				}
			}
		}

		company := Company{name, companyURL, location, emails}
		foundCompanies.Store(name, company)
	})
}

func getCompanyURL(companyLink string) string {
	doc, _ := GetPage(companyLink)
	URL, _ := doc.Find(".cc-web a").Attr("href")

	parsedURL, err := url.Parse(URL)

	if err != nil {
		return ""
	}

	return parsedURL.Scheme + "://" + parsedURL.Host
}
