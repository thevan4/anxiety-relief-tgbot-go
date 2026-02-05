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
	BackToMenu: "🏠 Menu",

	// Welcome holder — static message with bot info
	HolderText: `🌿 *Anxiety Relief Helper*

This bot helps manage anxiety using evidence-based techniques:

• 🌬️ Breathing 4-4-6 — quick calming
• 🌿 Grounding 5-4-3-2-1 — return to moment
• 🧘 Guided breathing — various patterns
• 💪 Muscle relaxation — release tension
• 🏷️ Thought labeling — work with anxiety
• 🌅 Visualization — relaxation through imagination

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
💪 Muscle relaxation — release tension
🏷️ Thought labeling — work with anxiety
🌅 Visualization — relaxation through imagination`,
	MenuBreathing:     "🌬️ Breathing 2 min",
	MenuGrounding:     "🌿 Grounding",
	MenuGuided:        "🧘 Guided",
	MenuPMR:           "💪 Muscles",
	MenuThought:       "🏷️ Thoughts",
	MenuVisualization: "🌅 Visualization",
	MenuLang:          "🌐 Language",

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
	PMRStopped: `💪 *Progressive Muscle Relaxation*

Exercise stopped.

Would you like to start over?`,
	PMRThanks: `✨ *Thank you for practicing!*

Progressive muscle relaxation reduces muscle tension and stress levels.

Choose a technique from the menu below.`,

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

	// Thought labeling categories
	ThoughtWorryName:    "Worry",
	ThoughtWorryDesc:    "Anxiety about future events or outcomes",
	ThoughtWorryReframe: "\"What if I fail the interview?\"",

	ThoughtCatastrophicName:    "Catastrophizing",
	ThoughtCatastrophicDesc:    "Imagining the worst-case scenario",
	ThoughtCatastrophicReframe: "\"If I make a mistake, everything will be terrible!\"",

	ThoughtSelfDoubtName:    "Self-Doubt",
	ThoughtSelfDoubtDesc:    "Doubting your own abilities",
	ThoughtSelfDoubtReframe: "\"I'm not good enough for this\"",

	ThoughtPerfectionistName:    "Perfectionism",
	ThoughtPerfectionistDesc:    "Unrealistically high standards",
	ThoughtPerfectionistReframe: "\"If it's not perfect, it's bad\"",

	ThoughtComparisonName:    "Comparison",
	ThoughtComparisonDesc:    "Comparing yourself to others",
	ThoughtComparisonReframe: "\"Everyone handles this better than me\"",

	ThoughtRuminationName:    "Rumination",
	ThoughtRuminationDesc:    "Constantly returning to past events",
	ThoughtRuminationReframe: "\"Why did I say that then?\"",

	ThoughtControlName:    "Control",
	ThoughtControlDesc:    "Wanting to control the uncontrollable",
	ThoughtControlReframe: "\"I must anticipate everything\"",

	ThoughtRejectionName:    "Fear of Rejection",
	ThoughtRejectionDesc:    "Fear of being rejected",
	ThoughtRejectionReframe: "\"They'll think I'm weird\"",

	ThoughtHealthName:    "Health Anxiety",
	ThoughtHealthDesc:    "Excessive worry about health",
	ThoughtHealthReframe: "\"Is this symptom a sign of illness?\"",

	ThoughtSocialName:    "Social Anxiety",
	ThoughtSocialDesc:    "Fear of social situations",
	ThoughtSocialReframe: "\"Everyone will be looking at me\"",

	ThoughtFinancialName:    "Financial Worries",
	ThoughtFinancialDesc:    "Stress about money",
	ThoughtFinancialReframe: "\"What if there's not enough money?\"",

	// Visualization
	VisualizationIntro: `🌅 *Peaceful Visualization*

A relaxation technique through imagining peaceful places.

*Choose a scene:*

`,
	VisualizationStopped: `🌅 *Peaceful Visualization*

Exercise stopped. Choose another scene:

`,
	VisualizationCompletion: `✅ *Great job!*

You've completed the visualization "%s".

Slowly return to reality. Wiggle your fingers, take a deep breath and open your eyes.

How do you feel?`,
	VisualizationThanks: `✨ *Thank you for practicing!*

Visualization is a powerful technique for reducing stress and anxiety. Regular practice enhances the effect.

Choose a technique from the menu below.`,
	VisualizationAtmosphere: "Atmosphere",
	VisualizationCloseEyes:  "Close your eyes and immerse yourself in this scene...",
	VisualizationStepFmt:    "Step %d/%d",

	// Scene: Mountain
	SceneMountainName:  "Mountain Peak",
	SceneMountainDesc:  "Sunrise at the mountain top",
	SceneMountainAtmo:  "Cool, fresh mountain air",
	SceneMountainStep1: "Close your eyes and imagine standing on a mountain peak early in the morning.",
	SceneMountainStep2: "Silence surrounds you. You feel the cool mountain air on your skin.",
	SceneMountainStep3: "The sun appears on the horizon, painting the sky in pink and orange hues.",
	SceneMountainStep4: "Warm rays of sunlight gently touch your face.",
	SceneMountainStep5: "You see endless expanses below. Everything seems so small and distant.",
	SceneMountainStep6: "Take a deep breath of pure mountain air. Feel the calmness.",
	SceneMountainStep7: "You are safe. This moment belongs only to you.",

	// Scene: Forest
	SceneForestName:  "Forest Clearing",
	SceneForestDesc:  "A quiet clearing among ancient trees",
	SceneForestAtmo:  "Warm sunlight, rustling leaves",
	SceneForestStep1: "Imagine yourself in a cozy clearing in an ancient forest.",
	SceneForestStep2: "Sunlight filters through the tree canopy, creating patterns on the grass.",
	SceneForestStep3: "You hear birds singing and leaves rustling in the wind.",
	SceneForestStep4: "Soft moss beneath your feet. You feel connected to the earth.",
	SceneForestStep5: "The scent of pine and flowers fills the air.",
	SceneForestStep6: "Sit down on the warm grass. Feel nature embracing you.",
	SceneForestStep7: "There is no rush here. Only peace and harmony with nature.",

	// Scene: Beach
	SceneBeachName:  "Ocean Beach",
	SceneBeachDesc:  "A calm beach with warm sand",
	SceneBeachAtmo:  "Sea breeze, sound of waves",
	SceneBeachStep1: "You walk barefoot on warm sand along the ocean.",
	SceneBeachStep2: "Waves gently roll onto the shore and retreat back.",
	SceneBeachStep3: "A light sea breeze refreshes your face and plays with your hair.",
	SceneBeachStep4: "You feel the warmth of sand under your feet with each step.",
	SceneBeachStep5: "Seagulls fly in the distance. Their cries blend with the sound of waves.",
	SceneBeachStep6: "Stop and look at the endless horizon.",
	SceneBeachStep7: "The ocean is infinite, just like your possibilities. Feel the freedom.",

	// Scene: Garden
	SceneGardenName:  "Blooming Garden",
	SceneGardenDesc:  "A beautiful garden with flowers and a fountain",
	SceneGardenAtmo:  "Fragrance of flowers, murmur of water",
	SceneGardenStep1: "You enter a beautiful garden full of blooming plants.",
	SceneGardenStep2: "Roses, lavender, jasmine — their fragrances mix in the air.",
	SceneGardenStep3: "A small fountain bubbles in the center of the garden.",
	SceneGardenStep4: "Butterflies flutter between flowers. Everything is full of life.",
	SceneGardenStep5: "You sit on a bench by the fountain and close your eyes.",
	SceneGardenStep6: "The sound of bubbling water calms your mind.",
	SceneGardenStep7: "This garden is your safe place. You can return here anytime.",

	// Scene: Starry
	SceneStarryName:  "Starry Night",
	SceneStarryDesc:  "A night meadow under the starry sky",
	SceneStarryAtmo:  "Cool night air, silence",
	SceneStarryStep1: "You lie on soft grass on a warm summer night.",
	SceneStarryStep2: "Above you is an endless sky dotted with millions of stars.",
	SceneStarryStep3: "The Milky Way stretches across the sky — a river of light.",
	SceneStarryStep4: "The night air is pleasantly cool. You hear crickets chirping.",
	SceneStarryStep5: "Each star is a sun in a distant galaxy.",
	SceneStarryStep6: "Feel your place in the universe. You are part of something vast.",
	SceneStarryStep7: "Your worries dissolve in the infinity of space.",

	// Language selection
	LangSelectTitle: "🌐 *Choose language:*",
	CycleWord:       "Cycle",
}
