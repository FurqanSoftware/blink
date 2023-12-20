package crawl

type Crawler interface {
	Run(entityCh chan<- Entity) error
}
