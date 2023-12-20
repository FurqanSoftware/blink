package site

import "errors"

var registry = map[string]Site{}

func Register(site Site) {
	key := site.key
	_, exists := registry[key]
	if exists {
		panic(ErrDuplicate)
	}
	registry[key] = site
}

func Get(key string) Site {
	return registry[key]
}

var (
	ErrDuplicate = errors.New("site: duplicate key")
)
