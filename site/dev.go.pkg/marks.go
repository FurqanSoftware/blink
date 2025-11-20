package pkggodev

import (
	"strings"

	"github.com/FurqanSoftware/blink/pipe"
	"github.com/PuerkitoBio/goquery"
)

func Marks() pipe.Filter {
	return pipe.FilterFunc(func(x pipe.Context, p pipe.Page) (pipe.Page, error) {
		kind := p.Doc.Find("#short-nav code").Text()
		kind = strings.TrimPrefix(kind, `import "`)
		kind = strings.TrimSuffix(kind, `"`)

		// Extract function, type, and constant names
		p.Doc.Find("#manual-nav a[href^='#']").Each(func(_ int, s *goquery.Selection) {
			href := s.AttrOr("href", "")
			if href == "" {
				return
			}
			p.Marks = append(p.Marks, pipe.Mark{
				Kind: kind,
				Name: strings.TrimSpace(s.Text()),
				Href: href,
			})
		})

		return p, nil
	})
}
