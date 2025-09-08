package service

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/snjrkn/generate-manga-feed/internal/generator"
	"github.com/snjrkn/generate-manga-feed/internal/site"
	"github.com/snjrkn/generate-manga-feed/internal/utils"
)

type ComicbunchkaiAwardExtractor struct {
	config site.Config
}

func NewComicbunchkaiAwardExtractor(cfg site.Config) *ComicbunchkaiAwardExtractor {
	return &ComicbunchkaiAwardExtractor{
		config: cfg,
	}
}

func ComicbunchkaiAward() *generator.Generator {
	cfg := site.Config{
		Title:       "コミックバンチKai 漫画賞",
		URL:         "https://comicbunch-kai.com/article/award",
		DateLayout:  "2006年01月02日",
		Description: "None",
	}
	return generator.NewGenerator(cfg, NewComicbunchkaiAwardExtractor(cfg))
}

func (extract ComicbunchkaiAwardExtractor) ExtractItems(doc *goquery.Document) ([]site.Item, error) {

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

func (extract ComicbunchkaiAwardExtractor) awardURLs(doc *goquery.Document) ([]string, error) {

	var urls []string
	doc.Find("ul.award-banner > li > a").Each(func(i int, sel *goquery.Selection) {
		if url, exist := sel.Attr("href"); exist {
			urls = append(urls, url)
		}
	})

	if len(urls) == 0 {
		return nil, fmt.Errorf("award URL not found")
	}

	return urls, nil
}

func (extract ComicbunchkaiAwardExtractor) productURLs(awUrls []string) ([]string, error) {

	var urls []string
	for _, awUrl := range awUrls {
		doc, err := utils.GetHtmlDoc(awUrl)
		if err != nil {
			return nil, fmt.Errorf("failed to GetHtmlDoc: %w", err)
		}
		doc.Find(".hatenablog-entry > p > a").Each(func(i int, sel *goquery.Selection) {
			if url, exist := sel.Attr("href"); exist && strings.Contains(url, "episode") {
				urls = append(urls, url)
			}
		})
	}

	if len(urls) == 0 {
		return nil, fmt.Errorf("product URL not found")
	}

	return urls, nil
}

func (extract ComicbunchkaiAwardExtractor) productItems(urls []string) ([]site.Item, error) {

	var items []site.Item
	for _, url := range urls {
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
	}

	return items, nil
}
