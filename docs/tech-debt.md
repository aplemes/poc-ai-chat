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
| BE-01 | 🔧 | `handlers/chat.go`, `handlers/form.go` | 75, 93, 31, 66 | Handler muta `session.PendingFormData` diretamente, sem passar pelo service. Viola encapsulamento e cria risco de race condition. |
| BE-02 | 🔧 | `handlers/form.go` | 42–64 | Lógica de negócio (montagem de `ToolCall` sintético, injeção de mensagens, reset de dados) dentro do handler. Pertence ao service. |
| BE-03 | 🔧 | `handlers/chat.go` linha 123, `handlers/form.go` linha 47 | — | `json.Marshal` com erro descartado via `_`. Em falha o cliente recebe evento SSE corrompido sem nenhum log. |
| BE-04 | 🔧 | `services/conversation.go` | 97 | `crypto/rand.Read` sem tratamento de erro. Em ambientes com entropia insuficiente produz IDs de sessão previsíveis ou repetidos. |
| BE-05 | 🔧 | `services/conversation.go` | — | `GetByID` retorna ponteiro direto ao objeto dentro do map. Qualquer acesso ao ponteiro fora do lock é uma data race. Mutex protege o map, não o objeto apontado. |
| BE-06 | 🔧 | `handlers/chat.go`, `handlers/field.go`, `handlers/form.go` | múltiplos | SSE headers (`Content-Type`, `Cache-Control`, `Connection`) copiados em 4 locais. Uma função `setSSEHeaders(c, sessionID)` eliminaria a duplicação. |
| BE-07 | 🔧 | Todo o backend | — | **Zero arquivos `_test.go`**. Nenhuma cobertura de testes em handlers, services ou models. |

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
| UX-08 | 🔧 | `ChatFab.vue` | — | Badge de notificação do FAB usa `aria-hidden="true"` mas o `aria-label` do botão é estático — não anuncia "nova mensagem". |
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
| BE-M1 | 🔧 | `services/prompts.go` | 280 linhas com 4 responsabilidades (resolução de língua, prompts principais, prompts por campo, prompt de análise). Candidatos: `prompts_main.go`, `prompts_field.go`. |
| BE-M2 | 🔧 | `handlers/chat.go`, `handlers/field.go` | `buildMessages` / `buildFieldMessages` estruturalmente idênticas. Unificar em `buildMessagesWithSystem(session, systemPrompt)`. |
| BE-M3 | 🔧 | `handlers/chat.go`, `handlers/field.go`, `handlers/form.go` | Bloco de adição de mensagens de tool call duplicado 3 vezes. |
| BE-M4 | 🔧 | `handlers/form.go` | Strings de histórico para o LLM (`"Proposal shown to user…"`, `"Form filled successfully."`) hardcoded sem constantes. Mudança quebra histórico silenciosamente. |
| BE-M5 | 🔧 | `services/conversation.go` | `Session.Status` inicializado como `"collecting"` mas `"complete"` nunca é atribuído. Campo é dead code. |
| BE-M6 | 🔧 | `services/conversation.go` | `bufio.Scanner` sem buffer customizado. Respostas SSE grandes (>64KB) causam `token too long` e encerramento silencioso do stream. |
| BE-M7 | 🔧 | `services/llama.go` | Chunks malformados descartados com `continue` sem nenhum log. Stream aparece incompleto sem indicação de erro no servidor. |
| BE-M8 | 🔧 | `models/types.go` | `ToolCall.Type` sempre `"function"` mas sem constante. `Session.Status` também sem constantes. |
| BE-M9 | 🔧 | `main.go` | `gin.Default()` produz logs sem estrutura JSON, dificultando observabilidade. |

### Frontend

| ID | Status | Arquivo | Problema |
|---|---|---|---|
| FE-M1 | ✅ | `ConfirmCard.vue` | `Object.entries` com filter extraído para `visibleFields` computed. |
| FE-M2 | 🔧 | Section*.vue (×10 campos) | Bloco `label-row` + badge + botão ✦ repetido para cada campo. Deveria ser componente `FieldLabel.vue`. |
| FE-M3 | ⏸ | `useAiFormFill.ts` | Dois watchers com semânticas distintas (chat global vs painel de campo) no mesmo composable. Reconhecido — separação requer refactor de store. |
| FE-M4 | ✅ | `views/HomeView.vue`, `router/index.ts` | Imports relativos `../` substituídos pelo alias `@/*`. |
| FE-M5 | ✅ | `chatService.ts`, `fieldChatService.ts` | `API_BASE` extraído para `services/config.ts` compartilhado. |
| FE-M6 | 🔧 | `FieldChatPanel.vue` | Mensagem de confirmação de fill (`"✓ Field filled with: …"`) hardcoded em inglês num componente multilíngue. |
| FE-M7 | ✅ | `stores/chat.ts` | `clearFormFill` agora reseta `formFilled` para `false`. |
| FE-M8 | ⏸ | `composables/useFormContext.ts` | `inject` lança erro sem fallback — comportamento intencional para composable que exige provider. Documentado como design decision. |
| FE-M9 | ✅ | `FieldChatPanel.vue` | 40 strings i18n extraídas para `data/fieldGreetings.ts`. |

### UX / Acessibilidade

| ID | Status | Arquivo | Problema |
|---|---|---|---|
| UX-M1 | 🔧 | `FloatingChat.vue`, `FieldChatPanel.vue`, `AnalysisModal.vue` | `font-size: 0.65rem`, `0.62rem`, `0.625rem` usados diretamente. Design system mínimo é `--font-size-xs: 0.75rem`. |
| UX-M2 | 🔧 | `FloatingChat.vue`, todos | Tokens `--space-*` existem mas **nunca usados**. Todos os paddings e gaps são valores `rem` inline. |
| UX-M3 | ✅ | `FloatingChat.vue` | `aria-label` de fechar e enviar usam strings do objeto `currentI18n()`, não mais hardcoded em PT. |
| UX-M4 | ✅ | `FloatingChat.vue` | Botões de idioma têm `aria-label` com nome completo do idioma via `langNames`. |
| UX-M5 | 🔧 | Seção de confirmação | Ordem de campos no card de confirmação depende de `Object.entries` — não controlada pelo design. |
| UX-M6 | 🔧 | `FieldChatPanel.vue` | Sem botão "Aceitar e fechar" após um `field_fill`. Usuário precisa fechar manualmente para ver o campo preenchido. |
| UX-M7 | 🔧 | `FloatingChat.vue` | Sem mecanismo de reset/nova conversa na UI. Único escape é limpar `localStorage` manualmente. |
| UX-M8 | 🔧 | `ChatFab.vue` | Badge de notificação usa `aria-hidden="true"` mas o `aria-label` do botão é estático — não anuncia "nova mensagem". |
| UX-M9 | ✅ | `AnalysisModal.vue` | "Submit anyway" tem `aria-describedby` explicando consequência da ação. |
| UX-M10 | 🔧 | `FieldChatPanel.vue` | Erros de conexão hardcoded em inglês (`'Could not process your message.'`). |

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
| LO-02 | 🔧 | Frontend | `FloatingChat.vue` | Classe `.focused` aplicada quando `loading === false` — lógica invertida, parece resíduo. |
| LO-03 | 🔧 | Frontend | Section*.vue | Caractere `✦` hardcoded em múltiplos locais. Deveria ser uma constante nomeada. |
| LO-04 | 🔧 | Frontend | `SectionBenefits.vue` | `<div class="field field--placeholder">` como espaçador de grid não-semântico. |
| LO-05 | ✅ | UX | `AnalysisModal.vue` | `aria-label="Close AI review dialog"` descritivo implementado. |
| LO-06 | 🔧 | UX | `FloatingChat.vue` | Dimensões do painel (400×620px) divergem da spec do design system (380×520px). |
| LO-07 | 🔧 | UX | `FloatingChat.vue` | FAB 52px (`3.25rem`) vs spec 56px. |
| LO-08 | 🔧 | Backend | `main.go` | `router.Run` sem `log.Fatal`. Falha de bind encerra o processo silenciosamente. |
| LO-09 | 🔧 | Backend | `services/tools.go` | `formProperties()` recria map completo a cada chamada. Deveria ser variável de pacote. |
| LO-10 | 🔧 | Backend | `services/llama.go` | Variável `body` reutilizada para bytes do request e bytes do erro no mesmo escopo. |
| LO-11 | ✅ | QA | `chatService.spec.ts` | Helper `makeSplitStream` morto removido. |
| LO-12 | 🔧 | QA | `chat-store.spec.ts` | Nomes de teste enganosos (`'stores a reference-equal copy'`, `'stores all required fields'` com 10 assertions). |

---

## Resumo Executivo

| Categoria | ✅ Resolvido | 🔧 Pendente | ⏸ Adiado | Total |
|---|---|---|---|---|
| Backend | 0 | 16 | 0 | 16 |
| Frontend | 13 | 4 | 2 | 19 |
| UX / Acessibilidade | 14 | 11 | 0 | 25 |
| QA / Testes | 2 | 16 | 0 | 18 |
| **Total** | **29** | **47** | **2** | **78** |

### Pendências de maior impacto restantes

1. **BE-05**: Race condition real em `ConversationService` — ponteiro compartilhado sem proteção fora do mutex.
2. **BE-07 + QA-01**: Zero testes em todo o backend — qualquer refactor é cego sem cobertura mínima.
3. **FE-M2**: `FieldLabel.vue` — bloco label/badge/botão repetido em 10 campos sem componente.
4. **UX-M6/M7**: Sem "Aceitar e fechar" após field_fill + sem reset de conversa na UI.
5. **QA-13**: Page Object Model para os testes E2E — seletores CSS espalhados em ~80 testes.
