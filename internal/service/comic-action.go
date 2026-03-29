package service

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/mmcdole/gofeed"

	"github.com/snjrkn/generate-manga-feed/internal/site"
	"github.com/snjrkn/generate-manga-feed/internal/util"
)

type comicActionExtractor struct {
	config site.Config
}

func newComicActionExtractor(cfg site.Config) *comicActionExtractor {
	return &comicActionExtractor{
		config: cfg,
	}
}

func ComicActionOneshot() site.Site {
	cfg := site.Config{
		Title:       "webアクション 読切作品",
		URL:         "https://comic-action.com/series/oneshot",
		Description: "None",
		DateLayout:  "2006/01/02",
	}
	return site.Site{
		Config:    cfg,
		Extractor: newComicActionExtractor(cfg),
	}
}

func (ext comicActionExtractor) ExtractItems(doc *goquery.Document) ([]site.Item, error) {

	oneshotURLs, err := ext.oneshotURLs(doc)
	if err != nil {
		return nil, fmt.Errorf("failed to oneshotURLs: %w", err)
	}

	rssURLs, productDescs, err := ext.rssURLsAndDescs(oneshotURLs)
	if err != nil {
		return nil, fmt.Errorf("failed to rssURLs: %w", err)
	}

	productItems, err := ext.productItems(rssURLs, productDescs)
	if err != nil {
		return nil, fmt.Errorf("failed to productItems: (Title='%v'): %w", ext.config.Title, err)
	}

	return productItems, nil
}

func (ext comicActionExtractor) oneshotURLs(doc *goquery.Document) ([]string, error) {

	var urls []string
	doc.Find("#oneshot a.SeriesListItem_thumb_link__kvQJN").Each(func(i int, sel *goquery.Selection) {
		if url, exists := sel.Attr("href"); exists && strings.Contains(url, "episode") {
			urls = append(urls, url)
		}
	})

	if len(urls) == 0 {
		return nil, fmt.Errorf("URL not found")
	}

	return urls, nil
}

func (ext comicActionExtractor) rssURLsAndDescs(urls []string) (rsUrls, descs []string, err error) {

	for i := range urls {
		doc, err := util.FetchHtmlDoc(urls[i])
		if err != nil {
			return nil, nil, fmt.Errorf("failed to FetchHtmlDoc: %w", err)
		}
		desc := strings.TrimSpace(doc.Find("p.series-header-description").First().Text())
		descs = append(descs, desc)

		url, exists := doc.Find("dd.rss > a").First().Attr("href")
		if exists {
			rsUrls = append(rsUrls, url)
		}

		util.ItemPerSleep(i, 9, 1)
	}

	if len(rsUrls) == 0 {
		return nil, nil, fmt.Errorf("URL not found")
	}

	return rsUrls, descs, nil
}

func (ext comicActionExtractor) productItems(urls, descs []string) ([]site.Item, error) {

	var items []site.Item
	for i := range urls {
		feed, err := gofeed.NewParser().ParseURL(urls[i])
		if err != nil {
			return nil, fmt.Errorf("failed to parsing feed: %w", err)
		}
		desc := descs[i]

		for i := range feed.Items {
			items = append(items, site.Item{
				Title:       feed.Items[i].Title,
				Link:        feed.Items[i].Link,
				Desc:        desc,
				CreatedDate: feed.Items[i].PublishedParsed.UTC(),
			})
		}
		util.ItemPerSleep(i, 9, 1)
	}

	if len(items) == 0 {
		return nil, fmt.Errorf("item not found")
	}

	return items, nil
}
