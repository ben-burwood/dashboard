package service

import "testing"

func TestValidateAndSetDefaults_ValidService(t *testing.T) {
	s := &Service{
		Link:  "http://example.com",
		Title: "Example",
		Group: "group1",
		Tags:  []string{"tag1", "tag2"},
	}
	if err := s.ValidateAndSetDefaults(); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestValidateAndSetDefaults_MissingLink(t *testing.T) {
	s := &Service{
		Title: "Example",
		Group: "group1",
	}
	if err := s.ValidateAndSetDefaults(); err == nil {
		t.Errorf("expected error for missing link, got nil")
	}
}

func TestValidateAndSetDefaults_MissingTitle(t *testing.T) {
	s := &Service{
		Link:  "http://example.com",
		Group: "group1",
	}
	if err := s.ValidateAndSetDefaults(); err == nil {
		t.Errorf("expected error for missing title, got nil")
	}
}

func TestValidateAndSetDefaults_DefaultGroup(t *testing.T) {
	s := &Service{
		Link:  "http://example.com",
		Title: "Example",
	}
	_ = s.ValidateAndSetDefaults()
	if s.Group != DefaultGroup {
		t.Errorf("expected default group, got %v", s.Group)
	}
}

func TestServiceIcon_DefaultColor(t *testing.T) {
	s := &Service{
		Link:  "http://example.com",
		Title: "Example",
		Icon:  Icon{Name: "home"},
	}
	err := s.ValidateAndSetDefaults()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if s.Icon.Color != DefaultIconColor {
		t.Errorf("expected default icon color '%s', got '%s'", DefaultIconColor, s.Icon.Color)
	}
}

func TestServiceIcon_CustomColor(t *testing.T) {
	s := &Service{
		Link:  "http://example.com",
		Title: "Example",
		Icon:  Icon{Name: "home", Color: "red"},
	}
	err := s.ValidateAndSetDefaults()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if s.Icon.Color != "red" {
		t.Errorf("expected icon color 'red', got '%s'", s.Icon.Color)
	}
}

func TestValidateAndSetDefaults_DuplicateTags(t *testing.T) {
	s := &Service{
		Link:  "http://example.com",
		Title: "Example",
		Tags:  []string{"tag1", "tag1"},
	}
	if err := s.ValidateAndSetDefaults(); err == nil {
		t.Errorf("expected error for duplicate tags, got nil")
	}
}
