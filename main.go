package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getAbsolute(link string, baseDomain string) string {
	urlLink, err := url.Parse(strings.TrimSpace(link))

	if err != nil {
		fmt.Println("Error while parsing url.\nError: ", err)
		os.Exit(1)
	}

	if !urlLink.IsAbs() {
		baseDomainLink, err := url.Parse(strings.TrimSpace(baseDomain))

		if err != nil {
			fmt.Println("Error while parsing url.\nError: ", err)
			os.Exit(1)
		}

		wholeLink := baseDomainLink.ResolveReference(urlLink).String()

		return wholeLink
	}

	return link
}

func getLinks(res *http.Response, baseURL string) []string {
	doc, err := goquery.NewDocumentFromReader(res.Body)

	if err != nil {
		fmt.Println("Error while parsing the reponse body for goquery documentation.\nError: ", err)
		os.Exit(1)
	}

	list := make([]string, 0)

	doc.Find("a").Each(func(i int, p *goquery.Selection) {
		val, exists := p.Attr("href")

		val = getAbsolute(val, baseURL)

		if exists {
			list = append(list, val)
		}
	})

	return list
}

func getSite(url string) (*http.Response, error) {
	return http.Get(url)
}

func crawler(url string) []string {
	fmt.Println("Message: Fetching the site \"", url, "\"")

	res, err := getSite(url)

	if err != nil {
		fmt.Println("Error encountered while fetching the site.\nError: ", err)
		os.Exit(1)
	}

	list := getLinks(res, url)

	if err != nil {
		fmt.Println("Error encountered while trying to get links.\nError: ", err)
		os.Exit(1)
	}

	return list
}

func main() {
	baseURL := "https://www.guardian.com"

	links := crawler(baseURL)

	track := make(map[string]bool)

	for _, link := range links {

		if !track[link] {
			track[link] = true

			links = append(links, crawler(link)...)
		}
	}
}
