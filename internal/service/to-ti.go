package service

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/snjrkn/generate-manga-feed/internal/generator"
	"github.com/snjrkn/generate-manga-feed/internal/site"
	"github.com/snjrkn/generate-manga-feed/internal/util"
)

type TotiExtractor struct {
	config site.Config
}

func NewTotiExtractor(cfg site.Config) *TotiExtractor {
	return &TotiExtractor{
		config: cfg,
	}
}

func Toti() *generator.Generator {
	cfg := site.Config{
		Title:       "トーチ",
		URL:         "https://to-ti.in/product",
		DateLayout:  "2006/01/02",
		Description: "None",
	}
	return generator.NewGenerator(cfg, NewTotiExtractor(cfg))
}

func (extract TotiExtractor) ExtractItems(doc *goquery.Document) ([]site.Item, error) {

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

func (extract TotiExtractor) productURLs(doc *goquery.Document) ([]string, error) {

	var urls []string
	doc.Find("article a").Each(func(i int, sel *goquery.Selection) {
		if url, exists := sel.Attr("href"); exists {
			urls = append(urls, url)
		}
	})

	if len(urls) == 0 {
		return nil, fmt.Errorf("URL not found")
	}

	return urls, nil
}

func (extract TotiExtractor) productItems(productURLs []string) ([]site.Item, error) {

	items := []site.Item{}
	for i := range productURLs {
		doc, err := util.FetchHtmlDoc(productURLs[i])
		if err != nil {
			return nil, fmt.Errorf("failed to FetchHtmlDoc: %w", err)
		}

		date := strings.TrimSpace(doc.Find("time").First().Text())
		product := strings.TrimSpace(doc.Find("header > h3").Text())
		desc := strings.TrimSpace(doc.Find("header > p").Text())
		story := strings.TrimSpace(doc.Find("p.next > a >  span").Text())
		link, exists := doc.Find("p.next a").Attr("href")
		// 第1話が公開された場合の対応
		if !exists {
			link = doc.Find("p.prev a").AttrOr("href", productURLs[i])
		}

		titleDate := date
		titleDate = strings.ReplaceAll(titleDate, "'", "’")
		title := fmt.Sprintf("%s %s %s", titleDate, product, story)
		date = "20" + strings.ReplaceAll(date, "'", "")
		date = strings.ReplaceAll(date, " UPDATE", "")

		items = append(items, site.Item{Title: title, Link: link, Desc: desc, Date: date})
	}

	if len(items) == 0 {
		return nil, fmt.Errorf("item not found")
	}

	return items, nil
}
