# Architecture: AI-Assisted Demand Form

## Visão Geral

O usuário descreve sua necessidade no chat flutuante. A IA faz perguntas até ter contexto suficiente, propõe um resumo para revisão e, após confirmação, preenche o formulário automaticamente.

```
┌─────────────────────────────────────────────────────────────┐
│                        FRONTEND (Vue 3)                      │
│                                                              │
│  ┌──────────────────┐         ┌──────────────────────────┐  │
│  │  FloatingChat    │         │      RequestForm          │  │
│  │                  │         │                           │  │
│  │  • Exibe msgs    │         │  • title                  │  │
│  │  • Envia msg     │         │  • businessLine           │  │
│  │  • Stream tokens │         │  • requesterBU            │  │
│  │  • Recebe        │         │  • busInterested          │  │
│  │    form_fill ────┼────────►│  • timeSensitive          │  │
│  │  • Recebe        │         │  • whyDemand              │  │
│  │    form_confirm  │         │  • whoIsImpacted          │  │
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
│  handlers/                                                   │
│  ┌──────────────────────────────────────────────────────┐   │
│  │  POST /api/chat/message       — coleta e streaming   │   │
│  │  POST /api/chat/confirm       — confirma proposta    │   │
│  │  POST /api/chat/field-message — preenche campo único │   │
│  │  POST /api/chat/analyze-form  — revisão de qualidade │   │
│  └──────────────────────────────────────────────────────┘   │
│                                                              │
│  services/llama.go              services/conversation.go     │
│  ┌────────────────────┐        ┌────────────────────────┐   │
│  │  HTTP client para  │        │  Gerencia sessões e    │   │
│  │  Groq API          │        │  histórico de msgs     │   │
│  │  (Llama 3.3 70B)   │        │  em memória            │   │
│  └────────────────────┘        └────────────────────────┘   │
│                                                              │
│  services/tools.go              services/prompts_main.go     │
│  ┌────────────────────┐        ┌────────────────────────┐   │
│  │  Schema das tools  │        │  System prompts e      │   │
│  │  (propose +        │        │  prompt de análise     │   │
│  │   fill_demand_form)│        └────────────────────────┘   │
│  └────────────────────┘                                      │
│                          services/prompts_field.go           │
│                          ┌────────────────────────┐         │
│                          │  Prompts por campo     │         │
│                          │  (fill_field tool)     │         │
│                          └────────────────────────┘         │
└─────────────────────────────────────────────────────────────┘
            │
            │ POST https://api.groq.com/openai/v1/chat/completions
            │ stream: true + tools: [propose_form_data, fill_demand_form]
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
{ "sessionId": "abc123", "message": "Preciso criar uma demanda de integração", "language": "pt" }
```

O campo `language` é opcional (padrão: `"en"`). Aceita: `pt`, `en`, `es`, `fr`.

### Eventos SSE emitidos pelo backend

| Tipo | Descrição | Payload |
|---|---|---|
| `token` | Trecho de texto da IA chegando em tempo real | `{ "type": "token", "content": "Qual é..." }` |
| `form_confirm` | IA coletou todos os campos e propõe um resumo para revisão | `{ "type": "form_confirm", "data": { "title": "...", ... } }` |
| `form_fill` | Formulário preenchido definitivamente (após confirmação) | `{ "type": "form_fill", "data": { "title": "...", ... } }` |
| `field_fill` | Campo único preenchido via `/chat/field-message` | `{ "type": "field_fill", "data": { "fieldName": "businessLine", "value": "18520" } }` |
| `done` | Fim do stream; `content` contém o sessionId (vazio em `/analyze-form`) | `{ "type": "done", "content": "abc123" }` |
| `error` | Erro no processamento | `{ "type": "error", "content": "mensagem" }` |

O `X-Session-ID` header na response devolve o ID da sessão (gerado automaticamente se não enviado).

---

## Lógica da IA: Três fases via Tool Calling

O modelo recebe um **system prompt** descrevendo os campos do formulário e suas opções válidas. Ele opera em três fases:

**Fase 1 — Coleta:** O modelo responde com texto (perguntas ao usuário). Cada token chega como evento `token`.

**Fase 2 — Proposta:** Quando o modelo considera que tem todos os campos obrigatórios, ele chama `propose_form_data`. O backend emite o evento `form_confirm` com os dados para o usuário revisar. O frontend exibe um resumo aguardando confirmação.

**Fase 3 — Preenchimento:** Após confirmação via `POST /api/chat/confirm`, o backend chama `fill_demand_form` e emite `form_fill`. Se o usuário solicitar correções após o preenchimento, o modelo pode chamar `fill_demand_form` diretamente com os valores corrigidos.

```
Tools disponíveis:

propose_form_data — chama quando todos os campos obrigatórios estão prontos
fill_demand_form  — chama após confirmação do usuário (ou para correções pós-fill)

Parâmetros obrigatórios (ambas as tools):
  - title               (string) — verbo de outcome + escopo de negócio
  - businessLine        (enum)   — IDs: 18518–18525, 19033
  - requesterBU         (enum)   — IDs ADEO-XXXX (32 opções)
  - busInterested       (array)  — IDs 20047–20078 (pode ser array vazio se nenhuma BU alinhada)
  - timeSensitive       (enum)   — "No" | "Legal" | "Security"
  - whyDemand           (string) — contexto, situação atual e pontos de dor
  - whoIsImpacted       (string) — personas e estimativa de usuários
  - benefitCategory     (enum)   — categoria do benefício esperado
  - benefitHypothesis   (string) — hipótese de como a demanda gera o benefício
  - measureBenefits     (string) — KPIs e timing para medir os benefícios

Parâmetros opcionais:
  - lowConfidenceFields (array)  — campos inferidos pelo modelo (não declarados pelo usuário)
```

---

## Endpoints adicionais

### `POST /api/chat/confirm`
Confirma a proposta exibida pelo `form_confirm`. O backend recupera os dados pendentes da sessão e emite `form_fill`.

```json
{ "sessionId": "abc123" }
```

### `POST /api/chat/field-message`
Preenche um único campo do formulário via conversa focada. Usa a tool `fill_field` com o schema do campo específico (preservando enums).

```json
{ "sessionId": "abc123", "fieldName": "businessLine", "message": "Trabalho na área de Supply Chain", "language": "pt" }
```

Emite SSE com eventos `token` e `field_fill` (contendo o campo e valor preenchido).

### `POST /api/chat/analyze-form`
Revisa a qualidade de um formulário já preenchido. Não usa tool calling — retorna tokens com feedback estruturado sobre cada campo. Não requer `sessionId` (sessão efêmera por request).

```json
{ "formData": { "title": "...", ... }, "language": "pt" }
```

---

## Gerenciamento de Sessão

Cada conversa tem um `sessionId` (UUID hex gerado no backend). O histórico de mensagens é mantido em memória no `ConversationService` e enviado ao modelo a cada request para manter contexto multi-turno.

```
sync.Map em memória — TTL de 2h, limpeza a cada 15min
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

Status possíveis: `"collecting"` | `"complete"`

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
│   │   ├── chat.go                 ← SendMessage (POST /api/chat/message)
│   │   ├── field.go                ← SendFieldMessage (POST /api/chat/field-message)
│   │   └── form.go                 ← ConfirmForm + AnalyzeForm
│   ├── services/
│   │   ├── llama.go                ← cliente Groq API com streaming
│   │   ├── conversation.go         ← gerencia sessões e histórico
│   │   ├── tools.go                ← schema das tools (propose + fill)
│   │   ├── prompts_main.go         ← system prompt principal e de análise
│   │   └── prompts_field.go        ← prompts por campo (fill_field)
│   └── models/
│       └── types.go                ← Message, Session, FormFillData
│
└── frontend/ai-chat/src/
    ├── services/
    │   └── chatService.ts          ← fetch + parse SSE stream
    ├── stores/
    │   └── chat.ts                 ← Pinia store (form_fill bridge)
    └── components/
        ├── FloatingChat.vue        ← conecta ao SSE, emite form_fill/form_confirm
        └── RequestForm.vue         ← observa store e auto-preenche
```

---

## Variáveis de Ambiente

```env
GROQ_API_KEY=gsk_...        # Chave da API Groq (https://console.groq.com)
PORT=8080                    # Porta do servidor (padrão: 8080)
ALLOWED_ORIGIN=*             # Origem permitida no CORS (padrão: *)
```
