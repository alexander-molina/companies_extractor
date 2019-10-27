package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/alexander-molina/companies_extractor/3rdparty/goquery"
	"github.com/alexander-molina/companies_extractor/3rdparty/xlsx"
)

// GetPage запрашивает и возвращает нужную веб страницу
func GetPage(pageURL string) (doc *goquery.Document, err error) {
	res, err := http.Get(pageURL)

	if err != nil {
		return nil, err
	}

	doc, err = goquery.NewDocumentFromReader(res.Body)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	defer res.Body.Close()

	return doc, nil
}

// NavigateSite находит и возвращает ссылки на другие страницы сайта
func NavigateSite(rootURL string) map[string]string {
	links := make(map[string]string)
	parsedURL, _ := url.Parse(rootURL)
	doc, err := GetPage(rootURL)
	if err == nil {
		doc.Find("a").Each(func(i int, s *goquery.Selection) {
			link, _ := s.Attr("href")
			parsedLink, err := url.Parse(link)
			if err == nil {
				if parsedLink.Host == parsedURL.Host && parsedLink.Scheme == parsedURL.Scheme {
					links[link] = ""
				}
			}
		})
	}

	return links
}

// SearchEmails ищет все email адреса на текущей станице
func SearchEmails(pageURL string, emails *map[string]string) {
	doc, err := GetPage(pageURL)
	if err == nil {
		doc.Find("a").Each(func(i int, s *goquery.Selection) {
			link, exists := s.Attr("href")
			if exists {
				parsedLink, err := url.Parse(link)
				if err == nil {
					if parsedLink.Scheme == "mailto" {
						(*emails)[parsedLink.Opaque] = ""
					}
				}
			}
		})
	}
}

// WriteToExcel записывает полученные данные в excel таблицу
func WriteToExcel(siteName string, data *sync.Map) {
	file := xlsx.NewFile()
	sheet, _ := file.AddSheet(siteName)

	row := sheet.AddRow()
	companyName := row.AddCell()
	companyURL := row.AddCell()
	companyLocation := row.AddCell()
	companyEmails := row.AddCell()

	companyName.Value = "Name"
	companyURL.Value = "URL"
	companyLocation.Value = "Location"
	companyEmails.Value = "Emails"

	fmt.Println("WRITING!")
	data.Range(func(k, v interface{}) bool {
		fmt.Println(v)
		company := v.(Company)

		row = sheet.AddRow()
		companyName = row.AddCell()
		companyURL = row.AddCell()
		companyLocation = row.AddCell()
		companyEmails = row.AddCell()

		companyName.Value = company.Name
		companyURL.Value = company.URL
		companyLocation.Value = company.Location

		emails := []string{}
		for key := range company.Emails {
			emails = append(emails, key)
		}

		companyEmails.Value = strings.Join(emails, ", ")

		return true
	})

	file.Save(siteName + ".xlsx")
}
