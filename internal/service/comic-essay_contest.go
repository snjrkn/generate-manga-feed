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

type ComicEssayContestExtractor struct {
	config site.Config
}

func NewComicEssayContestExtractor(cfg site.Config) *ComicEssayContestExtractor {
	return &ComicEssayContestExtractor{
		config: cfg,
	}
}

func ComicEssayContest() *generator.Generator {
	cfg := site.Config{
		Title:       "コミックエッセイ プチ大賞",
		URL:         "https://www.comic-essay.com/contest/winner/",
		DateLayout:  "20060102",
		Description: "None",
	}
	return generator.NewGenerator(cfg, NewComicEssayContestExtractor(cfg))
}

func (extract ComicEssayContestExtractor) ExtractItems(doc *goquery.Document) ([]site.Item, error) {

	ContestURLs, err := extract.ContestURLs(doc)
	if err != nil {
		return nil, fmt.Errorf("failed to ContestURLs: (Title='%v'): %w", extract.config.Title, err)
	}
	// 賞の最新分は賞のページにあるので追加
	ContestURLs = append(ContestURLs, extract.config.URL)

	productURLs, err := extract.productURLs(ContestURLs)
	if err != nil {
		return nil, fmt.Errorf("failed to productURLs: (Title='%v'): %w", extract.config.Title, err)
	}

	productItems, err := extract.productItems(productURLs)
	if err != nil {
		return nil, fmt.Errorf("failed to productItems: (Title='%v'): %w", extract.config.Title, err)
	}

	return productItems, nil
}

func (extract *ComicEssayContestExtractor) ContestURLs(doc *goquery.Document) ([]string, error) {

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

func (extract *ComicEssayContestExtractor) productURLs(awUrls []string) ([]string, error) {

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

func (extract *ComicEssayContestExtractor) productItems(urls []string) ([]site.Item, error) {

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

		items = append(items, site.Item{Title: title, Link: link, Desc: desc, Date: date})
	}

	if len(items) == 0 {
		return nil, fmt.Errorf("item not found")
	}

	return items, nil
}
