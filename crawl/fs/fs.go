package fs

import "github.com/FurqanSoftware/blink/crawl"

type Crawler struct{}

var _ crawl.Crawler = Crawler{}

func (c Crawler) Run(entityCh chan<- crawl.Entity) error {
	panic("unimplemented")
}
