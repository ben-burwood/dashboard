package service

import "errors"

const (
	DefaultIconColor       = "black"
	DefaultGroup     Group = "default"
)

var (
	ErrInvalidServiceConfiguration = errors.New("invalid Service Configuration")
)

type Group string

type Icon struct {
	Name  string `yaml:"name"`
	Color string `yaml:"color,omitempty"`
}

type Service struct {
	Link        string   `yaml:"link"`
	Icon        Icon     `yaml:"icon"`
	Title       string   `yaml:"title"`
	Description string   `yaml:"description,omitempty"`
	Group       Group    `yaml:"group,omitempty"`
	Tags        []string `yaml:"tags,omitempty"`
}

// ValidateAndSetDefaults validates the Service configuration and sets defaults where necessary
func (s *Service) ValidateAndSetDefaults() error {
	if len(s.Link) == 0 {
		return ErrInvalidServiceConfiguration
	}
	if len(s.Title) == 0 {
		return ErrInvalidServiceConfiguration
	}
	if len(s.Group) == 0 {
		s.Group = DefaultGroup
	}
	if len(s.Icon.Color) == 0 {
		s.Icon.Color = DefaultIconColor
	}
	// Check Tag Unique Set
	tagSet := make(map[string]struct{})
	for _, tag := range s.Tags {
		if _, exists := tagSet[tag]; exists {
			return ErrInvalidServiceConfiguration
		}
		tagSet[tag] = struct{}{}
	}
	return nil
}

// GroupedServices is a map of service groups to services
type GroupedServices map[Group][]*Service

// ByGroup groups the services by their group name
func ByGroup(services []*Service) GroupedServices {
	grouped := make(GroupedServices)
	for _, service := range services {
		grouped[service.Group] = append(grouped[service.Group], service)
	}
	return grouped
}
