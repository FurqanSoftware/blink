package pkggodev

import (
	"regexp"

	"github.com/FurqanSoftware/blink/crawl/web"
	"github.com/FurqanSoftware/blink/pipe"
	"github.com/FurqanSoftware/blink/site"
)

func Go() site.Site {
	baseURL := "http://localhost:6060/pkg/" // Run `godoc -http=:6060` first
	urlFilters := []*regexp.Regexp{
		regexp.MustCompile(`^http://localhost:6060/pkg/([a-zA-Z0-9]+/)*(#.+)?$`),
	}
	disallowedURLFilters := []*regexp.Regexp{}
	disallowedPaths := []string{}

	return site.New(
		"dev.golang.pkg/go",
		web.New(
			baseURL,
			web.AllowedDomains("localhost"),
			web.URLFilters(urlFilters...),
			web.DisallowedURLFilters(disallowedURLFilters...),
			web.DisallowedPaths(disallowedPaths...),
		),
		pipe.New(
			pipe.Filters(
				pipe.DefaultMeta(),
				pipe.SanitizeHTML(),
				Marks(),
				Meta(),
				CleanHTML(),
				pipe.Container("#page"),
				pipe.RewriteURLs(
					"dev.golang.pkg/go",
					baseURL,
				).
					WithURLFilters(urlFilters...).
					WithDisallowedURLFilters(disallowedURLFilters...).
					WithDisallowedPaths(disallowedPaths...),
				pipe.CleanClassName(),
				// pipe.CleanStyle(),
				pipe.SyntaxHighlight(),

				pipe.If(pipe.IsURL(baseURL)).
					Then(
						Heading("Go Standard Library"),
					),
			),
		),
		site.TrimPathPrefix("/pkg"),
		site.Title("Go Standard Library Reference"),
		site.Attribution(`© The Go Authors
Licensed under the BSD License.`),
	)
}

func init() {
	site.Register(Go())
}
