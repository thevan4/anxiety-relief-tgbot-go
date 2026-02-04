package localization

var messagesDE = Messages{
	// Buttons
	Start:      "▶️ Starten",
	Back:       "◀️ Zurück",
	Stop:       "🛑 Stopp",
	Next:       "➡️ Weiter",
	Done:       "✅ Fertig",
	FeelBetter: "✅ Besser",
	Repeat:     "🔄 Wiederholen",
	BackToMenu: "🏠 Menü",

	// Main menu
	MainMenuText: `Wählen Sie eine Technik zur Angstbewältigung:

*Schnell (2-5 Min):*
🌬️ Atmung — Beruhigung des Nervensystems
🌿 Erdung — Rückkehr in den Moment

*Fortgeschritten (5-15 Min):*
🧘 Geführte Atmung — verschiedene Muster
💪 Muskelentspannung — Spannung lösen
🏷️ Gedanken markieren — mit Angst arbeiten`,
	MenuBreathing: "🌬️ Atmung 2 Min",
	MenuGrounding: "🌿 Erdung",
	MenuGuided:    "🧘 Geführt",
	MenuPMR:       "💪 Muskeln",
	MenuThought:   "🏷️ Gedanken",
	MenuInfo:      "ℹ️ Info",
	MenuLang:      "🌐 Sprache",

	// Breathing
	BreathingIntro: `🌬️ *2-Minuten-Atmung*

Eine einfache Übung zur Beruhigung des Nervensystems.

*Anleitung:*
1️⃣ Setzen Sie sich bequem hin, schließen Sie die Augen
2️⃣ Einatmen durch die Nase — 4 Sekunden
3️⃣ Atem anhalten — 4 Sekunden
4️⃣ Ausatmen durch den Mund — 6 Sekunden
5️⃣ 8 Zyklen wiederholen (~2 Minuten)

Ich werde Sie durch jeden Schritt führen. Bereit zu beginnen?`,
	BreathingCompletion: `✅ *Sehr gut!*

Sie haben die Atemübung abgeschlossen.
Wie fühlen Sie sich?`,
	BreathingThanks: `✨ *Danke für die Übung!*

Regelmäßige Übungen helfen, das Angstniveau zu senken.

Wählen Sie eine Technik aus dem Menü unten.`,
	BreathingCycle:  "Zyklus %d/%d",
	BreathingInhale: "Einatmen",
	BreathingHold:   "Halten",
	BreathingExhale: "Ausatmen",
	BreathingPause:  "Pause",

	// Grounding
	GroundingIntro: `🌿 *5-4-3-2-1 Erdungstechnik*

Diese Technik hilft Ihnen, durch Ihre Sinne in den gegenwärtigen Moment zurückzukehren.

*Wie es funktioniert:*
Sie werden Dinge um sich herum benennen, die Sie mit verschiedenen Sinnen wahrnehmen.

Bereit zu beginnen?`,
	GroundingStep:       "Nennen Sie *%d %s*, die Sie *%s*",
	GroundingCompletion: "✅ *Sehr gut!*\n\nSie haben die Erdungstechnik abgeschlossen.\nWie fühlen Sie sich?",
	GroundingThanks: `✨ *Danke für die Übung!*

Die 5-4-3-2-1 Technik hilft, bei Angst schnell in den gegenwärtigen Moment zurückzukehren.

Wählen Sie eine Technik aus dem Menü unten.`,
	GroundingSee:        "sehen",
	GroundingHear:       "hören",
	GroundingFeel:       "fühlen",
	GroundingThings:     "Dinge",
	GroundingSounds:     "Geräusche",
	GroundingSensations: "Empfindungen",

	// Grounding steps
	GroundingStep1Title: "Schritt 1: Sehen",
	GroundingStep1Desc:  "Nennen Sie 5 Dinge, die Sie um sich herum sehen.\n\nZum Beispiel: Tisch, Lampe, Buch, Fenster, Stuhl.",
	GroundingStep2Title: "Schritt 2: Fühlen",
	GroundingStep2Desc:  "Nennen Sie 4 Dinge, die Sie mit Ihrem Körper fühlen.\n\nZum Beispiel: Füße auf dem Boden, Rücken am Stuhl, Hände auf dem Tisch, Kleidung am Körper.",
	GroundingStep3Title: "Schritt 3: Hören",
	GroundingStep3Desc:  "Nennen Sie 3 Geräusche, die Sie hören.\n\nZum Beispiel: Straßenlärm, Ticken der Uhr, Ihre Atmung.",
	GroundingStep4Title: "Schritt 4: Riechen",
	GroundingStep4Desc:  "Nennen Sie 2 Gerüche, die Sie wahrnehmen oder mögen.\n\nZum Beispiel: Kaffee, frische Luft, Blumenduft.",
	GroundingStep5Title: "Schritt 5: Schmecken",
	GroundingStep5Desc:  "Nennen Sie 1 Geschmack in Ihrem Mund.\n\nWenn nichts — erinnern Sie sich an Ihren Lieblingsgeschmack.",

	// Guided breathing
	GuidedIntro:      "🧘 *Geführte Atmung*\n\nWählen Sie eine Atemtechnik:\n\n",
	GuidedCompletion: "✅ *Sehr gut!*\n\nSie haben die Übung \"%s\" abgeschlossen.\n\nWie fühlen Sie sich?",
	GuidedThanks: `✨ *Danke für die Übung!*

Regelmäßige Atemübungen helfen, Angst zu reduzieren und die Konzentration zu verbessern.

Wählen Sie eine Technik aus dem Menü unten.`,
	GuidedStopped: "🧘 *Geführte Atmung*\n\nÜbung gestoppt. Wählen Sie eine Technik:\n\n",

	// Breathing patterns
	PatternBoxName:        "Box-Atmung 4-4-4-4",
	PatternBoxDesc:        "Fokus und Konzentration",
	PatternRelaxingName:   "Entspannend 4-7-8",
	PatternRelaxingDesc:   "Tiefe Entspannung und Schlaf",
	PatternEnergizingName: "Energetisierend 4-4-6",
	PatternEnergizingDesc: "Energie und geistige Klarheit",
	PatternQuickName:      "Schneller Reset 3-3-3",
	PatternQuickDesc:      "Schnelle Spannungslösung",

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
	PMRStopped: `💪 *Progressive Muskelentspannung*

Übung gestoppt.

Möchten Sie von vorne beginnen?`,
	PMRThanks: `✨ *Danke für die Übung!*

Progressive Muskelentspannung reduziert Muskelverspannungen und Stressniveau.

Wählen Sie eine Technik aus dem Menü unten.`,

	// Thought labeling
	ThoughtIntro: `🏷️ *Gedanken markieren*

Eine Achtsamkeitstechnik zur Arbeit mit ängstlichen Gedanken.

*Wie es funktioniert:*
1. Sie beschreiben einen ängstlichen Gedanken
2. Wählen Sie eine Verzerrungskategorie
3. Erhalten Sie einen Weg zur Neuformulierung

Dies hilft, sich von Gedanken zu distanzieren und sie objektiv zu sehen.

Bereit zu beginnen?`,
	ThoughtPrompt:     "📝 *Beschreiben Sie Ihren ängstlichen Gedanken*\n\nSchreiben Sie in einer Nachricht den Gedanken, der Sie beunruhigt.",
	ThoughtCategories: "🏷️ *Kategorisierung des Gedankens*\n\n_\"%s\"_\n\nWählen Sie die Art der kognitiven Verzerrung:",
	ThoughtResult:     "%s *%s*\n\n_%s_\n\n*Wie man umdenkt:*\n%s",
	ThoughtCompletion: "Möchten Sie einen weiteren Gedanken bearbeiten?",
	ThoughtThanks: `✨ *Danke für die Übung!*

Gedanken markieren hilft, kognitive Verzerrungen zu erkennen und ihre Auswirkungen zu reduzieren.

Wählen Sie eine Technik aus dem Menü unten.`,

	// Info
	InfoText: `ℹ️ *Über den Bot*

Dieser Bot hilft bei der Bewältigung von Angst mit evidenzbasierten Techniken.

*Verfügbare Techniken:*
• 🌬️ 4-4-6 Atmung — schnelle Beruhigung
• 🌿 5-4-3-2-1 Erdung — Rückkehr in den Moment
• 🧘 Geführte Atmung — verschiedene Muster
• 💪 Muskelentspannung — Spannung lösen
• 🏷️ Gedanken markieren — mit Angst arbeiten

*Datenschutz:*
Der Bot speichert keine persönlichen Daten. Sitzungen werden automatisch gelöscht.

Quellcode: github.com/thevan4/anxiety-relief-tgbot-go`,

	// Language selection
	LangSelectTitle: "🌐 *Sprache wählen:*",
	CycleWord:       "Zyklus",
}
