package service

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/snjrkn/generate-manga-feed/internal/generator"
	"github.com/snjrkn/generate-manga-feed/internal/site"
	"github.com/snjrkn/generate-manga-feed/internal/utils"
)

type ComicEssayExtractor struct {
	config site.Config
}

func NewComicEssayExtractor(cfg site.Config) *ComicEssayExtractor {
	return &ComicEssayExtractor{
		config: cfg,
	}
}

func ComicEssayGekijo() *generator.Generator {
	cfg := site.Config{
		Title:       "コミックエッセイ劇場",
		URL:         "https://www.comic-essay.com/comics/",
		DateLayout:  "2006.1.02",
		Description: "None",
	}
	return generator.NewGenerator(cfg, NewComicEssayExtractor(cfg))
}

func (extract ComicEssayExtractor) ExtractItems(doc *goquery.Document) ([]site.Item, error) {

	productURLs, productDates := extract.productURLsAndDates(doc)

	productItems, err := extract.productItems(productURLs, productDates)
	if err != nil {
		return nil, fmt.Errorf("failed to productItems: (Title='%v'): %w", extract.config.Title, err)
	}

	return productItems, nil
}

func (extract *ComicEssayExtractor) productURLsAndDates(doc *goquery.Document) (urls, dates []string) {

	doc.Find("li.thum-list__item").Each(func(i int, sel *goquery.Selection) {
		url, exist := sel.Find("a.thum-list__link").Attr("href")
		if exist {
			urls = append(urls, url)
		} else {
			return
		}
		dates = append(dates, strings.TrimSpace(sel.Find("div.thum-list__body--head__date").Text()))
	})

	return urls, dates
}

func (extract *ComicEssayExtractor) productItems(urls, dates []string) ([]site.Item, error) {

	items := []site.Item{}
	for i := range urls {
		doc, err := utils.GetHtmlDoc(urls[i])
		if err != nil {
			return nil, fmt.Errorf("failed to GetHtmlDoc: %w", err)
		}

		product := strings.TrimSpace(doc.Find("div.episode-info__title").First().Text())
		author := strings.TrimSpace(doc.Find("span.episode-info__author--item").First().Text())
		desc := strings.TrimSpace(doc.Find("div.episode-info__synopsis").First().Text())
		episode := strings.TrimSpace(doc.Find("div.episode-list__item--title").First().Text())
		link := doc.Find("a.episode-list__item--link").First().AttrOr("href", "")

		episode = strings.ReplaceAll(episode, "\t", " ")
		title := fmt.Sprintf("%s %s %s %s", dates[i], episode, product, author)
		desc = strings.ReplaceAll(desc, "\n", " ")
		date := dates[i]

		items = append(items, site.Item{Title: title, Link: link, Desc: desc, Date: date})
	}

	return items, nil
}
