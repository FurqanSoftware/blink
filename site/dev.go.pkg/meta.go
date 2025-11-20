package pkggodev

import (
	"strings"

	"github.com/FurqanSoftware/blink/pipe"
)

func Meta() pipe.Filter {
	return pipe.FilterFunc(func(x pipe.Context, p pipe.Page) (pipe.Page, error) {
		// Try to get title from the first mark (package name)
		if len(p.Marks) > 0 {
			p.Meta.Title = p.Marks[0].Name
		} else {
			// Fall back to extracting from the page title
			pageTitle := strings.TrimSpace(p.Doc.Find("title").First().Text())

			// Clean up the title
			if strings.Contains(pageTitle, " - ") {
				parts := strings.Split(pageTitle, " - ")
				if len(parts) > 0 {
					pageTitle = strings.TrimSpace(parts[0])
				}
			}

			// Remove "package " prefix if present
			if strings.HasPrefix(pageTitle, "package ") {
				pageTitle = strings.TrimPrefix(pageTitle, "package ")
			}

			// Remove " Go Packages" suffix if present
			if strings.HasSuffix(pageTitle, " Go Packages") {
				pageTitle = strings.TrimSuffix(pageTitle, " Go Packages")
			}

			p.Meta.Title = pageTitle
		}

		// If title is still empty, use a default based on the URL path
		if p.Meta.Title == "" {
			pathParts := strings.Split(strings.Trim(x.URL.Path, "/"), "/")
			if len(pathParts) > 0 {
				p.Meta.Title = pathParts[len(pathParts)-1]
			} else {
				p.Meta.Title = "Go Documentation"
			}
		}

		return p, nil
	})
}
