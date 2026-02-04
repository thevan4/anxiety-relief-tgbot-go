package localization

var messagesEN = Messages{
	// Buttons
	Start:      "▶️ Start",
	Back:       "◀️ Back",
	Stop:       "🛑 Stop",
	Next:       "➡️ Next",
	Done:       "✅ Done",
	FeelBetter: "✅ Better",
	Repeat:     "🔄 Repeat",
	BackToMenu: "🏠 Menu",

	// Main menu
	MainMenuText: `Choose a technique to manage anxiety:

*Quick (2-5 min):*
🌬️ Breathing — calm your nervous system
🌿 Grounding — return to the present moment

*Advanced (5-15 min):*
🧘 Guided breathing — various patterns
💪 Muscle relaxation — release tension
🏷️ Thought labeling — work with anxiety`,
	MenuBreathing: "🌬️ Breathing 2 min",
	MenuGrounding: "🌿 Grounding",
	MenuGuided:    "🧘 Guided",
	MenuPMR:       "💪 Muscles",
	MenuThought:   "🏷️ Thoughts",
	MenuInfo:      "ℹ️ Info",
	MenuLang:      "🌐 Language",

	// Breathing
	BreathingIntro: `🌬️ *2-Minute Breathing*

A simple exercise to calm your nervous system.

*Instructions:*
1️⃣ Sit comfortably, close your eyes
2️⃣ Inhale through nose — 4 seconds
3️⃣ Hold breath — 4 seconds
4️⃣ Exhale through mouth — 6 seconds
5️⃣ Repeat 8 cycles (~2 minutes)

I'll guide you through each step. Ready to begin?`,
	BreathingCompletion: `✅ *Great job!*

You've completed the breathing exercise.
How do you feel?`,
	BreathingThanks: `✨ *Thank you for practicing!*

Regular exercises help reduce anxiety levels.

Choose a technique from the menu below.`,
	BreathingCycle:  "Cycle %d/%d",
	BreathingInhale: "Inhale",
	BreathingHold:   "Hold",
	BreathingExhale: "Exhale",
	BreathingPause:  "Pause",

	// Grounding
	GroundingIntro: `🌿 *5-4-3-2-1 Grounding Technique*

This technique helps you return to the present moment through your senses.

*How it works:*
You'll name things around you that you perceive with different senses.

Ready to begin?`,
	GroundingStep:       "Name *%d %s* that you can *%s*",
	GroundingCompletion: "✅ *Great job!*\n\nYou've completed the grounding technique.\nHow do you feel?",
	GroundingThanks: `✨ *Thank you for practicing!*

The 5-4-3-2-1 technique helps quickly return to the present moment during anxiety.

Choose a technique from the menu below.`,
	GroundingSee:        "see",
	GroundingHear:       "hear",
	GroundingFeel:       "feel",
	GroundingThings:     "things",
	GroundingSounds:     "sounds",
	GroundingSensations: "sensations",

	// Grounding steps
	GroundingStep1Title: "Step 1: Sight",
	GroundingStep1Desc:  "Name 5 things you can see around you.\n\nFor example: table, lamp, book, window, chair.",
	GroundingStep2Title: "Step 2: Touch",
	GroundingStep2Desc:  "Name 4 things you can feel with your body.\n\nFor example: feet on the floor, back on the chair, hands on the table, clothes on your body.",
	GroundingStep3Title: "Step 3: Hearing",
	GroundingStep3Desc:  "Name 3 sounds you can hear.\n\nFor example: street noise, clock ticking, your breathing.",
	GroundingStep4Title: "Step 4: Smell",
	GroundingStep4Desc:  "Name 2 scents you can smell or enjoy.\n\nFor example: coffee, fresh air, flower fragrance.",
	GroundingStep5Title: "Step 5: Taste",
	GroundingStep5Desc:  "Name 1 taste you can sense in your mouth.\n\nIf nothing — recall your favorite taste.",

	// Guided breathing
	GuidedIntro:      "🧘 *Guided Breathing*\n\nChoose a breathing technique:\n\n",
	GuidedCompletion: "✅ *Great job!*\n\nYou've completed the \"%s\" exercise.\n\nHow do you feel?",
	GuidedThanks: `✨ *Thank you for practicing!*

Regular breathing exercises help reduce anxiety and improve focus.

Choose a technique from the menu below.`,
	GuidedStopped: "🧘 *Guided Breathing*\n\nExercise stopped. Choose a technique:\n\n",

	// Breathing patterns
	PatternBoxName:        "Box Breathing 4-4-4-4",
	PatternBoxDesc:        "Focus and concentration",
	PatternRelaxingName:   "Relaxing 4-7-8",
	PatternRelaxingDesc:   "Deep relaxation and sleep",
	PatternEnergizingName: "Energizing 4-4-6",
	PatternEnergizingDesc: "Energy and mental clarity",
	PatternQuickName:      "Quick Reset 3-3-3",
	PatternQuickDesc:      "Fast tension relief",

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
	PMRStopped: `💪 *Progressive Muscle Relaxation*

Exercise stopped.

Would you like to start over?`,
	PMRThanks: `✨ *Thank you for practicing!*

Progressive muscle relaxation reduces muscle tension and stress levels.

Choose a technique from the menu below.`,

	// Thought labeling
	ThoughtIntro: `🏷️ *Thought Labeling*

A mindfulness technique for working with anxious thoughts.

*How it works:*
1. You describe an anxious thought
2. Choose a distortion category
3. Get a way to reframe it

This helps separate yourself from thoughts and see them objectively.

Ready to begin?`,
	ThoughtPrompt:     "📝 *Describe your anxious thought*\n\nWrite in one message the thought that's bothering you.",
	ThoughtCategories: "🏷️ *Categorizing the thought*\n\n_\"%s\"_\n\nChoose the type of cognitive distortion:",
	ThoughtResult:     "%s *%s*\n\n_%s_\n\n*How to reframe:*\n%s",
	ThoughtCompletion: "Would you like to work through another thought?",
	ThoughtThanks: `✨ *Thank you for practicing!*

Thought labeling helps recognize cognitive distortions and reduce their impact.

Choose a technique from the menu below.`,

	// Info
	InfoText: `ℹ️ *About the bot*

This bot helps manage anxiety using evidence-based techniques.

*Available techniques:*
• 🌬️ 4-4-6 Breathing — quick calming
• 🌿 5-4-3-2-1 Grounding — return to moment
• 🧘 Guided breathing — various patterns
• 💪 Muscle relaxation — release tension
• 🏷️ Thought labeling — work with anxiety

*Privacy:*
The bot doesn't store personal data. Sessions are automatically deleted.

Source code: github.com/thevan4/anxiety-relief-tgbot-go`,

	// Language selection
	LangSelectTitle: "🌐 *Choose language:*",
	CycleWord:       "Cycle",
}
