package service

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/snjrkn/generate-manga-feed/internal/site"
	"github.com/snjrkn/generate-manga-feed/internal/util"
)

type comicBoostRensai struct {
	config site.Config
}

func NewComicBoostRensai(contentId string) site.Site {
	cfg := site.Config{
		Title:       "comicブースト【連載】",
		URL:         "https://comic-boost.com/content/" + contentId,
		DateLayout:  "2006/01/02",
		Description: "None",
	}

	ext := &comicBoostRensai{config: cfg}
	if err := ext.contentInfo(&cfg); err != nil {
		fmt.Printf("failed to contentInfo: %v\n", err)
	}

	return site.Site{
		Config:    cfg,
		Extractor: ext,
	}
}

func (ext comicBoostRensai) ExtractItems(doc *goquery.Document) ([]site.Item, error) {

	contentItems, err := ext.contentItems([]string{ext.config.URL})
	if err != nil {
		return nil, fmt.Errorf("failed to contentItems: (Title='%v'): %w", ext.config.Title, err)
	}

	return contentItems, nil
}

func (ext comicBoostRensai) contentInfo(cfg *site.Config) error {

	doc, err := util.FetchHtmlDoc(cfg.URL)
	if err != nil {
		return fmt.Errorf("failed to FetchHtmlDoc: (Site='%v'): %w", cfg.Title, err)
	}

	cfg.Description = strings.TrimSpace(doc.Find("p.comic-description-text").First().Text())
	cfg.Title += strings.TrimSpace(doc.Find("h1.comic-title").First().Text())

	return nil
}

func (ext comicBoostRensai) contentItems(contentURLs []string) ([]site.Item, error) {

	domain := util.GetFqdn(ext.config.URL)

	var items []site.Item
	processedIndex := 0 // contentURLsをキューとするインデックス
	for processedIndex < len(contentURLs) {

		contentURL := contentURLs[processedIndex]
		processedIndex++

		doc, err := util.FetchHtmlDoc(contentURL)
		if err != nil {
			return nil, fmt.Errorf("failed to FetchHtmlDoc: %w", err)
		}

		// 次のページがあればcontentURLsの末尾に追加
		nextPageURL, found := doc.Find(".to-next > a").First().Attr("href")
		if found && nextPageURL != "javascript:void(0);" {
			nextPageURL = domain + nextPageURL
			contentURLs = append(contentURLs, nextPageURL)
		}

		items = append(items, ext.extItems(doc, domain)...)

		util.ItemPerSleep(processedIndex, 9, 1)
	}

	if len(items) == 0 {
		return nil, fmt.Errorf("item not found")
	}

	return items, nil
}

func (ext comicBoostRensai) extItems(doc *goquery.Document, domain string) []site.Item {

	desc := "None"

	var items []site.Item
	doc.Find("a.book-product-list-item").Each(func(i int, sel *goquery.Selection) {
		episode := strings.TrimSpace(sel.Find("h4.title").Text())
		date := strings.TrimSpace(sel.Find("p.update-date").Text())
		link := sel.AttrOr("href", "")
		link = strings.Split(link, "?")[0]
		title := fmt.Sprintf("%s %s", episode, date)
		coin := strings.TrimSpace(sel.Find("div.book-product-list-item-meta-wrapper div.left div").Text())
		link = domain + link

		items = append(items, site.Item{Title: "【" + coin + "】" + title, Link: link, Desc: desc, Date: date})
	})

	return items
}
