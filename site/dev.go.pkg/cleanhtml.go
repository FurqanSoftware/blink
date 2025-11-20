package pkggodev

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/FurqanSoftware/blink/pipe"
	"github.com/PuerkitoBio/goquery"
)

func CleanHTML() pipe.Filter {
	return pipe.FilterFunc(func(x pipe.Context, p pipe.Page) (pipe.Page, error) {
		if x.Path == "/" {
			p.Doc.Find("#manual-nav").Remove()
		}

		p.Doc.Find(`
			#stdlib div:has(h2:contains("▹")),
			#pkg-overview div:has(h2:contains("▹")),
			#pkg-index div:has(h2:contains("▹")),
			#thirdparty,
			#other,
			#footer,
			span[title^="Added in"],
			#pkg-examples div:contains("(Expand All)")
		`).Remove()

		p.Doc.Find("#subrepo + p").Remove()
		p.Doc.Find("#subrepo + ul").Remove()
		p.Doc.Find("#subrepo").Remove()

		p.Doc.Find("#community + p").Remove()
		p.Doc.Find("#community + ul").Remove()
		p.Doc.Find("#community").Remove()

		p.Doc.Find("#short-nav, #manual-nav, #pkg-examples").Each(func(_ int, s *goquery.Selection) {
			s.Find("dl").Each(func(_ int, s *goquery.Selection) {
				s.Find("dd").Each(func(_ int, s *goquery.Selection) {
					html, _ := s.Html()
					s.ReplaceWithHtml("<div>" + html + "</div>")
				})
				html, _ := s.Html()
				s.ReplaceWithHtml("<div>" + html + "</div>")
			})
		})

		p.Doc.Find("h2").Each(func(_ int, s *goquery.Selection) {
			s.SetText(strings.TrimSpace(strings.ReplaceAll(s.Text(), "▾", "")))
		})

		p.Doc.Find(`h2:contains("¶")`).Each(func(_ int, s *goquery.Selection) {
			text := s.Text()
			i := strings.Index(text, "¶")
			text = text[:i]
			s.SetText(strings.TrimSpace(text))
		})

		p.Doc.Find(`a:contains("¶")`).Each(func(_ int, s *goquery.Selection) {
			if strings.TrimSpace(s.Text()) != "¶" {
				return
			}
			s.Remove()
		})

		p.Doc.Find(`div[id^="example_"]`).Each(func(_ int, s *goquery.Selection) {
			s.Find("div:first-child").Remove()
			s.Find(`p:contains("▾")`).Each(func(_ int, s *goquery.Selection) {
				s.SetText(strings.TrimSpace(strings.ReplaceAll(s.Text(), "▾", "")))
			})
		})

		paddingleftre, err := regexp.Compile(`padding-left:\s*(\d+)px`)
		if err != nil {
			return p, err
		}
		p.Doc.Find(".pkg-name[style]").Each(func(_ int, s *goquery.Selection) {
			m := paddingleftre.FindStringSubmatch(s.AttrOr("style", ""))
			if len(m) < 2 {
				return
			}
			padleft, _ := strconv.Atoi(m[1])
			if padleft > 0 {
				s.SetAttr("style", fmt.Sprintf("padding-left: %drem;", padleft/16+1))
			} else {
				s.RemoveAttr("style")
			}
		})

		p.Doc.Find("pre").Each(func(_ int, s *goquery.Selection) {
			s.SetAttr("data-language", "go")
			s.RemoveAttr("class")

			text := s.Text()
			text = strings.ReplaceAll(text, "\t", "    ")  // Convert tabs to spaces
			text = strings.ReplaceAll(text, "\u00a0", " ") // Convert nbsp to regular space
			s.SetText(text)
		})

		// // Handle code spans
		// p.Doc.Find("code").Each(func(_ int, s *goquery.Selection) {
		// 	// Remove any styling classes but keep the content
		// 	s.RemoveAttr("class")
		// 	text := s.Text()
		// 	text = strings.ReplaceAll(text, "\u00a0", " ")
		// 	s.SetText(text)
		// })

		// // Clean up tables by removing styling but keep structure
		// p.Doc.Find("table, th, td, tr").Each(func(_ int, s *goquery.Selection) {
		// 	s.RemoveAttr("style")
		// 	s.RemoveAttr("class")
		// })

		// // Clean up headings and preserve hierarchy
		// p.Doc.Find("h1, h2, h3, h4, h5, h6").Each(func(_ int, s *goquery.Selection) {
		// 	s.RemoveAttr("class")
		// 	// Keep id attributes for internal linking if they exist
		// 	text := strings.TrimSpace(s.Text())
		// 	if text != "" {
		// 		s.SetText(text)
		// 	}
		// })

		// // Clean up paragraphs and divs
		// p.Doc.Find("p, div").Each(func(_ int, s *goquery.Selection) {
		// 	s.RemoveAttr("class")
		// 	s.RemoveAttr("style")
		// })

		// // Clean up links - keep internal ones, convert external to text
		// p.Doc.Find("a").Each(func(_ int, s *goquery.Selection) {
		// 	href := s.AttrOr("href", "")
		// 	s.RemoveAttr("class")
		// 	s.RemoveAttr("style")

		// 	// Clean up the link text
		// 	text := strings.TrimSpace(s.Text())
		// 	if text != "" {
		// 		s.SetText(text)
		// 	}

		// 	// Keep only internal links
		// 	if href != "" && !strings.HasPrefix(href, "#") &&
		// 		!strings.HasPrefix(href, "/") && strings.Contains(href, "://") {
		// 		// Convert external links to plain text
		// 		s.RemoveAttr("href")
		// 	}
		// })

		// // Clean up lists
		// p.Doc.Find("ul, ol, li").Each(func(_ int, s *goquery.Selection) {
		// 	s.RemoveAttr("class")
		// 	s.RemoveAttr("style")
		// })

		// // Remove any remaining empty elements after all cleanup
		// p.Doc.Find("*:empty").Each(func(_ int, s *goquery.Selection) {
		// 	// Keep elements that should be empty like br, hr, img
		// 	if s.Nodes[0].DataAtom != atom.Br &&
		// 		s.Nodes[0].DataAtom != atom.Hr &&
		// 		s.Nodes[0].DataAtom != atom.Img {
		// 		s.Remove()
		// 	}
		// })

		return p, nil
	})
}
