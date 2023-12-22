package pipe

import (
	"github.com/PuerkitoBio/goquery"
)

type CleanStyleFilter struct {
	preserveSelectors []string
}

func CleanStyle() CleanStyleFilter {
	return CleanStyleFilter{}
}

func (f CleanStyleFilter) Apply(_ Context, p Page) (Page, error) {
	p.Doc.Find("[style]").Each(func(_ int, s *goquery.Selection) {
		if f.shouldPreserve(s) {
			return
		}
		s.RemoveAttr("style")
	})
	return p, nil
}

func (f CleanStyleFilter) shouldPreserve(s *goquery.Selection) bool {
	for _, v := range f.preserveSelectors {
		if s.Is(v) {
			return true
		}
	}
	return false
}

func (f CleanStyleFilter) WithPreserveSelectors(selectors ...string) CleanStyleFilter {
	f.preserveSelectors = append(f.preserveSelectors, selectors...)
	return f
}
