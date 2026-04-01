package service

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/snjrkn/generate-manga-feed/internal/site"
)

type matogrosso struct {
	config site.Config
}

func NewMatogrosso() site.Site {
	cfg := site.Config{
		Title:       "MATOGROSSO (マトグロッソ)",
		URL:         "https://matogrosso.jp",
		DateLayout:  "2006/01/02",
		Description: "None",
	}
	return site.Site{
		Config:    cfg,
		Extractor: &matogrosso{config: cfg},
	}
}

func (ext matogrosso) ExtractItems(doc *goquery.Document) ([]site.Item, error) {

	productItems, err := ext.productItems(doc)
	if err != nil {
		return nil, fmt.Errorf("failed to productItems: (Title='%v'): %w", ext.config.Title, err)
	}

	return productItems, nil
}

func (ext *matogrosso) productItems(doc *goquery.Document) ([]site.Item, error) {

	items := []site.Item{}
	doc.Find("div.serial_content").Each(func(i int, sel *goquery.Selection) {
		author := strings.TrimSpace(sel.Find("h2.entry-author").Text())
		product := strings.TrimSpace(sel.Find("h3.entry-title").Text())
		desc := strings.TrimSpace(sel.Find("div.asset-info > p").Text())
		date := strings.TrimSpace(sel.Find("p.published").Text())
		link := sel.Find("a").AttrOr("href", site.Config{}.URL)

		title := fmt.Sprintf("%s%s%s", author, product, date)
		date = strings.ReplaceAll(date, " 更新", "")

		items = append(items, site.Item{
			Title: title,
			Link:  link,
			Desc:  desc,
			Date:  date,
		})
	})

	if len(items) == 0 {
		return nil, fmt.Errorf("item not found")
	}

	return items, nil
}
