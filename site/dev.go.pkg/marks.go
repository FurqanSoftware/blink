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
			name := strings.TrimSpace(s.Text())
			sortkey := name
			switch {
			case strings.HasPrefix(name, "func "):
				sortkey = strings.TrimPrefix(name, "func ")
				if strings.HasPrefix(sortkey, "(") {
					typeright := strings.Index(sortkey, ")")
					typeleft := strings.LastIndexAny(sortkey[:typeright], "( *")
					nameright := strings.Index(sortkey[typeright+1:], "(")
					sortkey = sortkey[typeleft+1:typeright] + sortkey[typeright+2:typeright+1+nameright]
				} else {
					nameright := strings.Index(sortkey, "(")
					sortkey = sortkey[:nameright]
				}
				sortkey = "1" + sortkey
			case strings.HasPrefix(name, "type "):
				sortkey = "1" + strings.TrimPrefix(name, "type ")
			default:
				sortkey = "0" + sortkey
			}
			p.Marks = append(p.Marks, pipe.Mark{
				Kind:    kind,
				Name:    name,
				Href:    href,
				SortKey: sortkey,
			})
		})

		return p, nil
	})
}
