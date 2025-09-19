package service

import (
	"fmt"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/snjrkn/generate-manga-feed/internal/generator"
	"github.com/snjrkn/generate-manga-feed/internal/site"
	"github.com/snjrkn/generate-manga-feed/internal/util"
)

type ComicBunchKaiAwardExtractor struct {
	config site.Config
}

func NewComicBunchKaiAwardExtractor(cfg site.Config) *ComicBunchKaiAwardExtractor {
	return &ComicBunchKaiAwardExtractor{
		config: cfg,
	}
}

func ComicBunchKaiAward() *generator.Generator {
	cfg := site.Config{
		Title: "コミックバンチKai 漫画賞",
		// URL:         "https://comicbunch-kai.com/article/award",
		URL:         "https://comicbunch-kai.com/article/archive/category/%E6%BC%AB%E7%94%BB%E8%B3%9E_%E7%99%BA%E8%A1%A8",
		DateLayout:  "2006年01月02日",
		Description: "None",
	}
	return generator.NewGenerator(cfg, NewComicBunchKaiAwardExtractor(cfg))
}

func (extract ComicBunchKaiAwardExtractor) ExtractItems(doc *goquery.Document) ([]site.Item, error) {

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

func (extract ComicBunchKaiAwardExtractor) awardURLs(doc *goquery.Document) ([]string, error) {

	var urls []string
	doc.Find("a.entry-thumb-link").Each(func(i int, sel *goquery.Selection) {
		if url, exist := sel.Attr("href"); exist {
			urls = append(urls, url)
		}
	})

	if len(urls) == 0 {
		return nil, fmt.Errorf("URL not found")
	}

	return urls, nil
}

func (extract ComicBunchKaiAwardExtractor) productURLs(awUrls []string) ([]string, error) {

	var urls []string
	for _, awUrl := range awUrls {
		doc, err := util.FetchHtmlDoc(awUrl)
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

	time.Sleep(1 * time.Second)

	return urls, nil
}

func (extract ComicBunchKaiAwardExtractor) productItems(urls []string) ([]site.Item, error) {

	var items []site.Item
	for _, url := range urls {
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
	}

	if len(items) == 0 {
		return nil, fmt.Errorf("item not found")
	}

	return items, nil
}
