package techniques

import (
	"testing"
)

func TestGetMuscleGroups(t *testing.T) {
	t.Parallel()

	groups := GetMuscleGroups()

	if len(groups) == 0 {
		t.Error("expected at least one muscle group")
	}

	expectedCount := 10
	if len(groups) != expectedCount {
		t.Errorf("expected %d muscle groups, got %d", expectedCount, len(groups))
	}
}

func TestGetMuscleGroups_ValidData(t *testing.T) {
	t.Parallel()

	groups := GetMuscleGroups()

	for i, g := range groups {
		if g.Number != i+1 {
			t.Errorf("group %d: expected number %d, got %d", i, i+1, g.Number)
		}
		if g.ID == "" {
			t.Errorf("group %d: id should not be empty", i)
		}
		if g.Emoji == "" {
			t.Errorf("group %d: emoji should not be empty", i)
		}
		if g.TenseDuration <= 0 {
			t.Errorf("group %d: tense duration should be positive", i)
		}
		if g.RelaxDuration <= 0 {
			t.Errorf("group %d: relax duration should be positive", i)
		}
	}
}

func TestGetMuscleGroups_Order(t *testing.T) {
	t.Parallel()

	groups := GetMuscleGroups()

	expectedFirst := "hands"
	if groups[0].ID != expectedFirst {
		t.Errorf("expected first group to be %q, got %q", expectedFirst, groups[0].ID)
	}

	expectedLast := "calves"
	if groups[len(groups)-1].ID != expectedLast {
		t.Errorf("expected last group to be %q, got %q", expectedLast, groups[len(groups)-1].ID)
	}
}
