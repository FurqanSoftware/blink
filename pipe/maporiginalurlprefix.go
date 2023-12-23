package pipe

import "strings"

func MapOriginalURLPrefix(from, to string) Filter {
	return FilterFunc(func(x Context, p Page) (Page, error) {
		if strings.HasPrefix(p.Meta.OriginalURL, from) {
			p.Meta.OriginalURL = to + strings.TrimPrefix(p.Meta.OriginalURL, from)
		}
		return p, nil
	})
}
