package pipe

type Pipeline struct {
	filters []Filter
}

func New(options ...Option) Pipeline {
	p := Pipeline{}
	for _, o := range options {
		o.Apply(&p)
	}
	return p
}

func (p Pipeline) Process(x Context, page Page) (Page, error) {
	for _, f := range p.filters {
		var err error
		page, err = f.Apply(x, page)
		if err != nil {
			return page, err
		}
	}
	return page, nil
}
