package service

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/snjrkn/generate-manga-feed/internal/generator"
	"github.com/snjrkn/generate-manga-feed/internal/site"
)

type KurageFarmExtractor struct {
	config site.Config
}

func NewKurageFarmExtractor(cfg site.Config) *KurageFarmExtractor {
	return &KurageFarmExtractor{
		config: cfg,
	}
}

func KurageFarm() *generator.Generator {

	cfg := site.Config{
		Title:       "くらげファーム",
		URL:         "https://kuragebunch.com/farm",
		DateLayout:  "2006年1月2日",
		Description: "None",
	}

	return generator.NewGenerator(cfg, NewKurageFarmExtractor(cfg))
}

func (extract KurageFarmExtractor) ExtractItems(doc *goquery.Document) ([]site.Item, error) {

	productItems, err := extract.productItems(doc)
	if err != nil {
		return nil, fmt.Errorf("failed to productItems: (Title='%v'): %w", extract.config.Title, err)
	}

	return productItems, nil
}

func (extract KurageFarmExtractor) productItems(doc *goquery.Document) ([]site.Item, error) {

	var items []site.Item
	doc.Find("li.yomikiri-item-box").Each(func(i int, sel *goquery.Selection) {
		date := strings.TrimSpace(sel.Find("span.yomikiri-label-date").Text())
		product := strings.TrimSpace(sel.Find("div.yomikiri-link-title h4").Text())
		author := strings.TrimSpace(sel.Find("div.yomikiri-link-title h5").Text())
		link := sel.Find("a.yomikiri-link").AttrOr("href", "")

		title := fmt.Sprintf("%s %s %s", date, product, author)
		desc := "None" // 各作品のページに詳細情報あり

		items = append(items, site.Item{Title: title, Link: link, Desc: desc, Date: date})
	})

	return items, nil
}
