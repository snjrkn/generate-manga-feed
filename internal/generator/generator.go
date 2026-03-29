package generator

import (
	"fmt"

	"github.com/snjrkn/generate-manga-feed/internal/site"
	"github.com/snjrkn/generate-manga-feed/internal/util"
)

type Generator struct {
	site site.Site
}

func NewGenerator(st site.Site) *Generator {
	return &Generator{
		site: st,
	}
}

func (gen *Generator) MakeFeed() (string, error) {

	cfg := gen.site.Config

	doc, err := util.FetchHtmlDoc(cfg.URL)
	if err != nil {
		return "", fmt.Errorf("failed to FetchHtmlDoc (Site='%v'): %w", cfg.Title, err)
	}

	items, err := gen.site.Extractor.ExtractItems(doc)
	if err != nil {
		return "", fmt.Errorf("failed to ExtractItems (Site='%v', URL='%v'): %w", cfg.Title, cfg.URL, err)
	}

	if err := util.ValidateAndPrepare(cfg, items); err != nil {
		return "", fmt.Errorf("failed to ValidateAndPrepare: %w", err)
	}

	rss, err := util.GenerateFeed(cfg, items)
	if err != nil {
		return "", fmt.Errorf("failed to GenerateFeed: %w", err)
	}

	return rss, nil
}
