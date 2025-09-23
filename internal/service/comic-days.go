package service

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/snjrkn/generate-manga-feed/internal/generator"
	"github.com/snjrkn/generate-manga-feed/internal/site"
)

type ComicDaysExtractor struct {
	config site.Config
}

func NewComicDaysExtractor(cfg site.Config) *ComicDaysExtractor {
	return &ComicDaysExtractor{
		config: cfg,
	}
}

func ComicDaysNewcomer() *generator.Generator {
	cfg := site.Config{
		Title:       "コミックDAYS 新人賞",
		URL:         "https://comic-days.com/newcomer",
		DateLayout:  "2006年01月02日",
		Description: "None",
	}
	return generator.NewGenerator(cfg, NewComicDaysExtractor(cfg))
}

func ComicDaysOneshot() *generator.Generator {
	cfg := site.Config{
		Title:       "コミックDAYS 読み切り",
		URL:         "https://comic-days.com/oneshot",
		DateLayout:  "2006年01月02日",
		Description: "None",
	}
	return generator.NewGenerator(cfg, NewComicDaysExtractor(cfg))
}

func (extract ComicDaysExtractor) ExtractItems(doc *goquery.Document) ([]site.Item, error) {

	productItems, err := extract.productItems(doc)
	if err != nil {
		return nil, fmt.Errorf("failed to productItems: (Title='%v'): %w", extract.config.Title, err)
	}

	return productItems, nil
}

func (extract *ComicDaysExtractor) productItems(doc *goquery.Document) ([]site.Item, error) {

	items := []site.Item{}
	doc.Find("li.yomikiri-item-box").Each(func(i int, sel *goquery.Selection) {
		date := strings.TrimSpace(sel.Find("span.yomikiri-label-date").Text())
		product := strings.TrimSpace(sel.Find("div.yomikiri-link-title h4").Text())
		author := strings.TrimSpace(sel.Find("div.yomikiri-link-title h5").Text())
		link := sel.Find("a.yomikiri-link").AttrOr("href", site.Config{}.URL)

		title := fmt.Sprintf("%s %s %s", date, product, author)
		desc := "None"

		items = append(items, site.Item{Title: title, Link: link, Desc: desc, Date: date})
	})

	if len(items) == 0 {
		return nil, fmt.Errorf("item not found")
	}

	return items, nil
}
