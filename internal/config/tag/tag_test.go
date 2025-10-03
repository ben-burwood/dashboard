package tag

import "testing"

func TestEnsureUnique_NoDuplicates(t *testing.T) {
	tags := Tags{
		{Name: "A", Color: "red"},
		{Name: "B", Color: "blue"},
		{Name: "C", Color: "green"},
	}
	if err := tags.EnsureUnique(); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestEnsureUnique_WithDuplicates(t *testing.T) {
	tags := Tags{
		{Name: "A", Color: "red"},
		{Name: "B", Color: "blue"},
		{Name: "A", Color: "green"}, // duplicate name
	}
	if err := tags.EnsureUnique(); err == nil {
		t.Errorf("expected error for duplicate tag name, got nil")
	}
}

func TestEnsureUnique_EmptySlice(t *testing.T) {
	tags := Tags{}
	if err := tags.EnsureUnique(); err != nil {
		t.Errorf("expected no error for empty slice, got %v", err)
	}
}

func TestLookup_Found(t *testing.T) {
	tags := Tags{
		{Name: "A", Color: "red"},
		{Name: "B", Color: "blue"},
	}
	tag, err := tags.Lookup("A")
	if err != nil {
		t.Fatalf("expected to find tag 'A', got error: %v", err)
	}
	if tag == nil || tag.Name != "A" || tag.Color != "red" {
		t.Errorf("unexpected tag returned: %+v", tag)
	}
}

func TestLookup_NotFound(t *testing.T) {
	tags := Tags{
		{Name: "A", Color: "red"},
	}
	tag, err := tags.Lookup("X")
	if err == nil {
		t.Errorf("expected error for missing tag, got nil")
	}
	if tag != nil {
		t.Errorf("expected nil tag for missing name, got: %+v", tag)
	}
}

func TestLookup_EmptyTags(t *testing.T) {
	tags := Tags{}
	tag, err := tags.Lookup("A")
	if err == nil {
		t.Errorf("expected error for missing tag in empty slice, got nil")
	}
	if tag != nil {
		t.Errorf("expected nil tag for missing name in empty slice, got: %+v", tag)
	}
}
