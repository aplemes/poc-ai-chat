# Demand Assistant — AI-Powered Form Filler

**Live:** https://poc-ai-chat-1.onrender.com/

Internal Adeo tooling that lets users describe a demand in natural language. An AI assistant (Llama 3.3 70B via Groq) asks clarifying questions and, when it has enough context, auto-fills a structured demand request form and reviews it before submission.

---

## How it works

```
User types in FloatingChat
  → POST /api/chat/message (SSE)
  → Backend sends full history to Groq
  → Groq streams tokens back as {type:"token"} events
  → When AI has all fields → calls propose_form_data tool
  → Frontend shows ConfirmCard for user review
  → User confirms → POST /api/chat/confirm
  → Backend calls fill_demand_form → emits {type:"form_fill"}
  → Frontend auto-fills RequestForm fields with "AI" badge
  → User clicks Submit → POST /api/chat/analyze-form (SSE)
  → AI reviews form quality → AnalysisModal shows structured feedback
```

Per-field AI assistant (AI button next to each field):
```
User clicks AI button on a field
  → POST /api/chat/field-message (SSE)
  → Dedicated field-scoped AI fills that single field
```

---

## Stack

| Layer | Technology |
|---|---|
| Frontend | Vue 3 + TypeScript, Pinia, Vite |
| Backend | Go 1.22+, Gin |
| AI | Llama 3.3 70B via [Groq API](https://groq.com) |
| Streaming | Server-Sent Events (SSE) |
| Testing | Vitest + Playwright |

---

## Prerequisites

- **Node.js** `^20.19.0` or `>=22.12.0`
- **Go** `1.22+`
- **Groq API key** — get one at [console.groq.com](https://console.groq.com)

---

## Getting started

### 1. Backend

```bash
cd backend/gin-quickstart

# Option A — .env file
echo "GROQ_API_KEY=gsk_..." > .env

# Option B — environment variable
export GROQ_API_KEY=gsk_...

go run main.go
# Starts on :8080
```

### 2. Frontend

```bash
cd frontend/ai-chat

npm install
npm run dev
# Starts on :5173
```

Open [http://localhost:5173](http://localhost:5173) and click the chat button in the bottom-right corner.

---

## Project structure

```
chat-ai/
├── backend/gin-quickstart/
│   ├── handlers/
│   │   ├── chat.go        # POST /api/chat/message, /confirm — SSE stream handler
│   │   ├── field.go       # POST /api/chat/field-message — per-field AI handler
│   │   └── form.go        # POST /api/chat/analyze-form — quality review handler
│   ├── services/
│   │   ├── llama.go       # Groq client, streaming, tool call accumulation
│   │   ├── conversation.go# In-memory sessions (2h TTL, cleanup every 15min)
│   │   ├── prompts_main.go# System prompts for chat AI and analysis AI
│   │   ├── prompts_field.go# System prompt for per-field AI
│   │   └── tools.go       # Tool schemas (propose_form_data, fill_demand_form)
│   ├── models/
│   │   └── types.go       # Message, Session, ToolCall, FormFillData
│   └── main.go            # Routes + CORS middleware
│
├── frontend/ai-chat/src/
│   ├── components/
│   │   ├── FloatingChat.vue       # Global chat widget — SSE consumer, multi-language
│   │   ├── RequestForm.vue        # Demand form — AI fill + badge UX + submit flow
│   │   ├── AnalysisModal.vue      # Pre-submission AI quality review dialog
│   │   ├── AnalysisModal.css      # Styles for AnalysisModal
│   │   ├── FieldChatPanel.vue     # Per-field AI assistant panel
│   │   ├── chat/
│   │   │   ├── ChatFab.vue        # Floating action button
│   │   │   ├── ChatMessageList.vue# Message rendering + markdown
│   │   │   ├── ConfirmCard.vue    # Propose/confirm card in chat
│   │   │   └── FloatingChat.css   # Styles for FloatingChat
│   │   └── form/
│   │       ├── FieldLabel.vue     # Field label + AI badge
│   │       ├── SectionIdentity.vue# Title, BU, businessLine fields
│   │       ├── SectionContext.vue # Why demand, who is impacted fields
│   │       └── SectionBenefits.vue# Benefit category, hypothesis, KPIs fields
│   ├── composables/
│   │   ├── useAiFormFill.ts       # AI-fill state and badge logic
│   │   ├── useChatStream.ts       # Generic SSE stream hook
│   │   ├── useFormAnalysis.ts     # Pre-submission analysis flow
│   │   ├── useFormContext.ts      # Shared form state via provide/inject
│   │   └── useLanguage.ts         # Language persistence (localStorage)
│   ├── services/
│   │   ├── chatService.ts         # fetch + SSE parser for chat/analysis endpoints
│   │   ├── fieldChatService.ts    # fetch + SSE parser for field-message endpoint
│   │   └── config.ts              # Base URL and shared fetch config
│   ├── stores/
│   │   ├── chat.ts                # Pinia bridge: setFormFill / clearFormFill
│   │   └── fieldChat.ts           # Pinia store for per-field AI state
│   └── assets/
│       └── tokens.css             # Design system tokens
│
├── .claude/agents/                # Claude Code subagents
└── docs/                          # Architecture and feature docs
```

---

## API

### `POST /api/chat/message`

Sends a user message and receives an SSE stream.

**Request**
```json
{ "sessionId": "abc123" | null, "message": "I need to reduce supply chain costs", "language": "en" }
```

**Response** — SSE stream, `X-Session-ID` header
```
data: {"type":"token","content":"Sure, let me ask..."}
data: {"type":"form_confirm","data":{"title":"Reduce Supply Chain Costs",...}}
data: {"type":"form_fill","data":{"title":"Reduce Supply Chain Costs","businessLine":"18520",...}}
data: {"type":"done"}
data: {"type":"error","content":"something went wrong"}
```

| Event | When |
|---|---|
| `token` | AI is streaming a reply |
| `form_confirm` | AI proposes filled fields for user review (ConfirmCard) |
| `form_fill` | User confirmed — fields applied to the form |
| `done` | Stream finished |
| `error` | API or processing error |

### `POST /api/chat/confirm`

Confirms the proposed form data and triggers `fill_demand_form`.

**Request**
```json
{ "sessionId": "abc123" }
```

**Response** — SSE stream (same event types as above)

### `POST /api/chat/field-message`

Per-field AI assistant — fills a single form field.

**Request**
```json
{ "sessionId": "abc123" | null, "fieldName": "title", "message": "user input", "language": "en" }
```

**Response** — SSE stream

### `POST /api/chat/analyze-form`

Streams a quality review of the filled form before submission.

**Request**
```json
{ "formData": { "title": "...", "whyDemand": "...", ... }, "language": "en" }
```

**Response** — SSE stream of `token` events (markdown feedback)

### `GET /ping`

Health check — returns `{"message":"pong"}`.

---

## Form fields

| Field | Type | Required | Description |
|---|---|---|---|
| `title` | text | ✓ | Business need/outcome. Outcome verb + scope. e.g. "Reduce costs in the supply chain" |
| `businessLine` | select | ✓ | One of 9 organisation IDs (18518–18525) |
| `requesterBU` | select | ✓ | One of 32 BU IDs (ADEO-XXXX format) |
| `busInterested` | multi-select | — | Other aligned BUs (IDs 20047–20078) |
| `timeSensitive` | toggle | ✓ | `No` / `Legal` / `Security` |
| `whyDemand` | textarea | ✓ | Triggering event + current situation/tools + specific pain points |
| `whoIsImpacted` | textarea | ✓ | Personas and estimated user count |
| `benefitCategory` | select | ✓ | One of 6 benefit categories |
| `benefitHypothesis` | textarea | ✓ | Causal chain: "We believe X because Y" |
| `measureBenefits` | textarea | ✓ | Specific KPIs and measurement timeframe |

All field **values** are always stored and submitted in **English**, regardless of the chat language.

---

## Languages

The chat widget supports **PT · EN · ES · FR**. Language preference is persisted in `localStorage`. Feedback and conversation are in the selected language; form values are always in English.

---

## Development commands

### Frontend (`frontend/ai-chat/`)

```bash
npm run dev          # Dev server with HMR (port 5173)
npm run build        # Type-check + production build
npm run test:unit    # Vitest unit tests
npm run test:e2e     # Playwright end-to-end tests
npm run lint         # oxlint + eslint with auto-fix
npm run format       # Prettier auto-format
npm run type-check   # vue-tsc incremental check
```

### Backend (`backend/gin-quickstart/`)

```bash
go run main.go   # Start server
go test ./...    # Run all tests
```

---

## Environment variables

| Variable | Required | Default | Description |
|---|---|---|---|
| `GROQ_API_KEY` | ✓ | — | Groq API key |
| `PORT` | — | `8080` | Backend port |

---

## Design system

All visual tokens (colors, spacing, typography, shadows) are defined in `frontend/ai-chat/src/assets/tokens.css`. Primary color: Adeo green (`#00874A`).

---

## AI agents (Claude Code)

This project uses specialized Claude Code subagents for development. They live in `.claude/agents/`:

| Agent | Specialization |
|---|---|
| `go-backend-expert` | Go, Gin, goroutines, SSE server |
| `vue-frontend-expert` | Vue 3, TypeScript, Pinia, SSE consumer |
| `product-manager` | PRDs, user stories, acceptance criteria |
| `qa-specialist` | Vitest, Playwright, bug mapping |
| `ux-ui-designer` | Design system, UX review, accessibility |
| `requirements-analyst` | Idea → requirements, gap analysis, impact map |

---

## Known limitations

- Sessions are **in-memory only** — restarting the backend loses all conversation history
- The Submit button triggers an AI quality review but is not wired to a submission API
- No authentication or user identity
