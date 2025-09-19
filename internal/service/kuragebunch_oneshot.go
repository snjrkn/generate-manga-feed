package service

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/snjrkn/generate-manga-feed/internal/generator"
	"github.com/snjrkn/generate-manga-feed/internal/site"
	"github.com/snjrkn/generate-manga-feed/internal/util"
)

type KurageBunchOneshotExtractor struct {
	config site.Config
}

func NewKurageBunchOneshotExtractor(cfg site.Config) *KurageBunchOneshotExtractor {
	return &KurageBunchOneshotExtractor{
		config: cfg,
	}
}

func KurageBunchOneshot() *generator.Generator {
	cfg := site.Config{
		Title:       "くらげバンチ 読切",
		URL:         "https://kuragebunch.com/series/oneshot",
		DateLayout:  "2006年01月02日",
		Description: "None",
	}
	return generator.NewGenerator(cfg, NewKurageBunchOneshotExtractor(cfg))
}

func (extract KurageBunchOneshotExtractor) ExtractItems(doc *goquery.Document) ([]site.Item, error) {

	productURLs, err := extract.productURLs(doc)
	if err != nil {
		return nil, fmt.Errorf("failed to productURLs: (Title='%v'): %w", extract.config.Title, err)
	}

	productItems, err := extract.productItems(productURLs)
	if err != nil {
		return nil, fmt.Errorf("failed to productItems: (Title='%v'): %w", extract.config.Title, err)
	}

	return productItems, nil
}

func (extract KurageBunchOneshotExtractor) productURLs(doc *goquery.Document) ([]string, error) {

	var urls []string
	doc.Find(".item-box > a").Each(func(i int, sel *goquery.Selection) {
		if url, exist := sel.Attr("href"); exist && strings.Contains(url, "episode") {
			urls = append(urls, url)
		}
	})

	if len(urls) == 0 {
		return nil, fmt.Errorf("URL not found")
	}

	return urls, nil
}

func (extract KurageBunchOneshotExtractor) productItems(urls []string) ([]site.Item, error) {

	var items []site.Item
	for i, url := range urls {
		doc, err := util.FetchHtmlDoc(url)
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

		util.ItemPerSleep(i, 9, 1)
	}

	if len(items) == 0 {
		return nil, fmt.Errorf("item not found")
	}

	return items, nil
}
