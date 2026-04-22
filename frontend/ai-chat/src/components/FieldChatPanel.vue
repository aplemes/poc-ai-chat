<script setup lang="ts">
import { ref, watch, nextTick, onUnmounted, computed } from 'vue'
import { marked } from 'marked'
import DOMPurify from 'dompurify'
import { useFieldChatStore } from '@/stores/fieldChat'
import { sendFieldMessage } from '@/services/fieldChatService'

const fieldChatStore = useFieldChatStore()

interface Message {
  id: number
  role: 'user' | 'assistant'
  text: string
}

const fieldLabels: Record<string, string> = {
  title: 'Title',
  businessLine: 'Business Line',
  requesterBU: 'Requester BU',
  busInterested: 'BUs Interested',
  timeSensitive: 'Time Sensitive',
  whyDemand: 'Why Demand',
  whoIsImpacted: 'Who Is Impacted',
  benefitCategory: 'Benefit Category',
  benefitHypothesis: 'Benefit Hypothesis',
  measureBenefits: 'Measure Benefits',
}

const fieldGreetings: Record<string, Record<string, string>> = {
  en: {
    title:
      "I'll help you write a clear, well-formatted title for your demand. What do you want to achieve?",
    businessLine:
      "I'll help you identify the right Business Line. Which organisational area is responsible for this demand?",
    requesterBU: "I'll help you find your Business Unit. Which BU is making this request?",
    busInterested:
      "I'll help you identify interested BUs. Which Business Units have already shown interest in this demand?",
    timeSensitive:
      'Is this demand urgent due to a legal deadline or security risk? Or does it have no special urgency?',
    whyDemand:
      "I'll help you describe the context and problem. What is the current situation, and what pain points need solving?",
    whoIsImpacted:
      "I'll help you identify who is affected. Which user types are impacted, and roughly how many?",
    benefitCategory:
      "I'll help you pick the right benefit category. What is the main type of benefit this demand will bring?",
    benefitHypothesis:
      "I'll help you articulate your hypothesis. How do you believe this demand will achieve the expected benefit?",
    measureBenefits:
      "I'll help you define success metrics. What KPIs will you track, and over what timeframe?",
  },
  pt: {
    title: 'Vou te ajudar a escrever um título claro e bem formatado. O que você quer alcançar?',
    businessLine:
      'Vou te ajudar a identificar a Business Line correta. Qual área organizacional é responsável por esta demanda?',
    requesterBU:
      'Vou te ajudar a encontrar sua Unidade de Negócio. Qual BU está fazendo esta solicitação?',
    busInterested:
      'Vou te ajudar a identificar as BUs interessadas. Quais Unidades de Negócio já demonstraram interesse nesta demanda?',
    timeSensitive:
      'Esta demanda tem urgência por prazo legal ou risco de segurança? Ou não tem urgência especial?',
    whyDemand:
      'Vou te ajudar a descrever o contexto e o problema. Qual é a situação atual e quais são os pontos de dor?',
    whoIsImpacted:
      'Vou te ajudar a identificar quem é afetado. Quais tipos de usuários são impactados e quantos aproximadamente?',
    benefitCategory:
      'Vou te ajudar a escolher a categoria de benefício correta. Qual é o principal tipo de benefício que esta demanda trará?',
    benefitHypothesis:
      'Vou te ajudar a formular sua hipótese. Como você acredita que esta demanda vai gerar o benefício esperado?',
    measureBenefits:
      'Vou te ajudar a definir as métricas de sucesso. Quais KPIs você vai acompanhar e em qual prazo?',
  },
  es: {
    title: 'Te ayudaré a escribir un título claro y bien formateado. ¿Qué quieres lograr?',
    businessLine:
      'Te ayudaré a identificar la Business Line correcta. ¿Qué área organizacional es responsable de esta demanda?',
    requesterBU:
      'Te ayudaré a encontrar tu Unidad de Negocio. ¿Qué BU está haciendo esta solicitud?',
    busInterested:
      'Te ayudaré a identificar las BUs interesadas. ¿Qué Unidades de Negocio ya han mostrado interés en esta demanda?',
    timeSensitive:
      '¿Esta demanda tiene urgencia por plazo legal o riesgo de seguridad? ¿O no tiene urgencia especial?',
    whyDemand:
      'Te ayudaré a describir el contexto y el problema. ¿Cuál es la situación actual y cuáles son los puntos de dolor?',
    whoIsImpacted:
      'Te ayudaré a identificar quiénes se ven afectados. ¿Qué tipos de usuarios están impactados y cuántos aproximadamente?',
    benefitCategory:
      'Te ayudaré a elegir la categoría de beneficio correcta. ¿Cuál es el principal tipo de beneficio que aportará esta demanda?',
    benefitHypothesis:
      'Te ayudaré a formular tu hipótesis. ¿Cómo crees que esta demanda generará el beneficio esperado?',
    measureBenefits:
      'Te ayudaré a definir las métricas de éxito. ¿Qué KPIs seguirás y en qué plazo?',
  },
  fr: {
    title:
      'Je vais vous aider à rédiger un titre clair et bien formaté. Que souhaitez-vous accomplir ?',
    businessLine:
      'Je vais vous aider à identifier la bonne Business Line. Quelle zone organisationnelle est responsable de cette demande ?',
    requesterBU: 'Je vais vous aider à trouver votre Business Unit. Quelle BU fait cette demande ?',
    busInterested:
      "Je vais vous aider à identifier les BUs intéressées. Quelles Business Units ont déjà montré de l'intérêt pour cette demande ?",
    timeSensitive:
      "Cette demande est-elle urgente en raison d'une échéance légale ou d'un risque de sécurité ? Ou n'a-t-elle pas d'urgence particulière ?",
    whyDemand:
      'Je vais vous aider à décrire le contexte et le problème. Quelle est la situation actuelle et quels sont les points de douleur ?',
    whoIsImpacted:
      "Je vais vous aider à identifier qui est concerné. Quels types d'utilisateurs sont impactés et combien approximativement ?",
    benefitCategory:
      'Je vais vous aider à choisir la bonne catégorie de bénéfice. Quel est le principal type de bénéfice que cette demande apportera ?',
    benefitHypothesis:
      'Je vais vous aider à formuler votre hypothèse. Comment pensez-vous que cette demande générera le bénéfice attendu ?',
    measureBenefits:
      'Je vais vous aider à définir les indicateurs de succès. Quels KPIs suivrez-vous et dans quel délai ?',
  },
}

function getGreeting(field: string, lang: string): string {
  const byLang = fieldGreetings[lang] ?? fieldGreetings['en']
  return (
    byLang?.[field] ??
    `I'll help you fill the ${fieldLabels[field] ?? field} field. What information do you have?`
  )
}

function renderMarkdown(text: string): string {
  return DOMPurify.sanitize(marked.parse(text) as string)
}

const messages = ref<Message[]>([])
const input = ref('')
const loading = ref(false)
const messagesEl = ref<HTMLElement | null>(null)
const textareaEl = ref<HTMLTextAreaElement | null>(null)
let nextId = 1
let abortController: AbortController | null = null

const languages = [
  { code: 'pt', label: 'PT' },
  { code: 'en', label: 'EN' },
  { code: 'es', label: 'ES' },
  { code: 'fr', label: 'FR' },
]

const language = ref(localStorage.getItem('chat_language') ?? 'en')

function setLanguage(code: string) {
  language.value = code
  localStorage.setItem('chat_language', code)
}

const activeField = computed(() => fieldChatStore.activeField)
const isOpen = computed(() => activeField.value !== null)
const fieldLabel = computed(() =>
  activeField.value ? (fieldLabels[activeField.value] ?? activeField.value) : '',
)

// Reset conversation when field changes
watch(activeField, (field) => {
  if (!field) return
  abortController?.abort()
  abortController = null
  loading.value = false
  input.value = ''
  messages.value = [{ id: nextId++, role: 'assistant', text: getGreeting(field, language.value) }]
  nextTick(() => textareaEl.value?.focus())
})

// Update greeting when language changes while panel is open
watch(language, (lang) => {
  if (!activeField.value || messages.value.length !== 1) return
  messages.value[0]!.text = getGreeting(activeField.value, lang)
})

async function scrollToBottom() {
  await nextTick()
  if (messagesEl.value) messagesEl.value.scrollTop = messagesEl.value.scrollHeight
}

function autoResize() {
  const el = textareaEl.value
  if (!el) return
  el.style.height = 'auto'
  el.style.height = Math.min(el.scrollHeight, 120) + 'px'
}

async function send() {
  const text = input.value.trim()
  if (!text || loading.value || !activeField.value) return

  const field = activeField.value

  messages.value.push({ id: nextId++, role: 'user', text })
  input.value = ''
  loading.value = true
  nextTick(() => {
    if (textareaEl.value) textareaEl.value.style.height = 'auto'
  })
  await scrollToBottom()

  const assistantMsg: Message = { id: nextId++, role: 'assistant', text: '' }
  messages.value.push(assistantMsg)

  abortController?.abort()
  abortController = new AbortController()

  try {
    const sessionId = fieldChatStore.getSessionId(field)
    const newId = await sendFieldMessage(
      sessionId,
      field,
      text,
      language.value,
      (event) => {
        if (event.type === 'token' && event.content) {
          assistantMsg.text += event.content
          scrollToBottom()
        } else if (event.type === 'field_fill' && event.data) {
          fieldChatStore.setFieldFill({ fieldName: event.data.fieldName, value: event.data.value })
          if (!assistantMsg.text) {
            assistantMsg.text = `✓ Field filled with: **${event.data.value}**`
          }
          scrollToBottom()
        } else if (event.type === 'error') {
          assistantMsg.text = 'Could not process your message. Try again.'
        }
      },
      abortController.signal,
    )
    if (newId) fieldChatStore.setSessionId(field, newId)
  } catch (err) {
    assistantMsg.text =
      err instanceof Error && err.name === 'AbortError' ? '(cancelled)' : 'Connection error.'
  } finally {
    loading.value = false
    abortController = null
  }
}

function onKeydown(e: KeyboardEvent) {
  if (e.key === 'Enter' && !e.shiftKey) {
    e.preventDefault()
    send()
  }
  if (e.key === 'Escape') fieldChatStore.closePanel()
}

onUnmounted(() => abortController?.abort())
</script>

<template>
  <Teleport to="body">
    <Transition name="panel">
      <div v-if="isOpen" class="field-panel-backdrop" @click.self="fieldChatStore.closePanel()">
        <div class="field-panel" role="dialog" :aria-label="`AI assistant for ${fieldLabel}`">
          <!-- Header -->
          <div class="panel-header">
            <div class="panel-title">
              <span class="panel-sparkle">✦</span>
              <div>
                <p class="panel-subtitle">Helping with</p>
                <p class="panel-field-name">{{ fieldLabel }}</p>
              </div>
            </div>
            <div class="panel-header-right">
              <div class="lang-selector">
                <button
                  v-for="lang in languages"
                  :key="lang.code"
                  class="lang-btn"
                  :class="{ active: language === lang.code }"
                  :aria-label="lang.label"
                  @click="setLanguage(lang.code)"
                >
                  {{ lang.label }}
                </button>
              </div>
              <button
                class="panel-close"
                aria-label="Close assistant"
                @click="fieldChatStore.closePanel()"
              >
                <svg
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="2.5"
                  stroke-linecap="round"
                >
                  <path d="M18 6L6 18M6 6l12 12" />
                </svg>
              </button>
            </div>
          </div>

          <!-- Messages -->
          <div ref="messagesEl" class="panel-messages" role="log" aria-live="polite">
            <TransitionGroup name="msg">
              <div v-for="msg in messages" :key="msg.id" class="msg-row" :class="msg.role">
                <div v-if="msg.role === 'assistant'" class="msg-avatar">
                  <svg viewBox="0 0 20 20" fill="currentColor">
                    <path
                      d="M10 2a8 8 0 100 16A8 8 0 0010 2zm0 2a6 6 0 110 12A6 6 0 0110 4z"
                      opacity=".3"
                    />
                    <path d="M6.5 10a1 1 0 112 0 1 1 0 01-2 0zm5 0a1 1 0 112 0 1 1 0 01-2 0z" />
                  </svg>
                </div>
                <div class="msg-bubble">
                  <span v-if="msg.role === 'assistant' && !msg.text && loading" class="typing">
                    <span></span><span></span><span></span>
                  </span>
                  <span
                    v-else-if="msg.role === 'assistant'"
                    class="md"
                    v-html="renderMarkdown(msg.text)"
                  />
                  <template v-else>{{ msg.text }}</template>
                </div>
              </div>
            </TransitionGroup>
          </div>

          <!-- Input -->
          <div class="panel-footer">
            <div class="input-box">
              <textarea
                ref="textareaEl"
                v-model="input"
                rows="1"
                placeholder="Type your answer..."
                :disabled="loading"
                @keydown="onKeydown"
                @input="autoResize"
              />
              <button
                class="send-btn"
                :disabled="!input.trim() || loading"
                aria-label="Send"
                @click="send"
              >
                <svg viewBox="0 0 20 20" fill="currentColor">
                  <path
                    d="M3.105 2.288a.75.75 0 00-.826.95l1.337 4.01a.75.75 0 00.593.518l5.662.944-5.662.944a.75.75 0 00-.593.519l-1.337 4.01a.75.75 0 00.826.95 19.955 19.955 0 0016.233-8.568.75.75 0 000-.904A19.955 19.955 0 003.105 2.288z"
                  />
                </svg>
              </button>
            </div>
            <p class="footer-hint">Enter to send · Shift+Enter new line · Esc to close</p>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.field-panel-backdrop {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.25);
  z-index: 8000;
  display: flex;
  justify-content: flex-end;
}

.field-panel {
  width: 380px;
  height: 100%;
  display: flex;
  flex-direction: column;
  background: var(--color-neutral-0);
  box-shadow: -4px 0 24px rgba(0, 0, 0, 0.12);
  overflow: hidden;
}

/* ── Header ── */
.panel-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 1.1rem 1.25rem;
  background: var(--gradient-primary);
  flex-shrink: 0;
}

.panel-title {
  display: flex;
  align-items: center;
  gap: 0.65rem;
}

.panel-sparkle {
  font-size: 1.25rem;
  color: rgba(255, 255, 255, 0.9);
  line-height: 1;
}

.panel-subtitle {
  font-size: 0.65rem;
  color: rgba(255, 255, 255, 0.7);
  margin: 0;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.panel-field-name {
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-semibold);
  color: #fff;
  margin: 0;
  line-height: 1.2;
}

.panel-header-right {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  flex-shrink: 0;
}

.lang-selector {
  display: flex;
  gap: 2px;
  background: rgba(0, 0, 0, 0.15);
  border-radius: var(--radius-sm);
  padding: 2px;
}

.lang-btn {
  padding: 0.15rem 0.4rem;
  border-radius: var(--radius-xs);
  border: none;
  background: transparent;
  color: rgba(255, 255, 255, 0.65);
  font-size: 0.62rem;
  font-weight: var(--font-weight-semibold);
  letter-spacing: 0.04em;
  cursor: pointer;
  transition:
    background var(--transition-fast),
    color var(--transition-fast);
}

.lang-btn:hover {
  color: #fff;
  background: rgba(255, 255, 255, 0.12);
}

.lang-btn.active {
  background: rgba(255, 255, 255, 0.22);
  color: #fff;
}

.panel-close {
  width: 1.85rem;
  height: 1.85rem;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(0, 0, 0, 0.12);
  border: none;
  border-radius: var(--radius-sm);
  color: rgba(255, 255, 255, 0.85);
  cursor: pointer;
  transition: background var(--transition-fast);
  flex-shrink: 0;
}

.panel-close:hover {
  background: rgba(0, 0, 0, 0.22);
  color: #fff;
}
.panel-close:focus-visible {
  outline: 2px solid rgba(255, 255, 255, 0.6);
  outline-offset: 1px;
}
.panel-close svg {
  width: 0.9rem;
  height: 0.9rem;
}

/* ── Messages ── */
.panel-messages {
  flex: 1;
  overflow-y: auto;
  padding: 1.25rem 1rem;
  display: flex;
  flex-direction: column;
  gap: 1rem;
  background: var(--color-neutral-50);
  scroll-behavior: smooth;
  min-height: 0;
}

.panel-messages::-webkit-scrollbar {
  width: 4px;
}
.panel-messages::-webkit-scrollbar-track {
  background: transparent;
}
.panel-messages::-webkit-scrollbar-thumb {
  background: var(--color-neutral-300);
  border-radius: 4px;
}

.msg-row {
  display: flex;
  align-items: flex-end;
  gap: 0.5rem;
}

.msg-row.user {
  flex-direction: row-reverse;
}

.msg-avatar {
  width: 1.6rem;
  height: 1.6rem;
  border-radius: 50%;
  background: var(--gradient-primary);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.msg-avatar svg {
  width: 0.8rem;
  height: 0.8rem;
  color: #fff;
}

.msg-bubble {
  max-width: 80%;
  padding: 0.6rem 0.85rem;
  font-size: var(--font-size-sm);
  line-height: var(--line-height-relaxed);
  word-break: break-word;
  border-radius: var(--radius-xl);
}

.msg-row.assistant .msg-bubble {
  background: var(--color-neutral-0);
  color: var(--color-neutral-800);
  box-shadow: var(--shadow-sm);
  border: 1px solid var(--color-neutral-200);
  border-bottom-left-radius: var(--radius-xs);
}

.msg-row.user .msg-bubble {
  background: var(--gradient-primary);
  color: #fff;
  box-shadow: var(--shadow-primary);
  border-bottom-right-radius: var(--radius-xs);
}

/* Typing */
.typing {
  display: inline-flex;
  gap: 4px;
  align-items: center;
  height: 1.2rem;
  padding: 2px 0;
}
.typing span {
  width: 5px;
  height: 5px;
  border-radius: 50%;
  background: var(--color-neutral-400);
  animation: bounce 1.4s ease-in-out infinite;
}
.typing span:nth-child(2) {
  animation-delay: 0.18s;
}
.typing span:nth-child(3) {
  animation-delay: 0.36s;
}

@keyframes bounce {
  0%,
  60%,
  100% {
    transform: translateY(0);
    opacity: 0.6;
  }
  30% {
    transform: translateY(-5px);
    opacity: 1;
  }
}

/* Markdown */
.md :deep(p) {
  margin: 0 0 0.35rem;
}
.md :deep(p:last-child) {
  margin-bottom: 0;
}
.md :deep(ul),
.md :deep(ol) {
  margin: 0.25rem 0 0.35rem;
  padding-left: 1.2rem;
}
.md :deep(li) {
  margin-bottom: 0.15rem;
}
.md :deep(strong) {
  font-weight: var(--font-weight-semibold);
}
.md :deep(code) {
  background: rgba(0, 0, 0, 0.06);
  padding: 0.1em 0.3em;
  border-radius: 3px;
  font-size: 0.85em;
}

/* ── Footer ── */
.panel-footer {
  padding: 0.75rem 1rem 0.875rem;
  background: var(--color-neutral-0);
  border-top: 1px solid var(--color-neutral-200);
  display: flex;
  flex-direction: column;
  gap: 0.4rem;
  flex-shrink: 0;
}

.input-box {
  display: flex;
  align-items: flex-end;
  gap: 0.5rem;
  background: var(--color-neutral-50);
  border: 1.5px solid var(--color-neutral-200);
  border-radius: var(--radius-lg);
  padding: 0.5rem 0.5rem 0.5rem 0.75rem;
  transition:
    border-color var(--transition-fast),
    box-shadow var(--transition-fast);
}

.input-box:focus-within {
  border-color: var(--color-primary);
  box-shadow: 0 0 0 3px rgba(0, 135, 74, 0.12);
  background: var(--color-neutral-0);
}

.input-box textarea {
  flex: 1;
  resize: none;
  border: none;
  background: transparent;
  color: var(--color-neutral-800);
  font-size: var(--font-size-sm);
  font-family: var(--font-family);
  outline: none;
  max-height: 120px;
  line-height: var(--line-height-relaxed);
  padding: 0.2rem 0;
}

.input-box textarea::placeholder {
  color: var(--color-neutral-400);
}
.input-box textarea:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.send-btn {
  width: 2rem;
  height: 2rem;
  border-radius: var(--radius-md);
  background: var(--color-primary);
  color: #fff;
  border: none;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  transition: all var(--transition-fast);
  box-shadow: var(--shadow-primary);
}

.send-btn:hover:not(:disabled) {
  background: var(--color-primary-hover);
  transform: scale(1.05);
}
.send-btn:disabled {
  opacity: 0.35;
  cursor: not-allowed;
  box-shadow: none;
  transform: none;
}
.send-btn:focus-visible {
  outline: 2px solid var(--color-primary);
  outline-offset: 2px;
}
.send-btn svg {
  width: 0.85rem;
  height: 0.85rem;
}

.footer-hint {
  font-size: 0.62rem;
  color: var(--color-neutral-400);
  text-align: center;
  margin: 0;
}

/* ── Transitions ── */
.panel-enter-active {
  transition:
    transform 0.25s var(--ease-out),
    opacity 0.25s var(--ease-out);
}
.panel-leave-active {
  transition:
    transform 0.18s var(--ease-in),
    opacity 0.18s var(--ease-in);
}
.panel-enter-from,
.panel-leave-to {
  transform: translateX(100%);
  opacity: 0;
}

.msg-enter-active {
  transition:
    opacity 0.2s var(--ease-out),
    transform 0.2s var(--ease-out);
}
.msg-enter-from {
  opacity: 0;
  transform: translateY(6px);
}

@media (prefers-reduced-motion: reduce) {
  .panel-enter-active,
  .panel-leave-active,
  .msg-enter-active {
    transition: none;
  }
}

@media (max-width: 480px) {
  .field-panel {
    width: 100%;
  }
}
</style>
