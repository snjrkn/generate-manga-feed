package util

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/gorilla/feeds"
	"github.com/snjrkn/generate-manga-feed/internal/site"
)

// GenerateFeed はfeedデータを生成して文字列にして返す
func GenerateFeed(cfg site.Config, items []site.Item) (string, error) {

	feed := &feeds.Feed{
		Title:       cfg.Title,
		Link:        &feeds.Link{Href: cfg.URL},
		Created:     time.Now(),
		Description: cfg.Description,
	}

	for i := range items {
		feed.Add(&feeds.Item{
			Title:       items[i].Title,
			Link:        &feeds.Link{Href: items[i].Link},
			Description: items[i].Desc,
			Id:          items[i].Link,
			Created:     items[i].CreatedDate,
		})
	}

	rss, err := feed.ToRss()
	if err != nil {
		return "", fmt.Errorf("failed to convert: (Titile='%v'): %w ", cfg.Title, err)
	}

	return rss, nil
}

// ValidateAndPrepare はアイテムのデータの検証と事前処理を行う
func ValidateAndPrepare(cfg site.Config, items []site.Item) error {

	if len(items) == 0 {
		return fmt.Errorf("items slice is empty: (Titile='%v', URL='%v')", cfg.Title, cfg.URL)
	}

	for i, item := range items {

		items[i].Title = normalizeSpace(item.Title)
		items[i].Desc = normalizeSpace(item.Desc)

		// RSS2.0ではtitleかlinkのどちらかが必須
		if items[i].Title == "" && items[i].Link == "" {
			return fmt.Errorf("title or link is empty: (Titile='%v', Link='%v')", items[i].Title, items[i].Link)
		}

		// CreatedDateが空でDateの値がある場合、DateをパースしてCreatedDateに設定
		// CreatedDateの値がある場合はそのまま使用する（comiplexのみの対応）
		if items[i].CreatedDate.IsZero() && item.Date != "" {
			createdDate, err := parseTime(items[i].Date, cfg.DateLayout)
			if err != nil {
				return fmt.Errorf("failed to parse time: (Titile='%v'): %w", items[i].Title, err)
			}
			items[i].CreatedDate = createdDate
		}
	}

	sortItemsByCreated(items)

	return nil
}

// normalizeSpace は全角スペースを半角スペースに変換して連続した半角スペースを1つにして返す
func normalizeSpace(str string) string {
	// str = strings.ReplaceAll(str, "\t", " ")
	// str = strings.ReplaceAll(str, "\n", " ")
	str = strings.Join(strings.Fields(strings.ReplaceAll(str, "　", " ")), " ")
	return str
}

// sortItemsByCreated はCreatedで降順ソートする
func sortItemsByCreated(items []site.Item) {
	sort.Slice(items, func(i, j int) bool {
		return items[j].CreatedDate.Before(items[i].CreatedDate)
	})
}
