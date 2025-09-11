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

type ComicBunchKaiOneshotExtractor struct {
	config site.Config
}

func NewComicBunchKaiOneshotExtractor(cfg site.Config) *ComicBunchKaiOneshotExtractor {
	return &ComicBunchKaiOneshotExtractor{
		config: cfg,
	}
}

func ComicBunchKaiOneshot() *generator.Generator {
	cfg := site.Config{
		Title:       "コミックバンチKai 読切作品",
		URL:         "https://comicbunch-kai.com/series#oneshot",
		DateLayout:  "2006年01月02日",
		Description: "None",
	}
	return generator.NewGenerator(cfg, NewComicBunchKaiOneshotExtractor(cfg))
}

func (extract ComicBunchKaiOneshotExtractor) ExtractItems(doc *goquery.Document) ([]site.Item, error) {

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

func (extract ComicBunchKaiOneshotExtractor) productURLs(doc *goquery.Document) ([]string, error) {

	var urls []string
	doc.Find("#oneshot .SeriesList_series_item_wrapper__XHj7m > div > a").Each(func(i int, sel *goquery.Selection) {
		if url, exist := sel.Attr("href"); exist && strings.Contains(url, "episode") {
			urls = append(urls, url)
		}
	})

	if len(urls) == 0 {
		return nil, fmt.Errorf("product URL not found")
	}

	return urls, nil
}

func (extract ComicBunchKaiOneshotExtractor) productItems(urls []string) ([]site.Item, error) {

	var items []site.Item
	for i, url := range urls {
		doc, err := utils.GetHtmlDoc(url)
		if err != nil {
			return nil, fmt.Errorf("failed to GetHtmlDoc: %w", err)
		}

		product := strings.TrimSpace(doc.Find("h1.series-header-title").First().Text())
		date := strings.TrimSpace(doc.Find("p.episode-header-date").First().Text())
		author := strings.TrimSpace(doc.Find("h2.series-header-author").First().Text())

		title := fmt.Sprintf("%s %s %s", product, author, date)
		link := url
		desc := "None"

		items = append(items, site.Item{Title: title, Link: link, Desc: desc, Date: date})

		// 10作品毎に1秒スリープ
		if i/9 == 0 {
			time.Sleep(1 * time.Second)
		}
	}

	return items, nil
}
