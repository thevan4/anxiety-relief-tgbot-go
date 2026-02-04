package handlers

import (
	"testing"
)

func TestIsCallbackTooOldError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"nil error", nil, false},
		{"query too old", errWithMessageError("Bad Request: query is too old and response timeout expired"), true},
		{"query id invalid", errWithMessageError("Bad Request: query ID is invalid"), true},
		{"other error", errWithMessageError("some other error"), false},
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
		{"message not modified", errWithMessageError("Bad Request: message is not modified"), true},
		{"other error", errWithMessageError("some other error"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsMessageNotModifiedError(tt.err); got != tt.want {
				t.Errorf("IsMessageNotModifiedError() = %v, want %v", got, tt.want)
			}
		})
	}
}

// errWithMessageError is a simple error type for testing.
type errWithMessageError string

func (e errWithMessageError) Error() string {
	return string(e)
}
