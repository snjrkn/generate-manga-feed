package service

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/snjrkn/generate-manga-feed/internal/site"
	"github.com/snjrkn/generate-manga-feed/internal/util"
)

type toti struct {
	config site.Config
}

func NewToti() site.Site {
	cfg := site.Config{
		Title:       "トーチ",
		URL:         "https://to-ti.in/product",
		DateLayout:  "2006/01/02",
		Description: "None",
	}
	return site.Site{
		Config:    cfg,
		Extractor: &toti{config: cfg},
	}
}

func (ext toti) ExtractItems(doc *goquery.Document) ([]site.Item, error) {

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

func (ext toti) productURLs(doc *goquery.Document) ([]string, error) {

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

func (ext toti) productItems(productURLs []string) ([]site.Item, error) {

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
		// 日付にタイトルが含まれる場合の対応
		titleDate = regexp.MustCompile("^.*'").ReplaceAllString(date, "'")
		titleDate = strings.ReplaceAll(titleDate, "'", "’")
		title := fmt.Sprintf("%s %s %s", titleDate, product, story)

		// 日付にタイトルが含まれる場合の対応
		date = regexp.MustCompile("^.*'").ReplaceAllString(date, "'")

		date = "20" + strings.ReplaceAll(date, "'", "")
		date = strings.ReplaceAll(date, " UPDATE", "")

		items = append(items, site.Item{Title: title, Link: link, Desc: desc, Date: date})
	}

	if len(items) == 0 {
		return nil, fmt.Errorf("item not found")
	}

	return items, nil
}
