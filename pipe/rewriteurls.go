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
	prefixMappings       map[string]string
}

func RewriteURLs(siteKey string, rootURL string) RewriteURLsFilter {
	return RewriteURLsFilter{
		siteKey:        siteKey,
		rootURL:        rootURL,
		prefixMappings: map[string]string{},
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
		href := ""
		u, err := f.absoluteURL(baseURL, s.AttrOr("href", ""), x.URL.Scheme)
		if err != nil {
			return
		}
		if !f.isInternal(u) {
			s.AddClass("Ƀexternal")
			href = u.String()
		} else {
			u.Path = strings.TrimSuffix(u.Path, ".html")
			href = u.String()
			href = strings.TrimPrefix(href, f.rootURL)
			href = path.Join("/"+f.siteKey, href)
		}
		for from, to := range f.prefixMappings {
			if strings.HasPrefix(href, from) {
				href = to + strings.TrimPrefix(href, from)
				break
			}
		}
		s.SetAttr("href", href)
	})
	return p, nil
}

func (f RewriteURLsFilter) isInternal(u *url.URL) bool {
	if len(f.disallowedURLFilters) > 0 {
		for _, f := range f.disallowedURLFilters {
			if f.MatchString(u.String()) {
				return false
			}
		}
	}
	if len(f.urlFilters) > 0 {
		for _, f := range f.urlFilters {
			if f.MatchString(u.String()) {
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

func (f RewriteURLsFilter) WithPrefixMappings(fromto ...string) RewriteURLsFilter {
	for i := 0; i < len(fromto); i += 2 {
		f.prefixMappings[fromto[i]] = fromto[i+1]
	}
	return f
}

func (f RewriteURLsFilter) absoluteURL(baseURL *url.URL, ref string, scheme string) (*url.URL, error) {
	u, err := baseURL.Parse(ref)
	if err != nil {
		return nil, err
	}
	if u.Scheme == "//" {
		u.Scheme = scheme
	}
	return u, nil
}
