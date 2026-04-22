# Architecture: AI-Assisted Demand Form

## Visão Geral

O usuário descreve sua necessidade no chat flutuante. A IA faz perguntas até ter contexto suficiente e, ao final, retorna os valores estruturados que preenchem o formulário automaticamente.

```
┌─────────────────────────────────────────────────────────────┐
│                        FRONTEND (Vue 3)                      │
│                                                              │
│  ┌──────────────────┐         ┌──────────────────────────┐  │
│  │  FloatingChat    │         │      RequestForm          │  │
│  │                  │         │                           │  │
│  │  • Exibe msgs    │         │  • title                  │  │
│  │  • Envia msg     │         │  • demandScope            │  │
│  │  • Stream tokens │         │  • businessLine           │  │
│  │  • Recebe        │         │  • requesterBU            │  │
│  │    form_fill ────┼────────►│  • busInterested          │  │
│  │                  │         │  • timeSensitive          │  │
│  │                  │         │  • whyDemand              │  │
│  │                  │         │  • whoIsImpacted          │  │
│  │                  │         │  • benefitCategory        │  │
│  │                  │         │  • benefitHypothesis      │  │
│  │                  │         │  • measureBenefits        │  │
│  └────────┬─────────┘         └──────────────────────────┘  │
│           │        chatStore (Pinia)                         │
└───────────┼─────────────────────────────────────────────────┘
            │ POST /api/chat/message
            │ ◄── SSE stream (text/event-stream)
            │
┌───────────▼─────────────────────────────────────────────────┐
│                     BACKEND (Go / Gin)                       │
│                                                              │
│  handlers/chat.go                                            │
│  ┌──────────────────────────────────────────────────────┐   │
│  │  POST /api/chat/message                               │   │
│  │  1. Adiciona msg ao histórico da sessão               │   │
│  │  2. Chama LlamaService.StreamChat(...)                │   │
│  │  3. Faz stream dos tokens pro cliente via SSE         │   │
│  │  4. Se tool_call "fill_demand_form" → emite form_fill │   │
│  └──────────────────────────────────────────────────────┘   │
│                                                              │
│  services/llama.go              services/conversation.go     │
│  ┌────────────────────┐        ┌────────────────────────┐   │
│  │  HTTP client para  │        │  Gerencia sessões e    │   │
│  │  Groq API          │        │  histórico de msgs     │   │
│  │  (Llama 3.3 70B)   │        │  em memória (→ MongoDB)│   │
│  └────────────────────┘        └────────────────────────┘   │
│                                                              │
└─────────────────────────────────────────────────────────────┘
            │
            │ POST https://api.groq.com/openai/v1/chat/completions
            │ stream: true + tools: [fill_demand_form]
            ▼
┌─────────────────────────────────────────────────────────────┐
│                    Groq API (Llama 3.3 70B)                  │
└─────────────────────────────────────────────────────────────┘
```

---

## Protocolo SSE (Frontend ↔ Backend)

O frontend faz um `POST /api/chat/message` e lê a resposta como stream SSE.

### Request
```json
POST /api/chat/message
{ "sessionId": "abc123", "message": "Preciso criar uma demanda de integração" }
```

### Eventos SSE emitidos pelo backend

| Tipo | Descrição | Payload |
|---|---|---|
| `token` | Trecho de texto da IA chegando em tempo real | `{ "type": "token", "content": "Qual é..." }` |
| `form_fill` | IA coletou dados suficientes, preencher formulário | `{ "type": "form_fill", "data": { "title": "...", "businessLine": "18525", ... } }` |
| `done` | Fim do stream, inclui sessionId | `{ "type": "done", "content": "abc123" }` |
| `error` | Erro no processamento | `{ "type": "error", "content": "mensagem" }` |

O `X-Session-ID` header na response devolve o ID da sessão (gerado automaticamente se não enviado).

---

## Lógica da IA: Duas fases via Tool Calling

O modelo recebe um **system prompt** descrevendo os campos do formulário e suas opções válidas. Ele opera em dois modos:

**Fase 1 — Coleta:** O modelo responde com texto (perguntas ao usuário). Cada token chega como evento `token`.

**Fase 2 — Preenchimento:** Quando o modelo considera que tem contexto suficiente, ele chama a tool `fill_demand_form` em vez de responder com texto. O backend detecta o `tool_call`, parseia os argumentos JSON, e emite o evento `form_fill`.

```
Tool: fill_demand_form
Parâmetros obrigatórios:
  - title             (string) — verbo infinitivo + escopo
  - demandScope       (enum)   — "Intra-BU" | "Adeo Platform"
  - businessLine      (enum)   — IDs 18518–18525
  - requesterBU       (enum)   — IDs ADEO-XXXX (32 opções)
  - whyDemand         (string) — situação atual, dores e contexto motivador
  - benefitCategory   (enum)   — categoria do benefício esperado
  - benefitHypothesis (string) — como a demanda vai gerar o benefício
Parâmetros opcionais:
  - timeSensitive     (enum)   — "Yes" | "No"
  - busInterested     (array)  — IDs 20047–20078
  - whoIsImpacted     (string) — personas e estimativa de usuários
  - measureBenefits   (string) — métricas e timing
```

---

## Gerenciamento de Sessão

Cada conversa tem um `sessionId` (UUID hex gerado no backend). O histórico de mensagens é mantido em memória no `ConversationService` e enviado ao modelo a cada request para manter contexto multi-turno.

```
MongoDB (futuro) ← hoje: sync.Map em memória
```

Estrutura de uma sessão:
```json
{
  "id": "a3f8c2...",
  "status": "collecting",
  "messages": [
    { "role": "system", "content": "..." },
    { "role": "user",   "content": "Preciso criar uma demanda..." },
    { "role": "assistant", "content": "Qual BU está solicitando?" }
  ]
}
```

---

## Estrutura de Arquivos

```
chat-ai/
├── docs/
│   └── ARCHITECTURE.md             ← este arquivo
│
├── backend/gin-quickstart/
│   ├── main.go                     ← rotas, CORS, inicia servidor
│   ├── .env.example
│   ├── handlers/
│   │   └── chat.go                 ← POST /api/chat/message
│   ├── services/
│   │   ├── llama.go                ← cliente Groq API com streaming
│   │   └── conversation.go         ← gerencia sessões e histórico
│   └── models/
│       └── types.go                ← Message, Session, FormFillData
│
└── frontend/ai-chat/src/
    ├── services/
    │   └── chatService.ts          ← fetch + parse SSE stream
    ├── stores/
    │   └── chat.ts                 ← Pinia store (form_fill bridge)
    └── components/
        ├── FloatingChat.vue        ← conecta ao SSE, emite form_fill
        └── RequestForm.vue         ← observa store e auto-preenche
```

---

## Variáveis de Ambiente

```env
GROQ_API_KEY=gsk_...     # Chave da API Groq (https://console.groq.com)
PORT=8080                 # Porta do servidor (padrão: 8080)
```
