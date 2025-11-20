package pipe

import (
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

type Context struct {
	URL  *url.URL
	Path string
}

type Page struct {
	Meta  Meta
	Doc   *goquery.Document
	Marks []Mark
}

type Meta struct {
	Title       string
	OriginalURL string
}

type Mark struct {
	Kind string
	Name string
	Href string
}
