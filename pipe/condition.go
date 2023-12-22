package pipe

type Condition struct {
	clause                   Clause
	thenFilters, elseFilters []Filter
}

var _ Filter = Condition{}

func If(clause Clause) Condition {
	return Condition{clause: clause}
}

func (c Condition) Then(filters ...Filter) Condition {
	c.thenFilters = append(c.thenFilters, filters...)
	return c
}

func (c Condition) Else(filters ...Filter) Condition {
	c.elseFilters = append(c.elseFilters, filters...)
	return c
}

func (c Condition) Apply(x Context, p Page) (Page, error) {
	filters := c.elseFilters
	if c.clause.IsTrue(x, p) {
		filters = c.thenFilters
	}
	for _, f := range filters {
		var err error
		p, err = f.Apply(x, p)
		if err != nil {
			return p, err
		}
	}
	return p, nil
}

type Clause interface {
	IsTrue(Context, Page) bool
}

type ClauseFunc func(x Context, p Page) bool

func (f ClauseFunc) IsTrue(x Context, p Page) bool {
	return f(x, p)
}

func IsURL(url string) Clause {
	return ClauseFunc(func(x Context, p Page) bool {
		return x.URL.String() == url
	})
}
