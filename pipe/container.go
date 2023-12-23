package pipe

func Container(selectors ...string) Filter {
	return FilterFunc(func(x Context, p Page) (Page, error) {
		for _, selector := range selectors {
			s := p.Doc.Find(selector)
			if s.Length() == 0 {
				continue
			}
			html, _ := s.First().Html()
			p.Doc.Find("body").SetHtml(html)
			break
		}
		return p, nil
	})
}
