//nolint:dupl // Localization files have identical structure by design, only text values differ.
package localization

//nolint:gochecknoglobals // Localization bundle.
var messagesFR = Messages{
	// Buttons
	Start:  "▶️ Commencer",
	Back:   "◀️ Retour",
	Stop:   "🛑 Arrêter",
	Next:   "➡️ Suivant",
	Done:   "✨ Au menu",
	Repeat: "🔄 Répéter",

	// Welcome holder — static message with bot info
	HolderText: `🌿 *Aide contre l'anxiété*

Ce bot aide à gérer l'anxiété avec des techniques fondées sur des preuves :

• 🌬️ Respiration 4-4-6 — calme rapide
• 🌿 Ancrage 5-4-3-2-1 — retour au moment
• 🧘 Respiration guidée — différents schémas
• 💪 Relaxation musculaire — relâcher la tension

Confidentialité : le bot ne stocke pas de données personnelles.

👇 Appuyez sur le bouton ci-dessous pour commencer`,

	// Session expired
	SessionExpired: "Session expirée. Appuyez sur /start",

	// Main menu
	MainMenuText: `Choisissez une technique pour gérer l'anxiété :

*Rapides (2-5 min) :*
🌬️ Respiration — calmer le système nerveux
🌿 Ancrage — revenir au moment présent

*Avancées (5-15 min) :*
🧘 Respiration guidée — différents schémas
💪 Relaxation musculaire — relâcher la tension`,
	MenuBreathing: "🌬️ Respiration 2 min",
	MenuGrounding: "🌿 Ancrage 3-5 min",
	MenuGuided:    "🧘 Guidée 1-2 min",
	MenuPMR:       "💪 Muscles 4 min",
	MenuLang:      "🌐 Langue",

	// Breathing
	BreathingIntro: `🌬️ *Respiration de 2 minutes*

Un exercice simple pour calmer le système nerveux.

*Instructions :*
1️⃣ Asseyez-vous confortablement, détendez-vous
2️⃣ Inspirez par le nez — 4 secondes
3️⃣ Retenez votre souffle — 4 secondes
4️⃣ Expirez par la bouche — 6 secondes
5️⃣ Répétez 8 cycles (~2 minutes)

Je vous guiderai à chaque étape. Prêt à commencer ?`,
	BreathingCompletion: `✅ *Excellent !*

Vous avez terminé l'exercice de respiration.
Comment vous sentez-vous ?`,
	BreathingInhale: "Inspirer",
	BreathingHold:   "Retenir",
	BreathingExhale: "Expirer",

	// Grounding
	GroundingIntro: `🌿 *Technique d'ancrage 5-4-3-2-1*

Lors de l'anxiété et de la panique, l'attention se déplace vers la tête et les pensées négatives. Cette technique vous ramène au moment présent à travers 5 sens.

*Comment ça marche :*
5 étapes avec nombres décroissants :
• Étape 1 : Nommez 5 choses (vue)
• Étape 2 : Nommez 4 choses (toucher)
• Étape 3 : Nommez 3 sons (ouïe)
• Étape 4 : Nommez 2 odeurs (odorat)
• Étape 5 : Nommez 1 goût (goût)

*Temps :* ~3-5 minutes

Prêt à commencer ?`,

	// Grounding steps
	GroundingStep1Title: "Étape 1 : Vue",
	GroundingStep1Desc: "Nommez 5 choses que vous voyez autour de vous.\n\n" +
		"Par exemple : table, lampe, livre, fenêtre, chaise.",
	GroundingStep2Title: "Étape 2 : Toucher",
	GroundingStep2Desc: "Nommez 4 choses que vous ressentez avec votre corps.\n\n" +
		"Par exemple : pieds sur le sol, dos sur la chaise, mains sur la table, vêtements sur le corps.",
	GroundingStep3Title: "Étape 3 : Ouïe",
	GroundingStep3Desc: "Nommez 3 sons que vous entendez.\n\n" +
		"Par exemple : bruit de la rue, tic-tac de l'horloge, votre respiration.",
	GroundingStep4Title: "Étape 4 : Odorat",
	GroundingStep4Desc:  "Nommez 2 odeurs que vous sentez ou aimez.\n\nPar exemple : café, air frais, parfum de fleurs.",
	GroundingStep5Title: "Étape 5 : Goût",
	GroundingStep5Desc: "Nommez 1 goût que vous ressentez dans votre bouche.\n\n" +
		"Si rien — rappelez-vous votre goût préféré.",

	// Guided breathing
	GuidedIntro: `🧘 *Respiration guidée*

Différents schémas de respiration pour différents objectifs — de la concentration à la relaxation.

Choisissez une technique :

📦 *Respiration carrée 4-4-4-4* (~1.5 min)
Rythme régulier pour la concentration. Utilisée par les militaires et les athlètes.

😴 *Relaxant 4-7-8* (~1.5 min)
Longue expiration pour relaxation profonde et préparation au sommeil.

⚡ *Énergisant 4-4-6* (~1.5 min)
Énergie et clarté mentale, soulagement de la fatigue.

🚀 *Reset rapide 3-3-3* (~1 min)
Soulagement rapide de la tension et du stress aigus.
`,
	GuidedCompletion: "✅ *Excellent !*\n\nVous avez terminé l'exercice \"%s\".\n\nComment vous sentez-vous ?",
	GuidedStopped:    "🧘 *Respiration guidée*\n\nExercice arrêté. Choisissez une technique :\n\n",

	// Breathing patterns
	PatternBoxName:        "Respiration carrée 4-4-4-4",
	PatternBoxDesc:        "Focus et concentration",
	PatternRelaxingName:   "Relaxant 4-7-8",
	PatternRelaxingDesc:   "Relaxation profonde et sommeil",
	PatternEnergizingName: "Énergisant 4-4-6",
	PatternEnergizingDesc: "Énergie et clarté mentale",
	PatternQuickName:      "Reset rapide 3-3-3",
	PatternQuickDesc:      "Soulagement rapide des tensions",

	// Guided breathing phases
	GuidedInhale:  "Inspirez",
	GuidedHoldIn:  "Retenez",
	GuidedExhale:  "Expirez",
	GuidedHoldOut: "Pause",

	// PMR
	PMRIntro: `💪 *Relaxation musculaire progressive*

Relaxation profonde de tout le corps par tension et relâchement conscients des muscles. Aide à relâcher les tensions physiques du stress et de l'anxiété.

*Comment ça marche :*
1. Contractez un groupe musculaire pendant 7 secondes
2. Relâchez pendant 15 secondes — ressentez le contraste
3. Passez au groupe suivant

*%d groupes musculaires* — des mains aux pieds.
*Temps :* ~4 minutes

Prêt à commencer ?`,
	PMRTense: "🔴 *CONTRACTEZ*\n\n%s",
	PMRRelax: "🟢 *RELÂCHEZ*\n\n%s",
	PMRCompletion: `✅ *Excellent !*

Vous avez terminé la relaxation musculaire progressive.

Votre corps est maintenant complètement détendu. Restez assis encore une minute en profitant de cet état.`,

	// PMR muscle groups
	PMRHandsName:  "Mains",
	PMRHandsTense: "Serrez les poings aussi fort que possible. Sentez la tension dans vos doigts et vos paumes.",
	PMRHandsRelax: "Relâchez les poings et détendez les mains. Sentez la chaleur et la détente.",

	PMRForearmsName:  "Avant-bras et biceps",
	PMRForearmsTense: "Pliez les bras aux coudes et contractez les biceps. Gardez la tension.",
	PMRForearmsRelax: "Abaissez les bras et détendez-les complètement. Vos bras deviennent lourds.",

	PMRForeheadName:  "Front",
	PMRForeheadTense: "Levez les sourcils aussi haut que possible et plissez le front. Sentez la tension.",
	PMRForeheadRelax: "Détendez le front. Laissez les sourcils retomber. Le front devient lisse.",

	PMREyesName:  "Yeux et nez",
	PMREyesTense: "Fermez fortement les yeux et froncez le nez. Sentez la tension autour des yeux.",
	PMREyesRelax: "Détendez les yeux et le nez. Les paupières deviennent légères et calmes.",

	PMRJawName:  "Mâchoire",
	PMRJawTense: "Serrez la mâchoire et étirez les lèvres en un sourire tendu.",
	PMRJawRelax: "Détendez la mâchoire, entrouvrez légèrement la bouche. La langue est détendue.",

	PMRNeckName:  "Cou et épaules",
	PMRNeckTense: "Haussez les épaules vers les oreilles et contractez le cou. Gardez la tension.",
	PMRNeckRelax: "Abaissez les épaules et détendez le cou. Sentez le soulagement.",

	PMRChestName:  "Poitrine et dos",
	PMRChestTense: "Prenez une grande inspiration, retenez-la et contractez les muscles de la poitrine et du dos.",
	PMRChestRelax: "Expirez lentement et détendez la poitrine et le dos. Respirez calmement.",

	PMRStomachName:  "Ventre",
	PMRStomachTense: "Rentrez le ventre et contractez les abdos. Gardez la tension.",
	PMRStomachRelax: "Détendez le ventre. Laissez-le bouger librement avec la respiration.",

	PMRThighsName:  "Cuisses et fessiers",
	PMRThighsTense: "Contractez les fessiers et les cuisses, appuyez-les contre le siège.",
	PMRThighsRelax: "Détendez les fessiers et les cuisses. Sentez la lourdeur dans les jambes.",

	PMRCalvesName:  "Mollets et pieds",
	PMRCalvesTense: "Ramenez les orteils vers vous, contractez les mollets. Sentez l'étirement.",
	PMRCalvesRelax: "Détendez les pieds et les mollets. Vos jambes deviennent chaudes et lourdes.",

	// Language selection
	LangSelectTitle: "🌐 *Choisissez une langue :*",
	CycleWord:       "Cycle",
}
