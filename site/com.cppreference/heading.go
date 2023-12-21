package cppreference

import (
	"github.com/FurqanSoftware/blink/pipe"
)

func Heading(url string, heading string) pipe.Filter {
	return pipe.FilterFunc(func(x pipe.Context, p pipe.Page) (pipe.Page, error) {
		if x.URL.String() == url {
			p.Doc.Find("h1").SetText(heading)
		}
		return p, nil
	})
}
