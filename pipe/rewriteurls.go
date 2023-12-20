package pipe

import (
	"net/url"
	"path"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type RewriteURLsFilter struct {
	siteKey              string
	rootURL              string
	urlFilters           []*regexp.Regexp
	disallowedURLFilters []*regexp.Regexp
}

func RewriteURLs(siteKey string, rootURL string) RewriteURLsFilter {
	return RewriteURLsFilter{
		siteKey: siteKey,
		rootURL: rootURL,
	}
}

func (f RewriteURLsFilter) Apply(x Context, p Page) (Page, error) {
	var baseURL *url.URL
	if v, ok := p.Doc.Find("base").Attr("href"); ok {
		var err error
		baseURL, err = url.Parse(v)
		if err != nil {
			return p, err
		}
	} else {
		baseURL = x.URL
	}

	p.Doc.Find("a[href]").Each(func(_ int, s *goquery.Selection) {
		href := s.AttrOr("href", "")
		if href == "" {
			return
		}
		href = f.absoluteURL(baseURL, href, x.URL.Scheme)
		if !f.isInternal(href) {
			s.SetAttr("href", href).AddClass("Ƀexternal")
			return
		}
		href = strings.TrimPrefix(href, f.rootURL)
		href = path.Join("/"+f.siteKey, href)
		s.SetAttr("href", href)
	})
	return p, nil
}

func (f RewriteURLsFilter) isInternal(url string) bool {
	if len(f.disallowedURLFilters) > 0 {
		for _, f := range f.disallowedURLFilters {
			if f.MatchString(url) {
				return false
			}
		}
	}
	if len(f.urlFilters) > 0 {
		for _, f := range f.urlFilters {
			if f.MatchString(url) {
				return true
			}
		}
		return false
	}
	return true
}

func (f RewriteURLsFilter) WithURLFilters(filters ...*regexp.Regexp) RewriteURLsFilter {
	f.urlFilters = append(f.urlFilters, filters...)
	return f
}

func (f RewriteURLsFilter) WithDisallowedURLFilters(filters ...*regexp.Regexp) RewriteURLsFilter {
	f.disallowedURLFilters = append(f.disallowedURLFilters, filters...)
	return f
}

func (f RewriteURLsFilter) WithDisallowedPaths(paths ...string) RewriteURLsFilter {
	rootURL, _ := url.Parse(f.rootURL)
	filters := make([]*regexp.Regexp, 0, len(paths))
	for _, path := range paths {
		url, _ := rootURL.Parse(path)
		filters = append(filters, regexp.MustCompile("^"+regexp.QuoteMeta(url.String())))
	}
	f.disallowedURLFilters = append(f.disallowedURLFilters, filters...)
	return f
}

func (f RewriteURLsFilter) absoluteURL(baseURL *url.URL, ref string, scheme string) string {
	url, err := baseURL.Parse(ref)
	if err != nil {
		return ""
	}
	if url.Scheme == "//" {
		url.Scheme = scheme
	}
	return url.String()
}
