//nolint:dupl // Localization files have identical structure by design, only text values differ.
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

	// Welcome holder — static message with bot info
	HolderText: `🌿 *Helfer bei Angst*

Dieser Bot hilft bei der Bewältigung von Angst mit evidenzbasierten Techniken:

• 🌬️ Atmung 4-4-6 — schnelle Beruhigung
• 🌿 Erdung 5-4-3-2-1 — Rückkehr in den Moment
• 🧘 Geführte Atmung — verschiedene Muster
• 💪 Muskelentspannung — Spannung lösen
• 🏷️ Gedanken markieren — mit Angst arbeiten
• 🌅 Visualisierung — Entspannung durch Vorstellung

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
💪 Muskelentspannung — Spannung lösen
🏷️ Gedanken markieren — mit Angst arbeiten
🌅 Visualisierung — Entspannung durch Vorstellung`,
	MenuBreathing:     "🌬️ Atmung 2 Min",
	MenuGrounding:     "🌿 Erdung",
	MenuGuided:        "🧘 Geführt",
	MenuPMR:           "💪 Muskeln",
	MenuThought:       "🏷️ Gedanken",
	MenuVisualization: "🌅 Visualisierung",
	MenuLang:          "🌐 Sprache",

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

	// Visualization
	VisualizationIntro: `🌅 *Friedliche Visualisierung*

Eine Entspannungstechnik durch Vorstellung friedlicher Orte.

*Wählen Sie eine Szene:*

`,
	VisualizationStopped: `🌅 *Friedliche Visualisierung*

Übung gestoppt. Wählen Sie eine andere Szene:

`,
	VisualizationCompletion: `✅ *Sehr gut!*

Sie haben die Visualisierung "%s" abgeschlossen.

Kehren Sie langsam in die Realität zurück. Bewegen Sie Ihre Finger, atmen Sie tief ein und öffnen Sie die Augen.

Wie fühlen Sie sich?`,
	VisualizationThanks: `✨ *Danke für die Übung!*

Visualisierung ist eine kraftvolle Technik zur Reduzierung von Stress und Angst. Regelmäßige Übung verstärkt die Wirkung.

Wählen Sie eine Technik aus dem Menü unten.`,
	VisualizationAtmosphere: "Atmosphäre",
	VisualizationCloseEyes:  "Schließen Sie die Augen und tauchen Sie in diese Szene ein...",
	VisualizationStepFmt:    "Schritt %d/%d",

	// Scene: Mountain
	SceneMountainName:  "Berggipfel",
	SceneMountainDesc:  "Sonnenaufgang auf dem Berggipfel",
	SceneMountainAtmo:  "Kühle, frische Bergluft",
	SceneMountainStep1: "Schließen Sie die Augen und stellen Sie sich vor, Sie stehen am frühen Morgen auf einem Berggipfel.",
	SceneMountainStep2: "Stille umgibt Sie. Sie spüren die kühle Bergluft auf Ihrer Haut.",
	SceneMountainStep3: "Die Sonne erscheint am Horizont und taucht den Himmel in Rosa- und Orangetöne.",
	SceneMountainStep4: "Warme Sonnenstrahlen berühren sanft Ihr Gesicht.",
	SceneMountainStep5: "Sie sehen endlose Weiten unter sich. Alles erscheint so klein und fern.",
	SceneMountainStep6: "Atmen Sie tief die reine Bergluft ein. Spüren Sie die Ruhe.",
	SceneMountainStep7: "Sie sind sicher. Dieser Moment gehört nur Ihnen.",

	// Scene: Forest
	SceneForestName:  "Waldlichtung",
	SceneForestDesc:  "Eine ruhige Lichtung zwischen alten Bäumen",
	SceneForestAtmo:  "Warmes Sonnenlicht, raschelnde Blätter",
	SceneForestStep1: "Stellen Sie sich vor, Sie befinden sich auf einer gemütlichen Lichtung in einem alten Wald.",
	SceneForestStep2: "Sonnenlicht filtert durch die Baumkronen und erzeugt Muster im Gras.",
	SceneForestStep3: "Sie hören Vögel singen und Blätter im Wind rascheln.",
	SceneForestStep4: "Weiches Moos unter Ihren Füßen. Sie fühlen sich mit der Erde verbunden.",
	SceneForestStep5: "Der Duft von Kiefern und Blumen erfüllt die Luft.",
	SceneForestStep6: "Setzen Sie sich auf das warme Gras. Spüren Sie, wie die Natur Sie umarmt.",
	SceneForestStep7: "Hier gibt es keine Hektik. Nur Frieden und Harmonie mit der Natur.",

	// Scene: Beach
	SceneBeachName:  "Ozeanstrand",
	SceneBeachDesc:  "Ein ruhiger Strand mit warmem Sand",
	SceneBeachAtmo:  "Meeresbrise, Wellenrauschen",
	SceneBeachStep1: "Sie gehen barfuß über warmen Sand am Meer entlang.",
	SceneBeachStep2: "Wellen rollen sanft an den Strand und ziehen sich zurück.",
	SceneBeachStep3: "Eine leichte Meeresbrise erfrischt Ihr Gesicht und spielt mit Ihrem Haar.",
	SceneBeachStep4: "Sie spüren die Wärme des Sandes unter Ihren Füßen bei jedem Schritt.",
	SceneBeachStep5: "Möwen fliegen in der Ferne. Ihre Rufe vermischen sich mit dem Wellenrauschen.",
	SceneBeachStep6: "Bleiben Sie stehen und schauen Sie auf den endlosen Horizont.",
	SceneBeachStep7: "Der Ozean ist unendlich, genau wie Ihre Möglichkeiten. Spüren Sie die Freiheit.",

	// Scene: Garden
	SceneGardenName:  "Blühender Garten",
	SceneGardenDesc:  "Ein schöner Garten mit Blumen und Brunnen",
	SceneGardenAtmo:  "Blumenduft, plätscherndes Wasser",
	SceneGardenStep1: "Sie betreten einen wunderschönen Garten voller blühender Pflanzen.",
	SceneGardenStep2: "Rosen, Lavendel, Jasmin — ihre Düfte vermischen sich in der Luft.",
	SceneGardenStep3: "Ein kleiner Brunnen plätschert in der Mitte des Gartens.",
	SceneGardenStep4: "Schmetterlinge flattern zwischen den Blumen. Alles ist voller Leben.",
	SceneGardenStep5: "Sie setzen sich auf eine Bank am Brunnen und schließen die Augen.",
	SceneGardenStep6: "Das Geräusch des plätschernden Wassers beruhigt Ihren Geist.",
	SceneGardenStep7: "Dieser Garten ist Ihr sicherer Ort. Sie können jederzeit hierher zurückkehren.",

	// Scene: Starry
	SceneStarryName:  "Sternennacht",
	SceneStarryDesc:  "Eine nächtliche Wiese unter dem Sternenhimmel",
	SceneStarryAtmo:  "Kühle Nachtluft, Stille",
	SceneStarryStep1: "Sie liegen auf weichem Gras in einer warmen Sommernacht.",
	SceneStarryStep2: "Über Ihnen ist ein endloser Himmel, übersät mit Millionen von Sternen.",
	SceneStarryStep3: "Die Milchstraße erstreckt sich über den Himmel — ein Fluss aus Licht.",
	SceneStarryStep4: "Die Nachtluft ist angenehm kühl. Sie hören Grillen zirpen.",
	SceneStarryStep5: "Jeder Stern ist eine Sonne in einer fernen Galaxie.",
	SceneStarryStep6: "Spüren Sie Ihren Platz im Universum. Sie sind Teil von etwas Großem.",
	SceneStarryStep7: "Ihre Sorgen lösen sich in der Unendlichkeit des Weltraums auf.",

	// Language selection
	LangSelectTitle: "🌐 *Sprache wählen:*",
	CycleWord:       "Zyklus",
}
