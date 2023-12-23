package pythondocs

import (
	"regexp"
	"strings"

	"github.com/FurqanSoftware/blink/pipe"
)

func Marks() pipe.Filter {
	kindReplacements := map[string]string{
		"contextvars — Context Variables":          "Context Variables",
		"Cryptographic":                            "Cryptography",
		"Custom Interpreters":                      "Interpreters",
		"Data Compression & Archiving":             "Data Compression",
		"email — An email & MIME handling package": "Email",
		"Generic Operating System":                 "Operating System",
		"Graphical User Interfaces with Tk":        "Tk",
		"Internet Data Handling":                   "Internet Data",
		"Internet Protocols & Support":             "Internet",
		"Interprocess Communication & Networking":  "Networking",
		"Program Frameworks":                       "Frameworks",
		"Structured Markup Processing Tools":       "Structured Markup",
	}

	headingNumberRegexp := regexp.MustCompile(`\A[\d\.]+\s+`)
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
		case strings.HasPrefix(x.Path, "/reference"):
			kind = "Language Reference"
		case strings.HasPrefix(x.Path, "/c-api"):
			kind = "Python/C API"
		case strings.HasPrefix(x.Path, "/tutorial"):
			kind = "Tutorial"
		case strings.HasPrefix(x.Path, "/distributing"):
			kind = "Software Packaging & Distribution"
		case strings.HasPrefix(x.Path, "/distutils"):
			kind = "Software Packaging & Distribution"
		case strings.HasPrefix(x.Path, "/glossary"):
			kind = "Glossary"
		case !strings.HasPrefix(x.Path, "/library/") || strings.HasPrefix(x.Path, "/library/index"):
			kind = "Basics"
		case strings.HasPrefix(x.Path, "/library/logging"):
			kind = "Logging"
		case strings.HasPrefix(x.Path, "/library/asyncio"):
			kind = "Asynchronous I/O"
		}

		kind = p.Doc.Find(`.related a[accesskey="U"]`).Text()

		if kind == "The Python Standard Library" {
			kind = p.Doc.Find("h1").Text()
		} else if strings.Contains(kind, "I/O") || name == "select" || name == "selectors" {
			kind = "Input/output"
		} else if strings.HasPrefix(kind, "19") {
			kind = "Internet Data Handling"
		}

		kind = headingNumberRegexp.ReplaceAllString(kind, "")
		kind = strings.ReplaceAll(kind, "\u00b6", "")
		kind = strings.ReplaceAll(kind, " and ", " & ")
		for _, k := range []string{" Services", " Modules", " Specific", "Python "} {
			kind = strings.ReplaceAll(kind, k, "")
		}

		if v, ok := kindReplacements[kind]; ok {
			kind = v
		}

		p.Marks = append(p.Marks, pipe.Mark{
			Kind: kind,
			Name: name,
		})

		return p, nil
	})
}
