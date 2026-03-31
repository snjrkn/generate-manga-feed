package service

import (
	"fmt"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/snjrkn/generate-manga-feed/internal/site"
)

type andSofa struct {
	config site.Config
}

func NewAndSofa() site.Site {
	cfg := site.Config{
		Title:       "＆Sofa (アンドソファ)",
		URL:         "https://andsofa.com",
		DateLayout:  "2006年1月2日",
		Description: "None",
	}
	return site.Site{
		Config:    cfg,
		Extractor: &andSofa{config: cfg},
	}
}

func (ext andSofa) ExtractItems(doc *goquery.Document) ([]site.Item, error) {

	productItems, err := ext.productItems(doc)
	if err != nil {
		return nil, fmt.Errorf("failed to productItems: (Title='%v'): %w", ext.config.Title, err)
	}

	return productItems, nil
}

func (ext andSofa) productItems(doc *goquery.Document) ([]site.Item, error) {

	// 最新更新分の日付
	newestDate := strings.TrimSpace(doc.Find("div.updated-episodes-date").First().Text())
	newestDate = fmt.Sprint(time.Now().Local().Year()) + "年" + strings.ReplaceAll(newestDate, "日月", "日")

	// 前回更新分の日付
	t, err := time.Parse(ext.config.DateLayout, newestDate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse newest date: (Date='%v'): %w", newestDate, err)
	}
	lastDate := t.Add(-7 * 24 * time.Hour).Format(ext.config.DateLayout)

	items := []site.Item{}

	// 最新更新分のアイテム
	newestItems := ext.extItems(doc.Find("div.updated-episodes-wrapper").First(), newestDate)
	items = append(items, newestItems...)

	// 前回更新分のアイテム
	lastItems := ext.extItems(doc.Find("div.updated-episodes-wrapper").Next(), lastDate)
	items = append(items, lastItems...)

	if len(items) == 0 {
		return nil, fmt.Errorf("item not found")
	}

	return items, nil
}

func (ext andSofa) extItems(sel *goquery.Selection, date string) []site.Item {

	var items []site.Item
	sel.Find("li.updated-episodes-item").Each(func(i int, sel *goquery.Selection) {
		episode := strings.TrimSpace(sel.Find("p.episode-title").Text())
		product := strings.TrimSpace(sel.Find("div.updated-episodes-item-text h4").Text())
		author := strings.TrimSpace(sel.Find("div.updated-episodes-item-text h5").Text())
		link := sel.Find("a").AttrOr("href", ext.config.URL)
		desc := strings.TrimSpace(sel.Find("p.description").Text())

		title := fmt.Sprintf("%s (月) %s %s %s", date, episode, product, author)
		desc = strings.ReplaceAll(desc, "&", "＆")

		items = append(items, site.Item{Title: title, Link: link, Desc: desc, Date: date})
	})

	return items
}
