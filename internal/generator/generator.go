package generator

import (
	"fmt"

	"github.com/snjrkn/generate-manga-feed/internal/site"
	"github.com/snjrkn/generate-manga-feed/internal/utils"
)

type Generator struct {
	config    site.Config
	extractor site.Extractor
}

func NewGenerator(cfg site.Config, ext site.Extractor) *Generator {
	return &Generator{
		config:    cfg,
		extractor: ext,
	}
}

func (generator *Generator) MakeFeed() (string, error) {

	// 共通の前処理: HTMLドキュメントの取得
	doc, err := utils.GetHtmlDoc(generator.config.URL)
	if err != nil {
		return "", fmt.Errorf("failed to get HTML document: %w", err)
	}

	// サイト固有のデータ抽出
	items, err := generator.extractor.ExtractItems(doc)
	if err != nil {
		return "", fmt.Errorf("failed to extract items: %w", err)
	}

	// 共通の後処理: アイテムの検証と事前処理
	if err := utils.ValidateAndPrepare(generator.config, items); err != nil {
		return "", fmt.Errorf("failed to validate and prepare: %w", err)
	}

	// 共通の後処理: フィードの生成
	rss, err := utils.GenerateFeed(generator.config, items)
	if err != nil {
		return "", fmt.Errorf("failed to generate feed: %w", err)
	}

	return rss, nil
}
