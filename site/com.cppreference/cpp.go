package cppreference

import (
	"regexp"

	"github.com/FurqanSoftware/blink/crawl/web"
	"github.com/FurqanSoftware/blink/pipe"
	"github.com/FurqanSoftware/blink/site"
)

func Cpp() site.Site {
	baseURL := "https://en.cppreference.com/w/cpp"
	urlFilters := []*regexp.Regexp{
		regexp.MustCompile(`^https://en\.cppreference\.com/w/cpp($|/)`),
	}
	disallowedURLFilters := []*regexp.Regexp{
		regexp.MustCompile(`/experimental`),
	}
	disallowedPaths := []string{
		"/w/cpp/language/extending_std",
		"/w/cpp/language/history",
		"/w/cpp/regex/ecmascript",
		"/w/cpp/regex/regex_token_iterator/operator_cmp",
	}
	return site.New(
		"com.cppreference/cpp",
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
					"com.cppreference/cpp",
					baseURL,
				).
					WithURLFilters(urlFilters...).
					WithDisallowedURLFilters(disallowedURLFilters...).
					WithDisallowedPaths(disallowedPaths...),
				pipe.CleanClassName().
					WithPreserveClasses(
						"t-dcl-begin",
						"t-dsc-header",
						"t-mark-rev",
						"t-li1",
					),
				pipe.CleanStyle(),
				pipe.SyntaxHighlight(),

				pipe.If(pipe.IsURL(baseURL)).
					Then(
						Heading("C++ Programming Language"),
					),
			),
		),
		site.TrimPathPrefix("/w/cpp"),
		site.Title("C++ Reference"),
		site.Attribution(`© cppreference.com
Licensed under the Creative Commons Attribution-ShareAlike Unported License v3.0.`),
	)
}

func init() {
	site.Register(Cpp())
}
