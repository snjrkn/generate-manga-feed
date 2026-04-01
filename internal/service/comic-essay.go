package service

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/snjrkn/generate-manga-feed/internal/site"
	"github.com/snjrkn/generate-manga-feed/internal/util"
)

type comicEssay struct {
	config site.Config
}

func NewComicEssayGekijo() site.Site {
	cfg := site.Config{
		Title:       "コミックエッセイ劇場",
		URL:         "https://www.comic-essay.com/comics/",
		DateLayout:  "2006.1.02",
		Description: "None",
	}
	return site.Site{
		Config:    cfg,
		Extractor: &comicEssay{config: cfg},
	}
}

func (ext comicEssay) ExtractItems(doc *goquery.Document) ([]site.Item, error) {

	productURLs, productDates, err := ext.productURLsAndDates(doc)
	if err != nil {
		return nil, fmt.Errorf("failed to productURLs: %w", err)
	}

	productItems, err := ext.productItems(productURLs, productDates)
	if err != nil {
		return nil, fmt.Errorf("failed to productItems: (Title='%v'): %w", ext.config.Title, err)
	}

	return productItems, nil
}

func (ext *comicEssay) productURLsAndDates(doc *goquery.Document) (urls, dates []string, err error) {

	doc.Find("li.thum-list__item").Each(func(i int, sel *goquery.Selection) {
		url, exists := sel.Find("a.thum-list__link").Attr("href")
		if exists {
			urls = append(urls, url)
		} else {
			return
		}
		dates = append(dates, strings.TrimSpace(sel.Find("div.thum-list__body--head__date").Text()))
	})

	if len(urls) == 0 {
		return nil, nil, fmt.Errorf("URL not found")
	}

	return urls, dates, nil
}

func (ext *comicEssay) productItems(urls, dates []string) ([]site.Item, error) {

	items := []site.Item{}
	for i := range urls {
		doc, err := util.FetchHtmlDoc(urls[i])
		if err != nil {
			return nil, fmt.Errorf("failed to FetchHtmlDoc: %w", err)
		}

		product := strings.TrimSpace(doc.Find("div.episode-info__title").First().Text())
		author := strings.TrimSpace(doc.Find("span.episode-info__author--item").First().Text())
		desc := strings.TrimSpace(doc.Find("div.episode-info__synopsis").First().Text())
		episode := strings.TrimSpace(doc.Find("div.episode-list__item--title").First().Text())
		link := doc.Find("a.episode-list__item--link").First().AttrOr("href", urls[i])

		episode = strings.ReplaceAll(episode, "\t", " ")
		title := fmt.Sprintf("%s %s %s %s", dates[i], episode, product, author)
		date := dates[i]

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
