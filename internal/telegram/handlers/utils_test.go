package handlers

import (
	"testing"
)

func TestEmptyInlineKeyboard(t *testing.T) {
	kb := EmptyInlineKeyboard()
	if kb == nil {
		t.Fatal("EmptyInlineKeyboard returned nil")
	}
	if len(kb.InlineKeyboard) != 0 {
		t.Errorf("EmptyInlineKeyboard has %d rows, want 0", len(kb.InlineKeyboard))
	}
}

func TestGetMainMenuInline(t *testing.T) {
	kb := GetMainMenuInline()
	if kb == nil {
		t.Fatal("GetMainMenuInline returned nil")
	}
	// Should have 4 rows of buttons
	if len(kb.InlineKeyboard) != 4 {
		t.Errorf("GetMainMenuInline has %d rows, want 4", len(kb.InlineKeyboard))
	}
	// First row should have 2 buttons
	if len(kb.InlineKeyboard[0]) != 2 {
		t.Errorf("first row has %d buttons, want 2", len(kb.InlineKeyboard[0]))
	}
	// Check callback data
	if kb.InlineKeyboard[0][0].CallbackData != "menu_breathing" {
		t.Errorf("first button callback = %q, want %q", kb.InlineKeyboard[0][0].CallbackData, "menu_breathing")
	}
}

func TestIsCallbackTooOldError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"nil error", nil, false},
		{"query too old", errWithMessage("Bad Request: query is too old and response timeout expired"), true},
		{"query id invalid", errWithMessage("Bad Request: query ID is invalid"), true},
		{"other error", errWithMessage("some other error"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsCallbackTooOldError(tt.err); got != tt.want {
				t.Errorf("IsCallbackTooOldError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsMessageNotModifiedError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"nil error", nil, false},
		{"message not modified", errWithMessage("Bad Request: message is not modified"), true},
		{"other error", errWithMessage("some other error"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsMessageNotModifiedError(tt.err); got != tt.want {
				t.Errorf("IsMessageNotModifiedError() = %v, want %v", got, tt.want)
			}
		})
	}
}

// errWithMessage is a simple error type for testing
type errWithMessage string

func (e errWithMessage) Error() string {
	return string(e)
}
