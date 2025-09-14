package service

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/snjrkn/generate-manga-feed/internal/generator"
	"github.com/snjrkn/generate-manga-feed/internal/site"
	"github.com/snjrkn/generate-manga-feed/internal/utils"
)

type KurageAwardExtractor struct {
	config site.Config
}

func NewKurageAwardExtractor(cfg site.Config) *KurageAwardExtractor {
	return &KurageAwardExtractor{
		config: cfg,
	}
}

func KurageAward() *generator.Generator {
	cfg := site.Config{
		Title:       "くらげバンチ漫画賞",
		URL:         "https://kuragebunch.com/info/award",
		DateLayout:  "2006年01月02日",
		Description: "None",
	}
	return generator.NewGenerator(cfg, NewKurageAwardExtractor(cfg))
}

func (extract KurageAwardExtractor) ExtractItems(doc *goquery.Document) ([]site.Item, error) {

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

func (extract KurageAwardExtractor) awardURLs(doc *goquery.Document) ([]string, error) {

	var urls []string
	doc.Find("ul.award-banner > li > a").Each(func(i int, sel *goquery.Selection) {
		if url, exist := sel.Attr("href"); exist {
			urls = append(urls, url)
		}
	})

	if len(urls) == 0 {
		return nil, fmt.Errorf("URL not found")
	}

	return urls, nil
}

func (extract KurageAwardExtractor) productURLs(awUrls []string) ([]string, error) {

	var urls []string
	for _, awUrl := range awUrls {
		doc, err := utils.FetchHtmlDoc(awUrl)
		if err != nil {
			return nil, fmt.Errorf("failed to FetchHtmlDoc: %w", err)
		}
		doc.Find(".hatenablog-entry > p > a").Each(func(i int, sel *goquery.Selection) {
			if url, exist := sel.Attr("href"); exist && strings.Contains(url, "episode") {
				urls = append(urls, url)
			}
		})
	}

	if len(urls) == 0 {
		return nil, fmt.Errorf("URL not found")
	}

	return urls, nil
}

func (extract KurageAwardExtractor) productItems(urls []string) ([]site.Item, error) {

	var items []site.Item
	for _, url := range urls {
		doc, err := utils.FetchHtmlDoc(url)
		if err != nil {
			return nil, fmt.Errorf("failed to FetchHtmlDoc: %w", err)
		}

		product := strings.TrimSpace(doc.Find("h1.series-header-title").First().Text())
		date := strings.TrimSpace(doc.Find("p.episode-header-date").First().Text())
		author := strings.TrimSpace(doc.Find("h2.series-header-author").First().Text())

		// 作品名の最後に作者名が付いている場合の対応
		if strings.HasSuffix(product, author) {
			product = strings.ReplaceAll(product, author, "")
		}
		title := fmt.Sprintf("%s %s %s", product, author, date)
		link := url
		desc := "None"

		items = append(items, site.Item{Title: title, Link: link, Desc: desc, Date: date})
	}

	return items, nil
}
