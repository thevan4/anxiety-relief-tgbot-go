//nolint:dupl // Localization files have identical structure by design, only text values differ.
package localization

//nolint:gochecknoglobals // Localization bundle.
var messagesDE = Messages{
	// Buttons
	Start:  "▶️ Starten",
	Back:   "◀️ Zurück",
	Stop:   "🛑 Stopp",
	Next:   "➡️ Weiter",
	Done:   "✨ Zum Menü",
	Repeat: "🔄 Wiederholen",
	Pause:  "⏸ Pause",
	Resume: "▶️ Fortsetzen",

	PauseText: "⏸ *Pause*\n\nNehmen Sie sich Zeit. Fahren Sie fort, wenn Sie bereit sind.",

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
	SessionExpired: "Aktuelle Sitzung abgelaufen. Senden Sie /start oder tippen Sie auf die Schaltfläche „Start“ im Chat.",

	// Main menu
	MainMenuText: `Wählen Sie eine Technik:

*Schnell (2-5 Min):*
🌬️ Atmung — Beruhigung des Nervensystems
🌿 Erdung — Rückkehr in den Moment

*Fortgeschritten (5-15 Min):*
🧘 Atemmuster
💪 Muskelentspannung`,
	MenuBreathing: "🌬️ Einfache Atmung 2 Min",
	MenuGrounding: "🌿 Erdung 3-5 Min",
	MenuGuided:    "🧘 Atemmuster 1-2 Min",
	MenuPMR:       "💪 Muskelentspannung 4 Min",
	MenuLang:      "🌐 Sprache",

	// Breathing
	BreathingIntro: `🌬️ *Einfache Atmung*

Eine einfache Übung zur Beruhigung des Nervensystems.

*Anleitung:*
1️⃣ Setzen Sie sich bequem hin, entspannen Sie sich
2️⃣ Einatmen durch die Nase — 4 Sekunden
3️⃣ Atem anhalten — 4 Sekunden
4️⃣ Ausatmen durch den Mund — 6 Sekunden
5️⃣ 8 Zyklen wiederholen (~2 Minuten)

Ich werde Sie durch jeden Schritt führen. Bereit zu beginnen?`,
	BreathingCompletion: "✅ *Einfache Atmung*\n\nFertig. Sie können weitermachen.",
	BreathingInhale:     "Einatmen",
	BreathingHold:       "Halten",
	BreathingExhale:     "Ausatmen",

	// Grounding
	GroundingIntro: `🌿 *5-4-3-2-1 Erdungstechnik*

Bei Angst und Panik wandert die Aufmerksamkeit in den Kopf und zu negativen Gedanken.
Diese Technik bringt Sie durch 5 Sinne in den gegenwärtigen Moment zurück.

*Wie es funktioniert:*
5 Schritte mit abnehmenden Zahlen:
• Schritt 1: Nennen Sie 5 Dinge (Sehen)
• Schritt 2: Nennen Sie 4 Dinge (Fühlen)
• Schritt 3: Nennen Sie 3 Geräusche (Hören)
• Schritt 4: Nennen Sie 2 Gerüche (Riechen)
• Schritt 5: Nennen Sie 1 Geschmack (Schmecken)

*Zeit:* ~3-5 Minuten

Ich begleite Sie durch jeden Schritt. Bereit zu beginnen?`,
	GroundingCompletion: "✅ *Erdung 5-4-3-2-1*\n\nSie haben etwas Schwieriges geschafft.",

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
	GuidedIntro: `🧘 *Atemmuster*

Verschiedene Atemmuster für verschiedene Ziele — von Fokus bis Entspannung.

Wählen Sie eine Technik:

📦 *Box-Atmung 4-4-4-4* (~1.5 Min)
Gleichmäßiger Rhythmus für Fokus und Konzentration. Wird vom Militär und Sportlern verwendet.

😴 *Entspannend 4-7-8* (~1.5 Min)
Langer Ausatem für tiefe Entspannung und Schlafvorbereitung.

⚡ *Energetisierend 4-4-6* (~1.5 Min)
Energie und geistige Klarheit, Linderung von Müdigkeit.

🚀 *Schneller Reset 3-3-3* (~1 Min)
Schnelle Linderung von akuter Spannung und Stress.
`,
	GuidedCompletion: "✅ *Atemmuster*\n\nFertig. Sie können weitermachen.",

	// Guided breathing — pattern intros
	PatternBoxIntro: `📦 *Box-Atmung 4-4-4-4*

Gleichmäßiger Rhythmus: Einatmen, Halten, Ausatmen, Pause — je 4 Sekunden.
Wird vom Militär und Sportlern für schnellen Fokus verwendet.

*6 Zyklen, ~1.5 Minuten*

Ich begleite Sie durch jeden Zyklus. Bereit zu beginnen?`,
	PatternRelaxingIntro: `😴 *Entspannend 4-7-8*

Einatmen 4 Sek, Halten 7 Sek, langes Ausatmen 8 Sek.
Das lange Ausatmen aktiviert den Parasympathikus — tiefe Entspannung und Schlafvorbereitung.

*4 Zyklen, ~1.5 Minuten*

Ich begleite Sie durch jeden Zyklus. Bereit zu beginnen?`,
	PatternEnergizingIntro: `⚡ *Energetisierend 4-4-6*

Einatmen 4 Sek, Halten 4 Sek, Ausatmen 6 Sek.
Energie und geistige Klarheit, schnelle Linderung von Müdigkeit.

*6 Zyklen, ~1.5 Minuten*

Ich begleite Sie durch jeden Zyklus. Bereit zu beginnen?`,
	PatternQuickIntro: `🚀 *Schneller Reset 3-3-3*

Einatmen 3 Sek, Halten 3 Sek, Ausatmen 3 Sek.
Minimales Muster zur schnellen Linderung akuter Spannung.

*5 Zyklen, ~1 Minute*

Ich begleite Sie durch jeden Zyklus. Bereit zu beginnen?`,

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

Tiefe Entspannung des gesamten Körpers durch bewusstes Anspannen und Entspannen der Muskeln.
Hilft, körperliche Verspannungen durch Stress und Angst zu lösen.

*Wie es funktioniert:*
1. Spannen Sie eine Muskelgruppe für 7 Sekunden an
2. Entspannen Sie für 15 Sekunden — spüren Sie den Kontrast
3. Gehen Sie zur nächsten Gruppe über

*%d Muskelgruppen* — von den Händen bis zu den Füßen.
*Zeit:* ~4 Minuten

Ich begleite Sie durch jede Muskelgruppe. Bereit zu beginnen?`,
	PMRTense:      "🔴 *ANSPANNEN*\n\n%s",
	PMRRelax:      "🟢 *ENTSPANNEN*\n\n%s",
	PMRCompletion: "✅ *Muskelentspannung*\n\nDer Körper hat gearbeitet. Einfach ruhen jetzt.",

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
