package web

import (
	"net/url"
	"regexp"
)

type Option interface {
	Apply(*Crawler)
}

type OptionFunc func(*Crawler)

func (f OptionFunc) Apply(c *Crawler) {
	f(c)
}

func DisableCache() Option {
	return OptionFunc(func(c *Crawler) {
		c.disableCache = true
	})
}

func AllowedDomains(domains ...string) Option {
	return OptionFunc(func(c *Crawler) {
		c.allowedDomains = append(c.allowedDomains, domains...)
	})
}

func URLFilters(filters ...*regexp.Regexp) Option {
	return OptionFunc(func(c *Crawler) {
		c.urlFilters = append(c.urlFilters, filters...)
	})
}

func DisallowedURLFilters(filters ...*regexp.Regexp) Option {
	return OptionFunc(func(c *Crawler) {
		c.disallowedURLFilters = append(c.disallowedURLFilters, filters...)
	})
}

func DisallowedPaths(paths ...string) Option {
	return OptionFunc(func(c *Crawler) {
		startURL, _ := url.Parse(c.startURL)
		filters := make([]*regexp.Regexp, 0, len(paths))
		for _, path := range paths {
			url, _ := startURL.Parse(path)
			filters = append(filters, regexp.MustCompile("^"+regexp.QuoteMeta(url.String())))
		}
		c.disallowedURLFilters = append(c.disallowedURLFilters, filters...)
	})
}
