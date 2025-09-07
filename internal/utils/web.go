package utils

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

func GetHtmlDoc(url string) (*goquery.Document, error) {

	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("could not fetching URL: (URL='%v'): %w", url, err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("received non 200 status: (URL='%v'): %v", url, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, fmt.Errorf("could not create goquery doc: (URL='%v'): %w", url, err)
	}

	return doc, nil
}

func GetDomain(str string) (string, error) {
	u, err := url.Parse(str)
	if err != nil {
		return "", fmt.Errorf("could not parse domain: (URL='%v'): %w", str, err)
	}
	return fmt.Sprintf("%s://%s", u.Scheme, u.Host), nil
}
