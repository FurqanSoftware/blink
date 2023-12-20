package pipe

type Filter interface {
	Apply(Context, Page) (Page, error)
}

type FilterFunc func(Context, Page) (Page, error)

func (f FilterFunc) Apply(x Context, p Page) (Page, error) {
	return f(x, p)
}
