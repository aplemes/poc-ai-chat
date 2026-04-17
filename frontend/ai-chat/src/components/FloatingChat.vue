<script setup lang="ts">
import { ref, nextTick, onUnmounted } from 'vue'
import { marked } from 'marked'
import DOMPurify from 'dompurify'
import { sendMessage } from '@/services/chatService'
import { useChatStore } from '@/stores/chat'

function renderMarkdown(text: string): string {
  return DOMPurify.sanitize(marked.parse(text) as string)
}

interface Message {
  id: number
  role: 'user' | 'assistant'
  text: string
}

const chatStore = useChatStore()

const languages = [
  { code: 'pt', label: 'PT', placeholder: 'Escreva sua mensagem...', greeting: 'Olá! Como posso ajudar você hoje?', formFilled: 'Formulário preenchido! Revise os campos.', errorProcessing: 'Não consegui processar. Tente novamente.', cancelled: '(cancelado)', errorConnection: 'Sem conexão com o servidor.', footerHint: 'Enter para enviar · Shift+Enter nova linha' },
  { code: 'en', label: 'EN', placeholder: 'Write your message...', greeting: 'Hello! How can I help you today?', formFilled: 'Form filled! Please review the fields.', errorProcessing: 'Could not process your message. Try again.', cancelled: '(cancelled)', errorConnection: 'Could not connect to the server.', footerHint: 'Enter to send · Shift+Enter for new line' },
  { code: 'es', label: 'ES', placeholder: 'Escribe tu mensaje...', greeting: '¡Hola! ¿Cómo puedo ayudarte hoy?', formFilled: '¡Formulario completado! Revisa los campos.', errorProcessing: 'No se pudo procesar. Inténtalo de nuevo.', cancelled: '(cancelado)', errorConnection: 'No se pudo conectar al servidor.', footerHint: 'Enter para enviar · Shift+Enter nueva línea' },
  { code: 'fr', label: 'FR', placeholder: 'Écrivez votre message...', greeting: 'Bonjour ! Comment puis-je vous aider ?', formFilled: 'Formulaire rempli ! Vérifiez les champs.', errorProcessing: 'Impossible de traiter. Réessayez.', cancelled: '(annulé)', errorConnection: 'Connexion au serveur impossible.', footerHint: 'Entrée pour envoyer · Maj+Entrée nouvelle ligne' },
]

const langNames: Record<string, string> = { pt: 'Português', en: 'English', es: 'Español', fr: 'Français' }

const open = ref(false)
const input = ref('')
const messages = ref<Message[]>([])
const messagesEl = ref<HTMLElement | null>(null)
const textareaEl = ref<HTMLTextAreaElement | null>(null)
const loading = ref(false)
const sessionId = ref<string | null>(localStorage.getItem('chat_session_id'))
const storedLang = localStorage.getItem('chat_language')
const validCodes = languages.map((l) => l.code)
const language = ref<string>(storedLang && validCodes.includes(storedLang) ? storedLang : 'en')
const hasNewMessage = ref(false)
let abortController: AbortController | null = null
let nextId = 1

function currentLang() {
  return languages.find((l) => l.code === language.value) ?? languages[1]!
}

function setLanguage(code: string) {
  language.value = code
  localStorage.setItem('chat_language', code)
  if (messages.value.length === 1 && messages.value[0]?.role === 'assistant') {
    const lang = languages.find((l) => l.code === code) ?? languages[1]!
    messages.value[0].text = lang.greeting
  }
}

function toggle() {
  open.value = !open.value
  if (open.value) {
    hasNewMessage.value = false
    if (messages.value.length === 0) {
      messages.value.push({ id: nextId++, role: 'assistant', text: currentLang().greeting })
    }
    nextTick(() => textareaEl.value?.focus())
  }
}

async function scrollToBottom() {
  await nextTick()
  if (messagesEl.value) {
    messagesEl.value.scrollTop = messagesEl.value.scrollHeight
  }
}

function autoResize() {
  const el = textareaEl.value
  if (!el) return
  el.style.height = 'auto'
  el.style.height = Math.min(el.scrollHeight, 140) + 'px'
}

async function send() {
  const text = input.value.trim()
  if (!text || loading.value) return

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
    const newId = await sendMessage(sessionId.value, text, language.value, (event) => {
      if (event.type === 'token' && event.content) {
        assistantMsg.text += event.content
        scrollToBottom()
      } else if (event.type === 'form_fill' && event.data) {
        chatStore.setFormFill(event.data)
        if (!assistantMsg.text) assistantMsg.text = currentLang().formFilled
        if (!open.value) hasNewMessage.value = true
        scrollToBottom()
      } else if (event.type === 'error') {
        assistantMsg.text = currentLang().errorProcessing
      }
    }, abortController.signal)
    if (newId) {
      sessionId.value = newId
      localStorage.setItem('chat_session_id', newId)
    }
  } catch (err) {
    assistantMsg.text = err instanceof Error && err.name === 'AbortError'
      ? currentLang().cancelled
      : currentLang().errorConnection
  } finally {
    loading.value = false
    abortController = null
  }
}

onUnmounted(() => abortController?.abort())

function onKeydown(e: KeyboardEvent) {
  if (e.key === 'Enter' && !e.shiftKey) {
    e.preventDefault()
    send()
  }
}
</script>

<template>
  <div class="floating-chat">
    <Transition name="panel">
      <div v-if="open" class="chat-panel">

        <!-- Header -->
        <div class="chat-header">
          <div class="header-brand">
            <div class="brand-avatar">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M12 2a10 10 0 0 1 10 10c0 5.52-4.48 10-10 10S2 17.52 2 12 6.48 2 12 2z" />
                <path d="M8 12h.01M12 12h.01M16 12h.01" stroke-linecap="round" stroke-width="2.5" />
              </svg>
            </div>
            <div class="brand-info">
              <span class="brand-name">Demand Assistant</span>
              <span class="brand-status">
                <span class="status-dot"></span>Online
              </span>
            </div>
          </div>
          <div class="header-actions">
            <div class="lang-selector">
              <button
                v-for="lang in languages"
                :key="lang.code"
                class="lang-btn"
                :class="{ active: language === lang.code }"
                :aria-label="langNames[lang.code]"
                @click="setLanguage(lang.code)"
              >{{ lang.label }}</button>
            </div>
            <button class="icon-btn" aria-label="Close chat" @click="toggle">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round">
                <path d="M18 6L6 18M6 6l12 12" />
              </svg>
            </button>
          </div>
        </div>

        <!-- Messages -->
        <div ref="messagesEl" class="chat-messages" role="log" aria-live="polite" aria-label="Conversa">
          <TransitionGroup name="msg">
            <div v-for="msg in messages" :key="msg.id" class="msg-row" :class="msg.role">
              <div v-if="msg.role === 'assistant'" class="msg-avatar">
                <svg viewBox="0 0 20 20" fill="currentColor">
                  <path d="M10 2a8 8 0 100 16A8 8 0 0010 2zm0 2a6 6 0 110 12A6 6 0 0110 4z" opacity=".3"/>
                  <path d="M6.5 10a1 1 0 112 0 1 1 0 01-2 0zm5 0a1 1 0 112 0 1 1 0 01-2 0z"/>
                </svg>
              </div>
              <div class="msg-bubble">
                <span v-if="msg.role === 'assistant' && !msg.text && loading" class="typing">
                  <span></span><span></span><span></span>
                </span>
                <span v-else-if="msg.role === 'assistant'" class="md" v-html="renderMarkdown(msg.text)" />
                <template v-else>{{ msg.text }}</template>
              </div>
            </div>
          </TransitionGroup>
        </div>

        <!-- Input -->
        <div class="chat-footer">
          <div class="input-box" :class="{ focused: loading === false }">
            <textarea
              ref="textareaEl"
              v-model="input"
              rows="1"
              :placeholder="currentLang().placeholder"
              :disabled="loading"
              @keydown="onKeydown"
              @input="autoResize"
            />
            <button
              class="send-btn"
              :disabled="!input.trim() || loading"
              aria-label="Enviar"
              @click="send"
            >
              <svg viewBox="0 0 20 20" fill="currentColor">
                <path d="M3.105 2.288a.75.75 0 00-.826.95l1.337 4.01a.75.75 0 00.593.518l5.662.944-5.662.944a.75.75 0 00-.593.519l-1.337 4.01a.75.75 0 00.826.95 19.955 19.955 0 0016.233-8.568.75.75 0 000-.904A19.955 19.955 0 003.105 2.288z" />
              </svg>
            </button>
          </div>
          <p class="footer-hint">{{ currentLang().footerHint }}</p>
        </div>

      </div>
    </Transition>

    <!-- FAB -->
    <button class="fab" :class="{ open }" aria-label="Abrir assistente" @click="toggle">
      <Transition name="icon" mode="out-in">
        <svg v-if="!open" key="open" viewBox="0 0 24 24" fill="currentColor">
          <path d="M4.913 2.658c2.075-.27 4.19-.408 6.337-.408 2.147 0 4.262.139 6.337.408 1.922.25 3.291 1.861 3.405 3.727a4.403 4.403 0 00-1.032-.211 50.89 50.89 0 00-8.42 0c-2.358.196-4.04 2.19-4.04 4.434v4.286a4.47 4.47 0 002.433 3.984L7.28 21.53A.75.75 0 016 21v-4.03a48.527 48.527 0 01-1.087-.128C2.905 16.58 1.5 14.833 1.5 12.862V6.638c0-1.97 1.405-3.718 3.413-3.979z" />
          <path d="M15.75 7.5c-1.376 0-2.739.057-4.086.169C10.124 7.797 9 9.103 9 10.609v4.285c0 1.507 1.128 2.814 2.67 2.94 1.243.102 2.5.157 3.768.165l2.782 2.781a.75.75 0 001.28-.53v-2.39l.33-.026c1.542-.125 2.67-1.433 2.67-2.94v-4.286c0-1.505-1.125-2.811-2.664-2.94A49.392 49.392 0 0015.75 7.5z" />
        </svg>
        <svg v-else key="close" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round">
          <path d="M18 6L6 18M6 6l12 12" />
        </svg>
      </Transition>
      <span v-if="hasNewMessage" class="fab-badge" aria-hidden="true"></span>
    </button>
  </div>
</template>

<style scoped>
/* ── Layout ── */
.floating-chat {
  position: fixed;
  bottom: 1.75rem;
  right: 1.75rem;
  z-index: 9000;
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 1rem;
}

/* ── Panel ── */
.chat-panel {
  width: 400px;
  height: 620px;
  max-height: calc(100vh - 7rem);
  display: flex;
  flex-direction: column;
  background: var(--color-neutral-0);
  border-radius: var(--radius-2xl);
  box-shadow: var(--shadow-2xl);
  overflow: hidden;
  border: 1px solid rgba(0, 0, 0, 0.06);
}

/* ── Header ── */
.chat-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 1rem 1.25rem;
  background: var(--gradient-primary);
  gap: 0.75rem;
  flex-shrink: 0;
}

.header-brand {
  display: flex;
  align-items: center;
  gap: 0.65rem;
}

.brand-avatar {
  width: 2.25rem;
  height: 2.25rem;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.2);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  border: 1.5px solid rgba(255, 255, 255, 0.3);
}

.brand-avatar svg {
  width: 1.1rem;
  height: 1.1rem;
  color: #fff;
}

.brand-info {
  display: flex;
  flex-direction: column;
  gap: 1px;
}

.brand-name {
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-semibold);
  color: #fff;
  line-height: 1;
  letter-spacing: -0.01em;
}

.brand-status {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 0.7rem;
  color: rgba(255, 255, 255, 0.8);
}

.status-dot {
  width: 5px;
  height: 5px;
  border-radius: 50%;
  background: #a7f3d0;
  box-shadow: 0 0 0 2px rgba(167, 243, 208, 0.3);
  animation: pulse-dot 2s ease-in-out infinite;
}

@keyframes pulse-dot {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.lang-selector {
  display: flex;
  gap: 2px;
  background: rgba(0, 0, 0, 0.15);
  padding: 3px;
  border-radius: var(--radius-sm);
}

.lang-btn {
  background: transparent;
  border: none;
  color: rgba(255, 255, 255, 0.7);
  cursor: pointer;
  font-size: 0.65rem;
  font-weight: var(--font-weight-semibold);
  padding: 0.2rem 0.35rem;
  border-radius: 3px;
  transition: all var(--transition-fast);
  letter-spacing: 0.03em;
}

.lang-btn:hover { color: #fff; }

.lang-btn.active {
  background: rgba(255, 255, 255, 0.25);
  color: #fff;
}

.lang-btn:focus-visible {
  outline: 2px solid rgba(255, 255, 255, 0.6);
  outline-offset: 1px;
}

.icon-btn {
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
  transition: all var(--transition-fast);
}

.icon-btn:hover { background: rgba(0, 0, 0, 0.22); color: #fff; }

.icon-btn:focus-visible {
  outline: 2px solid rgba(255, 255, 255, 0.6);
  outline-offset: 1px;
}

.icon-btn svg { width: 0.9rem; height: 0.9rem; }

/* ── Messages ── */
.chat-messages {
  flex: 1;
  overflow-y: auto;
  padding: 1.25rem 1rem;
  display: flex;
  flex-direction: column;
  gap: 1rem;
  min-height: 0;
  background: var(--color-neutral-50);
  scroll-behavior: smooth;
}

.chat-messages::-webkit-scrollbar { width: 4px; }
.chat-messages::-webkit-scrollbar-track { background: transparent; }
.chat-messages::-webkit-scrollbar-thumb { background: var(--color-neutral-300); border-radius: 4px; }

.msg-row {
  display: flex;
  align-items: flex-end;
  gap: 0.5rem;
}

.msg-row.user { flex-direction: row-reverse; }

.msg-avatar {
  width: 1.75rem;
  height: 1.75rem;
  border-radius: 50%;
  background: var(--gradient-primary);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.msg-avatar svg { width: 0.85rem; height: 0.85rem; color: #fff; }

.msg-bubble {
  max-width: 78%;
  padding: 0.65rem 0.9rem;
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

/* Typing indicator */
.typing {
  display: inline-flex;
  gap: 4px;
  align-items: center;
  padding: 2px 2px;
  height: 1.2rem;
}

.typing span {
  width: 5px;
  height: 5px;
  border-radius: 50%;
  background: var(--color-neutral-400);
  animation: bounce 1.4s ease-in-out infinite;
}

.typing span:nth-child(2) { animation-delay: 0.18s; }
.typing span:nth-child(3) { animation-delay: 0.36s; }

@keyframes bounce {
  0%, 60%, 100% { transform: translateY(0); opacity: 0.6; }
  30% { transform: translateY(-5px); opacity: 1; }
}

/* ── Footer / Input ── */
.chat-footer {
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
  transition: border-color var(--transition-fast), box-shadow var(--transition-fast);
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
  max-height: 140px;
  line-height: var(--line-height-relaxed);
  padding: 0.2rem 0;
}

.input-box textarea::placeholder { color: var(--color-neutral-400); }
.input-box textarea:disabled { opacity: 0.5; cursor: not-allowed; }

.send-btn {
  width: 2.1rem;
  height: 2.1rem;
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

.send-btn svg { width: 0.9rem; height: 0.9rem; }

.footer-hint {
  font-size: 0.65rem;
  color: var(--color-neutral-400);
  text-align: center;
  margin: 0;
  letter-spacing: 0.01em;
}

/* ── FAB ── */
.fab {
  position: relative;
  width: 3.25rem;
  height: 3.25rem;
  border-radius: 50%;
  background: var(--gradient-primary);
  color: #fff;
  border: none;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: var(--shadow-xl), var(--shadow-primary);
  transition: all var(--transition-normal);
  flex-shrink: 0;
}

.fab:hover { transform: scale(1.07); box-shadow: 0 12px 32px rgba(0,0,0,0.15), 0 6px 16px rgba(0,135,74,0.4); }
.fab.open { box-shadow: var(--shadow-md); }

.fab:focus-visible {
  outline: 2px solid var(--color-primary);
  outline-offset: 3px;
}

.fab svg { width: 1.35rem; height: 1.35rem; }

.fab-badge {
  position: absolute;
  top: 2px;
  right: 2px;
  width: 10px;
  height: 10px;
  border-radius: 50%;
  background: #ef4444;
  border: 2px solid #fff;
  animation: badge-pop 0.3s var(--ease-out);
}

@keyframes badge-pop {
  from { transform: scale(0); opacity: 0; }
  to   { transform: scale(1); opacity: 1; }
}

/* ── Markdown ── */
.md :deep(p) { margin: 0 0 0.35rem; }
.md :deep(p:last-child) { margin-bottom: 0; }
.md :deep(ul), .md :deep(ol) { margin: 0.25rem 0 0.35rem; padding-left: 1.2rem; }
.md :deep(li) { margin-bottom: 0.15rem; }
.md :deep(strong) { font-weight: var(--font-weight-semibold); }
.md :deep(code) {
  background: rgba(0,0,0,0.06);
  padding: 0.1em 0.3em;
  border-radius: 3px;
  font-size: 0.85em;
}

/* ── Transitions ── */
.panel-enter-active { transition: opacity var(--transition-slow), transform var(--transition-slow); }
.panel-leave-active { transition: opacity 0.18s var(--ease-in), transform 0.18s var(--ease-in); }
.panel-enter-from, .panel-leave-to { opacity: 0; transform: translateY(16px) scale(0.96); }

.icon-enter-active, .icon-leave-active { transition: opacity var(--transition-fast), transform var(--transition-fast); }
.icon-enter-from, .icon-leave-to { opacity: 0; transform: rotate(25deg) scale(0.75); }

.msg-enter-active { transition: opacity 0.25s var(--ease-out), transform 0.25s var(--ease-out); }
.msg-enter-from { opacity: 0; transform: translateY(8px); }

@media (prefers-reduced-motion: reduce) {
  .panel-enter-active, .panel-leave-active,
  .icon-enter-active, .icon-leave-active,
  .msg-enter-active,
  .fab, .send-btn, .lang-btn {
    transition: none;
    animation: none;
  }
}

@media (max-width: 480px) {
  .chat-panel {
    width: calc(100vw - 2rem);
    height: calc(100vh - 6rem);
    border-radius: var(--radius-xl);
  }
}
</style>
