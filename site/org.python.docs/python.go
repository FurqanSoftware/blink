package pythondocs

import (
	"regexp"

	"github.com/FurqanSoftware/blink/crawl/fs"
	"github.com/FurqanSoftware/blink/crawl/web"
	"github.com/FurqanSoftware/blink/pipe"
	"github.com/FurqanSoftware/blink/site"
)

func Python(versionShort, zipName, basePath string) site.Site {
	urlFilters := []*regexp.Regexp{
		regexp.MustCompile(`^http://localhost:24488` + regexp.QuoteMeta(basePath) + `/`),
	}
	disallowedURLFilters := []*regexp.Regexp{
		regexp.MustCompile(`/whatsnew/`),
	}
	disallowedPaths := []string{
		basePath + "/library/2to3.html",
		basePath + "/library/formatter.html",
		basePath + "/library/intro.html",
		basePath + "/library/undoc.html",
		basePath + "/library/unittest.mock-examples.html",
		basePath + "/library/sunau.html",
	}
	crawler, err := fs.NewFromZIP(
		zipName,
		basePath+"/index.html",
		fs.WebOptions(
			web.AllowedDomains("localhost"),
			web.URLFilters(urlFilters...),
			web.DisallowedURLFilters(disallowedURLFilters...),
			web.DisallowedPaths(disallowedPaths...),
		),
	)
	if err != nil {
		panic(err)
	}
	return site.New(
		"org.python.docs/python-"+versionShort,
		crawler,
		pipe.New(
			pipe.Filters(
				pipe.DefaultMeta(),
				pipe.MapOriginalURLPrefix("http://localhost:24488"+basePath, "https://docs.python.org/"+versionShort),
				pipe.SanitizeHTML(),
				Marks(),
				Meta(),
				CleanHTML(),
				pipe.Container(`.body[role="main"]`),
				pipe.RewriteURLs(
					"org.python.docs/python-"+versionShort,
					"http://localhost:24488"+basePath,
				).
					WithURLFilters(urlFilters...).
					WithDisallowedURLFilters(disallowedURLFilters...).
					WithDisallowedPaths(disallowedPaths...).
					WithPrefixMappings(
						"http://localhost:24488"+basePath, "https://docs.python.org/"+versionShort,
					),
				pipe.SyntaxHighlight(),
				pipe.CleanClassName().
					WithPreserveClasses(
						"t-dcl-begin",
						"t-dsc-header",
						"t-mark-rev",
						"t-li1",
					),
				pipe.CleanStyle(),

				pipe.If(pipe.IsPath(basePath+"/")).
					Then(
						Heading("Python "+versionShort),
					),
			),
		),
		site.TrimPathPrefix(basePath),
		site.Title("Python "+versionShort+" Reference"),
		site.Attribution(`© 2001–2023 Python Software Foundation
Licensed under the Python Software Foundation License Version 2.`),
		site.Redirects(
			"", "/index",
		),
	)
}

func init() {
	site.Register(Python("3.12", "python-3.12.1-docs-html.zip", "/python-3.12.1-docs-html"))
	site.Register(Python("3.11", "python-3.11.7-docs-html.zip", "/python-3.11.7-docs-html"))
}
