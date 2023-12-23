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

func Redirects(fromto ...string) Option {
	return OptionFunc(func(s *Site) {
		for i := 0; i < len(fromto); i += 2 {
			s.redirects[fromto[i]] = fromto[i+1]
		}
	})
}
