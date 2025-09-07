package service

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/snjrkn/generate-manga-feed/internal/generator"
	"github.com/snjrkn/generate-manga-feed/internal/site"
	"github.com/snjrkn/generate-manga-feed/internal/utils"
)

type ComicEssayAwardExtractor struct {
	config site.Config
}

func NewComicEssayAwardExtractor(cfg site.Config) *ComicEssayAwardExtractor {
	return &ComicEssayAwardExtractor{
		config: cfg,
	}
}

func ComicEssayAward() *generator.Generator {
	cfg := site.Config{
		Title:       "コミックエッセイ プチ大賞",
		URL:         "https://www.comic-essay.com/contest/winner/",
		DateLayout:  "20060102",
		Description: "None",
	}
	return generator.NewGenerator(cfg, NewComicEssayAwardExtractor(cfg))
}

func (extract ComicEssayAwardExtractor) ExtractItems(doc *goquery.Document) ([]site.Item, error) {

	awardURLs, err := extract.awardURLs(doc)
	if err != nil {
		return nil, fmt.Errorf("failed to awardURLs: (Title='%v'): %w", extract.config.Title, err)
	}
	// 賞の最新分は賞のページにあるので追加
	awardURLs = append(awardURLs, extract.config.URL)

	productURLs, err := extract.productURLs(awardURLs)
	if err != nil {
		return nil, fmt.Errorf("failed to productURLs: (Title='%v'): %w", extract.config.Title, err)
	}

	productItems, err := extract.productItems(productURLs)
	if err != nil {
		return nil, fmt.Errorf("failed to productItems: (Title='%v'): %w", extract.config.Title, err)
	}

	return productItems, nil
}

func (extract *ComicEssayAwardExtractor) awardURLs(doc *goquery.Document) ([]string, error) {

	var urls []string
	doc.Find(".c-mt20 ._btn-500").Each(func(i int, sel *goquery.Selection) {
		if url, exist := sel.Attr("href"); exist && strings.Contains(url, "winner") {
			urls = append(urls, url)
		}
	})

	if len(urls) == 0 {
		return nil, fmt.Errorf("award URL not found")
	}

	return urls, nil
}

func (extract *ComicEssayAwardExtractor) productURLs(awUrls []string) ([]string, error) {

	var urls []string
	for _, awUrl := range awUrls {
		doc, err := utils.GetHtmlDoc(awUrl)
		if err != nil {
			return nil, fmt.Errorf("failed to GetHtmlDoc: %w", err)
		}
		doc.Find("._contest-btn-read").Each(func(i int, sel *goquery.Selection) {
			if url, exist := sel.Attr("href"); exist && strings.Contains(url, "episode") {
				urls = append(urls, url)
			}
		})
	}

	if len(urls) == 0 {
		return nil, fmt.Errorf("product URL not found")
	}

	return urls, nil
}

func (extract *ComicEssayAwardExtractor) productItems(urls []string) ([]site.Item, error) {

	items := []site.Item{}
	for _, url := range urls {
		doc, err := utils.GetHtmlDoc(url)
		if err != nil {
			return nil, fmt.Errorf("failed to GetHtmlDoc: %w", err)
		}

		product := strings.TrimSpace(doc.Find("div.episode-info__title").First().Text())
		author := strings.TrimSpace(doc.Find("span.episode-info__author--item").First().Text())
		desc := strings.TrimSpace(doc.Find("div.episode-info__synopsis").First().Text())
		imageUrl := doc.Find("div.detail-title-banner._episode > img").First().AttrOr("src", "")

		title := fmt.Sprintf("%s %s", product, author)
		link := url
		// 日付はページに明記されていないが、画像のディレクトリ名に西暦年と月があるので"01"を追加して日付とする
		date := strings.Split(imageUrl, "/")[3] + "01"

		items = append(items, site.Item{Title: title, Link: link, Desc: desc, Date: date})
	}

	return items, nil
}
