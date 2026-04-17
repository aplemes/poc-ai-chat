# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Full-stack AI assistant for internal Adeo tooling. The user describes their demand in a floating chat widget; the backend queries Llama 3.3 70B (via Groq) with tool calling, asks clarifying questions, and when it has enough context calls `fill_demand_form` — which auto-fills the `RequestForm` fields via a Pinia store bridge.

## Commands

Frontend — run from `frontend/ai-chat/`:

```bash
npm run dev           # Start dev server with HMR (port 5173)
npm run build         # Type-check + production build
npm run test:unit     # Vitest unit tests
npm run test:e2e      # Playwright end-to-end tests
npm run lint          # oxlint + eslint with auto-fix
npm run format        # Prettier auto-format
npm run type-check    # vue-tsc incremental check
```

Backend — run from `backend/gin-quickstart/`:

```bash
export GROQ_API_KEY=gsk_...
go run main.go        # Starts on :8080
go test ./...
```

## Architecture

See `docs/ARCHITECTURE.md` for the full system diagram and SSE event protocol.

### Chat flow

1. User types in `FloatingChat.vue` → `POST /api/chat/message` (SSE response)
2. Backend appends message to in-memory session, sends full history to Groq
3. Groq streams tokens back → backend forwards as `{type:"token"}` SSE events
4. When model has enough context it calls tool `fill_demand_form` → backend emits `{type:"form_fill", data:{...}}`
5. `FloatingChat` receives `form_fill` → writes to `useChatStore`
6. `RequestForm` watches the store and auto-fills its fields; filled fields get an "IA" badge that disappears on manual edit

### Form fields

| Field | Type | Required | Description |
|---|---|---|---|
| `title` | text | ✓ | Infinitive verb + scope. e.g. "Add Payment Method X on website" |
| `demandScope` | select | ✓ | `Intra-BU` or `Adeo Platform` |
| `businessLine` | select | ✓ | One of 9 organisation IDs (18518–18525) |
| `requesterBU` | select | ✓ | One of 32 BU IDs (ADEO-XXXX format) |
| `busInterested` | multi-select | — | Other aligned BUs (IDs 20047–20078) |
| `demandContext` | textarea | ✓ | Event/situation motivating the demand |
| `currentSituation` | textarea | ✓ | Current processes and tools in use |
| `problemsToSolve` | textarea | ✓ | Pain points and inefficiencies |
| `whoIsImpacted` | textarea | — | Personas and estimated user counts |
| `measureBenefits` | textarea | — | KPIs and timing to measure success |

### Backend (`backend/gin-quickstart/`)

```
handlers/chat.go          POST /api/chat/message — SSE stream handler
services/llama.go         Groq API client (streaming + tool call accumulation)
services/conversation.go  In-memory sessions with 2h TTL, cleanup every 15min
models/types.go           Message, Session, ToolCall, FormFillData
main.go                   Routes + CORS middleware
```

- Sessions are keyed by a hex UUID returned in `X-Session-ID` response header
- Tool schema enforces enums for `businessLine` (9 IDs) and `busInterested` (20047–20078)
- HTTP client has 90s timeout; API errors include the response body

### Frontend (`frontend/ai-chat/src/`)

```
components/FloatingChat.vue   Chat widget — SSE consumer, persists sessionId in localStorage
components/RequestForm.vue    Demand form — watches useChatStore, applies AI fill + badge UX
services/chatService.ts       fetch + SSE stream parser, accepts AbortSignal
stores/chat.ts                Pinia bridge: setFormFill / clearFormFill
```

- `sessionId` is stored in `localStorage` so conversation survives page navigation
- `AbortController` cancels the in-flight stream if a new message is sent

## Debugging

- If the backend logs `WARNING: GROQ_API_KEY is not set`, set the env var before starting
- Groq API errors include the full response body in the error message
- Session history is in-memory only — restart loses all sessions

## Code Style

- **Frontend:** Prettier (no semicolons, single quotes, 100-char line width), Oxlint + ESLint, `noUncheckedIndexedAccess: true`, path alias `@/*`, `<script setup lang="ts">`
- **Backend:** Standard Go formatting (`gofmt`), errors always handled and propagated
- Before editing frontend files run `npm run format` after changes
