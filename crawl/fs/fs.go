package fs

import (
	"archive/zip"
	"errors"
	"io/fs"
	"net/http"
	"net/url"

	"github.com/FurqanSoftware/blink/crawl"
	"github.com/FurqanSoftware/blink/crawl/web"
)

type Crawler struct {
	fs         fs.FS
	startPath  string
	webOptions []web.Option
}

var _ crawl.Crawler = Crawler{}

func New(fs fs.FS, startPath string, options ...Option) Crawler {
	c := Crawler{
		fs:        fs,
		startPath: startPath,
	}
	for _, o := range options {
		o.Apply(&c)
	}
	return c
}

func NewFromZIP(zipName string, startPath string, options ...Option) (Crawler, error) {
	zipFS, err := zip.OpenReader(zipName)
	if err != nil {
		return Crawler{}, err
	}
	return New(zipFS, startPath, options...), nil
}

func toHTTPError(err error) (msg string, httpStatus int) {
	if errors.Is(err, fs.ErrNotExist) {
		return "404 page not found", http.StatusNotFound
	}
	if errors.Is(err, fs.ErrPermission) {
		return "403 Forbidden", http.StatusForbidden
	}
	return "500 Internal Server Error", http.StatusInternalServerError
}

func (c Crawler) Run(entityCh chan<- crawl.Entity) error {
	go func() {
		err := http.ListenAndServe(":24488", FileServer{
			fs: c.fs,
		})
		if err != nil {
			panic(err)
		}
	}()
	url, err := url.Parse("http://localhost:24488")
	if err != nil {
		return err
	}
	url.Path = c.startPath
	webOptions := c.webOptions
	webOptions = append(webOptions, web.DisableCache())
	webCrawler := web.New(url.String(), webOptions...)
	return webCrawler.Run(entityCh)
}
