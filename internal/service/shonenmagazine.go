package service

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/snjrkn/generate-manga-feed/internal/generator"
	"github.com/snjrkn/generate-manga-feed/internal/site"
	"github.com/snjrkn/generate-manga-feed/internal/utils"
)

type ShonenMagazineExtractor struct {
	config site.Config
}

func NewShonenMagazineExtractor(cfg site.Config) *ShonenMagazineExtractor {
	return &ShonenMagazineExtractor{
		config: cfg,
	}
}

func ShonenMagazineAward() *generator.Generator {
	cfg := site.Config{
		Title:       "少年マガジン 新人漫画大賞",
		URL:         "https://debut.shonenmagazine.com/archive/#awards",
		DateLayout:  "2006/01/02",
		Description: "None",
	}
	return generator.NewGenerator(cfg, NewShonenMagazineExtractor(cfg))
}

func ShonenMagazineRise() *generator.Generator {
	cfg := site.Config{
		Title:       "少年マガジン ライズ",
		URL:         "https://debut.shonenmagazine.com/archive/#magazinerise",
		DateLayout:  "2006/01/02",
		Description: "None",
	}
	return generator.NewGenerator(cfg, NewShonenMagazineExtractor(cfg))
}

func (extract ShonenMagazineExtractor) ExtractItems(doc *goquery.Document) ([]site.Item, error) {

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

func (extract ShonenMagazineExtractor) productURLs(doc *goquery.Document) ([]string, error) {

	var findStr string
	index := strings.LastIndex(extract.config.URL, "#")
	if index != -1 {
		findStr = extract.config.URL[index:]
	}

	var urls []string
	doc.Find(findStr + " .works-list a").Each(func(i int, sel *goquery.Selection) {
		if url, exist := sel.Attr("href"); exist && (strings.Contains(url, "content") || strings.Contains(url, "episode")) {
			urls = append(urls, url)
		}
	})

	if len(urls) == 0 {
		return nil, fmt.Errorf("URL not found")
	}

	return urls, nil
}

func (extract ShonenMagazineExtractor) productItems(urls []string) ([]site.Item, error) {

	var items []site.Item
	for i, url := range urls {
		doc, err := utils.FetchHtmlDoc(url)
		if err != nil {
			return nil, fmt.Errorf("failed to FetchHtmlDoc: %w", err)
		}

		// 特選取って連載になった作品の特殊対応（受賞作ではなく、連載の第1話になっている作品があるため）
		str := strings.TrimSpace(doc.Find("h2.p-episode__header-ttl").First().Text())
		if strings.Contains(str, "第") && strings.Contains(str, "話") {
			continue
		}

		date := strings.TrimSpace(doc.Find("p.p-episode__header-date").First().Text())
		product := strings.TrimSpace(doc.Find("h1.p-episode__comic-ttl").First().Text())
		author := strings.TrimSpace(doc.Find("h3.p-episode__comic-name").First().Text())
		desc := strings.TrimSpace(doc.Find("div.p-episode__comic-description p").First().Text())

		title := fmt.Sprintf("%s %s %s", date, product, author)
		link := url

		items = append(items, site.Item{Title: title, Link: link, Desc: desc, Date: date})

		utils.ItemPerSleep(i, 9, 2)
	}

	return items, nil
}
