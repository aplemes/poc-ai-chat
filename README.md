# Demand Assistant — AI-Powered Form Filler

**Live:** https://poc-ai-chat-1.onrender.com/

Internal Adeo tooling that lets users describe a demand in natural language. An AI assistant (Llama 3.3 70B via Groq) asks clarifying questions and, when it has enough context, automatically fills a structured demand request form.

---

## How it works

```
User types in chat  →  POST /api/chat/message (SSE)
                    →  Backend sends history to Groq
                    →  Groq streams tokens back
                    →  When ready, calls fill_demand_form tool
                    →  Frontend auto-fills RequestForm fields
                    →  Fields show "AI" badge (disappears on manual edit)
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
│   │   └── chat.go          # POST /api/chat/message — SSE handler
│   ├── services/
│   │   ├── llama.go         # Groq client, streaming, tool call accumulation
│   │   └── conversation.go  # In-memory sessions (2h TTL, cleanup every 15min)
│   ├── models/
│   │   └── types.go         # Message, Session, ToolCall, FormFillData
│   └── main.go              # Routes + CORS middleware
│
├── frontend/ai-chat/src/
│   ├── components/
│   │   ├── FloatingChat.vue  # Chat widget — SSE consumer, multi-language
│   │   └── RequestForm.vue   # Demand form — AI fill + badge UX
│   ├── services/
│   │   └── chatService.ts   # fetch + SSE stream parser
│   ├── stores/
│   │   └── chat.ts          # Pinia bridge: setFormFill / clearFormFill
│   └── assets/
│       └── tokens.css       # Design system tokens
│
├── .claude/agents/          # Claude Code subagents
└── docs/                    # Architecture, QA reports, requirements
```

---

## API

### `POST /api/chat/message`

**Request body**
```json
{
  "sessionId": "abc123" | null,
  "message": "I need to add a new payment method",
  "language": "en"
}
```

**Response** — SSE stream, `X-Session-ID` header

```
data: {"type":"token","content":"Sure, let me ask..."}
data: {"type":"token","content":" a few questions."}
data: {"type":"form_fill","data":{"title":"Add Payment Method X","businessLine":"18518",...}}
data: {"type":"done"}
data: {"type":"error","data":"something went wrong"}
```

### `GET /ping`
Health check — returns `{"message":"pong"}`.

---

## Form fields

| Field | Type | Required | Description |
|---|---|---|---|
| `title` | text | ✓ | Infinitive verb + scope. e.g. "Add Payment Method X on website" |
| `businessLine` | select | ✓ | One of 9 organisation IDs (18518–18525) |
| `requesterBU` | select | ✓ | One of 32 BU IDs (ADEO-XXXX format) |
| `busInterested` | multi-select | — | Other aligned BUs (IDs 20047–20078) |
| `timeSensitive` | toggle | — | `No` / `Legal` / `Security` |
| `whyDemand` | textarea | ✓ | Current situation and pain points |
| `whoIsImpacted` | textarea | ✓ | Personas and estimated user count |
| `benefitCategory` | select | ✓ | One of 6 benefit categories |
| `benefitHypothesis` | textarea | ✓ | How the benefit will be achieved |
| `measureBenefits` | textarea | ✓ | KPIs and timeframe |

---

## Languages

The chat widget supports **PT · EN · ES · FR**. Language preference is persisted in `localStorage`.

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
- The Submit button is not yet wired to a submission API
- No authentication or user identity
