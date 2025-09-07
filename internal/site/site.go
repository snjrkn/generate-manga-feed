package site

import (
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Config struct {
	Title       string
	URL         string
	Description string
	DateLayout  string
}

type Item struct {
	Title       string
	Link        string
	Desc        string
	Date        string
	CreatedDate time.Time
}

type Extractor interface {
	ExtractItems(*goquery.Document) ([]Item, error)
}
