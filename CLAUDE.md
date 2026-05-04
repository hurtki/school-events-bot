# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Build
go build ./cmd/bot/

# Run all tests
go test ./...

# Run a single test
go test ./internal/parser/ -run TestName

# Run with race detector
go test -race ./...

# Vet
go vet ./cmd/bot/... ./internal/...

# Docker (production)
docker compose up -d
```

`cmd/dev/main.go` is a scratch file with a stub `main1` function — it intentionally doesn't compile as a standalone binary.

## Environment variables

Copy `.env.example` to `.env`. Required vars: `SCHEDULE_DOCUMENT_ID`, `TELEGRAM_BOT_TOKEN`, `UPDATES_TELEGRAM_CHANNEL_ID`, `SCHEDULE_FILE_REPOSIOTORY_PATH`, `PINNED_MESSAGE_STATE_FILE_REPOSIOTORY_PATH`, `SCHEDULE_WORKER_INTERVAL`, `UPCOMING_EVENTS_WORKER_INTERVAL`.

Optional AI vars (if absent, update notifications fall back to a basic template): `GEMINI_API_KEY`, `GEMINI_MODEL` (e.g. `gemini-2.5-flash`).

## Architecture

The bot serves an Israeli school. The schedule source is a Google Spreadsheets XLSX document where each sheet corresponds to a grade group (Hebrew sheet names map to `domain.Group` constants in `parser/parser.go`).

### Two main runtime loops

**ScheduleWorker** (polling interval from env) — fetches the XLSX, parses it, diffs against the persisted schedule, saves the new schedule, then publishes a `domain.ScheduleUpdate` to the event bus if anything changed.

**UpcomingEventsWorker** (polling interval from env) — independently refreshes the pinned Telegram message summarising the next N Bagrut/Magen events per group. The pinned message is re-sent from scratch if it is older than ~47.5 h; otherwise it is edited in place.

### Event bus

`evbus.ScheduleUpdateEventBus` wraps `asaskevich/EventBus` and adds synchronous fan-out: `Publish` blocks until every subscriber's goroutine calls `wg.Done()`. Two subscribers are wired in `main.go`:
- `updates.BotScheduleUpdatesService` → sends a change notification to the updates channel.
- `pinned.BotUpcomingEventsPinService` → refreshes the pinned message.

### Schedule diffing

`domain.NewScheduleUpdate(old, current)` builds `Added` / `Deleted` slices by comparing `Event.Hash()` values. `Hash()` concatenates `date + type + group + text`, so any text change (e.g. a room number) produces a delete+add pair. The AI summary layer in `bot.NotifyAboutUpdate` is responsible for grouping such pairs back into human-readable "modifications".

### AI summary (`internal/bot/update_notify.go`)

When `GEMINI_API_KEY` and `GEMINI_MODEL` are set, `generateAISummary` serialises the added/deleted events as JSON and calls Gemini with a prompt that instructs it to detect modifications (same date+type+group, different text) and produce a Telegram HTML message in Hebrew. `buildFallbackMessage` is used if AI is unavailable or returns an error.

### Dependency injection

All wiring happens in `cmd/bot/main.go`. Application-layer services (`internal/app/`) depend on interfaces, not concrete types. Infrastructure implementations (`internal/infrastructure/`) and `internal/bot/` are injected at the composition root. There is no framework — everything is manual.

### Persistence

Both repositories (`repository/schedule`, `repository/pinned_message`) are JSON-file-backed. The file paths are supplied via env vars and mounted as a Docker volume (`schedule-storage`) so state survives container restarts.

## Agent skills

### Issue tracker

Not used for this project. See `docs/agents/issue-tracker.md`.

### Triage labels

Default canonical label strings. See `docs/agents/triage-labels.md`.

### Domain docs

Single-context layout — `CONTEXT.md` + `docs/adr/` at repo root. See `docs/agents/domain.md`.
