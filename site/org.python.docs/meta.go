package pythondocs

import (
	"strings"

	"github.com/FurqanSoftware/blink/pipe"
)

func Meta() pipe.Filter {
	return pipe.FilterFunc(func(x pipe.Context, p pipe.Page) (pipe.Page, error) {
		if len(p.Marks) > 0 {
			p.Meta.Title = p.Marks[0].Name
		} else {
			p.Meta.Title = strings.TrimSpace(p.Doc.Find("title").First().Text())
		}
		return p, nil
	})
}
