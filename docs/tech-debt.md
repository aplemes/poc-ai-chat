# Levantamento de Débitos Técnicos

Auditoria realizada em 2026-04-22 por 4 agentes especializados (backend, frontend, QA, UX/acessibilidade).

---

## Legenda de status

| Status | Significado |
|---|---|
| ✅ Resolvido | Corrigido e commitado |
| 🔧 Pendente | Ainda não implementado |
| ⏸ Adiado | Reconhecido, sem prioridade imediata |

---

## Índice

- [Severidade Alta](#severidade-alta)
- [Severidade Média](#severidade-média)
- [Severidade Baixa](#severidade-baixa)

---

## Severidade Alta

### Backend

| ID | Status | Arquivo | Linha | Problema |
|---|---|---|---|---|
| BE-01 | ✅ | `services/conversation.go` | — | `SetPendingFormData`, `TakeAndClearPendingFormData` adicionados. Handlers não mutam mais `session.PendingFormData` diretamente. |
| BE-02 | ✅ | `handlers/form.go` | — | `addToolCallMessages` helper extraído para chat.go. `TakeAndClearPendingFormData` atomic no service. Lógica do handler simplificada. |
| BE-03 | ✅ | `handlers/chat.go` | — | `writeEvent` agora trata o erro de `json.Marshal` e loga via `log.Printf`. `argsBytes` em form.go também tratado. |
| BE-04 | ✅ | `services/conversation.go` | — | `rand.Read` verificado; falha causa `panic` imediato com mensagem descritiva. |
| BE-05 | ✅ | `services/conversation.go` | — | `GetByID` removido. `GetMessages` retorna cópia do slice. `TakeAndClearPendingFormData` acessa dados sob mutex. `sync.RWMutex` → `sync.Mutex` (escritas dominam). |
| BE-06 | ✅ | `handlers/chat.go` | — | `setSSEHeaders(c *gin.Context, sessionID string)` centraliza os 4 headers em todos os handlers. |
| BE-07 | ✅ | `services/conversation_test.go` | — | 15 testes cobrindo GetOrCreate, SessionExists, AddMessage, GetMessages (cópia), PendingFormData (set/take/clear), evictExpired e generateID. |

### Frontend

| ID | Status | Arquivo | Linha | Problema |
|---|---|---|---|---|
| FE-01 | ✅ | `composables/useChatStream.ts` | — | SSE stream extraído para composable compartilhado. `FloatingChat` e `FieldChatPanel` usam `useChatStream()`. |
| FE-02 | ✅ | `utils/markdown.ts` | — | `renderMarkdown` (DOMPurify + marked) extraído para utilitário compartilhado. |
| FE-03 | ✅ | `composables/useLanguage.ts` | — | Estado de linguagem extraído para composable com `localStorage` sincronizado. |
| FE-04 | ✅ | `FloatingChat.vue` | — | 249 linhas (era 1116). Responsabilidades extraídas: `ChatFab.vue`, `ChatMessageList.vue`, `ConfirmCard.vue`, `useChatStream.ts`, `FloatingChat.css`. |
| FE-05 | ✅ | `FieldChatPanel.vue` | — | 299 linhas (era 722). i18n extraído para `fieldGreetings.ts`, CSS para `FieldChatPanel.css`. |
| FE-06 | ✅ | `useAiFormFill.ts` | — | Tipagem corrigida via `[StringField, string \| undefined][]`. Sem mais double-cast `as unknown as`. |
| FE-07 | ✅ | `useFormAnalysis.ts` | 19 | `analysisAbort.value?.abort()` chamado antes de criar novo `AbortController`. |
| FE-08 | ✅ | `SectionIdentity.vue` | 28 | Handler inline extraído como função `handleBusInterestedSelect`. |
| FE-09 | ✅ | `useAiFormFill.ts` | — | `aiFilledFields` não é mais mutado diretamente por componentes externos. Toda mutação ocorre dentro do composable. |

### UX / Acessibilidade

| ID | Status | Arquivo | Linha | Problema |
|---|---|---|---|---|
| UX-01 | ✅ | `AnalysisModal.vue` | — | Focus trap implementado: Tab/Shift+Tab aprisionados, Escape fecha, foco restaurado ao fechar. |
| UX-02 | ✅ | `FieldChatPanel.vue` | 42–115 | Focus trap implementado no drawer com `trapFocus` e restauração de foco. |
| UX-03 | ✅ | `SectionIdentity.vue` | múltiplos | Botões ✦ têm `aria-label` descritivo em cada campo ("Open AI assistant for Title", etc.). |
| UX-04 | ✅ | `SectionIdentity.vue` | 304–305 | `timeSensitive` segmented control tem `role="group"` e `aria-labelledby`. |
| UX-05 | ✅ | `RequestForm.vue` | — | Live region `aria-live="polite"` anuncia quantos campos foram preenchidos pela IA. |
| UX-06 | ✅ | `SectionIdentity.vue` | 254 | `<select>` de BUs Interessadas tem `aria-labelledby="busInterested-label"`. |
| UX-07 | ✅ | `tokens.css` | — | Todos os tokens referenciados (`--radius-2xl`, `--radius-xl`, `--radius-xs`, `--shadow-2xl`, `--shadow-primary`) declarados no design system. |
| UX-08 | ✅ | `FloatingChat.vue` | — | `aria-label` do FAB dinâmico: usa `ariaOpenChatNew` quando `hasNewMessage=true`. |
| UX-09 | ✅ | `AnalysisModal.vue` | 129–140 | Estado de erro exibe apenas "Return to form" — "Submit anyway" oculto em erro. |
| UX-10 | ✅ | `AnalysisModal.vue` | — | Botão renomeado de "Fix issues" para "Return to form". |
| UX-11 | ✅ | `useFormAnalysis.ts` | 52–56 | `confirmSubmit()` chama `onConfirmedSubmit()` que executa `submitForm()` real. |
| UX-12 | ✅ | `ConfirmCard.vue` | — | Chaves camelCase substituídas por labels legíveis via `FIELD_LABELS` map. Loop extraído para `visibleFields` computed. |
| UX-13 | ✅ | `AnalysisModal.vue` | — | `aria-modal="true"` adicionado e foco inicial atribuído ao abrir via `panelRef.focus()`. |
| UX-14 | ✅ | `tokens.css` | 24, 97 | `--color-neutral-500` e `--shadow-primary` declarados. |
| UX-15 | ✅ | `FloatingChat.css`, `ChatFab.vue` | — | `bounce`, `pulse-dot` e `badge-pop` cobertos pelo bloco `prefers-reduced-motion`. |

### QA / Testes

| ID | Status | Área | Problema |
|---|---|---|---|
| QA-01 | 🔧 | Backend inteiro | Zero testes. `services/conversation.go`, `services/llama.go`, todos os handlers e services sem `_test.go`. |
| QA-02 | 🔧 | `src/utils/sse.ts` | `readSSEStream` nunca testado em isolamento. Primitive central de todo o parsing SSE. |
| QA-03 | 🔧 | `chatService.ts` | `analyzeForm` e `confirmForm` sem nenhum teste. |
| QA-04 | 🔧 | `fieldChatService.ts` | Zero testes. Evento `field_fill` nunca exercitado em teste. |
| QA-05 | 🔧 | `useAiFormFill.ts` | Lógica de badge, `badgeClass`, coerção de `busInterested` para array, e watcher de `pendingFieldFill` sem cobertura. |
| QA-06 | 🔧 | `useFormAnalysis.ts` | `handleSubmit`, `closeAnalysis` (abort), `renderMd` (XSS sanitisation) sem testes. |
| QA-07 | 🔧 | `stores/fieldChat.ts` | Todas as actions sem testes: `openPanel`, `closePanel`, `setSessionId`, `getSessionId`, `setFieldFill`. |
| QA-08 | 🔧 | E2E — fluxo confirm | Nenhum E2E cobre `form_confirm` → confirm card → clique em "Confirm & fill" → `form_fill` → form preenchido. |
| QA-09 | 🔧 | E2E — FieldChatPanel | Nenhum E2E abre painel de campo, envia mensagem, recebe `field_fill` e verifica o campo no form. |
| QA-10 | 🔧 | E2E — AnalysisModal | Nenhum E2E clica em Submit, verifica modal, aguarda markdown, confirma ou cancela. |
| QA-11 | 🔧 | E2E — `lowConfidenceFields` | Badge `ai-badge--uncertain` nunca exercitado em E2E. |
| QA-12 | 🔧 | `chat-form-flow.spec.ts` | **Asserção tautológica** em testes de AbortError: `expect(typeof hasCancelled).toBe('boolean')` — sempre passa, nunca valida cancelamento real. |
| QA-13 | 🔧 | `chat-form-flow.spec.ts` | **Sem Page Object Model**. ~80 testes com seletores CSS literais espalhados. Qualquer rename de classe quebra o arquivo inteiro. |
| QA-14 | 🔧 | `chat-store.spec.ts` | `formFilled` nunca testado. `clearFormFill` não reseta a flag — bug confirmado, sem teste que o capture. |
| QA-15 | 🔧 | Backend — `services/llama.go` | Acumulação de tool call arguments em múltiplos chunks SSE sem nenhum teste. |
| QA-16 | 🔧 | Backend — `handlers/form.go` | `ConfirmForm` com sessão não encontrada (404) e sem dados pendentes (400) sem testes. |
| QA-17 | 🔧 | Backend — `handlers/field.go` | `parseFieldValue` para `busInterested` (string → `[]string`) e campo escalar sem testes. |
| QA-18 | 🔧 | Backend — `services/conversation.go` | `evictExpired` e `AddMessage` com sessão inexistente sem cobertura. Race condition em `GetOrCreate` concorrente. |

---

## Severidade Média

### Backend

| ID | Status | Arquivo | Problema |
|---|---|---|---|
| BE-M1 | ✅ | `services/prompts_main.go`, `services/prompts_field.go` | `prompts.go` (280 linhas) dividido em dois arquivos por responsabilidade. |
| BE-M2 | ✅ | `handlers/chat.go` | `buildMessages` e `buildFieldMessages` unificados em `buildMessagesWithSystem(msgs, systemPrompt)`. |
| BE-M3 | ✅ | `handlers/chat.go` | `addToolCallMessages(h, sessionID, toolCall, result)` elimina o bloco duplicado em 3 handlers. |
| BE-M4 | ✅ | `handlers/chat.go` | Constantes `msgProposalShown`, `msgFormFilled`, `msgFieldFilled`, `msgUserConfirmed` declaradas e usadas nos 3 handlers. |
| BE-M5 | ⏸ | `services/conversation.go` | `Session.Status` mantido; constantes `SessionStatusCollecting/Complete` adicionadas para documentar os valores válidos. |
| BE-M6 | ✅ | `services/llama.go` | `scanner.Buffer(make([]byte, 1<<20), 1<<20)` evita `token too long` em chunks >64KB. |
| BE-M7 | ✅ | `services/llama.go` | Chunks malformados agora logados via `log.Printf("WARN: skipping malformed SSE chunk: %v", err)`. |
| BE-M8 | ✅ | `models/types.go` | `ToolCallTypeFunction = "function"`, `SessionStatusCollecting`, `SessionStatusComplete` declarados como constantes. |
| BE-M9 | 🔧 | `main.go` | `gin.Default()` produz logs sem estrutura JSON. Pendente para quando houver requisito de observabilidade estruturada. |

### Frontend

| ID | Status | Arquivo | Problema |
|---|---|---|---|
| FE-M1 | ✅ | `ConfirmCard.vue` | `Object.entries` com filter extraído para `visibleFields` computed. |
| FE-M2 | ✅ | `components/form/FieldLabel.vue` | Componente criado; `SectionIdentity`, `SectionContext` e `SectionBenefits` refatorados para usá-lo. `✦` centralizado em `SPARKLE`. |
| FE-M3 | ⏸ | `useAiFormFill.ts` | Dois watchers com semânticas distintas (chat global vs painel de campo) no mesmo composable. Reconhecido — separação requer refactor de store. |
| FE-M4 | ✅ | `views/HomeView.vue`, `router/index.ts` | Imports relativos `../` substituídos pelo alias `@/*`. |
| FE-M5 | ✅ | `chatService.ts`, `fieldChatService.ts` | `API_BASE` extraído para `services/config.ts` compartilhado. |
| FE-M6 | ✅ | `FieldChatPanel.vue` | `fieldFilledPrefix` extraído para `fieldChatI18n` — resolvido junto com UX-M10. |
| FE-M7 | ✅ | `stores/chat.ts` | `clearFormFill` agora reseta `formFilled` para `false`. |
| FE-M8 | ⏸ | `composables/useFormContext.ts` | `inject` lança erro sem fallback — comportamento intencional para composable que exige provider. Documentado como design decision. |
| FE-M9 | ✅ | `FieldChatPanel.vue` | 40 strings i18n extraídas para `data/fieldGreetings.ts`. |

### UX / Acessibilidade

| ID | Status | Arquivo | Problema |
|---|---|---|---|
| UX-M1 | ✅ | `FloatingChat.css`, `FieldChatPanel.css`, `form.css` | `0.65rem`, `0.62rem`, `0.7rem` substituídos por `var(--font-size-xs)` em lang-btn, footer-hint, brand-status e ai-badge. |
| UX-M2 | ✅ | `FloatingChat.css`, `FieldChatPanel.css` | Padding e gap dos containers principais (`chat-header`, `chat-messages`, `panel-header`, `panel-messages`, `panel-footer`) substituídos por `var(--space-*)`. |
| UX-M3 | ✅ | `FloatingChat.vue` | `aria-label` de fechar e enviar usam strings do objeto `currentI18n()`, não mais hardcoded em PT. |
| UX-M4 | ✅ | `FloatingChat.vue` | Botões de idioma têm `aria-label` com nome completo do idioma via `langNames`. |
| UX-M5 | ✅ | `ConfirmCard.vue` | `FIELD_ORDER` array garante ordem determinística: title → businessLine → requesterBU → busInterested → … |
| UX-M6 | ✅ | `FieldChatPanel.vue` | Botão "Accept & close" (i18n) aparece com transição após `field_fill`; fecha o painel imediatamente. |
| UX-M7 | ✅ | `FloatingChat.vue` | Botão de reset na header limpa mensagens, sessionId e localStorage, reiniciando a conversa. |
| UX-M8 | ✅ | `FloatingChat.vue` | Duplicata de UX-08 — resolvida junto. |
| UX-M9 | ✅ | `AnalysisModal.vue` | "Submit anyway" tem `aria-describedby` explicando consequência da ação. |
| UX-M10 | ✅ | `FieldChatPanel.vue`, `data/fieldChatI18n.ts` | Todas as strings extraídas para `fieldChatI18n` com suporte a PT/EN/ES/FR. |

### QA / Testes

| ID | Status | Área | Problema |
|---|---|---|---|
| QA-M1 | 🔧 | E2E | Testes de AbortError com `waitForTimeout(500)` — flakiness em CI lento. |
| QA-M2 | 🔧 | E2E | `'send button is disabled while loading'` usa `setTimeout(3000)` dentro do route handler — atraso real, sensível ao tempo. |
| QA-M3 | 🔧 | E2E | 52 testes idênticos repetidos por idioma. Cobertura real é a mesma; 4× overhead de manutenção. |
| QA-M4 | 🔧 | `chatService.spec.ts` | Asserção de AbortError: `expect(err instanceof Error || err instanceof DOMException).toBe(true)` — passa para qualquer erro lançado. |
| QA-M5 | 🔧 | Backend | `services/conversation.go` `AddMessage` com sessão inexistente faz no-op silencioso sem teste. |
| QA-M6 | 🔧 | Backend | `bufio.Scanner` sem buffer — chunks >64KB causam término silencioso, sem teste. |
| QA-M7 | 🔧 | Backend | `sse.ts` linha 6: `response.body!` non-null assertion sem guarda. `body === null` lança erro não descritivo. |

---

## Severidade Baixa

| ID | Status | Área | Arquivo | Problema |
|---|---|---|---|---|
| LO-01 | ✅ | Frontend | `stores/counter.ts` | Store de scaffolding Vite deletada. |
| LO-02 | ✅ | Frontend | `FloatingChat.vue` | Binding `:class="{ focused: … }"` morto removido do `.input-box`. |
| LO-03 | ✅ | Frontend | `FieldLabel.vue` | `✦` centralizado na constante `SPARKLE` dentro de `FieldLabel.vue`. |
| LO-04 | ✅ | Frontend | `SectionBenefits.vue` | `<div class="field field--placeholder">` removido — CSS Grid auto-avança `field--full` corretamente. |
| LO-05 | ✅ | UX | `AnalysisModal.vue` | `aria-label="Close AI review dialog"` descritivo implementado. |
| LO-06 | ✅ | UX | `FloatingChat.css` | Painel corrigido para `width: 380px; height: 520px` conforme spec. |
| LO-07 | ✅ | UX | `ChatFab.vue` | FAB corrigido para `3.5rem` (56px) conforme spec. |
| LO-08 | ✅ | Backend | `main.go` | `router.Run` envolto em `if err := ...; err != nil { log.Fatal(err) }`. |
| LO-09 | ✅ | Backend | `services/tools.go` | `formProperties()` removida; substituída por `var formProps` calculado uma vez no startup. |
| LO-10 | ✅ | Backend | `services/llama.go` | Variável de request renomeada para `reqBytes`; `errBody` usada para leitura de erro de resposta. |
| LO-11 | ✅ | QA | `chatService.spec.ts` | Helper `makeSplitStream` morto removido. |
| LO-12 | 🔧 | QA | `chat-store.spec.ts` | Nomes de teste enganosos (`'stores a reference-equal copy'`, `'stores all required fields'` com 10 assertions). |

---

## Resumo Executivo

| Categoria | ✅ Resolvido | 🔧 Pendente | ⏸ Adiado | Total |
|---|---|---|---|---|
| Backend | 14 | 2 | 2 | 18 |
| Frontend | 16 | 1 | 2 | 19 |
| UX / Acessibilidade | 25 | 0 | 0 | 25 |
| QA / Testes | 2 | 16 | 0 | 18 |
| **Total** | **57** | **19** | **4** | **80** |

### Pendências de maior impacto restantes

1. **BE-M9**: Observabilidade estruturada — `gin.Default()` sem logs JSON.
2. **FE-M3**: Dois watchers com semânticas distintas no mesmo composable `useAiFormFill`.
3. **QA-03/04/05/06/07**: `analyzeForm`, `confirmForm`, `useAiFormFill`, `useFormAnalysis`, stores sem cobertura de testes.
4. **QA-13**: Page Object Model para os testes E2E — seletores CSS espalhados em ~80 testes.
