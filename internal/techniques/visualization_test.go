package techniques

import (
	"testing"
)

func TestGetVisualizationScenes(t *testing.T) {
	t.Parallel()

	scenes := GetVisualizationScenes()

	if len(scenes) == 0 {
		t.Error("expected at least one visualization scene")
	}

	expectedCount := 5
	if len(scenes) != expectedCount {
		t.Errorf("expected %d scenes, got %d", expectedCount, len(scenes))
	}
}

func TestGetVisualizationScenes_ValidData(t *testing.T) {
	t.Parallel()

	scenes := GetVisualizationScenes()

	for i, s := range scenes {
		if s.ID == "" {
			t.Errorf("scene %d: ID should not be empty", i)
		}
		if s.Name == "" {
			t.Errorf("scene %d: name should not be empty", i)
		}
		if s.Emoji == "" {
			t.Errorf("scene %d: emoji should not be empty", i)
		}
		if s.Description == "" {
			t.Errorf("scene %d: description should not be empty", i)
		}
		if s.Atmosphere == "" {
			t.Errorf("scene %d: atmosphere should not be empty", i)
		}
		if len(s.Steps) == 0 {
			t.Errorf("scene %d: should have at least one step", i)
		}
	}
}

func TestGetVisualizationScenes_UniqueIDs(t *testing.T) {
	t.Parallel()

	scenes := GetVisualizationScenes()
	ids := make(map[string]bool)

	for _, s := range scenes {
		if ids[s.ID] {
			t.Errorf("duplicate scene ID: %s", s.ID)
		}
		ids[s.ID] = true
	}
}

func TestGetVisualizationScenes_StepsNumbering(t *testing.T) {
	t.Parallel()

	scenes := GetVisualizationScenes()

	for _, s := range scenes {
		for i, step := range s.Steps {
			expectedNum := i + 1
			if step.Number != expectedNum {
				t.Errorf("scene %s, step %d: expected number %d, got %d",
					s.ID, i, expectedNum, step.Number)
			}
			if step.Instruction == "" {
				t.Errorf("scene %s, step %d: instruction should not be empty", s.ID, i)
			}
		}
	}
}

func TestGetVisualizationScenes_ExpectedScenes(t *testing.T) {
	t.Parallel()

	scenes := GetVisualizationScenes()

	expectedIDs := []string{"mountain", "forest", "beach", "garden", "starry"}

	for _, expectedID := range expectedIDs {
		found := false
		for _, s := range scenes {
			if s.ID == expectedID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected scene with ID %q not found", expectedID)
		}
	}
}
