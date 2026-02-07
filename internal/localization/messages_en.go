//nolint:dupl // Localization files have identical structure by design, only text values differ.
package localization

//nolint:gochecknoglobals // Localization bundle.
var messagesEN = Messages{
	// Buttons
	Start:      "▶️ Start",
	Back:       "◀️ Back",
	Stop:       "🛑 Stop",
	Next:       "➡️ Next",
	Done:       "✅ Done",
	FeelBetter: "✅ Better",
	Repeat:     "🔄 Repeat",

	// Welcome holder — static message with bot info
	HolderText: `🌿 *Anxiety Relief Helper*

This bot helps manage anxiety using evidence-based techniques:

• 🌬️ Breathing 4-4-6 — quick calming
• 🌿 Grounding 5-4-3-2-1 — return to moment
• 🧘 Guided breathing — various patterns
• 💪 Muscle relaxation — release tension

Privacy: the bot doesn't store personal data.

👇 Press the button below to start`,

	// Session expired
	SessionExpired: "Session expired. Press /start",

	// Main menu
	MainMenuText: `Choose a technique to manage anxiety:

*Quick (2-5 min):*
🌬️ Breathing — calm your nervous system
🌿 Grounding — return to the present moment

*Advanced (5-15 min):*
🧘 Guided breathing — various patterns
💪 Muscle relaxation — release tension`,
	MenuBreathing: "🌬️ Breathing 2 min",
	MenuGrounding: "🌿 Grounding",
	MenuGuided:    "🧘 Guided",
	MenuPMR:       "💪 Muscles",
	MenuLang:      "🌐 Language",

	// Breathing
	BreathingIntro: `🌬️ *2-Minute Breathing*

A simple exercise to calm your nervous system.

*Instructions:*
1️⃣ Sit comfortably, relax
2️⃣ Inhale through nose — 4 seconds
3️⃣ Hold breath — 4 seconds
4️⃣ Exhale through mouth — 6 seconds
5️⃣ Repeat 8 cycles (~2 minutes)

I'll guide you through each step. Ready to begin?`,
	BreathingCompletion: `✅ *Great job!*

You've completed the breathing exercise.
How do you feel?`,
	BreathingInhale: "Inhale",
	BreathingHold:   "Hold",
	BreathingExhale: "Exhale",

	// Grounding
	GroundingIntro: `🌿 *5-4-3-2-1 Grounding Technique*

This technique helps you return to the present moment through your senses.

*How it works:*
You'll name things around you that you perceive with different senses.

Ready to begin?`,

	// Grounding steps
	GroundingStep1Title: "Step 1: Sight",
	GroundingStep1Desc:  "Name 5 things you can see around you.\n\nFor example: table, lamp, book, window, chair.",
	GroundingStep2Title: "Step 2: Touch",
	GroundingStep2Desc: "Name 4 things you can feel with your body.\n\n" +
		"For example: feet on the floor, back on the chair, hands on the table, clothes on your body.",
	GroundingStep3Title: "Step 3: Hearing",
	GroundingStep3Desc:  "Name 3 sounds you can hear.\n\nFor example: street noise, clock ticking, your breathing.",
	GroundingStep4Title: "Step 4: Smell",
	GroundingStep4Desc:  "Name 2 scents you can smell or enjoy.\n\nFor example: coffee, fresh air, flower fragrance.",
	GroundingStep5Title: "Step 5: Taste",
	GroundingStep5Desc:  "Name 1 taste you can sense in your mouth.\n\nIf nothing — recall your favorite taste.",

	// Guided breathing
	GuidedIntro:      "🧘 *Guided Breathing*\n\nChoose a breathing technique:\n\n",
	GuidedCompletion: "✅ *Great job!*\n\nYou've completed the \"%s\" exercise.\n\nHow do you feel?",
	GuidedStopped:    "🧘 *Guided Breathing*\n\nExercise stopped. Choose a technique:\n\n",

	// Breathing patterns
	PatternBoxName:        "Box Breathing 4-4-4-4",
	PatternBoxDesc:        "Focus and concentration",
	PatternRelaxingName:   "Relaxing 4-7-8",
	PatternRelaxingDesc:   "Deep relaxation and sleep",
	PatternEnergizingName: "Energizing 4-4-6",
	PatternEnergizingDesc: "Energy and mental clarity",
	PatternQuickName:      "Quick Reset 3-3-3",
	PatternQuickDesc:      "Fast tension relief",

	// Guided breathing phases
	GuidedInhale:  "Inhale",
	GuidedHoldIn:  "Hold",
	GuidedExhale:  "Exhale",
	GuidedHoldOut: "Pause",

	// PMR
	PMRIntro: `💪 *Progressive Muscle Relaxation*

A deep relaxation technique through tensing and relaxing muscles.

*How it works:*
1. Tense a muscle group for 7 seconds
2. Relax for 15 seconds
3. Move to the next group

*%d muscle groups* — from hands to feet.

Ready to begin?`,
	PMRTense: "🔴 *TENSE*\n\n%s",
	PMRRelax: "🟢 *RELAX*\n\n%s",
	PMRCompletion: `✅ *Great job!*

You've completed progressive muscle relaxation.

Your body is now fully relaxed. Sit for another minute, enjoying this state.`,

	// PMR muscle groups
	PMRHandsName:  "Hands",
	PMRHandsTense: "Clench your fists as hard as you can. Feel the tension in your fingers and palms.",
	PMRHandsRelax: "Release your fists and relax your hands. Feel the warmth and relaxation.",

	PMRForearmsName:  "Forearms and Biceps",
	PMRForearmsTense: "Bend your arms at the elbows and tense your biceps. Hold the tension.",
	PMRForearmsRelax: "Lower your arms and completely relax them. Your arms become heavy.",

	PMRForeheadName:  "Forehead",
	PMRForeheadTense: "Raise your eyebrows as high as possible, wrinkle your forehead. Feel the tension.",
	PMRForeheadRelax: "Relax your forehead. Let your eyebrows drop. Your forehead becomes smooth.",

	PMREyesName:  "Eyes and Nose",
	PMREyesTense: "Squeeze your eyes shut tightly and wrinkle your nose. Feel the tension around your eyes.",
	PMREyesRelax: "Relax your eyes and nose. Your eyelids become light and calm.",

	PMRJawName:  "Jaw",
	PMRJawTense: "Clench your jaw and stretch your lips in a tense smile.",
	PMRJawRelax: "Relax your jaw, slightly open your mouth. Your tongue is relaxed.",

	PMRNeckName:  "Neck and Shoulders",
	PMRNeckTense: "Raise your shoulders to your ears and tense your neck. Hold the tension.",
	PMRNeckRelax: "Lower your shoulders down and relax your neck. Feel the relief.",

	PMRChestName:  "Chest and Back",
	PMRChestTense: "Take a deep breath, hold it and tense your chest and back muscles.",
	PMRChestRelax: "Slowly exhale and relax your chest and back. Breathe calmly.",

	PMRStomachName:  "Stomach",
	PMRStomachTense: "Pull in your stomach and tense your abs. Hold the tension.",
	PMRStomachRelax: "Relax your stomach. Let it move freely with your breathing.",

	PMRThighsName:  "Thighs and Glutes",
	PMRThighsTense: "Tense your glutes and thighs, press them to the seat.",
	PMRThighsRelax: "Relax your glutes and thighs. Feel the heaviness in your legs.",

	PMRCalvesName:  "Calves and Feet",
	PMRCalvesTense: "Pull your toes toward you, tense your calves. Feel the stretch.",
	PMRCalvesRelax: "Relax your feet and calves. Your legs become warm and heavy.",

	// Language selection
	LangSelectTitle: "🌐 *Choose language:*",
	CycleWord:       "Cycle",
}
