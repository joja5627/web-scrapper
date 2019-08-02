package scrape

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/joja5627/webscrapper/internal/async"
)

var (
	selectors = []string{".result-row .result-image", "#sortable-results > ul > li:nth-child(1) > p > a"}
)

func ScrapeCL(urls []string) []string {

	var links []string
	c := async.NewClient(http.DefaultClient, len(urls))
	async.FetchAll(urls, c)

	for i := 0; i < len(urls); i++ {
		select {
		case resp := <-c.Resp:

			document, err := goquery.NewDocumentFromReader(resp.Body)
			if err != nil {
				log.Fatal("Error loading HTTP response body. ", err)
			}
			for _, selector := range selectors {
				document.Find(selector).Each(func(i int, s *goquery.Selection) {

					if href, ok := s.Attr("href"); ok {
						links = append(links, href)
					} else {
						fmt.Printf("No link found %s\n", s.Text())
					}

				})
			}

		case err := <-c.Err:
			fmt.Printf("Error received: %s\n", err)
		}
	}

	return links
}
