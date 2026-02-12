package handlers

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/thevan4/anxiety-relief-tgbot-go/internal/localization"
	"github.com/thevan4/anxiety-relief-tgbot-go/internal/session"
	"github.com/thevan4/anxiety-relief-tgbot-go/internal/statistic"
	"github.com/thevan4/anxiety-relief-tgbot-go/internal/techniques"
	"github.com/thevan4/anxiety-relief-tgbot-go/internal/telegram/messages"
)

// PMRHandler handles Progressive Muscle Relaxation exercise.
type PMRHandler struct {
	ctx               context.Context
	bot               *telego.Bot
	localizer         *localization.Localizer
	statistics        statistic.Stats
	sessionStorage    session.Storage
	sessionManager    *session.Manager
	callbackProcessor *CallbackProcessor
}

// NewPMRHandler creates a new PMR exercise handler.
func NewPMRHandler(
	ctx context.Context,
	bot *telego.Bot,
	localizer *localization.Localizer,
	statistics statistic.Stats,
	sessionStorage session.Storage,
	sessionManager *session.Manager,
	callbackProcessor *CallbackProcessor,
) *PMRHandler {
	return &PMRHandler{
		ctx:               ctx,
		bot:               bot,
		localizer:         localizer,
		statistics:        statistics,
		sessionStorage:    sessionStorage,
		sessionManager:    sessionManager,
		callbackProcessor: callbackProcessor,
	}
}

// getLang returns user's language from session or default.
func (h *PMRHandler) getLang(ctx context.Context, userID int64) string {
	return GetLang(ctx, h.sessionStorage, userID)
}

// HandleMenuSelect handles selection from main menu.
func (h *PMRHandler) HandleMenuSelect(_ *th.Context, cb telego.CallbackQuery) error {
	info := h.callbackProcessor.Extract(cb)
	if info == nil {
		return nil
	}

	h.statistics.IncreaseRequestsStatisticForUser(info.UserID, cb.From.Username, cb.From.IsPremium, cb.From.IsBot)
	h.showIntroRecreate(h.ctx, info.ChatID, info.UserID)
	return nil
}

// HandleCallback handles PMR exercise callbacks.
func (h *PMRHandler) HandleCallback(_ *th.Context, cb telego.CallbackQuery) error {
	info := h.callbackProcessor.Extract(cb)
	if info == nil {
		return nil
	}

	h.processCallback(info.ChatID, info.UserID, info.MessageID, info.Data)
	return nil
}

func (h *PMRHandler) processCallback(chatID, userID int64, messageID int, data string) {
	switch data {
	case "pmr_start":
		// Create new session context (cancels previous if any)
		sessionCtx := h.sessionManager.StartSession(h.ctx, userID)
		go h.startExercise(sessionCtx, chatID, userID, messageID)
	case "pmr_pause":
		if err := h.sessionStorage.SetState(h.ctx, userID, session.StatePMRPaused); err != nil {
			log.Printf("ERROR: set state pmr paused: %v", err)
		}
	case "pmr_resume":
		if err := h.sessionStorage.SetState(h.ctx, userID, session.StatePMRRunning); err != nil {
			log.Printf("ERROR: set state pmr running: %v", err)
		}
	case "pmr_stop":
		h.sessionManager.CancelSession(userID)
		h.stopExercise(h.ctx, chatID, userID, messageID)
	case "pmr_complete":
		h.sessionManager.CancelSession(userID)
		h.completeExercise(h.ctx, chatID, userID, messageID)
	case "pmr_cancel":
		h.sessionManager.CancelSession(userID)
		h.cancelExercise(h.ctx, chatID, userID, messageID)
	}
}

func (h *PMRHandler) showIntroRecreate(ctx context.Context, chatID, userID int64) {
	lang := h.getLang(ctx, userID)
	m := h.localizer.Get(lang)
	muscleGroups := techniques.GetMuscleGroups()
	text := fmt.Sprintf(m.PMRIntro, len(muscleGroups))
	ShowIntroRecreate(
		ctx, h.bot, h.sessionStorage, h.localizer, lang,
		chatID, userID, session.StatePMRActive, text,
		"pmr_cancel", "pmr_start", "pmr intro",
	)
}

func (h *PMRHandler) startExercise(ctx context.Context, chatID, userID int64, messageID int) {
	h.runMuscleGroup(ctx, chatID, userID, messageID, 0)
}

func (h *PMRHandler) runMuscleGroup(ctx context.Context, chatID, userID int64, messageID int, muscleIdx int) {
	muscleGroups := techniques.GetMuscleGroups()
	if muscleIdx < 0 || muscleIdx >= len(muscleGroups) {
		h.sendCompletion(ctx, chatID, userID, messageID)
		return
	}

	if err := h.sessionStorage.SetState(ctx, userID, session.StatePMRRunning); err != nil {
		log.Printf("ERROR: set state pmr running: %v", err)
	}

	muscle := muscleGroups[muscleIdx]
	totalGroups := len(muscleGroups)

	// Tense phase with progress
	if !h.runPhaseWithProgress(ctx, chatID, userID, messageID, muscle, totalGroups, true, muscle.TenseDuration) {
		return
	}

	// Relax phase with progress
	if !h.runPhaseWithProgress(ctx, chatID, userID, messageID, muscle, totalGroups, false, muscle.RelaxDuration) {
		return
	}

	h.runMuscleGroup(ctx, chatID, userID, messageID, muscleIdx+1)
}

func (h *PMRHandler) getLocalizedMuscle(muscleID string, m localization.Messages) (name, tense, relax string) {
	type muscleTexts struct {
		name  string
		tense string
		relax string
	}
	texts := map[string]muscleTexts{
		"hands":    {m.PMRHandsName, m.PMRHandsTense, m.PMRHandsRelax},
		"forearms": {m.PMRForearmsName, m.PMRForearmsTense, m.PMRForearmsRelax},
		"forehead": {m.PMRForeheadName, m.PMRForeheadTense, m.PMRForeheadRelax},
		"eyes":     {m.PMREyesName, m.PMREyesTense, m.PMREyesRelax},
		"jaw":      {m.PMRJawName, m.PMRJawTense, m.PMRJawRelax},
		"neck":     {m.PMRNeckName, m.PMRNeckTense, m.PMRNeckRelax},
		"chest":    {m.PMRChestName, m.PMRChestTense, m.PMRChestRelax},
		"stomach":  {m.PMRStomachName, m.PMRStomachTense, m.PMRStomachRelax},
		"thighs":   {m.PMRThighsName, m.PMRThighsTense, m.PMRThighsRelax},
		"calves":   {m.PMRCalvesName, m.PMRCalvesTense, m.PMRCalvesRelax},
	}
	if text, ok := texts[muscleID]; ok {
		return text.name, text.tense, text.relax
	}
	return muscleID, "", ""
}

func (h *PMRHandler) runPhaseWithProgress(
	ctx context.Context, chatID, userID int64, messageID int,
	muscle techniques.MuscleGroup, totalGroups int, isTense bool, duration time.Duration,
) bool {
	totalSeconds := int(duration.Seconds())
	m := h.localizer.Get(h.getLang(ctx, userID))

	keyboard := &telego.InlineKeyboardMarkup{
		InlineKeyboard: [][]telego.InlineKeyboardButton{
			{{Text: m.Pause, CallbackData: "pmr_pause"}},
		},
	}

	for elapsed := 0; elapsed < totalSeconds; elapsed++ {
		if !h.checkRunningOrPause(ctx, chatID, messageID, userID) {
			return false
		}

		progress := messages.TimerCountdown(elapsed, totalSeconds)
		var phaseText string
		name, tenseInstr, relaxInstr := h.getLocalizedMuscle(muscle.ID, m)
		if isTense {
			phaseText = fmt.Sprintf(m.PMRTense, tenseInstr) + "\n\n" + progress
		} else {
			phaseText = fmt.Sprintf(m.PMRRelax, relaxInstr) + "\n\n" + progress
		}

		text := fmt.Sprintf("%s *%s* (%d/%d)\n\n%s", muscle.Emoji, name, muscle.Number, totalGroups, phaseText)

		if _, err := h.bot.EditMessageText(ctx, &telego.EditMessageTextParams{
			ChatID:      tu.ID(chatID),
			MessageID:   messageID,
			Text:        text,
			ParseMode:   "Markdown",
			ReplyMarkup: keyboard,
		}); err != nil {
			if IsMessageNotFoundError(err) {
				return false
			}
			if !IsMessageNotModifiedError(err) {
				log.Printf("ERROR: edit pmr phase message: %v", err)
			}
		}

		timer := time.NewTimer(1 * time.Second)
		select {
		case <-ctx.Done():
			timer.Stop()
			return false
		case <-timer.C:
		}
	}
	return true
}

func (h *PMRHandler) checkRunningOrPause(
	ctx context.Context, chatID int64, messageID int, userID int64,
) bool {
	if ctx.Err() != nil {
		return false
	}
	state, err := h.sessionStorage.GetState(ctx, userID)
	if err != nil {
		return false
	}
	if state == session.StatePMRPaused {
		h.showPauseScreen(ctx, chatID, messageID, userID)
		return WaitForResume(
			ctx, h.sessionStorage, userID,
			session.StatePMRRunning, session.StatePMRPaused,
		)
	}
	return state == session.StatePMRRunning
}

func (h *PMRHandler) showPauseScreen(ctx context.Context, chatID int64, messageID int, userID int64) {
	ShowPauseScreen(ctx, h.bot, h.localizer, h.getLang(ctx, userID), chatID, messageID, "pmr_stop", "pmr_resume")
}

func (h *PMRHandler) sendCompletion(ctx context.Context, chatID, userID int64, messageID int) {
	m := h.localizer.Get(h.getLang(ctx, userID))
	SendCompletionScreen(
		ctx, h.bot, h.sessionStorage, chatID, userID, messageID,
		session.StatePMRRunning, m.PMRCompletion,
		"pmr_start", "pmr_complete", m, "send pmr completion",
	)
}

func (h *PMRHandler) stopExercise(ctx context.Context, chatID, userID int64, messageID int) {
	m := h.localizer.Get(h.getLang(ctx, userID))
	muscleGroups := techniques.GetMuscleGroups()
	text := fmt.Sprintf(m.PMRIntro, len(muscleGroups))
	SetStateAndEditIntro(
		ctx, h.bot, h.sessionStorage, chatID, userID, messageID,
		session.StatePMRActive, text,
		"pmr_cancel", "pmr_start", m, "pmr stop",
	)
}

func (h *PMRHandler) completeExercise(ctx context.Context, chatID, userID int64, messageID int) {
	ClearAndShowMainMenu(
		ctx, h.bot, h.sessionStorage, h.localizer, h.getLang(ctx, userID),
		chatID, userID, messageID, "edit pmr complete",
	)
}

func (h *PMRHandler) cancelExercise(ctx context.Context, chatID, userID int64, messageID int) {
	CancelAndShowMainMenu(
		ctx, h.bot, h.sessionStorage, h.localizer, h.getLang(ctx, userID),
		chatID, userID, messageID, "edit pmr cancel",
	)
}
