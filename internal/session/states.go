package session

type State string

const (
	StateUnknown          State = ""
	StateIdle             State = "idle"
	StateBreathingActive  State = "breathing_active"
	StateBreathingRunning State = "breathing_running"
	StateBreathingPaused  State = "breathing_paused"
	StateGroundingStep1   State = "grounding_step_1"
	StateGroundingStep2   State = "grounding_step_2"
	StateGroundingStep3   State = "grounding_step_3"
	StateGroundingStep4   State = "grounding_step_4"
	StateGroundingStep5   State = "grounding_step_5"
)
