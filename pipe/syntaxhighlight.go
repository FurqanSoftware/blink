package pipe

import (
	"bytes"
	"log"

	"github.com/PuerkitoBio/goquery"
	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
)

func SyntaxHighlight() Filter {
	return FilterFunc(func(_ Context, p Page) (Page, error) {
		p.Doc.Find("pre[data-language]").Each(func(_ int, s *goquery.Selection) {
			l := s.AttrOr("data-language", "")
			lexer := lexers.Get(l)
			if lexer == nil {
				lexer = lexers.Fallback
			}
			style := styles.Get("friendly")
			formatter := html.New(html.Standalone(true))
			iterator, err := lexer.Tokenise(nil, s.Text())
			if err != nil {
				log.Printf("[P] SyntaxHighlight: %s", err)
				return
			}
			b := bytes.Buffer{}
			err = formatter.Format(&b, style, iterator)
			if err != nil {
				log.Printf("[P] SyntaxHighlight: %s", err)
				return
			}
			s.ReplaceWithHtml(b.String())
		})
		return p, nil
	})
}
