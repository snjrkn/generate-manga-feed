package service

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/mmcdole/gofeed"

	"github.com/snjrkn/generate-manga-feed/internal/site"
	"github.com/snjrkn/generate-manga-feed/internal/util"
)

type comiplexExtractor struct {
	config site.Config
}

func newComiplexExtractor(cfg site.Config) *comiplexExtractor {
	return &comiplexExtractor{
		config: cfg,
	}
}

func ComiplexOneshot() site.Site {
	cfg := site.Config{
		Title:       "コミプレ 読切作品",
		URL:         "https://viewer.heros-web.com/series/oneshot",
		Description: "None",
		DateLayout:  "2006/01/02",
	}
	return site.Site{
		Config:    cfg,
		Extractor: newComiplexExtractor(cfg),
	}
}

func (ext comiplexExtractor) ExtractItems(doc *goquery.Document) ([]site.Item, error) {

	oneshotURLs, err := ext.oneshotURLs(doc)
	if err != nil {
		return nil, fmt.Errorf("failed to oneshotURLs: %w", err)
	}

	rssURLs, err := ext.rssURLs(oneshotURLs)
	if err != nil {
		return nil, fmt.Errorf("failed to rssURLs: %w", err)
	}

	productItems, err := ext.productItems(rssURLs)
	if err != nil {
		return nil, fmt.Errorf("failed to productItems: (Title='%v'): %w", ext.config.Title, err)
	}

	return productItems, nil
}

func (ext comiplexExtractor) oneshotURLs(doc *goquery.Document) ([]string, error) {

	var urls []string
	doc.Find("a.series-item-updated-link").Each(func(i int, sel *goquery.Selection) {
		if url, exists := sel.Attr("href"); exists && strings.Contains(url, "episode") {
			urls = append(urls, url)
		}
	})

	if len(urls) == 0 {
		return nil, fmt.Errorf("URL not found")
	}

	return urls, nil
}

func (ext comiplexExtractor) rssURLs(urls []string) ([]string, error) {

	var rsUrls []string
	for i := range urls {
		doc, err := util.FetchHtmlDoc(urls[i])
		if err != nil {
			return nil, fmt.Errorf("failed to FetchHtmlDoc: %w", err)
		}

		url, exists := doc.Find("dd.rss > a").First().Attr("href")
		if exists {
			rsUrls = append(rsUrls, url)
		}
	}

	if len(rsUrls) == 0 {
		return nil, fmt.Errorf("URL not found")
	}

	return rsUrls, nil
}

func (ext comiplexExtractor) productItems(urls []string) ([]site.Item, error) {

	var items []site.Item
	for i := range urls {
		feed, err := gofeed.NewParser().ParseURL(urls[i])
		if err != nil {
			return nil, fmt.Errorf("failed to parsing feed: %w", err)
		}

		for i := range feed.Items {
			items = append(items, site.Item{
				Title:       feed.Items[i].Title,
				Link:        feed.Items[i].Link,
				Desc:        feed.Items[i].Description,
				CreatedDate: feed.Items[i].PublishedParsed.UTC(),
			})
		}
	}

	if len(items) == 0 {
		return nil, fmt.Errorf("item not found")
	}

	return items, nil
}
