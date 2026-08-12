package render

import "strings"

// faviconHref resolves the configured page icon to an href for <link rel="icon">.
//
// It accepts, in order of detection:
//   - a URL, absolute/relative path, or data: URI -> used verbatim;
//   - an Iconify name (prefix:name, as used for service icons) -> served
//     self-contained from the Iconify SVG API (no local file needed);
//   - anything else -> treated as a filename relative to the page.
//
// It returns "" when unset (no favicon link is emitted).
func faviconHref(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	if strings.Contains(value, "://") ||
		strings.HasPrefix(value, "/") ||
		strings.HasPrefix(value, "./") ||
		strings.HasPrefix(value, "../") ||
		strings.HasPrefix(value, "data:") {
		return value
	}
	if prefix, name, found := strings.Cut(value, ":"); found {
		return "https://api.iconify.design/" + prefix + "/" + name + ".svg"
	}
	return value
}
