package techniques

import "testing"

func TestGetGroundingSteps(t *testing.T) {
	t.Parallel()
	steps := GetGroundingSteps()

	if len(steps) != 5 {
		t.Errorf("GetGroundingSteps() returned %d steps, want 5", len(steps))
	}

	expected := []struct {
		number int
		sense  string
		count  int
		emoji  string
	}{
		{1, "sight", 5, "👁️"},
		{2, "touch", 4, "🤚"},
		{3, "hearing", 3, "👂"},
		{4, "smell", 2, "👃"},
		{5, "taste", 1, "👅"},
	}

	for i, exp := range expected {
		step := steps[i]
		if step.Number != exp.number {
			t.Errorf("steps[%d].Number = %d, want %d", i, step.Number, exp.number)
		}
		if step.Sense != exp.sense {
			t.Errorf("steps[%d].Sense = %q, want %q", i, step.Sense, exp.sense)
		}
		if step.Count != exp.count {
			t.Errorf("steps[%d].Count = %d, want %d", i, step.Count, exp.count)
		}
		if step.Emoji != exp.emoji {
			t.Errorf("steps[%d].Emoji = %q, want %q", i, step.Emoji, exp.emoji)
		}
	}
}

func TestGroundingStepsPattern(t *testing.T) {
	t.Parallel()
	steps := GetGroundingSteps()

	// 5-4-3-2-1 pattern
	expectedCounts := []int{5, 4, 3, 2, 1}
	for i, step := range steps {
		if step.Count != expectedCounts[i] {
			t.Errorf("step %d count = %d, want %d (5-4-3-2-1 pattern)", i+1, step.Count, expectedCounts[i])
		}
	}
}
