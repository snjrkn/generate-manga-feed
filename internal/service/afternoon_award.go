package service

import (
	"fmt"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/snjrkn/generate-manga-feed/internal/generator"
	"github.com/snjrkn/generate-manga-feed/internal/site"
	"github.com/snjrkn/generate-manga-feed/internal/utils"
)

type AfternoonAwardExtractor struct {
	config site.Config
}

func NewAfternoonAwardExtractor(cfg site.Config) *AfternoonAwardExtractor {
	return &AfternoonAwardExtractor{
		config: cfg,
	}
}

func AfternoonAward() *generator.Generator {
	cfg := site.Config{
		Title:       "アフタヌーン 四季賞",
		URL:         "https://afternoon.kodansha.co.jp/award/",
		DateLayout:  "2006年01月02日",
		Description: "None",
	}
	return generator.NewGenerator(cfg, NewAfternoonAwardExtractor(cfg))
}

func (extract AfternoonAwardExtractor) ExtractItems(doc *goquery.Document) ([]site.Item, error) {

	awardURLs, err := extract.awardURLs(doc)
	if err != nil {
		return nil, fmt.Errorf("failed to awardURLs: (Title='%v'): %w", extract.config.Title, err)
	}

	productURLs, err := extract.productURLs(awardURLs)
	if err != nil {
		return nil, fmt.Errorf("failed to productURLs: (Title='%v'): %w", extract.config.Title, err)
	}

	productItems, err := extract.productItems(productURLs)
	if err != nil {
		return nil, fmt.Errorf("failed to productItems: (Title='%v'): %w", extract.config.Title, err)
	}

	return productItems, nil
}

func (extract AfternoonAwardExtractor) awardURLs(doc *goquery.Document) ([]string, error) {

	var urls []string
	doc.Find(".mB50 a").Each(func(i int, sel *goquery.Selection) {
		if url, exist := sel.Attr("href"); exist {
			urls = append(urls, utils.GetFqdn(extract.config.URL)+"/"+url)
		}
	})

	if len(urls) == 0 {
		return nil, fmt.Errorf("award URL not found")
	}

	return urls, nil
}

func (extract AfternoonAwardExtractor) productURLs(awUrls []string) ([]string, error) {

	var urls []string
	for _, awUrl := range awUrls {
		doc, err := utils.FetchHtmlDoc(awUrl)
		if err != nil {
			return nil, fmt.Errorf("failed to FetchHtmlDoc: %w", err)
		}
		doc.Find(".viewOnWEB > a").Each(func(i int, sel *goquery.Selection) {
			if url, exist := sel.Attr("href"); exist && strings.Contains(url, "episode") {
				urls = append(urls, url)
			}
		})
	}

	if len(urls) == 0 {
		return nil, fmt.Errorf("product URL not found")
	}

	time.Sleep(3 * time.Second)

	return urls, nil
}

func (extract AfternoonAwardExtractor) productItems(urls []string) ([]site.Item, error) {

	var items []site.Item
	for i, url := range urls {
		doc, err := utils.FetchHtmlDoc(url)
		if err != nil {
			return nil, fmt.Errorf("failed to FetchHtmlDoc: %w", err)
		}

		product := strings.TrimSpace(doc.Find("h1.series-header-title").First().Text())
		date := strings.TrimSpace(doc.Find("p.episode-header-date").First().Text())
		author := strings.TrimSpace(doc.Find("h2.series-header-author").First().Text())
		desc := strings.TrimSpace(doc.Find("p.series-header-description").First().Text())

		title := fmt.Sprintf("%s %s %s", product, author, date)
		link := url

		items = append(items, site.Item{Title: title, Link: link, Desc: desc, Date: date})

		utils.ItemPerSleep(i, 9, 2)
	}

	return items, nil
}
