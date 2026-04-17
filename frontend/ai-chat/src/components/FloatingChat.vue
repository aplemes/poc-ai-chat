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
  { code: 'pt', label: 'PT', placeholder: 'Descreva sua demanda...', emptyState: 'Me conte o que você precisa criar e vou te ajudar a preencher o formulário.' },
  { code: 'en', label: 'EN', placeholder: 'Describe your demand...', emptyState: 'Tell me what you need to create and I will help you fill in the form.' },
  { code: 'es', label: 'ES', placeholder: 'Describe tu demanda...', emptyState: 'Cuéntame qué necesitas crear y te ayudaré a rellenar el formulario.' },
  { code: 'fr', label: 'FR', placeholder: 'Décrivez votre demande...', emptyState: 'Dites-moi ce que vous souhaitez créer et je vous aiderai à remplir le formulaire.' },
]

const open = ref(false)
const input = ref('')
const messages = ref<Message[]>([])
const messagesEl = ref<HTMLElement | null>(null)
const loading = ref(false)
const sessionId = ref<string | null>(localStorage.getItem('chat_session_id'))
const language = ref<string>(localStorage.getItem('chat_language') ?? 'en')
let abortController: AbortController | null = null
let nextId = 1

function setLanguage(code: string) {
  language.value = code
  localStorage.setItem('chat_language', code)
}

function toggle() {
  open.value = !open.value
}

async function scrollToBottom() {
  await nextTick()
  if (messagesEl.value) {
    messagesEl.value.scrollTop = messagesEl.value.scrollHeight
  }
}

async function send() {
  const text = input.value.trim()
  if (!text || loading.value) return

  messages.value.push({ id: nextId++, role: 'user', text })
  input.value = ''
  loading.value = true
  await scrollToBottom()

  // Placeholder for streaming assistant message
  const assistantMsg: Message = { id: nextId++, role: 'assistant', text: '' }
  messages.value.push(assistantMsg)

  abortController = new AbortController()

  try {
    const newId = await sendMessage(sessionId.value, text, language.value, (event) => {
      if (event.type === 'token' && event.content) {
        assistantMsg.text += event.content
        scrollToBottom()
      } else if (event.type === 'form_fill' && event.data) {
        chatStore.setFormFill(event.data)
        if (!assistantMsg.text) {
          assistantMsg.text = 'Formulário preenchido! Confira os campos acima.'
        }
        scrollToBottom()
      } else if (event.type === 'error') {
        assistantMsg.text = 'Erro ao processar sua mensagem. Tente novamente.'
      }
    }, abortController.signal)
    sessionId.value = newId
    localStorage.setItem('chat_session_id', newId)
  } catch (err) {
    if (err instanceof Error && err.name === 'AbortError') {
      assistantMsg.text = '(cancelado)'
    } else {
      assistantMsg.text = 'Não foi possível conectar ao servidor.'
    }
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
    <Transition name="chat">
      <div v-if="open" class="chat-window">
        <div class="chat-header">
          <span>Demand Assistant</span>
          <div class="lang-selector">
            <button
              v-for="lang in languages"
              :key="lang.code"
              class="lang-btn"
              :class="{ 'lang-btn--active': language === lang.code }"
              @click="setLanguage(lang.code)"
            >
              {{ lang.label }}
            </button>
          </div>
          <button class="close-btn" @click="toggle" aria-label="Close chat">
            <svg viewBox="0 0 20 20" fill="currentColor">
              <path d="M6.28 5.22a.75.75 0 00-1.06 1.06L8.94 10l-3.72 3.72a.75.75 0 101.06 1.06L10 11.06l3.72 3.72a.75.75 0 101.06-1.06L11.06 10l3.72-3.72a.75.75 0 00-1.06-1.06L10 8.94 6.28 5.22z" />
            </svg>
          </button>
        </div>

        <div ref="messagesEl" class="chat-messages">
          <div v-if="messages.length === 0" class="empty-state">
            {{ languages.find(l => l.code === language)?.emptyState }}
          </div>
          <div
            v-for="msg in messages"
            :key="msg.id"
            class="message"
            :class="msg.role"
          >
            <span v-if="msg.role === 'assistant' && !msg.text && loading" class="typing-indicator">
              <span></span><span></span><span></span>
            </span>
            <span v-else-if="msg.role === 'assistant'" class="md" v-html="renderMarkdown(msg.text)" />
            <template v-else>{{ msg.text }}</template>
          </div>
        </div>

        <div class="chat-input-area">
          <textarea
            v-model="input"
            rows="1"
            :placeholder="languages.find(l => l.code === language)?.placeholder"
            :disabled="loading"
            @keydown="onKeydown"
          />
          <button class="send-btn" :disabled="!input.trim() || loading" @click="send" aria-label="Enviar">
            <svg viewBox="0 0 20 20" fill="currentColor">
              <path d="M3.105 2.288a.75.75 0 00-.826.95l1.337 4.01a.75.75 0 00.593.518l5.662.944-5.662.944a.75.75 0 00-.593.519l-1.337 4.01a.75.75 0 00.826.95 19.955 19.955 0 0016.233-8.568.75.75 0 000-.904A19.955 19.955 0 003.105 2.288z" />
            </svg>
          </button>
        </div>
      </div>
    </Transition>

    <button class="fab" @click="toggle" aria-label="Abrir chat">
      <Transition name="icon" mode="out-in">
        <svg v-if="!open" key="chat" viewBox="0 0 24 24" fill="currentColor">
          <path d="M4.913 2.658c2.075-.27 4.19-.408 6.337-.408 2.147 0 4.262.139 6.337.408 1.922.25 3.291 1.861 3.405 3.727a4.403 4.403 0 00-1.032-.211 50.89 50.89 0 00-8.42 0c-2.358.196-4.04 2.19-4.04 4.434v4.286a4.47 4.47 0 002.433 3.984L7.28 21.53A.75.75 0 016 21v-4.03a48.527 48.527 0 01-1.087-.128C2.905 16.58 1.5 14.833 1.5 12.862V6.638c0-1.97 1.405-3.718 3.413-3.979z" />
          <path d="M15.75 7.5c-1.376 0-2.739.057-4.086.169C10.124 7.797 9 9.103 9 10.609v4.285c0 1.507 1.128 2.814 2.67 2.94 1.243.102 2.5.157 3.768.165l2.782 2.781a.75.75 0 001.28-.53v-2.39l.33-.026c1.542-.125 2.67-1.433 2.67-2.94v-4.286c0-1.505-1.125-2.811-2.664-2.94A49.392 49.392 0 0015.75 7.5z" />
        </svg>
        <svg v-else key="close" viewBox="0 0 24 24" fill="currentColor">
          <path fill-rule="evenodd" d="M5.47 5.47a.75.75 0 011.06 0L12 10.94l5.47-5.47a.75.75 0 111.06 1.06L13.06 12l5.47 5.47a.75.75 0 11-1.06 1.06L12 13.06l-5.47 5.47a.75.75 0 01-1.06-1.06L10.94 12 5.47 6.53a.75.75 0 010-1.06z" clip-rule="evenodd" />
        </svg>
      </Transition>
    </button>
  </div>
</template>

<style scoped>
.floating-chat {
  position: fixed;
  bottom: 1.5rem;
  right: 1.5rem;
  z-index: 100;
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 0.75rem;
}

.fab {
  width: 3.25rem;
  height: 3.25rem;
  border-radius: 50%;
  background: hsla(160, 100%, 37%, 1);
  color: #fff;
  border: none;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.2);
  transition: background 0.2s, transform 0.15s;
  flex-shrink: 0;
}

.fab:hover {
  background: hsla(160, 100%, 30%, 1);
  transform: scale(1.05);
}

.fab svg {
  width: 1.4rem;
  height: 1.4rem;
}

.chat-window {
  width: 560px;
  height: 700px;
  max-height: calc(100vh - 6rem);
  display: flex;
  flex-direction: column;
  background: var(--color-background-soft);
  border: 1px solid var(--color-border);
  border-radius: 12px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.15);
  overflow: hidden;
}

.chat-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0.75rem 1rem;
  background: hsla(160, 100%, 37%, 1);
  color: #fff;
  font-weight: 600;
  font-size: 0.9rem;
}

.lang-selector {
  display: flex;
  gap: 0.25rem;
}

.lang-btn {
  background: rgba(255, 255, 255, 0.15);
  border: 1px solid rgba(255, 255, 255, 0.25);
  color: #fff;
  cursor: pointer;
  font-size: 0.7rem;
  font-weight: 600;
  padding: 0.15rem 0.4rem;
  border-radius: 4px;
  opacity: 0.7;
  transition: opacity 0.15s, background 0.15s;
}

.lang-btn:hover {
  opacity: 1;
}

.lang-btn--active {
  background: rgba(255, 255, 255, 0.3);
  border-color: rgba(255, 255, 255, 0.6);
  opacity: 1;
}

.close-btn {
  background: none;
  border: none;
  color: #fff;
  cursor: pointer;
  display: flex;
  padding: 0.1rem;
  border-radius: 4px;
  opacity: 0.8;
  transition: opacity 0.15s;
}

.close-btn:hover {
  opacity: 1;
}

.close-btn svg {
  width: 1.1rem;
  height: 1.1rem;
}

.chat-messages {
  flex: 1;
  overflow-y: auto;
  padding: 0.75rem;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  min-height: 0;
}

.empty-state {
  text-align: center;
  color: var(--color-text);
  opacity: 0.4;
  font-size: 0.85rem;
  margin: auto;
  padding: 0 1rem;
  line-height: 1.5;
}

.message {
  max-width: 85%;
  padding: 0.6rem 0.9rem;
  border-radius: 12px;
  font-size: 0.9rem;
  line-height: 1.7;
  word-break: break-word;
}

.message.user {
  align-self: flex-end;
  background: hsla(160, 100%, 37%, 1);
  color: #fff;
  border-bottom-right-radius: 4px;
}

.message.assistant {
  align-self: flex-start;
  background: var(--color-background);
  color: var(--color-text);
  border: 1px solid var(--color-border);
  border-bottom-left-radius: 4px;
}

.typing-indicator {
  display: inline-flex;
  gap: 3px;
  align-items: center;
  padding: 2px 0;
}

.typing-indicator span {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: var(--color-text);
  opacity: 0.4;
  animation: bounce 1.2s infinite;
}

.typing-indicator span:nth-child(2) { animation-delay: 0.2s; }
.typing-indicator span:nth-child(3) { animation-delay: 0.4s; }

@keyframes bounce {
  0%, 60%, 100% { transform: translateY(0); }
  30% { transform: translateY(-4px); }
}

.chat-input-area {
  display: flex;
  align-items: flex-end;
  gap: 0.5rem;
  padding: 0.6rem 0.75rem;
  border-top: 1px solid var(--color-border);
  background: var(--color-background);
}

.chat-input-area textarea {
  flex: 1;
  resize: none;
  border: 1px solid var(--color-border);
  border-radius: 8px;
  padding: 0.5rem 0.75rem;
  font-size: 0.9rem;
  background: var(--color-background-soft);
  color: var(--color-text);
  outline: none;
  max-height: 120px;
  line-height: 1.6;
  transition: border-color 0.2s;
}

.chat-input-area textarea:focus {
  border-color: hsla(160, 100%, 37%, 1);
}

.chat-input-area textarea:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.send-btn {
  width: 2rem;
  height: 2rem;
  border-radius: 8px;
  background: hsla(160, 100%, 37%, 1);
  color: #fff;
  border: none;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  transition: background 0.2s;
}

.send-btn:hover:not(:disabled) {
  background: hsla(160, 100%, 30%, 1);
}

.send-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.send-btn svg {
  width: 1rem;
  height: 1rem;
}

/* Transitions */
.chat-enter-active,
.chat-leave-active {
  transition: opacity 0.2s, transform 0.2s;
}

.chat-enter-from,
.chat-leave-to {
  opacity: 0;
  transform: translateY(12px) scale(0.97);
}

.icon-enter-active,
.icon-leave-active {
  transition: opacity 0.15s, transform 0.15s;
}

.icon-enter-from,
.icon-leave-to {
  opacity: 0;
  transform: rotate(30deg) scale(0.8);
}

.md :deep(p) {
  margin: 0 0 0.4rem;
}
.md :deep(p:last-child) {
  margin-bottom: 0;
}
.md :deep(ol),
.md :deep(ul) {
  margin: 0.3rem 0 0.4rem;
  padding-left: 1.2rem;
}
.md :deep(li) {
  margin-bottom: 0.2rem;
}
.md :deep(strong) {
  font-weight: 600;
}

@media (max-width: 480px) {
  .chat-window {
    width: calc(100vw - 3rem);
    height: calc(100vh - 6rem);
  }
}
</style>
