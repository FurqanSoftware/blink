package pipe

type Option interface {
	Apply(*Pipeline)
}

type OptionFunc func(*Pipeline)

func (f OptionFunc) Apply(c *Pipeline) {
	f(c)
}

func Filters(filters ...Filter) Option {
	return OptionFunc(func(p *Pipeline) {
		p.filters = append(p.filters, filters...)
	})
}
