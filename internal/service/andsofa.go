package service

import (
	"fmt"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/snjrkn/generate-manga-feed/internal/generator"
	"github.com/snjrkn/generate-manga-feed/internal/site"
)

type AndSofaExtractor struct {
	config site.Config
}

func NewAndSofaExtractor(cfg site.Config) *AndSofaExtractor {
	return &AndSofaExtractor{
		config: cfg,
	}
}

func AndSofa() *generator.Generator {
	cfg := site.Config{
		Title:       "＆Sofa (アンドソファ)",
		URL:         "https://andsofa.com",
		DateLayout:  "2006年1月2日",
		Description: "None",
	}
	return generator.NewGenerator(cfg, NewAndSofaExtractor(cfg))
}

func (extract AndSofaExtractor) ExtractItems(doc *goquery.Document) ([]site.Item, error) {

	productItems, err := extract.productItems(doc)
	if err != nil {
		return nil, fmt.Errorf("failed to productItems: (Title='%v'): %w", extract.config.Title, err)
	}

	return productItems, nil
}

func (extract AndSofaExtractor) productItems(doc *goquery.Document) ([]site.Item, error) {

	// 最新更新分の日付
	newestDate := strings.TrimSpace(doc.Find("div.updated-episodes-date").First().Text())
	newestDate = fmt.Sprint(time.Now().Local().Year()) + "年" + strings.ReplaceAll(newestDate, "日月", "日")

	// 前回更新分の日付
	t, err := time.Parse(extract.config.DateLayout, newestDate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse newest date: (Date='%v'): %w", newestDate, err)
	}
	lastDate := t.Add(-7 * 24 * time.Hour).Format(extract.config.DateLayout)

	items := []site.Item{}

	// 最新更新分のアイテム
	newestItems := extract.extractItems(doc.Find("div.updated-episodes-wrapper").First(), newestDate)
	items = append(items, newestItems...)

	// 前回更新分のアイテム
	lastItems := extract.extractItems(doc.Find("div.updated-episodes-wrapper").Next(), lastDate)
	items = append(items, lastItems...)

	if len(items) == 0 {
		return nil, fmt.Errorf("item not found")
	}

	return items, nil
}

func (extract AndSofaExtractor) extractItems(sel *goquery.Selection, date string) []site.Item {

	var items []site.Item
	sel.Find("li.updated-episodes-item").Each(func(i int, sel *goquery.Selection) {
		episode := strings.TrimSpace(sel.Find("p.episode-title").Text())
		product := strings.TrimSpace(sel.Find("div.updated-episodes-item-text h4").Text())
		author := strings.TrimSpace(sel.Find("div.updated-episodes-item-text h5").Text())
		link := sel.Find("a").AttrOr("href", site.Config{}.URL)
		desc := strings.TrimSpace(sel.Find("p.description").Text())

		title := fmt.Sprintf("%s (月) %s %s %s", date, episode, product, author)
		desc = strings.ReplaceAll(desc, "&", "＆")

		items = append(items, site.Item{Title: title, Link: link, Desc: desc, Date: date})
	})

	return items
}
