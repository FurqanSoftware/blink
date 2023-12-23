package pythondocs

import (
	"github.com/FurqanSoftware/blink/pipe"
)

func Heading(heading string) pipe.Filter {
	return pipe.FilterFunc(func(x pipe.Context, p pipe.Page) (pipe.Page, error) {
		p.Doc.Find("h1").SetText(heading)
		return p, nil
	})
}
