package service

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/snjrkn/generate-manga-feed/internal/site"
)

type comicDays struct {
	config site.Config
}

func NewComicDaysNewcomer() site.Site {
	cfg := site.Config{
		Title:       "コミックDAYS 新人賞",
		URL:         "https://comic-days.com/newcomer",
		DateLayout:  "2006年01月02日",
		Description: "None",
	}
	return site.Site{
		Config:    cfg,
		Extractor: &comicDays{config: cfg},
	}
}

func NewComicDaysOneshot() site.Site {
	cfg := site.Config{
		Title:       "コミックDAYS 読み切り",
		URL:         "https://comic-days.com/oneshot",
		DateLayout:  "2006年01月02日",
		Description: "None",
	}
	return site.Site{
		Config:    cfg,
		Extractor: &comicDays{config: cfg},
	}
}

func (ext comicDays) ExtractItems(doc *goquery.Document) ([]site.Item, error) {

	productItems, err := ext.productItems(doc)
	if err != nil {
		return nil, fmt.Errorf("failed to productItems: (Title='%v'): %w", ext.config.Title, err)
	}

	return productItems, nil
}

func (ext *comicDays) productItems(doc *goquery.Document) ([]site.Item, error) {

	items := []site.Item{}
	doc.Find("li.yomikiri-item-box").Each(func(i int, sel *goquery.Selection) {
		date := strings.TrimSpace(sel.Find("span.yomikiri-label-date").Text())
		product := strings.TrimSpace(sel.Find("div.yomikiri-link-title h4").Text())
		author := strings.TrimSpace(sel.Find("div.yomikiri-link-title h5").Text())
		link := sel.Find("a.yomikiri-link").AttrOr("href", ext.config.URL)

		title := fmt.Sprintf("%s %s %s", date, product, author)
		desc := "None"

		items = append(items, site.Item{Title: title, Link: link, Desc: desc, Date: date})
	})

	if len(items) == 0 {
		return nil, fmt.Errorf("item not found")
	}

	return items, nil
}
