package site

import (
	"github.com/FurqanSoftware/blink/crawl"
	"github.com/FurqanSoftware/blink/pipe"
)

type Site struct {
	key                  string
	crawler              crawl.Crawler
	pipeline             pipe.Pipeline
	trimPathPrefix       string
	title                string
	attributionContentMd string
}

func New(key string, crawler crawl.Crawler, pipeline pipe.Pipeline, options ...Option) Site {
	s := Site{
		key:      key,
		crawler:  crawler,
		pipeline: pipeline,
	}
	for _, o := range options {
		o.Apply(&s)
	}
	return s
}

func (s Site) Key() string {
	return s.key
}
