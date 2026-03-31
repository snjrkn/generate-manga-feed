package site

import (
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Site struct {
	Config    Config
	Extractor Extractor
}

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
