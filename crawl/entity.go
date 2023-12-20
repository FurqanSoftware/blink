package crawl

import (
	"net/url"

	"github.com/gocolly/colly/v2"
)

type EntityKind int

const (
	EntityInvalid EntityKind = iota
	EntityPage
	EntityRedirect
)

type Entity interface {
	Kind() EntityKind
}

type Page struct {
	Request *colly.Request
	Content []byte
}

func (Page) Kind() EntityKind {
	return EntityPage
}

type Redirect struct {
	Request *colly.Request
	To      *url.URL
}

func (Redirect) Kind() EntityKind {
	return EntityRedirect
}
