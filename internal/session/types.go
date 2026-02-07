// Package session provides user session state management and storage.
package session

import (
	"errors"
	"time"
)

// =============================================================================
// State — состояние пользователя в flow бота
// =============================================================================

// State represents user's current session state in the bot flow.
type State string

const (
	// StateUnknown represents no active session.
	StateUnknown State = ""
	// StateIdle represents idle state when user is not in any exercise.
	StateIdle State = "idle"
	// StateBreathingActive represents breathing technique menu is active.
	StateBreathingActive State = "breathing_active"
	// StateBreathingRunning represents breathing exercise is currently running.
	StateBreathingRunning State = "breathing_running"
	// StateGroundingStep1 represents grounding step 1 (see 5 things).
	StateGroundingStep1 State = "grounding_step_1"
	// StateGroundingStep2 represents grounding step 2 (touch 4 things).
	StateGroundingStep2 State = "grounding_step_2"
	// StateGroundingStep3 represents grounding step 3 (hear 3 sounds).
	StateGroundingStep3 State = "grounding_step_3"
	// StateGroundingStep4 represents grounding step 4 (smell 2 things).
	StateGroundingStep4 State = "grounding_step_4"
	// StateGroundingStep5 represents grounding step 5 (taste 1 thing).
	StateGroundingStep5 State = "grounding_step_5"

	// StateGuidedBreathingSelect indicates user is selecting breathing pattern.
	StateGuidedBreathingSelect State = "guided_breathing_select"
	// StateGuidedBreathingRunning represents guided breathing exercise is running.
	StateGuidedBreathingRunning State = "guided_breathing_running"

	// StatePMRActive indicates PMR technique menu is active.
	StatePMRActive State = "pmr_active"
	// StatePMRRunning represents PMR exercise is running.
	StatePMRRunning State = "pmr_running"
)

// IsRunning returns true if user is in an active exercise that should not be interrupted.
func (s State) IsRunning() bool {
	switch s {
	case StateBreathingRunning,
		StateGuidedBreathingRunning,
		StatePMRRunning:
		return true
	}
	return false
}

// IsActive returns true if user has any active session (including selection screens).
func (s State) IsActive() bool {
	return s != StateUnknown && s != StateIdle
}

// =============================================================================
// Screen — типы экранов для логики переходов
// =============================================================================

// Screen represents a UI screen type for transition logic.
type Screen int

const (
	// ScreenUnknown — неизвестный экран.
	ScreenUnknown Screen = iota
	// ScreenHolder — приветственное сообщение после /start.
	ScreenHolder
	// ScreenMenu — главное меню с кнопками упражнений.
	ScreenMenu
	// ScreenExerciseIntro — вступительный экран упражнения.
	ScreenExerciseIntro
	// ScreenExerciseRunning — упражнение в процессе выполнения.
	ScreenExerciseRunning
	// ScreenExerciseComplete — упражнение завершено.
	ScreenExerciseComplete
	// ScreenLanguageSelect — экран выбора языка.
	ScreenLanguageSelect
	// ScreenError — экран ошибки.
	ScreenError
)

// String returns human-readable screen name.
func (s Screen) String() string {
	switch s {
	case ScreenHolder:
		return "Holder"
	case ScreenMenu:
		return "Menu"
	case ScreenExerciseIntro:
		return "ExerciseIntro"
	case ScreenExerciseRunning:
		return "ExerciseRunning"
	case ScreenExerciseComplete:
		return "ExerciseComplete"
	case ScreenLanguageSelect:
		return "LanguageSelect"
	case ScreenError:
		return "Error"
	default:
		return "Unknown"
	}
}

// =============================================================================
// TransitionOp — операции перехода между экранами
// =============================================================================

// TransitionOp represents the operation to perform when transitioning between screens.
type TransitionOp int

const (
	// OpEdit — редактировать существующее сообщение.
	OpEdit TransitionOp = iota
	// OpRecreate — удалить старое и отправить новое сообщение.
	OpRecreate
	// OpSend — отправить новое сообщение (нет существующего).
	OpSend
)

// String returns human-readable operation name.
func (op TransitionOp) String() string {
	switch op {
	case OpEdit:
		return "Edit"
	case OpRecreate:
		return "Recreate"
	case OpSend:
		return "Send"
	default:
		return "Unknown"
	}
}

// =============================================================================
// Cleanup — типы для логики очистки ресурсов
// =============================================================================

// CleanupRetry stores retry information for cleanup decisions.
type CleanupRetry struct {
	State   string `json:"state"`
	Attempt int    `json:"attempt"`
}

// CleanupDecision represents the result of cleanup logic evaluation.
type CleanupDecision int

const (
	// DecisionDelete means the resource should be deleted.
	DecisionDelete CleanupDecision = iota
	// DecisionKeep means the resource should be kept (user became active).
	DecisionKeep
	// DecisionRetry means we should check again later.
	DecisionRetry
)

// String returns human-readable decision name.
func (d CleanupDecision) String() string {
	switch d {
	case DecisionDelete:
		return "Delete"
	case DecisionKeep:
		return "Keep"
	case DecisionRetry:
		return "Retry"
	default:
		return "Unknown"
	}
}

// CleanupResult contains the decision and optional retry info.
type CleanupResult struct {
	Decision  CleanupDecision
	Retry     *CleanupRetry // set if Decision == DecisionRetry
	NextCheck time.Time     // when to check next (for DecisionRetry)
	Reason    string        // human-readable reason for logging
}

// ErrCleanupRetryNotFound indicates missing cleanup retry entry.
//
//nolint:gochecknoglobals // Sentinel error for storage.
var ErrCleanupRetryNotFound = errors.New("cleanup retry not found")

// =============================================================================
// Config — конфигурация тайминга
// =============================================================================

const (
	defaultStateTTL      = 5 * time.Minute
	defaultResourceTTL   = 48 * time.Hour
	defaultCleanupDelay  = 47 * time.Hour
	defaultRetryInterval = 5 * time.Minute
	defaultMaxRetries    = 6
)

// Config holds timing configuration for session management.
type Config struct {
	// StateTTL is how long state lives without refresh (activity indicator).
	StateTTL time.Duration

	// ResourceTTL is max lifetime of managed resource (e.g., 48h for Telegram messages).
	ResourceTTL time.Duration

	// CleanupDelay is time before first cleanup check after resource creation.
	CleanupDelay time.Duration

	// RetryInterval is time between retry checks.
	RetryInterval time.Duration

	// MaxRetries is maximum retry attempts before force cleanup.
	MaxRetries int
}

// DefaultConfig returns sensible defaults for Telegram-like platforms.
func DefaultConfig() Config {
	return Config{
		StateTTL:      defaultStateTTL,
		ResourceTTL:   defaultResourceTTL,
		CleanupDelay:  defaultCleanupDelay,
		RetryInterval: defaultRetryInterval,
		MaxRetries:    defaultMaxRetries,
	}
}
