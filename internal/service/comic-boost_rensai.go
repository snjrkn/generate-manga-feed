package service

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/snjrkn/generate-manga-feed/internal/generator"
	"github.com/snjrkn/generate-manga-feed/internal/site"
	"github.com/snjrkn/generate-manga-feed/internal/util"
)

type ComicBoostRensaiExtractor struct {
	config site.Config
}

func NewComicBoostRensaiExtractor(cfg site.Config) *ComicBoostRensaiExtractor {
	return &ComicBoostRensaiExtractor{
		config: cfg,
	}
}

func ComicBoostRensai(productId string) *generator.Generator {
	cfg := site.Config{
		Title:       "comicブースト【連載】",
		URL:         "https://comic-boost.com/content/" + productId,
		DateLayout:  "2006/01/02",
		Description: "None",
	}

	extract := NewComicBoostRensaiExtractor(cfg)
	if err := extract.productInfo(&cfg); err != nil {
		fmt.Printf("failed to get product info: %v\n", err)
	}

	return generator.NewGenerator(cfg, extract)
}

func (extract ComicBoostRensaiExtractor) ExtractItems(doc *goquery.Document) ([]site.Item, error) {

	productItems, err := extract.productItems([]string{extract.config.URL})
	if err != nil {
		return nil, fmt.Errorf("failed to productItems: (Title='%v'): %w", extract.config.Title, err)
	}

	return productItems, nil
}

func (extract ComicBoostRensaiExtractor) productInfo(cfg *site.Config) error {

	doc, err := util.FetchHtmlDoc(cfg.URL)
	if err != nil {
		return fmt.Errorf("failed to get HTML document (Site='%v'): %w", cfg.Title, err)
	}

	cfg.Description = strings.TrimSpace(doc.Find("p.comic-description-text").First().Text())
	cfg.Title += strings.TrimSpace(doc.Find("h1.comic-title").First().Text())

	return nil
}

func (extract ComicBoostRensaiExtractor) productItems(productURLs []string) ([]site.Item, error) {

	domain := util.GetFqdn(extract.config.URL)

	var items []site.Item
	processedIndex := 0 // productURLsをキューとするインデックス
	for processedIndex < len(productURLs) {

		productURL := productURLs[processedIndex]
		processedIndex++

		doc, err := util.FetchHtmlDoc(productURL)
		if err != nil {
			return nil, fmt.Errorf("failed to FetchHtmlDoc: %w", err)
		}

		// 次のページがあればproductURLsの末尾に追加
		nextPageURL, found := doc.Find(".to-next > a").First().Attr("href")
		if found && nextPageURL != "javascript:void(0);" {
			nextPageURL = domain + nextPageURL
			productURLs = append(productURLs, nextPageURL)
		}

		items = append(items, extract.extractItems(doc, domain)...)

		util.ItemPerSleep(processedIndex, 9, 1)
	}

	if len(items) == 0 {
		return nil, fmt.Errorf("item not found")
	}

	return items, nil
}

func (extract ComicBoostRensaiExtractor) extractItems(doc *goquery.Document, domain string) []site.Item {

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
