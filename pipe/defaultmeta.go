package pipe

func DefaultMeta() Filter {
	return FilterFunc(func(x Context, p Page) (Page, error) {
		p.Meta.Title = p.Doc.Find("head title").First().Text()
		p.Meta.OriginalURL = x.URL.String()
		return p, nil
	})
}
