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

type ChampionCrossExtractor struct {
	config site.Config
}

func NewChampionCrossExtractor(cfg site.Config) *ChampionCrossExtractor {
	return &ChampionCrossExtractor{
		config: cfg,
	}
}

func ChampionCrossOneshot() *generator.Generator {
	cfg := site.Config{
		Title:       "チャンピオンクロス 読み切り",
		URL:         "https://championcross.jp/category/manga?type=%E8%AA%AD%E3%81%BF%E5%88%87%E3%82%8A",
		DateLayout:  "2006年1月2日",
		Description: "None",
	}
	return generator.NewGenerator(cfg, NewChampionCrossExtractor(cfg))
}

func (extract ChampionCrossExtractor) ExtractItems(doc *goquery.Document) ([]site.Item, error) {

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

func (extract ChampionCrossExtractor) productURLs(doc *goquery.Document) ([]string, error) {

	var urls []string
	doc.Find(".category-box-vertical > a").Each(func(i int, sel *goquery.Selection) {
		if url, exists := sel.Attr("href"); exists && strings.Contains(url, "series") {
			urls = append(urls, url)
		}
	})

	if len(urls) == 0 {
		return nil, fmt.Errorf("URL not found")
	}

	return urls, nil
}

func (extract ChampionCrossExtractor) productItems(urls []string) ([]site.Item, error) {

	var items []site.Item
	for i := range urls {
		doc, err := util.FetchHtmlDoc(urls[i])
		if err != nil {
			return nil, fmt.Errorf("failed to FetchHtmlDoc: %w", err)
		}

		product := strings.TrimSpace(doc.Find("span.series-ep-list-item-h-text").First().Text())
		author := strings.TrimSpace(doc.Find("span.article-text").First().Text())
		desc := strings.TrimSpace(doc.Find(".series-h-credit-info-text-text span").First().Text())
		date := strings.TrimSpace(doc.Find("time.series-ep-list-date-time").First().Text())
		link := doc.Find("a.series-act-read-btn").AttrOr("href", urls[i])

		title := fmt.Sprintf("%s %s %s", product, author, date)

		// 今年分の作品には西暦年が付いていないための対応（去年分以前は付いているので不要）
		if !strings.Contains(date, "年") {
			date = time.Now().In(util.GetTokyoLocation()).Format("2006") + "年" + date
		}

		items = append(items, site.Item{Title: title, Link: link, Desc: desc, Date: date})

		util.ItemPerSleep(i, 9, 1)
	}

	if len(items) == 0 {
		return nil, fmt.Errorf("item not found")
	}

	return items, nil
}
