package web

import (
	"log"
	"net/http"
	"net/url"
	"regexp"

	"github.com/FurqanSoftware/blink/crawl"
	"github.com/gocolly/colly/v2"
)

type Crawler struct {
	startURL             string
	disableCache         bool
	allowedDomains       []string
	urlFilters           []*regexp.Regexp
	disallowedURLFilters []*regexp.Regexp
}

var _ crawl.Crawler = Crawler{}

func New(startURL string, options ...Option) Crawler {
	c := Crawler{
		startURL: startURL,
	}
	for _, o := range options {
		o.Apply(&c)
	}
	return c
}

func (c Crawler) Run(entityCh chan<- crawl.Entity) error {
	options := []colly.CollectorOption{}
	if !c.disableCache {
		options = append(options, colly.CacheDir(".colly-cache"))
	}
	if c.allowedDomains != nil {
		options = append(options, colly.AllowedDomains(c.allowedDomains...))
	}
	if c.urlFilters != nil {
		options = append(options, colly.URLFilters(c.urlFilters...))
	}
	if c.disallowedURLFilters != nil {
		options = append(options, colly.DisallowedURLFilters(c.disallowedURLFilters...))
	}

	collector := colly.NewCollector(options...)

	collector.OnError(func(r *colly.Response, err error) {
		switch r.StatusCode {
		case http.StatusMovedPermanently, http.StatusFound, http.StatusTemporaryRedirect, http.StatusPermanentRedirect:
			url, _ := url.Parse(r.Request.AbsoluteURL(r.Headers.Get("Location")))
			entityCh <- crawl.Redirect{
				Request: r.Request,
				To:      url,
			}
			collector.Visit(url.String())
		}
	})

	collector.SetRedirectHandler(func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	})

	collector.OnRequest(func(r *colly.Request) {
		log.Printf("[C] Get %s", r.URL)
	})

	collector.OnResponse(func(r *colly.Response) {
		entityCh <- crawl.Page{
			Request: r.Request,
			Content: r.Body,
		}
	})

	collector.OnHTML("a[href]", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
	})

	return collector.Visit(c.startURL)
}
