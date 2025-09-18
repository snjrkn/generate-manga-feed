package service

import (
	"fmt"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/snjrkn/generate-manga-feed/internal/generator"
	"github.com/snjrkn/generate-manga-feed/internal/site"
	"github.com/snjrkn/generate-manga-feed/internal/util"
)

type YoungAnimalExtractor struct {
	config site.Config
}

func NewYoungAnimalExtractor(cfg site.Config) *YoungAnimalExtractor {
	return &YoungAnimalExtractor{
		config: cfg,
	}
}

func YoungAnimalOneshot() *generator.Generator {
	cfg := site.Config{
		Title:       "ヤングアニマル 読み切り",
		URL:         "https://younganimal.com/category/manga?type=%E8%AA%AD%E3%81%BF%E5%88%87%E3%82%8A",
		DateLayout:  "2006年1月2日",
		Description: "None",
	}
	return generator.NewGenerator(cfg, NewYoungAnimalExtractor(cfg))
}

func (extract YoungAnimalExtractor) ExtractItems(doc *goquery.Document) ([]site.Item, error) {

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

func (extract YoungAnimalExtractor) productURLs(doc *goquery.Document) ([]string, error) {

	var urls []string
	doc.Find(".ranking-item.category-box-vertical > a").Each(func(i int, sel *goquery.Selection) {
		if url, exist := sel.Attr("href"); exist && strings.Contains(url, "series") {
			urls = append(urls, url)
		}
	})

	if len(urls) == 0 {
		return nil, fmt.Errorf("URL not found")
	}

	return urls, nil
}

func (extract YoungAnimalExtractor) productItems(urls []string) ([]site.Item, error) {

	var items []site.Item
	for i, url := range urls {
		doc, err := util.FetchHtmlDoc(url)
		if err != nil {
			return nil, fmt.Errorf("failed to FetchHtmlDoc: %w", err)
		}

		product := strings.TrimSpace(doc.Find("span.series-ep-list-item-h-text").First().Text())
		author := strings.TrimSpace(doc.Find("span.article-text").First().Text())
		desc := doc.Find(".series-h-credit-info-text-text span").First().Text()
		date := doc.Find("time.series-ep-list-date-time").First().Text()
		link := doc.Find("a.series-act-read-btn").AttrOr("href", "")

		title := fmt.Sprintf("%s %s %s", product, author, date)

		// 今年分の作品には西暦年が付いていないための対応（去年分以前は付いているので不要）
		if !strings.Contains(date, "年") {
			date = time.Now().In(util.GetTokyoLocation()).Format("2006") + "年" + date
		}

		items = append(items, site.Item{Title: title, Link: link, Desc: desc, Date: date})

		util.ItemPerSleep(i, 9, 1)
	}

	return items, nil
}
