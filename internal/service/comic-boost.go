package service

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/snjrkn/generate-manga-feed/internal/generator"
	"github.com/snjrkn/generate-manga-feed/internal/site"
	"github.com/snjrkn/generate-manga-feed/internal/utils"
)

type ComicBoostExtractor struct {
	config site.Config
}

func NewComicBoostExtractor(cfg site.Config) *ComicBoostExtractor {
	return &ComicBoostExtractor{
		config: cfg,
	}
}

func ComicBoostOneshot() *generator.Generator {
	cfg := site.Config{
		Title:       "comicブースト 読み切り",
		URL:         "https://comic-boost.com/genre/3",
		DateLayout:  "2006/01/02",
		Description: "None",
	}
	return generator.NewGenerator(cfg, NewComicBoostExtractor(cfg))
}

func (extract ComicBoostExtractor) ExtractItems(doc *goquery.Document) ([]site.Item, error) {

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
func (extract ComicBoostExtractor) productURLs(doc *goquery.Document) ([]string, error) {

	var urls []string
	doc.Find("a.book-list-item-thum-wrapper").Each(func(i int, sel *goquery.Selection) {
		link := utils.GetFqdn(extract.config.URL) + sel.AttrOr("href", "")
		urls = append(urls, link)
	})

	if len(urls) == 0 {
		return nil, fmt.Errorf("URL not found")
	}

	return urls, nil
}

func (extract ComicBoostExtractor) productItems(productURLs []string) ([]site.Item, error) {

	domain := utils.GetFqdn(extract.config.URL)

	var items []site.Item
	processedIndex := 0 // productURLsをキューとするインデックス
	for processedIndex < len(productURLs) {

		productURL := productURLs[processedIndex]
		processedIndex++

		doc, err := utils.FetchHtmlDoc(productURL)
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

		utils.ItemPerSleep(processedIndex, 9, 1)
	}

	return items, nil
}

func (extract ComicBoostExtractor) extractItems(doc *goquery.Document, domain string) []site.Item {

	product := strings.TrimSpace(doc.Find("h1.comic-title").Text())
	author := strings.TrimSpace(doc.Find("li.author a").First().Text())
	desc := strings.TrimSpace(doc.Find("p.comic-description-text").First().Text())

	// タイトルに余計な文字列が付いてる作品対応
	if strings.Contains(product, "／読切版") {
		product = strings.ReplaceAll(product, "／読切版", "")
	}

	var items []site.Item
	doc.Find("a.book-product-list-item").Each(func(i int, sel *goquery.Selection) {

		episode := strings.TrimSpace(sel.Find("h4.title").Text())
		date := strings.TrimSpace(sel.Find("p.update-date").Text())
		link := sel.AttrOr("href", "")

		// 作品名の最後に「／作者名」が付いている場合の作者名取得対応（ブーストShorts対応）
		if strings.Contains(episode, "／") {
			// ブーストShortsの場合
			author = episode[strings.LastIndex(episode, "／")+len("／"):]
			if strings.HasSuffix(episode, author) {
				episode = strings.ReplaceAll(episode, "／"+author, "")
			}
		} else {
			// 通常作品の場合
			author = strings.TrimSpace(doc.Find("li.author a").First().Text())
		}

		// 作品ページのエピソードにも作品名が含まれているか、"読み切り"とだけなっている場合の対応
		if strings.Contains(episode, product) || strings.Contains(episode, "読み切り") {
			episode = ""
		}

		title := fmt.Sprintf("%s %s %s", product+" "+episode, author, date)
		link = domain + link

		items = append(items, site.Item{Title: title, Link: link, Desc: desc, Date: date})
	})

	return items
}
