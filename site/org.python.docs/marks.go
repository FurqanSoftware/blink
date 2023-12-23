package pythondocs

import (
	"regexp"
	"strings"

	"github.com/FurqanSoftware/blink/pipe"
)

func Marks() pipe.Filter {
	headingNumberRegexp := regexp.MustCompile(`\A[\d\.]+ `)
	dashSuffixRegexp := regexp.MustCompile(` [` + "\u2013\u2014" + `].+\z`)

	return pipe.FilterFunc(func(x pipe.Context, p pipe.Page) (pipe.Page, error) {

		name := p.Doc.Find("h1").Text()
		name = headingNumberRegexp.ReplaceAllString(name, "")
		name = strings.ReplaceAll(name, "\u00b6", "")
		name = dashSuffixRegexp.ReplaceAllString(name, "")
		name = strings.ReplaceAll(name, "Built-in", "")
		name = strings.TrimSpace(name)

		kind := ""
		switch {
		case strings.HasPrefix(x.URL.Path, "/reference"):
			kind = "Language Reference"
		case strings.HasPrefix(x.URL.Path, "/c-api"):
			kind = "Python/C API"
		case strings.HasPrefix(x.URL.Path, "/tutorial"):
			kind = "Tutorial"
		case strings.HasPrefix(x.URL.Path, "/distributing"):
			kind = "Software Packaging & Distribution"
		case strings.HasPrefix(x.URL.Path, "/distutils"):
			kind = "Software Packaging & Distribution"
		case strings.HasPrefix(x.URL.Path, "/glossary"):
			kind = "Glossary"
		case !strings.HasPrefix(x.URL.Path, "/library/") || strings.HasPrefix(x.URL.Path, "/library/index"):
			kind = "Basics"
		case strings.HasPrefix(x.URL.Path, "/library/logging"):
			kind = "Logging"
		case strings.HasPrefix(x.URL.Path, "/library/asyncio"):
			kind = "Asynchronous I/O"
		default:
			kind = "???"
		}

		// type = at_css('.related a[accesskey="U"]').content

		// if type == 'The Python Standard Library'
		//   type = at_css('h1').content
		// elsif type.include?('I/O') || %w(select selectors).include?(name)
		//   type = 'Input/ouput'
		// elsif type.start_with? '19'
		//   type = 'Internet Data Handling'
		// end

		// type.remove! %r{\A\d+\.\s+} unless include_h2? # remove list number
		// type.remove! "\u{00b6}" # remove paragraph character
		// type.sub! ' and ', ' & '
		// [' Services', ' Modules', ' Specific', 'Python '].each { |str| type.remove!(str) }

		p.Marks = append(p.Marks, pipe.Mark{
			Kind: kind,
			Name: name,
		})

		return p, nil
	})
}
