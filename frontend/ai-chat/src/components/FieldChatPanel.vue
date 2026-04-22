<script setup lang="ts">
import { ref, watch, nextTick, onBeforeUnmount, computed } from 'vue'
import { useFieldChatStore } from '@/stores/fieldChat'
import { sendFieldMessage } from '@/services/fieldChatService'
import type { FieldChatEvent } from '@/services/fieldChatService'
import { renderMarkdown } from '@/utils/markdown'
import { useLanguage } from '@/composables/useLanguage'
import { useChatStream } from '@/composables/useChatStream'
import { fieldGreetings } from '@/data/fieldGreetings'
import { fieldChatI18n } from '@/data/fieldChatI18n'

const fieldChatStore = useFieldChatStore()

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

function getGreeting(field: string, lang: string): string {
  const byLang = fieldGreetings[lang] ?? fieldGreetings['en']
  return (
    byLang?.[field] ??
    `I'll help you fill the ${fieldLabels[field] ?? field} field. What information do you have?`
  )
}

const panelRef = ref<HTMLElement | null>(null)
const messagesEl = ref<HTMLElement | null>(null)
const textareaEl = ref<HTMLTextAreaElement | null>(null)
let previouslyFocused: HTMLElement | null = null

const FOCUSABLE =
  'a[href], button:not([disabled]), input, select, textarea, [tabindex]:not([tabindex="-1"])'

function trapFocus(e: KeyboardEvent) {
  if (e.key !== 'Tab' || !panelRef.value) return
  const focusable = Array.from(panelRef.value.querySelectorAll<HTMLElement>(FOCUSABLE))
  const first = focusable[0]
  const last = focusable[focusable.length - 1]
  if (!first || !last) return
  if (e.shiftKey) {
    if (document.activeElement === first) {
      e.preventDefault()
      last.focus()
    }
  } else {
    if (document.activeElement === last) {
      e.preventDefault()
      first.focus()
    }
  }
}

const { language, setLanguage, languages } = useLanguage()

function i18n() {
  return fieldChatI18n[language.value] ?? fieldChatI18n['en']!
}

const activeField = computed(() => fieldChatStore.activeField)
const isOpen = computed(() => activeField.value !== null)
const fieldLabel = computed(() =>
  activeField.value ? (fieldLabels[activeField.value] ?? activeField.value) : '',
)

const fieldFilled = ref(false)

const chatStream = useChatStream<FieldChatEvent>({
  messagesEl: () => messagesEl.value,
  textareaEl: () => textareaEl.value,
  onEvent(event, assistantMsg, { scrollToBottom }) {
    if (event.type === 'token' && event.content) {
      assistantMsg.text += event.content
      scrollToBottom()
    } else if (event.type === 'field_fill' && event.data) {
      fieldChatStore.setFieldFill({ fieldName: event.data.fieldName, value: event.data.value })
      if (!assistantMsg.text) {
        assistantMsg.text = `${i18n().fieldFilledPrefix} **${event.data.value}**`
      }
      fieldFilled.value = true
      scrollToBottom()
    } else if (event.type === 'error') {
      assistantMsg.text = i18n().errorProcessing
    }
  },
  sendFn: async (text, lang, onEvent, signal) => {
    const field = activeField.value
    if (!field) return null
    try {
      const sessionId = fieldChatStore.getSessionId(field)
      const newId = await sendFieldMessage(sessionId, field, text, lang, onEvent, signal)
      if (newId) fieldChatStore.setSessionId(field, newId)
      return newId
    } catch (err) {
      const msgs = chatStream.messages.value
      const lastMsg = msgs[msgs.length - 1]
      if (lastMsg?.role === 'assistant') {
        lastMsg.text =
          err instanceof Error && err.name === 'AbortError'
            ? i18n().cancelled
            : i18n().connectionError
      }
      return null
    }
  },
})

watch(isOpen, (open) => {
  if (open) {
    previouslyFocused = document.activeElement as HTMLElement
    document.addEventListener('keydown', trapFocus)
  } else {
    document.removeEventListener('keydown', trapFocus)
    previouslyFocused?.focus()
  }
})

onBeforeUnmount(() => document.removeEventListener('keydown', trapFocus))

watch(activeField, (field) => {
  if (!field) return
  fieldFilled.value = false
  chatStream.abort()
  chatStream.loading.value = false
  chatStream.input.value = ''
  chatStream.messages.value = [
    { id: 1, role: 'assistant', text: getGreeting(field, language.value) },
  ]
  nextTick(() => textareaEl.value?.focus())
})

watch(language, (lang) => {
  if (!activeField.value || chatStream.messages.value.length !== 1) return
  const first = chatStream.messages.value[0]
  if (first) first.text = getGreeting(activeField.value, lang)
})

async function send() {
  await chatStream.send(language.value)
}

function onKeydown(e: KeyboardEvent) {
  if (e.key === 'Enter' && !e.shiftKey) {
    e.preventDefault()
    send()
  }
  if (e.key === 'Escape') fieldChatStore.closePanel()
}
</script>

<template>
  <Teleport to="body">
    <Transition name="panel">
      <div v-if="isOpen" class="field-panel-backdrop" @click.self="fieldChatStore.closePanel()">
        <div
          ref="panelRef"
          class="field-panel"
          role="dialog"
          aria-modal="true"
          :aria-label="`AI assistant for ${fieldLabel}`"
          tabindex="-1"
        >
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
              <div
                v-for="msg in chatStream.messages.value"
                :key="msg.id"
                class="msg-row"
                :class="msg.role"
              >
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
                  <span
                    v-if="msg.role === 'assistant' && !msg.text && chatStream.loading.value"
                    class="typing"
                  >
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
            <Transition name="accept">
              <button
                v-if="fieldFilled"
                class="accept-close-btn"
                @click="fieldChatStore.closePanel()"
              >
                <svg viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
                  <path
                    fill-rule="evenodd"
                    d="M16.707 5.293a1 1 0 010 1.414L8.414 15l-4.121-4.121a1 1 0 011.414-1.414L8.414 12.172l6.879-6.879a1 1 0 011.414 0z"
                    clip-rule="evenodd"
                  />
                </svg>
                {{ i18n().acceptClose }}
              </button>
            </Transition>
            <div class="input-box">
              <textarea
                ref="textareaEl"
                v-model="chatStream.input.value"
                rows="1"
                :placeholder="i18n().placeholder"
                :disabled="chatStream.loading.value"
                @keydown="onKeydown"
                @input="chatStream.autoResize()"
              />
              <button
                class="send-btn"
                :disabled="!chatStream.input.value.trim() || chatStream.loading.value"
                :aria-label="i18n().ariaSend"
                @click="send"
              >
                <svg viewBox="0 0 20 20" fill="currentColor">
                  <path
                    d="M3.105 2.288a.75.75 0 00-.826.95l1.337 4.01a.75.75 0 00.593.518l5.662.944-5.662.944a.75.75 0 00-.593.519l-1.337 4.01a.75.75 0 00.826.95 19.955 19.955 0 0016.233-8.568.75.75 0 000-.904A19.955 19.955 0 003.105 2.288z"
                  />
                </svg>
              </button>
            </div>
            <p class="footer-hint">{{ i18n().footerHint }}</p>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style src="@/components/chat/FieldChatPanel.css"></style>

<style scoped>
/* ── Markdown (uses :deep — must remain scoped) ── */
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

/* ── Accept & close button ── */
.accept-close-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--space-2);
  width: 100%;
  padding: var(--space-2) var(--space-4);
  background: var(--color-primary);
  color: #fff;
  border: none;
  border-radius: var(--radius-md);
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-semibold);
  font-family: var(--font-family);
  cursor: pointer;
  box-shadow: var(--shadow-primary);
  transition: background var(--transition-fast);
}

.accept-close-btn:hover {
  background: var(--color-primary-hover);
}

.accept-close-btn svg {
  width: 0.875rem;
  height: 0.875rem;
  flex-shrink: 0;
}

.accept-enter-active,
.accept-leave-active {
  transition:
    opacity var(--transition-fast),
    transform var(--transition-fast);
}
.accept-enter-from,
.accept-leave-to {
  opacity: 0;
  transform: translateY(-4px);
}
</style>
