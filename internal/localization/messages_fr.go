//nolint:dupl // Localization files have identical structure by design, only text values differ.
package localization

var messagesFR = Messages{
	// Buttons
	Start:      "▶️ Commencer",
	Back:       "◀️ Retour",
	Stop:       "🛑 Arrêter",
	Next:       "➡️ Suivant",
	Done:       "✅ Terminé",
	FeelBetter: "✅ Mieux",
	Repeat:     "🔄 Répéter",
	BackToMenu: "🏠 Menu",

	// Welcome holder — static message with bot info
	HolderText: `🌿 *Aide contre l'anxiété*

Ce bot aide à gérer l'anxiété avec des techniques fondées sur des preuves :

• 🌬️ Respiration 4-4-6 — calme rapide
• 🌿 Ancrage 5-4-3-2-1 — retour au moment
• 🧘 Respiration guidée — différents schémas
• 💪 Relaxation musculaire — relâcher la tension
• 🏷️ Étiquetage des pensées — travailler l'anxiété
• 🌅 Visualisation — relaxation par l'imagination

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
💪 Relaxation musculaire — relâcher la tension
🏷️ Étiquetage des pensées — travailler l'anxiété
🌅 Visualisation — relaxation par l'imagination`,
	MenuBreathing:     "🌬️ Respiration 2 min",
	MenuGrounding:     "🌿 Ancrage",
	MenuGuided:        "🧘 Guidée",
	MenuPMR:           "💪 Muscles",
	MenuThought:       "🏷️ Pensées",
	MenuVisualization: "🌅 Visualisation",
	MenuInfo:          "ℹ️ Info",
	MenuLang:          "🌐 Langue",

	// Breathing
	BreathingIntro: `🌬️ *Respiration de 2 minutes*

Un exercice simple pour calmer le système nerveux.

*Instructions :*
1️⃣ Asseyez-vous confortablement, fermez les yeux
2️⃣ Inspirez par le nez — 4 secondes
3️⃣ Retenez votre souffle — 4 secondes
4️⃣ Expirez par la bouche — 6 secondes
5️⃣ Répétez 8 cycles (~2 minutes)

Je vous guiderai à chaque étape. Prêt à commencer ?`,
	BreathingCompletion: `✅ *Excellent !*

Vous avez terminé l'exercice de respiration.
Comment vous sentez-vous ?`,
	BreathingThanks: `✨ *Merci pour la pratique !*

Des exercices réguliers aident à réduire le niveau d'anxiété.

Choisissez une technique dans le menu ci-dessous.`,
	BreathingCycle:  "Cycle %d/%d",
	BreathingInhale: "Inspirer",
	BreathingHold:   "Retenir",
	BreathingExhale: "Expirer",
	BreathingPause:  "Pause",

	// Grounding
	GroundingIntro: `🌿 *Technique d'ancrage 5-4-3-2-1*

Cette technique vous aide à revenir au moment présent à travers vos sens.

*Comment ça marche :*
Vous nommerez des choses autour de vous que vous percevez avec différents sens.

Prêt à commencer ?`,
	GroundingStep:       "Nommez *%d %s* que vous pouvez *%s*",
	GroundingCompletion: "✅ *Excellent !*\n\nVous avez terminé la technique d'ancrage.\nComment vous sentez-vous ?",
	GroundingThanks: `✨ *Merci pour la pratique !*

La technique 5-4-3-2-1 aide à revenir rapidement au moment présent lors de l'anxiété.

Choisissez une technique dans le menu ci-dessous.`,
	GroundingSee:        "voir",
	GroundingHear:       "entendre",
	GroundingFeel:       "sentir",
	GroundingThings:     "choses",
	GroundingSounds:     "sons",
	GroundingSensations: "sensations",

	// Grounding steps
	GroundingStep1Title: "Étape 1 : Vue",
	GroundingStep1Desc:  "Nommez 5 choses que vous voyez autour de vous.\n\nPar exemple : table, lampe, livre, fenêtre, chaise.",
	GroundingStep2Title: "Étape 2 : Toucher",
	GroundingStep2Desc:  "Nommez 4 choses que vous ressentez avec votre corps.\n\nPar exemple : pieds sur le sol, dos sur la chaise, mains sur la table, vêtements sur le corps.",
	GroundingStep3Title: "Étape 3 : Ouïe",
	GroundingStep3Desc:  "Nommez 3 sons que vous entendez.\n\nPar exemple : bruit de la rue, tic-tac de l'horloge, votre respiration.",
	GroundingStep4Title: "Étape 4 : Odorat",
	GroundingStep4Desc:  "Nommez 2 odeurs que vous sentez ou aimez.\n\nPar exemple : café, air frais, parfum de fleurs.",
	GroundingStep5Title: "Étape 5 : Goût",
	GroundingStep5Desc:  "Nommez 1 goût que vous ressentez dans votre bouche.\n\nSi rien — rappelez-vous votre goût préféré.",

	// Guided breathing
	GuidedIntro:      "🧘 *Respiration guidée*\n\nChoisissez une technique de respiration :\n\n",
	GuidedCompletion: "✅ *Excellent !*\n\nVous avez terminé l'exercice \"%s\".\n\nComment vous sentez-vous ?",
	GuidedThanks: `✨ *Merci pour la pratique !*

Des exercices de respiration réguliers aident à réduire l'anxiété et améliorer la concentration.

Choisissez une technique dans le menu ci-dessous.`,
	GuidedStopped: "🧘 *Respiration guidée*\n\nExercice arrêté. Choisissez une technique :\n\n",

	// Breathing patterns
	PatternBoxName:        "Respiration carrée 4-4-4-4",
	PatternBoxDesc:        "Focus et concentration",
	PatternRelaxingName:   "Relaxant 4-7-8",
	PatternRelaxingDesc:   "Relaxation profonde et sommeil",
	PatternEnergizingName: "Énergisant 4-4-6",
	PatternEnergizingDesc: "Énergie et clarté mentale",
	PatternQuickName:      "Reset rapide 3-3-3",
	PatternQuickDesc:      "Soulagement rapide des tensions",

	// PMR
	PMRIntro: `💪 *Relaxation musculaire progressive*

Une technique de relaxation profonde par la tension et le relâchement des muscles.

*Comment ça marche :*
1. Contractez un groupe musculaire pendant 7 secondes
2. Relâchez pendant 15 secondes
3. Passez au groupe suivant

*%d groupes musculaires* — des mains aux pieds.

Prêt à commencer ?`,
	PMRTense: "🔴 *CONTRACTEZ*\n\n%s",
	PMRRelax: "🟢 *RELÂCHEZ*\n\n%s",
	PMRCompletion: `✅ *Excellent !*

Vous avez terminé la relaxation musculaire progressive.

Votre corps est maintenant complètement détendu. Restez assis encore une minute en profitant de cet état.`,
	PMRStopped: `💪 *Relaxation musculaire progressive*

Exercice arrêté.

Voulez-vous recommencer ?`,
	PMRThanks: `✨ *Merci pour la pratique !*

La relaxation musculaire progressive réduit la tension musculaire et le niveau de stress.

Choisissez une technique dans le menu ci-dessous.`,

	// Thought labeling
	ThoughtIntro: `🏷️ *Étiquetage des pensées*

Une technique de pleine conscience pour travailler avec les pensées anxieuses.

*Comment ça marche :*
1. Vous décrivez une pensée anxieuse
2. Choisissez une catégorie de distorsion
3. Obtenez une façon de la reformuler

Cela aide à se séparer des pensées et à les voir objectivement.

Prêt à commencer ?`,
	ThoughtPrompt:     "📝 *Décrivez votre pensée anxieuse*\n\nÉcrivez en un message la pensée qui vous préoccupe.",
	ThoughtCategories: "🏷️ *Catégorisation de la pensée*\n\n_\"%s\"_\n\nChoisissez le type de distorsion cognitive :",
	ThoughtResult:     "%s *%s*\n\n_%s_\n\n*Comment reformuler :*\n%s",
	ThoughtCompletion: "Voulez-vous travailler sur une autre pensée ?",
	ThoughtThanks: `✨ *Merci pour la pratique !*

L'étiquetage des pensées aide à reconnaître les distorsions cognitives et à réduire leur impact.

Choisissez une technique dans le menu ci-dessous.`,

	// Visualization
	VisualizationIntro: `🌅 *Visualisation paisible*

Une technique de relaxation par l'imagination de lieux paisibles.

*Choisissez une scène :*

`,
	VisualizationStopped: `🌅 *Visualisation paisible*

Exercice arrêté. Choisissez une autre scène :

`,
	VisualizationCompletion: `✅ *Excellent !*

Vous avez terminé la visualisation "%s".

Revenez lentement à la réalité. Bougez vos doigts, inspirez profondément et ouvrez les yeux.

Comment vous sentez-vous ?`,
	VisualizationThanks: `✨ *Merci pour la pratique !*

La visualisation est une technique puissante pour réduire le stress et l'anxiété. La pratique régulière renforce l'effet.

Choisissez une technique dans le menu ci-dessous.`,
	VisualizationAtmosphere: "Atmosphère",
	VisualizationCloseEyes:  "Fermez les yeux et plongez dans cette scène...",
	VisualizationStepFmt:    "Étape %d/%d",

	// Scene: Mountain
	SceneMountainName:  "Sommet de montagne",
	SceneMountainDesc:  "Lever de soleil au sommet de la montagne",
	SceneMountainAtmo:  "Air frais et pur de la montagne",
	SceneMountainStep1: "Fermez les yeux et imaginez-vous debout sur un sommet de montagne tôt le matin.",
	SceneMountainStep2: "Le silence vous entoure. Vous sentez l'air frais de la montagne sur votre peau.",
	SceneMountainStep3: "Le soleil apparaît à l'horizon, peignant le ciel de teintes roses et orangées.",
	SceneMountainStep4: "Les rayons chauds du soleil touchent doucement votre visage.",
	SceneMountainStep5: "Vous voyez des étendues infinies en dessous. Tout semble si petit et lointain.",
	SceneMountainStep6: "Prenez une grande inspiration d'air pur de montagne. Ressentez le calme.",
	SceneMountainStep7: "Vous êtes en sécurité. Ce moment n'appartient qu'à vous.",

	// Scene: Forest
	SceneForestName:  "Clairière forestière",
	SceneForestDesc:  "Une clairière tranquille parmi les arbres anciens",
	SceneForestAtmo:  "Lumière chaude du soleil, bruissement des feuilles",
	SceneForestStep1: "Imaginez-vous dans une clairière confortable au milieu d'une forêt ancienne.",
	SceneForestStep2: "La lumière du soleil filtre à travers la canopée, créant des motifs sur l'herbe.",
	SceneForestStep3: "Vous entendez les oiseaux chanter et les feuilles bruisser dans le vent.",
	SceneForestStep4: "Mousse douce sous vos pieds. Vous vous sentez connecté à la terre.",
	SceneForestStep5: "Le parfum des pins et des fleurs emplit l'air.",
	SceneForestStep6: "Asseyez-vous sur l'herbe chaude. Sentez la nature vous envelopper.",
	SceneForestStep7: "Il n'y a pas de précipitation ici. Seulement la paix et l'harmonie avec la nature.",

	// Scene: Beach
	SceneBeachName:  "Plage océanique",
	SceneBeachDesc:  "Une plage calme avec du sable chaud",
	SceneBeachAtmo:  "Brise marine, bruit des vagues",
	SceneBeachStep1: "Vous marchez pieds nus sur le sable chaud le long de l'océan.",
	SceneBeachStep2: "Les vagues roulent doucement sur le rivage et se retirent.",
	SceneBeachStep3: "Une légère brise marine rafraîchit votre visage et joue avec vos cheveux.",
	SceneBeachStep4: "Vous sentez la chaleur du sable sous vos pieds à chaque pas.",
	SceneBeachStep5: "Des mouettes volent au loin. Leurs cris se mêlent au bruit des vagues.",
	SceneBeachStep6: "Arrêtez-vous et regardez l'horizon infini.",
	SceneBeachStep7: "L'océan est infini, tout comme vos possibilités. Ressentez la liberté.",

	// Scene: Garden
	SceneGardenName:  "Jardin fleuri",
	SceneGardenDesc:  "Un beau jardin avec des fleurs et une fontaine",
	SceneGardenAtmo:  "Parfum des fleurs, murmure de l'eau",
	SceneGardenStep1: "Vous entrez dans un magnifique jardin plein de plantes fleuries.",
	SceneGardenStep2: "Roses, lavande, jasmin — leurs parfums se mélangent dans l'air.",
	SceneGardenStep3: "Une petite fontaine murmure au centre du jardin.",
	SceneGardenStep4: "Des papillons voltigent entre les fleurs. Tout est plein de vie.",
	SceneGardenStep5: "Vous vous asseyez sur un banc près de la fontaine et fermez les yeux.",
	SceneGardenStep6: "Le son de l'eau qui coule apaise votre esprit.",
	SceneGardenStep7: "Ce jardin est votre lieu sûr. Vous pouvez y revenir à tout moment.",

	// Scene: Starry
	SceneStarryName:  "Nuit étoilée",
	SceneStarryDesc:  "Une prairie nocturne sous le ciel étoilé",
	SceneStarryAtmo:  "Air frais de la nuit, silence",
	SceneStarryStep1: "Vous êtes allongé sur l'herbe douce par une chaude nuit d'été.",
	SceneStarryStep2: "Au-dessus de vous, un ciel infini parsemé de millions d'étoiles.",
	SceneStarryStep3: "La Voie lactée s'étend à travers le ciel — une rivière de lumière.",
	SceneStarryStep4: "L'air nocturne est agréablement frais. Vous entendez les grillons chanter.",
	SceneStarryStep5: "Chaque étoile est un soleil dans une galaxie lointaine.",
	SceneStarryStep6: "Ressentez votre place dans l'univers. Vous faites partie de quelque chose d'immense.",
	SceneStarryStep7: "Vos soucis se dissolvent dans l'infini de l'espace.",

	// Info
	InfoText: `ℹ️ *À propos du bot*

Ce bot aide à gérer l'anxiété avec des techniques basées sur des preuves.

*Techniques disponibles :*
• 🌬️ Respiration 4-4-6 — calme rapide
• 🌿 Ancrage 5-4-3-2-1 — retour au moment
• 🧘 Respiration guidée — différents schémas
• 💪 Relaxation musculaire — relâcher la tension
• 🏷️ Étiquetage des pensées — travailler l'anxiété
• 🌅 Visualisation — relaxation par l'imagination

*Confidentialité :*
Le bot ne stocke pas de données personnelles. Les sessions sont automatiquement supprimées.

Code source : github.com/thevan4/anxiety-relief-tgbot-go`,

	// Language selection
	LangSelectTitle: "🌐 *Choisissez une langue :*",
	CycleWord:       "Cycle",
}
