// Package render turns a validated config.Config into the dashboard's static HTML.
// It is the single renderer behind both the HTTP server and the static export.
package render

import (
	"bytes"
	"embed"
	"html/template"

	"github.com/ben-burwood/dashboard/internal/config"
	"github.com/ben-burwood/dashboard/internal/config/service"
	"github.com/ben-burwood/dashboard/internal/config/tag"
)

//go:embed templates/*.gohtml
var templatesFS embed.FS

var tmpl = template.Must(template.ParseFS(templatesFS, "templates/*.gohtml"))

// tagView is a resolved tag (name + color variant) for a service card.
type tagView struct {
	Name  string
	Color string
}

// serviceView is a single service card with its tag names already resolved.
type serviceView struct {
	Link string
	// IconName is trusted config. html/template treats the `icon` attribute as
	// URL-typed, so an Iconify name like "simple-icons:proxmox" would be rejected
	// as an unknown URL scheme; template.URL bypasses that filtering.
	IconName  template.URL
	IconColor template.CSS // trusted config value, emitted verbatim in a style attribute
	Title     string
	Tags      []tagView
}

// groupView is a titled section of service cards.
type groupView struct {
	Title    string
	Services []serviceView
}

// pageView is the top-level template model.
type pageView struct {
	Title   string
	Favicon template.URL // trusted, pre-resolved href ("" -> no <link>)
	Groups  []groupView
}

// resolveTags resolves a service's tag names to Tag objects, dropping unknown names.
func resolveTags(tags tag.Tags, svc *service.Service) []tagView {
	var resolved []tagView
	for _, name := range svc.Tags {
		t, err := tags.Lookup(name)
		if err != nil {
			continue
		}
		resolved = append(resolved, tagView{Name: t.Name, Color: t.Color})
	}
	return resolved
}

// buildGroups groups services by their group, preserving first-occurrence order.
func buildGroups(cfg *config.Config) []groupView {
	var groups []groupView
	index := make(map[service.Group]int)
	for _, svc := range cfg.Services {
		sv := serviceView{
			Link:      svc.Link,
			IconName:  template.URL(svc.Icon.Name),
			IconColor: template.CSS(svc.Icon.Color),
			Title:     svc.Title,
			Tags:      resolveTags(cfg.Tags, svc),
		}
		i, ok := index[svc.Group]
		if !ok {
			index[svc.Group] = len(groups)
			groups = append(groups, groupView{Title: string(svc.Group)})
			i = len(groups) - 1
		}
		groups[i].Services = append(groups[i].Services, sv)
	}
	return groups
}

// Render renders the dashboard page for the given config to an HTML byte slice.
func Render(cfg *config.Config) ([]byte, error) {
	page := pageView{
		Title:   cfg.Title,
		Favicon: template.URL(faviconHref(cfg.Favicon)),
		Groups:  buildGroups(cfg),
	}
	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, "index.html.gohtml", page); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
