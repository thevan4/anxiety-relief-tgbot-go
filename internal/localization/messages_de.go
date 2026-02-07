//nolint:dupl // Localization files have identical structure by design, only text values differ.
package localization

//nolint:gochecknoglobals // Localization bundle.
var messagesDE = Messages{
	// Buttons
	Start:      "▶️ Starten",
	Back:       "◀️ Zurück",
	Stop:       "🛑 Stopp",
	Next:       "➡️ Weiter",
	Done:       "✅ Fertig",
	FeelBetter: "✅ Besser",
	Repeat:     "🔄 Wiederholen",

	// Welcome holder — static message with bot info
	HolderText: `🌿 *Helfer bei Angst*

Dieser Bot hilft bei der Bewältigung von Angst mit evidenzbasierten Techniken:

• 🌬️ Atmung 4-4-6 — schnelle Beruhigung
• 🌿 Erdung 5-4-3-2-1 — Rückkehr in den Moment
• 🧘 Geführte Atmung — verschiedene Muster
• 💪 Muskelentspannung — Spannung lösen

Datenschutz: Der Bot speichert keine persönlichen Daten.

👇 Drücken Sie die Taste unten, um zu beginnen`,

	// Session expired
	SessionExpired: "Sitzung abgelaufen. Drücken Sie /start",

	// Main menu
	MainMenuText: `Wählen Sie eine Technik zur Angstbewältigung:

*Schnell (2-5 Min):*
🌬️ Atmung — Beruhigung des Nervensystems
🌿 Erdung — Rückkehr in den Moment

*Fortgeschritten (5-15 Min):*
🧘 Geführte Atmung — verschiedene Muster
💪 Muskelentspannung — Spannung lösen`,
	MenuBreathing: "🌬️ Atmung 2 Min",
	MenuGrounding: "🌿 Erdung",
	MenuGuided:    "🧘 Geführt",
	MenuPMR:       "💪 Muskeln",
	MenuLang:      "🌐 Sprache",

	// Breathing
	BreathingIntro: `🌬️ *2-Minuten-Atmung*

Eine einfache Übung zur Beruhigung des Nervensystems.

*Anleitung:*
1️⃣ Setzen Sie sich bequem hin, entspannen Sie sich
2️⃣ Einatmen durch die Nase — 4 Sekunden
3️⃣ Atem anhalten — 4 Sekunden
4️⃣ Ausatmen durch den Mund — 6 Sekunden
5️⃣ 8 Zyklen wiederholen (~2 Minuten)

Ich werde Sie durch jeden Schritt führen. Bereit zu beginnen?`,
	BreathingCompletion: `✅ *Sehr gut!*

Sie haben die Atemübung abgeschlossen.
Wie fühlen Sie sich?`,
	BreathingInhale: "Einatmen",
	BreathingHold:   "Halten",
	BreathingExhale: "Ausatmen",

	// Grounding
	GroundingIntro: `🌿 *5-4-3-2-1 Erdungstechnik*

Diese Technik hilft Ihnen, durch Ihre Sinne in den gegenwärtigen Moment zurückzukehren.

*Wie es funktioniert:*
Sie werden Dinge um sich herum benennen, die Sie mit verschiedenen Sinnen wahrnehmen.

Bereit zu beginnen?`,

	// Grounding steps
	GroundingStep1Title: "Schritt 1: Sehen",
	GroundingStep1Desc: "Nennen Sie 5 Dinge, die Sie um sich herum sehen.\n\n" +
		"Zum Beispiel: Tisch, Lampe, Buch, Fenster, Stuhl.",
	GroundingStep2Title: "Schritt 2: Fühlen",
	GroundingStep2Desc: "Nennen Sie 4 Dinge, die Sie mit Ihrem Körper fühlen.\n\n" +
		"Zum Beispiel: Füße auf dem Boden, Rücken am Stuhl, Hände auf dem Tisch, Kleidung am Körper.",
	GroundingStep3Title: "Schritt 3: Hören",
	GroundingStep3Desc: "Nennen Sie 3 Geräusche, die Sie hören.\n\n" +
		"Zum Beispiel: Straßenlärm, Ticken der Uhr, Ihre Atmung.",
	GroundingStep4Title: "Schritt 4: Riechen",
	GroundingStep4Desc: "Nennen Sie 2 Gerüche, die Sie wahrnehmen oder mögen.\n\n" +
		"Zum Beispiel: Kaffee, frische Luft, Blumenduft.",
	GroundingStep5Title: "Schritt 5: Schmecken",
	GroundingStep5Desc: "Nennen Sie 1 Geschmack in Ihrem Mund.\n\n" +
		"Wenn nichts — erinnern Sie sich an Ihren Lieblingsgeschmack.",

	// Guided breathing
	GuidedIntro:      "🧘 *Geführte Atmung*\n\nWählen Sie eine Atemtechnik:\n\n",
	GuidedCompletion: "✅ *Sehr gut!*\n\nSie haben die Übung \"%s\" abgeschlossen.\n\nWie fühlen Sie sich?",
	GuidedStopped:    "🧘 *Geführte Atmung*\n\nÜbung gestoppt. Wählen Sie eine Technik:\n\n",

	// Breathing patterns
	PatternBoxName:        "Box-Atmung 4-4-4-4",
	PatternBoxDesc:        "Fokus und Konzentration",
	PatternRelaxingName:   "Entspannend 4-7-8",
	PatternRelaxingDesc:   "Tiefe Entspannung und Schlaf",
	PatternEnergizingName: "Energetisierend 4-4-6",
	PatternEnergizingDesc: "Energie und geistige Klarheit",
	PatternQuickName:      "Schneller Reset 3-3-3",
	PatternQuickDesc:      "Schnelle Spannungslösung",

	// Guided breathing phases
	GuidedInhale:  "Einatmen",
	GuidedHoldIn:  "Halten",
	GuidedExhale:  "Ausatmen",
	GuidedHoldOut: "Pause",

	// PMR
	PMRIntro: `💪 *Progressive Muskelentspannung*

Eine Tiefenentspannungstechnik durch Anspannen und Entspannen der Muskeln.

*Wie es funktioniert:*
1. Spannen Sie eine Muskelgruppe für 7 Sekunden an
2. Entspannen Sie für 15 Sekunden
3. Gehen Sie zur nächsten Gruppe über

*%d Muskelgruppen* — von den Händen bis zu den Füßen.

Bereit zu beginnen?`,
	PMRTense: "🔴 *ANSPANNEN*\n\n%s",
	PMRRelax: "🟢 *ENTSPANNEN*\n\n%s",
	PMRCompletion: `✅ *Sehr gut!*

Sie haben die progressive Muskelentspannung abgeschlossen.

Ihr Körper ist jetzt vollständig entspannt. Sitzen Sie noch eine Minute und genießen Sie diesen Zustand.`,

	// PMR muscle groups
	PMRHandsName:  "Hände",
	PMRHandsTense: "Ballen Sie die Fäuste so fest wie möglich. Spüren Sie die Spannung in Fingern und Handflächen.",
	PMRHandsRelax: "Öffnen Sie die Fäuste und entspannen Sie die Hände. Spüren Sie Wärme und Entspannung.",

	PMRForearmsName:  "Unterarme und Bizeps",
	PMRForearmsTense: "Beugen Sie die Arme an den Ellbogen und spannen Sie den Bizeps an. Halten Sie die Spannung.",
	PMRForearmsRelax: "Senken Sie die Arme und entspannen Sie sie vollständig. Ihre Arme werden schwer.",

	PMRForeheadName:  "Stirn",
	PMRForeheadTense: "Heben Sie die Augenbrauen so hoch wie möglich und runzeln Sie die Stirn. Spüren Sie die Spannung.",
	PMRForeheadRelax: "Entspannen Sie die Stirn. Lassen Sie die Augenbrauen sinken. Die Stirn wird glatt.",

	PMREyesName:  "Augen und Nase",
	PMREyesTense: "Drücken Sie die Augen fest zusammen und runzeln Sie die Nase. Spüren Sie die Spannung um die Augen.",
	PMREyesRelax: "Entspannen Sie Augen und Nase. Ihre Augenlider werden leicht und ruhig.",

	PMRJawName:  "Kiefer",
	PMRJawTense: "Pressen Sie den Kiefer zusammen und ziehen Sie die Lippen zu einem angespannten Lächeln.",
	PMRJawRelax: "Entspannen Sie den Kiefer, öffnen Sie den Mund leicht. Die Zunge ist entspannt.",

	PMRNeckName:  "Nacken und Schultern",
	PMRNeckTense: "Ziehen Sie die Schultern zu den Ohren und spannen Sie den Nacken an. Halten Sie die Spannung.",
	PMRNeckRelax: "Senken Sie die Schultern und entspannen Sie den Nacken. Spüren Sie die Erleichterung.",

	PMRChestName:  "Brust und Rücken",
	PMRChestTense: "Atmen Sie tief ein, halten Sie den Atem an und spannen Sie Brust- und Rückenmuskeln an.",
	PMRChestRelax: "Atmen Sie langsam aus und entspannen Sie Brust und Rücken. Atmen Sie ruhig.",

	PMRStomachName:  "Bauch",
	PMRStomachTense: "Ziehen Sie den Bauch ein und spannen Sie die Bauchmuskeln an. Halten Sie die Spannung.",
	PMRStomachRelax: "Entspannen Sie den Bauch. Lassen Sie ihn sich frei mit Ihrer Atmung bewegen.",

	PMRThighsName:  "Oberschenkel und Gesäß",
	PMRThighsTense: "Spannen Sie Gesäß und Oberschenkel an, drücken Sie sie gegen die Sitzfläche.",
	PMRThighsRelax: "Entspannen Sie Gesäß und Oberschenkel. Spüren Sie die Schwere in den Beinen.",

	PMRCalvesName:  "Waden und Füße",
	PMRCalvesTense: "Ziehen Sie die Zehen zu sich heran und spannen Sie die Waden an. Spüren Sie die Dehnung.",
	PMRCalvesRelax: "Entspannen Sie Füße und Waden. Ihre Beine werden warm und schwer.",

	// Language selection
	LangSelectTitle: "🌐 *Sprache wählen:*",
	CycleWord:       "Zyklus",
}
