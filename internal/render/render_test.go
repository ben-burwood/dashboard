package render

import (
	"strings"
	"testing"

	"github.com/ben-burwood/dashboard/internal/config"
	"github.com/ben-burwood/dashboard/internal/config/service"
	"github.com/ben-burwood/dashboard/internal/config/tag"
)

func sampleConfig(t *testing.T) *config.Config {
	t.Helper()
	cfg := &config.Config{
		Title:   "Configured Dashboard",
		Favicon: "mdi:view-dashboard",
		Tags: tag.Tags{
			{Name: "tag1", Color: "primary"},
			{Name: "tag2", Color: "secondary"},
			{Name: "tag3", Color: "success"},
		},
		Services: []*service.Service{
			{Title: "Service One", Link: "https://service1.example.com", Group: "Group A",
				Icon: service.Icon{Name: "simple-icons:proxmox", Color: "#E57334"},
				Tags: []string{"tag1", "tag2", "unknown"}},
			{Title: "Service Two", Link: "https://service2.example.com", Group: "Group B",
				Icon: service.Icon{Name: "mdi:server"}, Tags: []string{"tag2"}},
			{Title: "Service Three", Link: "https://service3.example.com", Group: "Group A",
				Icon: service.Icon{Name: "mdi:cloud"}, Tags: []string{"tag3"}},
			{Title: "Service Four", Link: "https://service4.example.com",
				Icon: service.Icon{Name: "mdi:database"}},
		},
	}
	for _, svc := range cfg.Services {
		if err := svc.ValidateAndSetDefaults(); err != nil {
			t.Fatalf("invalid service %q: %v", svc.Title, err)
		}
	}
	return cfg
}

func TestRenderContainsCoreContent(t *testing.T) {
	html, err := Render(sampleConfig(t))
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	out := string(html)

	wants := []string{
		"<title>Configured Dashboard</title>",
		`<h2 class="page-title">Configured Dashboard</h2>`,
		// favicon resolved from an Iconify name to the SVG API URL
		`<link rel="icon" href="https://api.iconify.design/mdi/view-dashboard.svg">`,
		// service links
		`href="https://service1.example.com"`,
		`href="https://service4.example.com"`,
		// icon + trusted color emitted verbatim (no ZgotmplZ)
		`icon="simple-icons:proxmox" style="color: #E57334"`,
		// default group name for a service with no group
		`>default</h2>`,
		// tag pills resolved to color variants
		`<span class="badge badge-primary">tag1</span>`,
		`<span class="badge badge-secondary">tag2</span>`,
	}
	for _, w := range wants {
		if !strings.Contains(out, w) {
			t.Errorf("rendered HTML missing %q", w)
		}
	}

	if strings.Contains(out, "ZgotmplZ") {
		t.Error("rendered HTML contains ZgotmplZ (a value was filtered by html/template escaping)")
	}
	// unknown tag name must be dropped, not rendered
	if strings.Contains(out, "unknown") {
		t.Error("unknown tag name should have been dropped from the output")
	}
}

func TestRenderGroupsPreserveFirstOccurrenceOrder(t *testing.T) {
	html, err := Render(sampleConfig(t))
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	out := string(html)

	// Group A appears first, Group B second, then the default group — even though
	// Service Three (Group A) comes after Service Two (Group B) in the service list.
	a := strings.Index(out, `>Group A</h2>`)
	b := strings.Index(out, `>Group B</h2>`)
	d := strings.Index(out, `>default</h2>`)
	if a < 0 || b < 0 || d < 0 {
		t.Fatalf("missing a group heading: A=%d B=%d default=%d", a, b, d)
	}
	if !(a < b && b < d) {
		t.Errorf("group order wrong: Group A=%d, Group B=%d, default=%d", a, b, d)
	}

	// Service Three should be grouped under Group A (before Group B renders).
	three := strings.Index(out, "Service Three")
	if three < 0 || !(a < three && three < b) {
		t.Errorf("Service Three not grouped under Group A: A=%d three=%d B=%d", a, three, b)
	}
}

func TestFaviconHref(t *testing.T) {
	cases := map[string]string{
		"":                        "",
		"  ":                      "",
		"mdi:view-dashboard":      "https://api.iconify.design/mdi/view-dashboard.svg",
		"https://x.example/i.png": "https://x.example/i.png",
		"/favicon.ico":            "/favicon.ico",
		"./icon.svg":              "./icon.svg",
		"data:image/png;base64,A": "data:image/png;base64,A",
		"favicon.ico":             "favicon.ico",
	}
	for in, want := range cases {
		if got := faviconHref(in); got != want {
			t.Errorf("faviconHref(%q) = %q, want %q", in, got, want)
		}
	}
}
