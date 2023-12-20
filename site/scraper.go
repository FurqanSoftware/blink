package site

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/FurqanSoftware/blink/crawl"
	"github.com/FurqanSoftware/blink/pipe"
	"github.com/PuerkitoBio/goquery"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/html"
	"github.com/yuin/goldmark"
	goldmarkhtml "github.com/yuin/goldmark/renderer/html"
	"golang.org/x/sync/errgroup"
	"gopkg.in/yaml.v2"
)

type Scraper struct{}

func (r Scraper) Run(ctx context.Context, s Site) error {
	log.Printf("[S] Starting scrape %s", s.key)

	eg := errgroup.Group{}

	entityCh := make(chan crawl.Entity)

	eg.Go(func() error {
		err := s.crawler.Run(entityCh)
		close(entityCh)
		return err
	})

	m := minify.New()
	m.Add("text/html", &html.Minifier{})

	marktree := map[string][]Mark{}
	redirects := map[string]string{}

L:
	for {
		select {
		case entity, ok := <-entityCh:
			if !ok {
				break L
			}

			switch entity := entity.(type) {
			case crawl.Page:
				log.Printf("[S] Scrape %s", entity.Request.URL)
				d, err := goquery.NewDocumentFromReader(bytes.NewReader(entity.Content))
				if err != nil {
					return err
				}

				x := pipe.Context{
					URL: entity.Request.URL,
				}
				p, err := s.pipeline.Process(x, pipe.Page{
					Doc: d,
				})
				if err != nil {
					return err
				}

				for _, m := range p.Marks {
					marktree[m.Kind] = append(marktree[m.Kind], Mark{
						Name: m.Name,
						Path: strings.TrimPrefix(x.URL.Path, s.trimPathPrefix),
					})
				}

				err = r.writePage(s, x, p, m)
				if err != nil {
					return err
				}

			case crawl.Redirect:
				redirects[strings.TrimPrefix(entity.Request.URL.Path, s.trimPathPrefix)] = strings.TrimPrefix(entity.To.Path, s.trimPathPrefix)
			}

		case <-ctx.Done():
			break L
		}
	}

	marks := []Mark{}
	for k, v := range marktree {
		if k == "" {
			continue
		}
		slices.SortStableFunc(v, func(a, b Mark) int {
			c := strings.Compare(strings.ToLower(a.Name), strings.ToLower(b.Name))
			if c == 0 {
				return strings.Compare(a.Path, b.Path)
			}
			return c
		})
		v = slices.CompactFunc(v, func(a, b Mark) bool {
			return a.Name == b.Name && a.Path == b.Path
		})
		marks = append(marks, Mark{
			Name:  k,
			Marks: v,
		})
	}
	slices.SortStableFunc(marks, func(a, b Mark) int {
		return strings.Compare(strings.ToLower(a.Name), strings.ToLower(b.Name))
	})

	err := r.writeSiteJSON(s, marks, redirects)
	if err != nil {
		return err
	}

	return nil
}

func (r Scraper) writePage(s Site, x pipe.Context, p pipe.Page, m *minify.M) error {
	outPath := filepath.Join("out", s.key, strings.TrimPrefix(x.URL.Path, s.trimPathPrefix), "index.html")
	dirPath := filepath.Dir(outPath)
	err := os.MkdirAll(dirPath, 0755)
	if err != nil {
		return err
	}
	html, err := p.Doc.Find("body").Html()
	if err != nil {
		return err
	}
	f, err := os.Create(outPath)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(f, "---")
	if err != nil {
		return err
	}
	err = yaml.NewEncoder(f).Encode(p.Meta)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(f, "---")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(f)
	if err != nil {
		return err
	}
	err = m.Minify("text/html", f, strings.NewReader(html))
	if err != nil {
		return err
	}
	return f.Close()
}

func (r Scraper) writeSiteJSON(s Site, marks []Mark, redirects map[string]string) error {
	outPath := filepath.Join("out", s.key, "site.json")
	dirPath := filepath.Dir(outPath)
	err := os.MkdirAll(dirPath, 0755)
	if err != nil {
		return err
	}

	f, err := os.Create(outPath)
	if err != nil {
		return err
	}

	b := bytes.Buffer{}
	err = goldmark.New(
		goldmark.WithRendererOptions(
			goldmarkhtml.WithHardWraps(),
		),
	).Convert(
		[]byte(s.attributionContentMd),
		&b,
	)
	if err != nil {
		return err
	}

	err = json.NewEncoder(f).Encode(struct {
		Key                string            `json:"key"`
		Title              string            `json:"title"`
		AttributionContent string            `json:"attributionContent"`
		Marks              []Mark            `json:"marks"`
		Redirects          map[string]string `json:"redirects"`
		UpdatedAt          time.Time         `json:"updatedAt"`
	}{
		Key:                s.key,
		Title:              s.title,
		AttributionContent: b.String(),
		Marks:              marks,
		Redirects:          redirects,
		UpdatedAt:          time.Now(),
	})
	if err != nil {
		return err
	}

	return f.Close()
}
