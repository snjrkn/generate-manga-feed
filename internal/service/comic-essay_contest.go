package service

import (
	"fmt"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/snjrkn/generate-manga-feed/internal/site"
	"github.com/snjrkn/generate-manga-feed/internal/util"
)

type comicEssayContest struct {
	config site.Config
}

func NewComicEssayContest() site.Site {
	cfg := site.Config{
		Title:       "コミックエッセイ プチ大賞",
		URL:         "https://www.comic-essay.com/contest/winner/",
		DateLayout:  "20060102",
		Description: "None",
	}
	return site.Site{
		Config:    cfg,
		Extractor: &comicEssayContest{config: cfg},
	}
}

func (ext comicEssayContest) ExtractItems(doc *goquery.Document) ([]site.Item, error) {

	ContestURLs, err := ext.ContestURLs(doc)
	if err != nil {
		return nil, fmt.Errorf("failed to ContestURLs: %w", err)
	}
	// 賞の最新分は賞のページにあるので追加
	ContestURLs = append(ContestURLs, ext.config.URL)

	productURLs, err := ext.productURLs(ContestURLs)
	if err != nil {
		return nil, fmt.Errorf("failed to productURLs: %w", err)
	}

	productItems, err := ext.productItems(productURLs)
	if err != nil {
		return nil, fmt.Errorf("failed to productItems: (Title='%v'): %w", ext.config.Title, err)
	}

	return productItems, nil
}

func (ext *comicEssayContest) ContestURLs(doc *goquery.Document) ([]string, error) {

	var urls []string
	doc.Find(".c-mt20 ._btn-500").Each(func(i int, sel *goquery.Selection) {
		if url, exists := sel.Attr("href"); exists && strings.Contains(url, "winner") {
			urls = append(urls, url)
		}
	})

	if len(urls) == 0 {
		return nil, fmt.Errorf("URL not found")
	}

	return urls, nil
}

func (ext *comicEssayContest) productURLs(awUrls []string) ([]string, error) {

	var urls []string
	for _, awUrl := range awUrls {
		doc, err := util.FetchHtmlDoc(awUrl)
		if err != nil {
			return nil, fmt.Errorf("failed to FetchHtmlDoc: %w", err)
		}
		doc.Find("._contest-btn-read").Each(func(i int, sel *goquery.Selection) {
			if url, exists := sel.Attr("href"); exists && strings.Contains(url, "episode") {
				urls = append(urls, url)
			}
		})
	}

	if len(urls) == 0 {
		return nil, fmt.Errorf("URL not found")
	}

	time.Sleep(1 * time.Second)

	return urls, nil
}

func (ext *comicEssayContest) productItems(urls []string) ([]site.Item, error) {

	items := []site.Item{}
	for i := range urls {
		doc, err := util.FetchHtmlDoc(urls[i])
		if err != nil {
			return nil, fmt.Errorf("failed to FetchHtmlDoc: %w", err)
		}

		product := strings.TrimSpace(doc.Find("div.episode-info__title").First().Text())
		author := strings.TrimSpace(doc.Find("span.episode-info__author--item").First().Text())
		desc := strings.TrimSpace(doc.Find("div.episode-info__synopsis").First().Text())
		imageUrl := doc.Find("div.detail-title-banner._episode > img").First().AttrOr("src", "")

		title := fmt.Sprintf("%s %s", product, author)
		link := urls[i]
		// 日付はページに明記されていないが、画像のディレクトリ名に西暦年と月があるので"01"を追加して日付とする
		date := strings.Split(imageUrl, "/")[3] + "01"

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
