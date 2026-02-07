package session

import (
	"testing"
	"time"
)

func TestNeedsRecreate(t *testing.T) {
	t.Parallel()
	ttl := 48 * time.Hour

	tests := []struct {
		name      string
		createdAt time.Time
		want      bool
	}{
		{"zero time", time.Time{}, true},
		{"fresh", time.Now(), false},
		{"old", time.Now().Add(-49 * time.Hour), true},
		{"exactly at limit", time.Now().Add(-48 * time.Hour), true},
		{"just under limit", time.Now().Add(-47 * time.Hour), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := NeedsRecreate(tt.createdAt, ttl); got != tt.want {
				t.Errorf("NeedsRecreate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetTransitionOp(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		from Screen
		to   Screen
		want TransitionOp
	}{
		// Unknown source
		{"unknown to menu", ScreenUnknown, ScreenMenu, OpSend},
		{"unknown to holder", ScreenUnknown, ScreenHolder, OpSend},
		{"unknown to exercise", ScreenUnknown, ScreenExerciseIntro, OpSend},

		// From Holder
		{"holder to menu", ScreenHolder, ScreenMenu, OpRecreate},
		{"holder to exercise", ScreenHolder, ScreenExerciseIntro, OpRecreate},
		{"holder to holder", ScreenHolder, ScreenHolder, OpRecreate},

		// Same screen
		{"menu to menu", ScreenMenu, ScreenMenu, OpEdit},
		{"exercise to exercise", ScreenExerciseRunning, ScreenExerciseRunning, OpEdit},
		{"holder to holder edit", ScreenHolder, ScreenHolder, OpRecreate},

		// Menu transitions
		{"menu to exercise intro", ScreenMenu, ScreenExerciseIntro, OpEdit},
		{"menu to exercise running", ScreenMenu, ScreenExerciseRunning, OpEdit},
		{"menu to language", ScreenMenu, ScreenLanguageSelect, OpEdit},
		{"menu to holder", ScreenMenu, ScreenHolder, OpRecreate},
		{"menu to error", ScreenMenu, ScreenError, OpEdit}, // default case

		// Exercise transitions
		{"exercise to menu", ScreenExerciseComplete, ScreenMenu, OpEdit},
		{"running to complete", ScreenExerciseRunning, ScreenExerciseComplete, OpEdit},
		{"exercise to holder", ScreenExerciseIntro, ScreenHolder, OpRecreate},
		{"intro to language", ScreenExerciseIntro, ScreenLanguageSelect, OpEdit},
		{"running to intro", ScreenExerciseRunning, ScreenExerciseIntro, OpEdit},
		{"complete to language", ScreenExerciseComplete, ScreenLanguageSelect, OpEdit},

		// Language transitions
		{"language to menu", ScreenLanguageSelect, ScreenMenu, OpEdit},
		{"language to holder", ScreenLanguageSelect, ScreenHolder, OpRecreate},
		{"language to exercise intro", ScreenLanguageSelect, ScreenExerciseIntro, OpEdit},
		{"language to exercise running", ScreenLanguageSelect, ScreenExerciseRunning, OpEdit},
		{"language to exercise complete", ScreenLanguageSelect, ScreenExerciseComplete, OpEdit},
		{"language to language", ScreenLanguageSelect, ScreenLanguageSelect, OpEdit},

		// Error screen
		{"error to menu", ScreenError, ScreenMenu, OpEdit},
		{"error to holder", ScreenError, ScreenHolder, OpRecreate},
		{"error to exercise", ScreenError, ScreenExerciseIntro, OpEdit},
		{"error to language", ScreenError, ScreenLanguageSelect, OpEdit},
		{"error to error", ScreenError, ScreenError, OpEdit},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := GetTransitionOp(tt.from, tt.to); got != tt.want {
				t.Errorf("GetTransitionOp(%v, %v) = %v, want %v", tt.from, tt.to, got, tt.want)
			}
		})
	}
}

func TestShouldRecreateForAge(t *testing.T) {
	t.Parallel()
	ttl := 48 * time.Hour
	oldTime := time.Now().Add(-49 * time.Hour)
	freshTime := time.Now()

	tests := []struct {
		name      string
		op        TransitionOp
		createdAt time.Time
		want      TransitionOp
	}{
		{"edit fresh stays edit", OpEdit, freshTime, OpEdit},
		{"edit old becomes recreate", OpEdit, oldTime, OpRecreate},
		{"recreate stays recreate", OpRecreate, freshTime, OpRecreate},
		{"send stays send", OpSend, oldTime, OpSend},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := ShouldRecreateForAge(tt.op, tt.createdAt, ttl); got != tt.want {
				t.Errorf("ShouldRecreateForAge() = %v, want %v", got, tt.want)
			}
		})
	}
}
