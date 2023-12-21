package cppreference

import (
	"regexp"
	"strings"

	"github.com/FurqanSoftware/blink/pipe"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func Marks() pipe.Filter {
	titleCaser := cases.Title(language.English)
	noiseReplacer := strings.NewReplacer(
		"C++ concepts:", "",
		"C++ keywords:", "",
	)
	removeRegexps := []*regexp.Regexp{
		regexp.MustCompile(`\(<.+?>\)`),
	}
	stdlibHeaderRegexp := regexp.MustCompile(`\AStandard library header <(.+)>\z`)
	operatorHeaderRegexp := regexp.MustCompile(`(operator).+?(\(.+?\)|)$`)
	repeatedSpaceRegexp := regexp.MustCompile(`\s+`)

	return pipe.FilterFunc(func(x pipe.Context, p pipe.Page) (pipe.Page, error) {
		kind := ""
		if strings.Contains(p.Doc.Find("#firstHeading").Text(), "C++ keyword") {
			kind = "Keywords"
		} else if strings.Contains(x.URL.Path, "experimental") {
			kind = "Experimental libraries"
		} else if strings.Contains(x.URL.Path, "/language/") {
			kind = "Language"
		} else if strings.Contains(x.URL.Path, "/freestanding/") {
			kind = "Utilities"
		} else {
			kind = p.Doc.Find(".t-navbar > div:nth-child(4) > :first-child").Text()
			kind = strings.ReplaceAll(kind, "library", "")
			kind = strings.ReplaceAll(kind, "utilities", "")
			kind = strings.ReplaceAll(kind, "C++", "")
			kind = strings.TrimSpace(kind)
			kind = titleCaser.String(kind)
		}

		name := p.Doc.Find("#firstHeading").Text()
		name = noiseReplacer.Replace(name)
		if name != "C++ language" {
			name = strings.ReplaceAll(name, "C++", "")
		}
		name = strings.TrimSpace(name)
		for _, re := range removeRegexps {
			name = re.ReplaceAllString(name, "")
		}
		name = repeatedSpaceRegexp.ReplaceAllString(name, " ")

		if stdlibHeaderRegexp.MatchString(name) {
			m := stdlibHeaderRegexp.FindStringSubmatch(name)
			name = m[1]
		}

		if strings.Contains(name, "operator") && strings.Contains(name, ",") {
			name = operatorHeaderRegexp.ReplaceAllString(name, "operators $2")
			if strings.Contains(name, "(") && !strings.HasSuffix(name, ")") {
				name += ")"
			}
			name = strings.ReplaceAll(name, "()", "")
			name = repeatedSpaceRegexp.ReplaceAllString(name, " ")
		}

		for _, name := range strings.Split(name, ",") {
			name = strings.TrimSpace(name)
			if name == "" {
				continue
			}
			p.Marks = append(p.Marks, pipe.Mark{
				Kind: kind,
				Name: name,
			})
		}

		return p, nil
	})
}
