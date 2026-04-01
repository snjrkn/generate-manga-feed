package service

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/snjrkn/generate-manga-feed/internal/site"
	"github.com/snjrkn/generate-manga-feed/internal/util"
)

type comicBunchKaiOneshot struct {
	config site.Config
}

func NewComicBunchKaiOneshot() site.Site {
	cfg := site.Config{
		Title:       "コミックバンチKai 読切作品",
		URL:         "https://comicbunch-kai.com/series#oneshot",
		DateLayout:  "2006年01月02日",
		Description: "None",
	}
	return site.Site{
		Config:    cfg,
		Extractor: &comicBunchKaiOneshot{config: cfg},
	}
}

func (ext comicBunchKaiOneshot) ExtractItems(doc *goquery.Document) ([]site.Item, error) {

	productURLs, err := ext.productURLs(doc)
	if err != nil {
		return nil, fmt.Errorf("failed to productURLs: %w", err)
	}

	productItems, err := ext.productItems(productURLs)
	if err != nil {
		return nil, fmt.Errorf("failed to productItems: (Title='%v'): %w", ext.config.Title, err)
	}

	return productItems, nil
}

func (ext comicBunchKaiOneshot) productURLs(doc *goquery.Document) ([]string, error) {

	var urls []string
	doc.Find("#oneshot .SeriesList_series_item_wrapper__XHj7m > div > a").Each(func(i int, sel *goquery.Selection) {
		if url, exists := sel.Attr("href"); exists && strings.Contains(url, "episode") {
			urls = append(urls, url)
		}
	})

	if len(urls) == 0 {
		return nil, fmt.Errorf("URL not found")
	}

	return urls, nil
}

func (ext comicBunchKaiOneshot) productItems(urls []string) ([]site.Item, error) {

	var items []site.Item
	for i := range urls {
		doc, err := util.FetchHtmlDoc(urls[i])
		if err != nil {
			return nil, fmt.Errorf("failed to FetchHtmlDoc: %w", err)
		}

		product := strings.TrimSpace(doc.Find("h1.series-header-title").First().Text())
		date := strings.TrimSpace(doc.Find("p.episode-header-date").First().Text())
		author := strings.TrimSpace(doc.Find("h2.series-header-author").First().Text())
		desc := strings.TrimSpace(doc.Find("p.series-header-description").First().Text())

		title := fmt.Sprintf("%s %s %s", product, author, date)
		link := urls[i]

		items = append(items, site.Item{
			Title: title,
			Link:  link,
			Desc:  desc,
			Date:  date,
		})

		util.ItemPerSleep(i, 9, 1)
	}

	if len(items) == 0 {
		return nil, fmt.Errorf("item not found")
	}

	return items, nil
}
