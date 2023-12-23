package pipe

func Container(selector string) Filter {
	return FilterFunc(func(x Context, p Page) (Page, error) {
		html, _ := p.Doc.Find(selector).First().Html()
		p.Doc.Find("body").SetHtml(html)
		return p, nil
	})
}
