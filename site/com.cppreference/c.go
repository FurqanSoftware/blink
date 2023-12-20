package cppreference

import (
	"regexp"

	"github.com/FurqanSoftware/blink/crawl/web"
	"github.com/FurqanSoftware/blink/pipe"
	"github.com/FurqanSoftware/blink/site"
)

func C() site.Site {
	baseURL := "https://en.cppreference.com/w/c"
	urlFilters := []*regexp.Regexp{
		regexp.MustCompile(`^https://en\.cppreference\.com/w/c($|/)`),
	}
	disallowedURLFilters := []*regexp.Regexp{
		regexp.MustCompile(`/experimental`),
	}
	disallowedPaths := []string{
		"/w/c/language/history",
	}
	return site.New(
		"com.cppreference/c",
		web.New(
			baseURL,
			web.AllowedDomains("en.cppreference.com"),
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
				pipe.Container("#content"),
				pipe.RewriteURLs(
					"com.cppreference/c",
					baseURL,
				).
					WithURLFilters(urlFilters...).
					WithDisallowedURLFilters(disallowedURLFilters...).
					WithDisallowedPaths(disallowedPaths...),
				pipe.SyntaxHighlight(),
				pipe.CleanClassNames().
					WithPreserveClasses(
						"t-dcl-begin",
						"t-dsc-header",
						"t-mark-rev",
						"t-li1",
					),
				// pipe.CleanStyles(),
			),
		),
		site.TrimPathPrefix("/w/c"),
		site.Title("C Reference"),
		site.Attribution(`© cppreference.com
Licensed under the Creative Commons Attribution-ShareAlike Unported License v3.0.`),
	)
}

func init() {
	site.Register(C())
}
