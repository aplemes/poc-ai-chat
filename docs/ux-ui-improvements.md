# UX/UI Improvements — Demand Assistant

**Auditoria realizada por:** ux-ui-designer agent  
**Data:** 2026-04-17  
**Arquivos analisados:** `FloatingChat.vue`, `RequestForm.vue`, `stores/chat.ts`, `services/chatService.ts`, `assets/base.css`, `assets/main.css`

---

## Sumário executivo

O sistema tem uma base funcional sólida, mas sofre de **três problemas centrais**:

1. **Design system fragmentado** — cores hardcoded em HSL/hex espalhadas por dois componentes sem tokens CSS centralizados. Qualquer mudança de marca exige busca e substituição manual.
2. **Identidade visual genérica** — o verde Vue padrão (`hsla(160, 100%, 37%, 1)`) foi reaproveitado diretamente. Não há identidade visual Adeo.
3. **Lacunas de acessibilidade e feedback** — ausência de `aria-live` na lista de mensagens, foco invisível em vários elementos, e estados de erro sem feedback visual adequado.

---

## 1. Design System — Criar arquivo de tokens

### Problema
Não existe um arquivo central de tokens CSS. As cores aparecem em pelo menos **12 lugares diferentes** hardcoded:

| Valor | Ocorrências |
|---|---|
| `hsla(160, 100%, 37%, 1)` | FAB, header, user message, send button, input focus, submit button, toggle active, ai-filled border |
| `#1e293b` | form-title, input color |
| `#374151` | label color |
| `#d1d5db` | input border |
| `#f9fafb` | input background |
| `#3b82f6` | input focus (azul — cor diferente do verde primário!) |
| `#0d9488` | toggle active (teal — terceira cor primária!) |

### Ação necessária
Criar `frontend/ai-chat/src/assets/tokens.css` e importar em `main.css`:

```css
:root {
  /* === Marca Adeo === */
  --adeo-green-500: hsla(160, 100%, 37%, 1);   /* cor primária atual */
  --adeo-green-600: hsla(160, 100%, 30%, 1);   /* hover */
  --adeo-green-50:  hsla(160, 100%, 37%, 0.08);

  /* === Paleta primária (substituir o azul #3b82f6 e teal #0d9488) === */
  --color-primary:        var(--adeo-green-500);
  --color-primary-hover:  var(--adeo-green-600);
  --color-primary-subtle: var(--adeo-green-50);

  /* === Neutros === */
  --color-neutral-900: #0f172a;
  --color-neutral-800: #1e293b;
  --color-neutral-700: #374151;
  --color-neutral-400: #9ca3af;
  --color-neutral-300: #d1d5db;
  --color-neutral-200: #e2e8f0;
  --color-neutral-100: #f1f5f9;
  --color-neutral-50:  #f9fafb;
  --color-neutral-0:   #ffffff;

  /* === Semânticos === */
  --color-error:   #dc2626;
  --color-warning: #d97706;
  --color-success: #16a34a;

  /* === Badge IA === */
  --color-ai-bg:     hsla(160, 100%, 37%, 0.12);
  --color-ai-text:   hsla(160, 100%, 28%, 1);
  --color-ai-border: hsla(160, 100%, 37%, 0.3);

  /* === Tipografia === */
  --font-size-xs:   0.75rem;
  --font-size-sm:   0.875rem;
  --font-size-base: 0.9rem;
  --font-size-md:   1rem;
  --font-size-lg:   1.25rem;

  /* === Espaçamento === */
  --space-1: 0.25rem;
  --space-2: 0.5rem;
  --space-3: 0.75rem;
  --space-4: 1rem;
  --space-6: 1.5rem;
  --space-8: 2rem;

  /* === Bordas === */
  --radius-sm:   4px;
  --radius-md:   6px;
  --radius-lg:   8px;
  --radius-xl:   12px;
  --radius-full: 9999px;

  /* === Sombras === */
  --shadow-sm: 0 1px 2px rgba(0,0,0,0.05);
  --shadow-md: 0 4px 16px rgba(0,0,0,0.12);
  --shadow-lg: 0 8px 32px rgba(0,0,0,0.15);

  /* === Transições === */
  --transition-fast:   0.15s ease-out;
  --transition-normal: 0.2s ease-out;
  --transition-slow:   0.25s ease-out;
}
```

Após criar os tokens, **substituir todos os valores hardcoded** nos dois componentes pelas variáveis CSS correspondentes.

---

## 2. Inconsistência de cores primárias

### Problema — Severidade: CRÍTICA
O sistema usa **três cores diferentes** para elementos de destaque:
- `hsla(160, 100%, 37%)` — verde (FAB, header, botões)
- `#3b82f6` — azul (foco dos inputs no `RequestForm`)
- `#0d9488` — teal (toggle ativo no `RequestForm`)

Isso cria uma experiência visual inconsistente e confusa.

### Arquivo: `RequestForm.vue:394-398` e `RequestForm.vue:447-449`

### Ação necessária
Substituir `#3b82f6` e `#0d9488` pela variável `--color-primary`. O foco dos inputs e o toggle ativo devem usar a mesma cor primária da marca.

---

## 3. FloatingChat — Melhorias de UX

### 3.1 Ausência de `aria-live` na lista de mensagens — Severidade: CRÍTICA (Acessibilidade)

**Arquivo:** `FloatingChat.vue:142`

Leitores de tela não anunciam novas mensagens quando chegam. A `div.chat-messages` precisa de `role="log"` e `aria-live="polite"`.

```html
<!-- Atual -->
<div ref="messagesEl" class="chat-messages">

<!-- Correto -->
<div ref="messagesEl" class="chat-messages" role="log" aria-live="polite" aria-label="Conversa com o assistente">
```

### 3.2 Textarea não cresce automaticamente — Severidade: MAIOR

**Arquivo:** `FloatingChat.vue:158-164`

O `textarea` tem `rows="1"` e `max-height: 120px`, mas não cresce conforme o usuário digita. Usuários que escrevem mensagens longas precisam rolar dentro do campo, o que é frustrante.

**Ação:** Implementar auto-resize via `@input` handler:

```typescript
function autoResize(e: Event) {
  const el = e.target as HTMLTextAreaElement
  el.style.height = 'auto'
  el.style.height = Math.min(el.scrollHeight, 120) + 'px'
}
```

### 3.3 Mensagem de confirmação de form fill hardcoded em PT — Severidade: MAIOR

**Arquivo:** `FloatingChat.vue:87`

```typescript
assistantMsg.text = 'Formulário preenchido! Confira os campos acima.'
```

O sistema suporta 4 idiomas, mas esta mensagem está hardcoded em português. Deve entrar no objeto `languages`.

**Ação:** Adicionar campo `formFilled` ao objeto de cada idioma:
```typescript
{ code: 'pt', ..., formFilled: 'Formulário preenchido! Confira os campos acima.' },
{ code: 'en', ..., formFilled: 'Form filled! Please review the fields above.' },
{ code: 'es', ..., formFilled: '¡Formulario completado! Revisa los campos arriba.' },
{ code: 'fr', ..., formFilled: 'Formulaire rempli ! Vérifiez les champs ci-dessus.' },
```

### 3.4 Mensagem de erro hardcoded em PT — Severidade: MAIOR

**Arquivo:** `FloatingChat.vue:91` e `FloatingChat.vue:98-101`

Mesmas mensagens `'Erro ao processar...'`, `'(cancelado)'` e `'Não foi possível conectar...'` estão hardcoded em português. Aplicar o mesmo tratamento multilíngue.

### 3.5 Botão FAB sem indicador de notificação — Severidade: MENOR

Quando o assistente preenche o formulário e o chat está fechado, o usuário não recebe nenhum feedback visual de que algo aconteceu. Um badge/dot de notificação no FAB melhoraria a percepção do estado da IA.

**Ação:** Adicionar um `ref` booleano `hasNewMessage` que ativa um dot vermelho no FAB enquanto o chat estiver fechado e houver mensagem nova não lida.

### 3.6 Foco invisível no botão de fechar e FAB — Severidade: MAIOR (Acessibilidade)

**Arquivo:** `FloatingChat.vue:278-292` (`.close-btn`) e `FloatingChat.vue:200-224` (`.fab`)

Nenhum dos dois tem estilo de `:focus-visible`. Navegação por teclado não mostra onde o foco está.

**Ação:** Adicionar em ambos:
```css
.fab:focus-visible,
.close-btn:focus-visible {
  outline: 2px solid var(--color-primary);
  outline-offset: 2px;
}
```

### 3.7 Seletor de idioma acessível — Severidade: MENOR

**Arquivo:** `FloatingChat.vue:125-133`

Os botões de idioma não têm `aria-label` descritivo. `aria-label="Português"` seria mais informativo que só o código `"PT"`.

---

## 4. RequestForm — Melhorias de UX

### 4.1 Campo `busInterested` só aceita uma BU — Severidade: CRÍTICA (Funcional)

**Arquivo:** `RequestForm.vue:200-216` e `RequestForm.vue:125`

O campo `busInterested` é declarado como `string` no estado do form, mas o schema do backend suporta um array de BUs (`busInterested: string[]`). O `<select>` atual não permite seleção múltipla.

**Ação:** Mudar o tipo para `string[]`, usar `<select multiple>` ou um componente de multi-select com chips/tags.

```typescript
// Atual
busInterested: '',

// Correto
busInterested: [] as string[],
```

### 4.2 Ausência de validação visual de campos obrigatórios — Severidade: CRÍTICA

**Arquivo:** `RequestForm.vue:138-140`

O `handleSubmit` só faz `console.log`. Não há validação, não há feedback de erro nos campos obrigatórios. O usuário pode submeter um formulário vazio.

**Ação:** Implementar validação com estados de erro por campo:
```typescript
const errors = ref<Record<string, string>>({})

function validate(): boolean {
  errors.value = {}
  if (!form.value.title.trim()) errors.value.title = 'Campo obrigatório'
  if (!form.value.businessLine) errors.value.businessLine = 'Selecione uma opção'
  // ... demais campos obrigatórios
  return Object.keys(errors.value).length === 0
}
```

Mostrar mensagem de erro abaixo de cada campo com `role="alert"`.

### 4.3 Label do `whyDemand` é confusa — Severidade: MENOR

**Arquivo:** `RequestForm.vue:241-243`

O label atual é uma instrução longa no lugar de um nome de campo claro. Melhor separar o nome do campo da instrução usando `.field-hint`.

```html
<!-- Atual: label mistura título e instrução -->
<label>Describe current situation, the pain points, comparison with competitors... - required</label>

<!-- Melhor -->
<label>Why are you making this demand? <span class="required">- required</span></label>
<p class="field-hint">Describe the current situation, pain points, and comparison with competitors.</p>
```

### 4.4 Transição de remoção do badge IA ausente — Severidade: MENOR

**Arquivo:** `RequestForm.vue` — badge `.ai-badge` e classe `.ai-filled`

Quando o usuário edita um campo preenchido pela IA, o badge some instantaneamente. Uma transição de 150ms melhora a percepção de que foi o usuário quem fez a mudança.

**Ação:**
```css
.ai-badge {
  transition: opacity var(--transition-fast), transform var(--transition-fast);
}

/* Usar v-show com transition wrapper ao invés de v-if para animar a saída */
```

### 4.5 Foco dos inputs sem `outline` customizado — Severidade: MENOR (Acessibilidade)

**Arquivo:** `RequestForm.vue:392-398`

O `outline: none` remove o foco nativo, mas o substituto `box-shadow` funciona. Porém ao usar `prefers-reduced-motion`, a transição da borda deveria ser suprimida.

**Ação:** Adicionar:
```css
@media (prefers-reduced-motion: reduce) {
  input, select, textarea, .toggle-btn, .ai-badge {
    transition: none;
  }
}
```

### 4.6 Submit button sem estado `disabled` durante loading — Severidade: MAIOR

**Arquivo:** `RequestForm.vue:321-323`

O botão Submit não tem estado desabilitado. O usuário pode clicar múltiplas vezes e submeter o formulário várias vezes enquanto um envio já está em curso.

**Ação:** Adicionar `ref<boolean> submitting` e desabilitar o botão + mostrar spinner durante o envio.

### 4.7 Formulário sem seções visuais — Severidade: MENOR (UX)

O formulário tem 10 campos apresentados em sequência sem agrupamento visual. Adicionar separadores e títulos de seção melhora a escaneabilidade:

- **Seção 1:** "About the demand" (title, businessLine, requesterBU, busInterested, timeSensitive)
- **Seção 2:** "Context & impact" (whyDemand, whoIsImpacted)
- **Seção 3:** "Expected benefits" (benefitCategory, benefitHypothesis, measureBenefits)

---

## 5. Paleta de cores recomendada

A paleta atual usa o verde Vue padrão. Para um produto interno Adeo, recomenda-se alinhar à identidade da marca. Opções:

### Opção A — Manter verde, refinar tons
```
Primário:  #00875A  (verde Adeo mais sóbrio)
Hover:     #006644
Foco:      #00875A  (unificar com primário — eliminar o azul e teal)
```

### Opção B — Azul corporativo (mais neutro para ferramenta interna)
```
Primário:  #2563EB
Hover:     #1D4ED8
Foco:      #2563EB
```

> **Recomendação:** confirmar com o time de design Adeo qual cor primária deve ser usada. Enquanto isso, a ação mais segura é **unificar** as três cores atuais (verde, azul, teal) em uma única variável `--color-primary`.

---

## 6. Lista de ações priorizadas

| # | Prioridade | Ação | Arquivo |
|---|---|---|---|
| 1 | CRÍTICA | Criar `tokens.css` e substituir todos os valores hardcoded | `tokens.css`, ambos componentes |
| 2 | CRÍTICA | Unificar cor primária (eliminar `#3b82f6` e `#0d9488`) | `RequestForm.vue:394`, `RequestForm.vue:448` |
| 3 | CRÍTICA | Campo `busInterested` → multi-select com array | `RequestForm.vue:125`, `RequestForm.vue:200-216` |
| 4 | CRÍTICA | Implementar validação de formulário com estados de erro | `RequestForm.vue:138` |
| 5 | CRÍTICA | Adicionar `role="log"` e `aria-live="polite"` no chat | `FloatingChat.vue:142` |
| 6 | MAIOR | Auto-resize do textarea do chat | `FloatingChat.vue:158-164` |
| 7 | MAIOR | Internacionalizar mensagens de confirmação e erro do chat | `FloatingChat.vue:87`, `FloatingChat.vue:91-101` |
| 8 | MAIOR | Focus ring visível no FAB e close button | `FloatingChat.vue:200`, `FloatingChat.vue:278` |
| 9 | MAIOR | Desabilitar submit durante envio (estado loading) | `RequestForm.vue:321` |
| 10 | MENOR | Badge IA com transição de saída animada | `RequestForm.vue` — `.ai-badge` |
| 11 | MENOR | Separar label de instrução no campo `whyDemand` | `RequestForm.vue:241` |
| 12 | MENOR | Adicionar seções visuais no formulário | `RequestForm.vue` |
| 13 | MENOR | FAB com badge de notificação quando há mensagem nova | `FloatingChat.vue` |
| 14 | MENOR | `aria-label` descritivo nos botões de idioma | `FloatingChat.vue:125` |
| 15 | MENOR | `prefers-reduced-motion` para todas as transições | Ambos componentes |

---

## 7. Instruções para o agente vue-frontend-expert

Ao implementar as melhorias acima, seguir esta ordem:

1. **Primeiro:** criar `frontend/ai-chat/src/assets/tokens.css` com todos os tokens e importar em `main.css`
2. **Segundo:** refatorar `FloatingChat.vue` substituindo hardcodes pelos tokens (itens 2, 5, 6, 7, 8, 13, 14, 15)
3. **Terceiro:** refatorar `RequestForm.vue` substituindo hardcodes pelos tokens e corrigindo bugs (itens 2, 3, 4, 9, 10, 11, 12, 15)
4. **Por fim:** rodar `npm run format`, `npm run lint` e `npm run type-check` para garantir que tudo passa

Não alterar lógica de negócio (SSE, store, sessionId) — apenas UX, estilos e acessibilidade.
