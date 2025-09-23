package service

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/mmcdole/gofeed"

	"github.com/snjrkn/generate-manga-feed/internal/generator"
	"github.com/snjrkn/generate-manga-feed/internal/site"
	"github.com/snjrkn/generate-manga-feed/internal/util"
)

type ComicActionExtractor struct {
	config site.Config
}

func NewComicActionExtractor(cfg site.Config) *ComicActionExtractor {
	return &ComicActionExtractor{
		config: cfg,
	}
}

func ComicActionOneshot() *generator.Generator {
	cfg := site.Config{
		Title:       "webアクション 読切作品",
		URL:         "https://comic-action.com/series/oneshot",
		Description: "None",
		DateLayout:  "2006/01/02",
	}
	return generator.NewGenerator(cfg, NewComicActionExtractor(cfg))
}

func (extract ComicActionExtractor) ExtractItems(doc *goquery.Document) ([]site.Item, error) {

	oneshotURLs, err := extract.oneshotURLs(doc)
	if err != nil {
		return nil, fmt.Errorf("failed to oneshotURLs: (Title='%v'): %w", extract.config.Title, err)
	}

	rssURLs, productDescs, err := extract.rssURLsAndDescs(oneshotURLs)
	if err != nil {
		return nil, fmt.Errorf("failed to rssURLs: (Title='%v'): %w", extract.config.Title, err)
	}

	productItems, err := extract.productItems(rssURLs, productDescs)
	if err != nil {
		return nil, fmt.Errorf("failed to productItems: (Title='%v'): %w", extract.config.Title, err)
	}

	return productItems, nil
}

func (extract ComicActionExtractor) oneshotURLs(doc *goquery.Document) ([]string, error) {

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

func (extract ComicActionExtractor) rssURLsAndDescs(urls []string) (rsUrls, descs []string, err error) {

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

func (extract ComicActionExtractor) productItems(urls, descs []string) ([]site.Item, error) {

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
