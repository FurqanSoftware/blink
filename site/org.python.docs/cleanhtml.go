package pythondocs

import (
	"regexp"
	"strings"

	"github.com/FurqanSoftware/blink/pipe"
	"github.com/PuerkitoBio/goquery"
)

func CleanHTML() pipe.Filter {
	codeLangRegexpes := []*regexp.Regexp{
		regexp.MustCompile(`code (\w+) highlight`),
		regexp.MustCompile(`highlight\-([\w\+]+)`),
		regexp.MustCompile(`hl\-(\w+)`),
	}

	return pipe.FilterFunc(func(_ pipe.Context, p pipe.Page) (pipe.Page, error) {
		p.Doc.Find(".headerlink, hr, #contents .topic-title, #topics .topic-title, colgroup, .line-block, .anchor-link, #reference-index").Remove()

		p.Doc.Find("#the-python-language-reference").Each(func(_ int, s *goquery.Selection) {
			s.BeforeSelection(s.Contents())
			s.Remove()
		})

		p.Doc.Find(`div[class*="highlight-"], div[class*="hl-"]`).Each(func(_ int, s *goquery.Selection) {
			pre := s.Find("pre")
			pre.SetText(pre.Text())
			class := s.AttrOr("class", "")
			lang := ""
			for _, r := range codeLangRegexpes {
				m := r.FindStringSubmatch(class)
				if len(m) == 2 {
					lang = m[1]
					break
				}
			}
			if lang == "default" || strings.HasPrefix(lang, "python") || strings.HasPrefix(lang, "ipython") {
				lang = "python"
			}
			pre.SetAttr("data-language", lang)
			s.ReplaceWithSelection(pre)
		})

		return p, nil
	})
}

// Some CSS selectors and transformaton rules in this file were adapted from DevDocs:
// - https://github.com/freeCodeCamp/devdocs/blob/ab9aeb2622838131574023ea3a01e933e2c770df/lib/docs/filters/sphinx/clean_html.rb
// - https://github.com/freeCodeCamp/devdocs/blob/ab9aeb2622838131574023ea3a01e933e2c770df/lib/docs/filters/python/clean_html.rb
//
// DevDocs is licensed under the terms of the Mozilla Public License v2.0.
