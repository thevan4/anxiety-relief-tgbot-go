package techniques

import (
	"testing"
)

func TestGetThoughtCategories(t *testing.T) {
	t.Parallel()

	categories := GetThoughtCategories()

	if len(categories) == 0 {
		t.Error("expected at least one thought category")
	}

	expectedCount := 11
	if len(categories) != expectedCount {
		t.Errorf("expected %d categories, got %d", expectedCount, len(categories))
	}
}

func TestGetThoughtCategories_ValidData(t *testing.T) {
	t.Parallel()

	categories := GetThoughtCategories()

	for i, c := range categories {
		if c.ID == "" {
			t.Errorf("category %d: ID should not be empty", i)
		}
		if c.Name == "" {
			t.Errorf("category %d: name should not be empty", i)
		}
		if c.Emoji == "" {
			t.Errorf("category %d: emoji should not be empty", i)
		}
		if c.Description == "" {
			t.Errorf("category %d: description should not be empty", i)
		}
		if c.Example == "" {
			t.Errorf("category %d: example should not be empty", i)
		}
	}
}

func TestGetThoughtCategories_UniqueIDs(t *testing.T) {
	t.Parallel()

	categories := GetThoughtCategories()
	ids := make(map[string]bool)

	for _, c := range categories {
		if ids[c.ID] {
			t.Errorf("duplicate category ID: %s", c.ID)
		}
		ids[c.ID] = true
	}
}

func TestNewLabelingSession(t *testing.T) {
	t.Parallel()

	session := NewLabelingSession()

	if session.ThoughtsLabeled != 0 {
		t.Errorf("expected 0 thoughts labeled, got %d", session.ThoughtsLabeled)
	}
	if session.Categories == nil {
		t.Error("expected non-nil categories map")
	}
}

func TestLabelingSession_AddLabel(t *testing.T) {
	t.Parallel()

	session := NewLabelingSession()

	session.AddLabel("worry")
	session.AddLabel("worry")
	session.AddLabel("catastrophic")

	if session.ThoughtsLabeled != 3 {
		t.Errorf("expected 3 thoughts labeled, got %d", session.ThoughtsLabeled)
	}
	if session.Categories["worry"] != 2 {
		t.Errorf("expected 2 worry labels, got %d", session.Categories["worry"])
	}
	if session.Categories["catastrophic"] != 1 {
		t.Errorf("expected 1 catastrophic label, got %d", session.Categories["catastrophic"])
	}
}

func TestLabelingSession_GetMostFrequent(t *testing.T) {
	t.Parallel()

	session := NewLabelingSession()

	session.AddLabel("worry")
	session.AddLabel("worry")
	session.AddLabel("catastrophic")

	mostFrequent := session.GetMostFrequent()
	if mostFrequent != "worry" {
		t.Errorf("expected most frequent to be 'worry', got %q", mostFrequent)
	}
}

func TestLabelingSession_GetMostFrequent_Empty(t *testing.T) {
	t.Parallel()

	session := NewLabelingSession()

	mostFrequent := session.GetMostFrequent()
	if mostFrequent != "" {
		t.Errorf("expected empty string for empty session, got %q", mostFrequent)
	}
}
