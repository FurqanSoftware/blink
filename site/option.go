package site

type Option interface {
	Apply(*Site)
}

type OptionFunc func(*Site)

func (f OptionFunc) Apply(c *Site) {
	f(c)
}

func Title(title string) Option {
	return OptionFunc(func(s *Site) {
		s.title = title
	})
}

func Attribution(contentMd string) Option {
	return OptionFunc(func(s *Site) {
		s.attributionContentMd = contentMd
	})
}

func TrimPathPrefix(prefix string) Option {
	return OptionFunc(func(s *Site) {
		s.trimPathPrefix = prefix
	})
}
