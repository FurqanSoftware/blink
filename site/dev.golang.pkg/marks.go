package pkggodev

import (
	"regexp"
	"strings"

	"github.com/FurqanSoftware/blink/pipe"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func Marks() pipe.Filter {
	titleCaser := cases.Title(language.English)

	// Patterns to clean up package names and function signatures
	genericParamsRegexp := regexp.MustCompile(`\[[^\]]+\]`)
	paramsRegexp := regexp.MustCompile(`\([^)]*\)`)
	repeatedSpaceRegexp := regexp.MustCompile(`\s+`)

	return pipe.FilterFunc(func(x pipe.Context, p pipe.Page) (pipe.Page, error) {
		kind := ""

		// Determine the kind of documentation page
		if strings.HasPrefix(x.URL.Path, "/std") {
			kind = "Standard Library"
		} else if strings.Contains(x.URL.Path, "/builtin") {
			kind = "Built-in"
		} else {
			// Extract package category from breadcrumbs or header
			breadcrumb := p.Doc.Find(".go-Breadcrumb").Text()
			if breadcrumb != "" {
				parts := strings.Split(breadcrumb, "/")
				if len(parts) > 1 {
					kind = titleCaser.String(parts[0])
				}
			}
		}

		if kind == "" {
			kind = "Packages"
		}

		// Extract package name from the page title or header
		packageName := ""
		pageTitle := p.Doc.Find("h1").First().Text()
		if pageTitle != "" {
			packageName = strings.TrimSpace(pageTitle)
			// Remove "package" prefix if present
			if strings.HasPrefix(packageName, "package ") {
				packageName = strings.TrimPrefix(packageName, "package ")
			}
		}

		// If we have a package name, add it as a mark
		if packageName != "" && packageName != "standard library" {
			p.Marks = append(p.Marks, pipe.Mark{
				Kind: kind,
				Name: packageName,
			})
		}

		// Extract function, type, and constant names
		p.Doc.Find("#pkg-index a[href^='#']").Each(func(_ int, s *goquery.Selection) {
			name := strings.TrimSpace(s.Text())
			if name == "" {
				return
			}

			// Determine the type based on the link context
			markKind := kind
			href := s.AttrOr("href", "")

			if strings.Contains(href, "func.") {
				markKind = "Functions"
			} else if strings.Contains(href, "type.") {
				markKind = "Types"
			} else if strings.Contains(href, "const.") || strings.Contains(href, "var.") {
				markKind = "Constants"
			}

			// Clean up the name
			name = genericParamsRegexp.ReplaceAllString(name, "")
			name = paramsRegexp.ReplaceAllString(name, "")
			name = repeatedSpaceRegexp.ReplaceAllString(name, " ")
			name = strings.TrimSpace(name)

			// Skip empty names, single characters, or unwanted symbols
			if name == "" || len(name) <= 1 || name == "¶" || strings.Contains(name, "¶") {
				return
			}

			p.Marks = append(p.Marks, pipe.Mark{
				Kind: markKind,
				Name: name,
			})
		})

		// Extract type methods
		p.Doc.Find("#pkg-index").Find("a[href*='method.']").Each(func(_ int, s *goquery.Selection) {
			name := strings.TrimSpace(s.Text())
			if name == "" {
				return
			}

			// Clean up method names
			if strings.Contains(name, ".") {
				parts := strings.Split(name, ".")
				if len(parts) > 1 {
					name = parts[len(parts)-1] // Get method name after the type
				}
			}

			name = genericParamsRegexp.ReplaceAllString(name, "")
			name = paramsRegexp.ReplaceAllString(name, "")
			name = repeatedSpaceRegexp.ReplaceAllString(name, " ")
			name = strings.TrimSpace(name)

			if name != "" && len(name) > 1 && name != "¶" && !strings.Contains(name, "¶") {
				p.Marks = append(p.Marks, pipe.Mark{
					Kind: "Methods",
					Name: name,
				})
			}
		})

		return p, nil
	})
}
