package session

type State string

const (
	StateUnknown          State = ""
	StateIdle             State = "idle"
	StateBreathingActive  State = "breathing_active"
	StateBreathingRunning State = "breathing_running"
	StateGroundingStep1   State = "grounding_step_1"
	StateGroundingStep2   State = "grounding_step_2"
	StateGroundingStep3   State = "grounding_step_3"
	StateGroundingStep4   State = "grounding_step_4"
	StateGroundingStep5   State = "grounding_step_5"

	// Guided Breathing states.
	StateGuidedBreathingSelect  State = "guided_breathing_select"
	StateGuidedBreathingRunning State = "guided_breathing_running"

	// PMR states.
	StatePMRActive  State = "pmr_active"
	StatePMRTense   State = "pmr_tense"
	StatePMRRelax   State = "pmr_relax"
	StatePMRRunning State = "pmr_running"

	// Thought Labeling states.
	StateThoughtLabelingActive   State = "thought_labeling_active"
	StateThoughtLabelingInput    State = "thought_labeling_input"
	StateThoughtLabelingCategory State = "thought_labeling_category"

	// Visualization states.
	StateVisualizationSelect  State = "visualization_select"
	StateVisualizationRunning State = "visualization_running"
)
