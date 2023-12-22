package cppreference

import (
	"strconv"
	"strings"

	"github.com/FurqanSoftware/blink/pipe"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html/atom"
)

func CleanHTML() pipe.Filter {
	codeLineBreakReplacer := strings.NewReplacer(
		"<br>\n", "\n",
		"<br>", "\n",
		"\n</p>\n", "</p>\n",
	)

	return pipe.FilterFunc(func(_ pipe.Context, p pipe.Page) (pipe.Page, error) {
		p.Doc.Find("div > span.source-c, div > span.source-cpp").Each(func(_ int, s *goquery.Selection) {
			p := s.Parent()
			p.AddClass(s.AttrOr("class", ""))
			n := p.Nodes[0]
			n.DataAtom = atom.Pre
			n.Data = "pre"
			html, err := s.Html()
			if err != nil {
				return
			}
			s.SetHtml(codeLineBreakReplacer.Replace(html))
		})

		p.Doc.Find("h1").SetText(p.Doc.Find("h1").Text())

		p.Doc.Find(".t-dcl-rev-aux td[rowspan]").Each(func(_ int, s *goquery.Selection) {
			rowspan, err := strconv.Atoi(s.AttrOr("rowspan", ""))
			if err != nil {
				return
			}
			if rowspan > 3 {
				n := s.ParentsFiltered("tbody").First().Find("tr").Length()
				if n > 1 {
					s.SetAttr("rowspan", strconv.Itoa(n))
				}
			}
		})

		p.Doc.Find("#top, #mw-js-message, #siteSub, #contentSub, .printfooter, .t-navbar, .editsection, #toc, .t-dsc-sep, .t-dcl-sep, #catlinks, .ambox-notice, .mw-cite-backlink, .t-sdsc-sep:first-child:last-child, .t-example-live-link, .t-dcl-rev-num > .t-dcl-rev-aux ~ tr:not(.t-dcl-rev-aux) > td:nth-child(2)").Remove()

		p.Doc.Find(`#bodyContent, .mw-content-ltr, span[style], div[class^="t-ref"], .t-image, th > div, td > div, .t-dsc-see, .mainpagediv, code > b, tbody`).Each(func(_ int, s *goquery.Selection) {
			s.BeforeSelection(s.Contents())
			s.Remove()
		})

		p.Doc.Find("h2 > span[id], h3 > span[id], h4 > span[id], h5 > span[id], h6 > span[id]").Each(func(_ int, s *goquery.Selection) {
			s.Parent().SetAttr("id", s.AttrOr("id", ""))
			s.BeforeSelection(s.Contents())
			s.Remove()
		})

		p.Doc.Find("table[style], th[style], td[style]").RemoveAttr("style")

		p.Doc.Find("pre").Each(func(_ int, s *goquery.Selection) {
			lang := "c"
			if strings.Contains(s.AttrOr("class", ""), "cpp") || strings.Contains(s.Parent().AttrOr("class", ""), "cpp") {
				lang = "cpp"
			}
			s.SetAttr("data-language", lang)
			s.RemoveAttr("class")
			text := s.Text()
			text = strings.ReplaceAll(text, "\t", "        ")
			text = strings.ReplaceAll(text, "\u00a0", " ")
			s.SetText(text)
		})

		return p, nil
	})
}
