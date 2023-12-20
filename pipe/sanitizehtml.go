package pipe

func SanitizeHTML() Filter {
	return FilterFunc(func(_ Context, p Page) (Page, error) {
		p.Doc.Find("link, style, script").Remove()
		return p, nil
	})
}
