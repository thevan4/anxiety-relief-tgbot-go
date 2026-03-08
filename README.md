<div align="center">

<img src="assets/logo.jpg" width="160" alt="Anxiety Relief Bot" />

# Anxiety Relief Bot

*Anxiety support — in a Telegram bot*

https://t.me/anxietyhelp_bot

</div>

## Techniques

| | Technique | Duration | What it does |
|:---:|:---|:---:|:---|
| 🌬️ | **Simple Breathing** | 2 min | Calms the nervous system |
| 🌿 | **Grounding** | 3-5 min | Returns attention to the present moment |
| 🧘 | **Breathing Patterns** | 1-2 min | 4 patterns: box 4-4-4-4, relaxing 4-7-8, energizing 4-4-6, quick reset 3-3-3 |
| 💪 | **Muscle Relaxation** | 4 min | Deep relaxation of the entire body |

## Principles

**Free** — no subscriptions, no paid features, no limits.

**No ads** — no banners or integrations.

**No tracking** — no personal data. Only the current session state is stored (Redis, TTL 5 min). Nothing remains after the session ends.

**No notifications** — the bot never messages first and never sends reminders.

**Open source** — everything can be verified.

## Languages

🇬🇧 EN · 🇩🇪 DE · 🇫🇷 FR · 🇷🇺 RU · 🇧🇾 BE · 🇺🇦 UK

## Technical Info

### Project structure

```
internal/
├── localization/   # Translations (6 languages)
├── session/        # State machine — user state (Redis)
├── techniques/     # Exercise data: phases, timings, texts
├── statistic/      # Statistics
├── cleanup/        # Stale session cleanup
└── telegram/
    ├── handlers/   # Handlers: breathing, grounding, guided_breathing, pmr
    └── messages/   # Message formatting
```

### Run locally

```bash
# Set BOT_TOKEN in .env
go run ./cmd/main.go
```

### Linter

```bash
golangci-lint run --config .golangci.pipeline.yaml ./...
```

## Disclaimer

This bot is not a substitute for professional mental health support. If anxiety is affecting your life — please consult a specialist.

# About

Based on https://github.com/msgnoki/anxiety-aid-tools — a fork of https://anxietyaidtools.com. The original author made the repository private, abandoning the open-source principles the project was built on. 
