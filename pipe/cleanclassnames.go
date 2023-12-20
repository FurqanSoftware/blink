package pipe

import (
	"slices"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type CleanClassNamesFilter struct {
	preserveSelectors []string
	preserveClasses   []string
}

func CleanClassNames() CleanClassNamesFilter {
	return CleanClassNamesFilter{}
}

func (f CleanClassNamesFilter) Apply(_ Context, p Page) (Page, error) {
	p.Doc.Find("[class]").Each(func(_ int, s *goquery.Selection) {
		if f.shouldPreserve(s) {
			return
		}
		classes := strings.Split(s.AttrOr("class", ""), " ")
		for _, class := range classes {
			if strings.HasPrefix(class, "Ƀ") {
				continue
			}
			if slices.Contains(f.preserveClasses, class) {
				s.AddClass("Ƀkeep_" + class)
			}
			s.RemoveClass(class)
		}
		if strings.TrimSpace(s.AttrOr("class", "")) == "" {
			s.RemoveAttr("class")
		}
	})
	return p, nil
}

func (f CleanClassNamesFilter) shouldPreserve(s *goquery.Selection) bool {
	for _, v := range f.preserveSelectors {
		if s.Is(v) {
			return true
		}
	}
	return false
}

func (f CleanClassNamesFilter) WithPreserveSelectors(selectors ...string) CleanClassNamesFilter {
	f.preserveSelectors = append(f.preserveSelectors, selectors...)
	return f
}

func (f CleanClassNamesFilter) WithPreserveClasses(classes ...string) CleanClassNamesFilter {
	f.preserveClasses = append(f.preserveClasses, classes...)
	return f
}
