package fs

import "github.com/FurqanSoftware/blink/crawl/web"

type Option interface {
	Apply(*Crawler)
}

type OptionFunc func(*Crawler)

func (f OptionFunc) Apply(c *Crawler) {
	f(c)
}

func WebOptions(options ...web.Option) Option {
	return OptionFunc(func(c *Crawler) {
		c.webOptions = append(c.webOptions, options...)
	})
}
