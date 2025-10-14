package tag

import (
	"errors"
	"slices"
)

var (
	ErrInvalidTagConfiguration = errors.New("invalid Tag Configuration")
)

var colorVariants = []string{
	"primary",
	"secondary",
	"accent",
	"info",
	"success",
	"warning",
	"error",
	"neutral",
}

type Tag struct {
	Name  string `yaml:"name"`
	Color string `yaml:"color"`
}

type Tags []*Tag

// EnsureColorVariant checks that the tag color is one of the predefined variants
func (tag *Tag) EnsureColorVariant() bool {
	return slices.Contains(colorVariants, tag.Color)
}

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
