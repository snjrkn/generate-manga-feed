package service

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/snjrkn/generate-manga-feed/internal/site"
	"github.com/snjrkn/generate-manga-feed/internal/util"
)

type toti struct {
	config site.Config
}

func NewToti(productId string) site.Site {
	cfg := site.Config{
		Title:       "トーチ",
		URL:         "https://to-ti.in/product/" + strings.TrimSpace(productId),
		DateLayout:  "2006/01/02",
		Description: "None",
	}

	ext := &toti{config: cfg}
	if err := ext.productInfo(&cfg); err != nil {
		fmt.Printf("failed to contentInfo: %v\n", err)
	}

	return site.Site{
		Config:    cfg,
		Extractor: &toti{config: cfg},
	}

}

func (ext toti) ExtractItems(doc *goquery.Document) ([]site.Item, error) {

	items := []site.Item{}
	var err error
	if ext.config.URL == "https://to-ti.in/product/" {
		productURLs, err := ext.productURLs(doc)
		if err != nil {
			return nil, fmt.Errorf("failed to productURLs: %w", err)
		}
		items, err = ext.productItems(productURLs)
		if err != nil {
			return nil, fmt.Errorf("failed to productItems: (Title='%v'): %w", ext.config.Title, err)
		}
	} else {
		items, err = ext.storyItems(doc)
		if err != nil {
			return nil, fmt.Errorf("failed to storyItems: (Title='%v'): %w", ext.config.Title, err)
		}
	}

	return items, nil
}

func (ext toti) productURLs(doc *goquery.Document) ([]string, error) {

	var urls []string
	doc.Find("article a").Each(func(i int, sel *goquery.Selection) {
		if url, exists := sel.Attr("href"); exists {
			urls = append(urls, url)
		}
	})

	if len(urls) == 0 {
		return nil, fmt.Errorf("URL not found")
	}

	return urls, nil
}

func (ext toti) productItems(productURLs []string) ([]site.Item, error) {

	items := []site.Item{}
	for i := range productURLs {
		doc, err := util.FetchHtmlDoc(productURLs[i])
		if err != nil {
			return nil, fmt.Errorf("failed to FetchHtmlDoc: %w", err)
		}

		date := strings.TrimSpace(doc.Find("time").First().Text())
		product := strings.TrimSpace(doc.Find("header > h3").Text())
		desc := strings.TrimSpace(doc.Find("header > p").Text())
		story := strings.TrimSpace(doc.Find("p.next > a >  span").Text())
		link, exists := doc.Find("p.next a").Attr("href")
		// 第1話が公開された場合の対応
		if !exists {
			link = doc.Find("p.prev a").AttrOr("href", productURLs[i])
		}
		// 書店様へが更新された場合の対応（各話のリンクではなくページのリンクを設定）
		if strings.Contains(productURLs[i], "store_information") && strings.Contains(product, "書店様へ") {
			link = productURLs[i]
		}

		titleDate := date
		// 日付にタイトルが含まれる場合の対応
		titleDate = regexp.MustCompile("^.*'").ReplaceAllString(date, "'")
		titleDate = strings.ReplaceAll(titleDate, "'", "’")
		title := fmt.Sprintf("%s %s %s", titleDate, product, story)

		// 日付にタイトルが含まれる場合の対応
		date = regexp.MustCompile("^.*'").ReplaceAllString(date, "'")

		date = "20" + strings.ReplaceAll(date, "'", "")
		date = strings.ReplaceAll(date, " UPDATE", "")

		items = append(items, site.Item{
			Title: title,
			Link:  link,
			Desc:  desc,
			Date:  date,
		})
	}

	if len(items) == 0 {
		return nil, fmt.Errorf("item not found")
	}

	return items, nil
}

func (ext toti) productInfo(cfg *site.Config) error {

	doc, err := util.FetchHtmlDoc(cfg.URL)
	if err != nil {
		return fmt.Errorf("failed to FetchHtmlDoc: (Site='%v'): %w", cfg.Title, err)
	}

	if ext.config.URL != "https://to-ti.in/product/" {
		cfg.Title += "【連載】" + strings.TrimSpace(doc.Find("header > h3").Text())

		str := strings.TrimSpace(doc.Find("header > p").Text())
		str += strings.TrimSpace(doc.Find(".description").Text())
		str = strings.Join(strings.Fields(str), " ")
		cfg.Description = str
	}

	return nil
}

func (ext toti) storyItems(doc *goquery.Document) ([]site.Item, error) {

	items := []site.Item{}
	doc.Find(".episode li").Each(func(i int, sel *goquery.Selection) {
		title := strings.TrimSpace(sel.Find("span").Text())
		link, exists := sel.Find("a").Attr("href")
		if !exists {
			link = ext.config.URL
		}
		desc := "none"
		date := ""

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
