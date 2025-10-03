package tag

import "errors"

var (
	ErrInvalidTagConfiguration = errors.New("invalid Tag Configuration")
)

type Tag struct {
	Name  string `yaml:"name"`
	Color string `yaml:"color"`
}

type Tags []*Tag

// EnsureUnique checks that all tags have unique names
func (tags Tags) EnsureUnique() error {
	uniqueTags := make(map[string]Tag)
	for _, tag := range tags {
		if _, exists := uniqueTags[tag.Name]; exists {
			return ErrInvalidTagConfiguration
		}
		uniqueTags[tag.Name] = *tag
	}
	return nil
}

// Lookup returns the Tag with the given name
func (tags Tags) Lookup(name string) (*Tag, error) {
	for _, tag := range tags {
		if tag.Name == name {
			return tag, nil
		}
	}
	return nil, errors.New("tag not found")
}
