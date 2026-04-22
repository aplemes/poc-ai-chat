<script setup lang="ts">
import { ref, nextTick, onUnmounted } from 'vue'
import { sendMessage, confirmForm } from '@/services/chatService'
import type { ChatEvent } from '@/services/chatService'
import type { ChatMessage } from '@/components/chat/ChatMessageList.vue'
import { useChatStore } from '@/stores/chat'
import { useLanguage } from '@/composables/useLanguage'
import { useChatStream } from '@/composables/useChatStream'
import ChatFab from '@/components/chat/ChatFab.vue'
import ChatMessageList from '@/components/chat/ChatMessageList.vue'
import { chatI18n, langNames } from '@/data/chatI18n'

const chatStore = useChatStore()

const open = ref(false)
const messagesEl = ref<HTMLElement | null>(null)
const textareaEl = ref<HTMLTextAreaElement | null>(null)
const sessionId = ref<string | null>(localStorage.getItem('chat_session_id'))
const hasNewMessage = ref(false)
let confirmAbortController: AbortController | null = null

const { language, setLanguage: setLang, languages } = useLanguage()

function currentI18n() {
  return chatI18n.find((l) => l.code === language.value) ?? chatI18n[1]!
}

function setLanguage(code: string) {
  setLang(code)
  const stream = chatStream
  if (stream.messages.value.length === 1 && stream.messages.value[0]?.role === 'assistant') {
    const strings = chatI18n.find((l) => l.code === code) ?? chatI18n[1]!
    stream.messages.value[0].text = strings.greeting
  }
}

const chatStream = useChatStream<ChatEvent>({
  messagesEl: () => messagesEl.value,
  textareaEl: () => textareaEl.value,
  onEvent(event, assistantMsg, { addMessage, scrollToBottom }) {
    if (event.type === 'token' && event.content) {
      assistantMsg.text += event.content
      scrollToBottom()
    } else if (event.type === 'form_confirm' && event.data) {
      chatStream.messages.value = chatStream.messages.value.filter((m) => m.id !== assistantMsg.id)
      addMessage({
        role: 'assistant',
        text: '',
        type: 'confirm',
        confirmData: event.data,
        confirmed: false,
      })
      if (!open.value) hasNewMessage.value = true
      scrollToBottom()
    } else if (event.type === 'form_fill' && event.data) {
      chatStore.setFormFill(event.data)
      if (!assistantMsg.text) assistantMsg.text = currentI18n().formFilled
      addMessage({ role: 'assistant', text: currentI18n().correctionHint })
      if (!open.value) hasNewMessage.value = true
      scrollToBottom()
    } else if (event.type === 'error') {
      assistantMsg.text = currentI18n().errorProcessing
    }
  },
  sendFn: async (text, lang, onEvent, signal) => {
    try {
      const newId = await sendMessage(sessionId.value, text, lang, onEvent, signal)
      if (newId) {
        sessionId.value = newId
        localStorage.setItem('chat_session_id', newId)
      }
      return newId
    } catch (err) {
      const lastMsg = chatStream.messages.value[chatStream.messages.value.length - 1]
      if (lastMsg?.role === 'assistant') {
        lastMsg.text =
          err instanceof Error && err.name === 'AbortError'
            ? currentI18n().cancelled
            : currentI18n().errorConnection
      }
      return null
    }
  },
})

function toggle() {
  open.value = !open.value
  if (open.value) {
    hasNewMessage.value = false
    if (chatStream.messages.value.length === 0) {
      chatStream.addMessage({ role: 'assistant', text: currentI18n().greeting })
    }
    nextTick(() => textareaEl.value?.focus())
  }
}

async function send() {
  await chatStream.send(language.value)
}

async function handleConfirm(msg: ChatMessage) {
  if (!sessionId.value || msg.confirmed) return
  msg.confirmed = true
  chatStream.loading.value = true

  confirmAbortController?.abort()
  confirmAbortController = new AbortController()

  try {
    await confirmForm(
      sessionId.value,
      (event) => {
        if (event.type === 'form_fill' && event.data) {
          chatStore.setFormFill(event.data)
          chatStream.addMessage({ role: 'assistant', text: currentI18n().formFilled })
          chatStream.addMessage({ role: 'assistant', text: currentI18n().correctionHint })
          if (!open.value) hasNewMessage.value = true
          chatStream.scrollToBottom()
        }
      },
      confirmAbortController.signal,
    )
  } catch (err) {
    msg.confirmed = false
    if (!(err instanceof Error && err.name === 'AbortError')) {
      chatStream.addMessage({ role: 'assistant', text: currentI18n().errorProcessing })
      chatStream.scrollToBottom()
    }
  } finally {
    chatStream.loading.value = false
    confirmAbortController = null
  }
}

function handleCorrect() {
  nextTick(() => textareaEl.value?.focus())
}

onUnmounted(() => confirmAbortController?.abort())

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
              <span class="brand-status"> <span class="status-dot"></span>Online </span>
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
              >
                {{ lang.label }}
              </button>
            </div>
            <button class="icon-btn" :aria-label="currentI18n().ariaCloseChat" @click="toggle">
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
        <div ref="messagesEl" class="chat-messages">
          <ChatMessageList
            :messages="chatStream.messages.value"
            :loading="chatStream.loading.value"
            :confirm-title="currentI18n().confirmTitle"
            :confirm-btn="currentI18n().confirmBtn"
            :correct-btn="currentI18n().correctBtn"
            :form-filled="currentI18n().formFilled"
            :conversation-label="currentI18n().ariaConversation"
            @confirm="handleConfirm"
            @correct="handleCorrect"
          />
        </div>

        <!-- Input -->
        <div class="chat-footer">
          <div class="input-box" :class="{ focused: chatStream.loading.value === false }">
            <textarea
              ref="textareaEl"
              v-model="chatStream.input.value"
              rows="1"
              :placeholder="currentI18n().placeholder"
              :disabled="chatStream.loading.value"
              @keydown="onKeydown"
              @input="chatStream.autoResize()"
            />
            <button
              class="send-btn"
              :disabled="!chatStream.input.value.trim() || chatStream.loading.value"
              :aria-label="currentI18n().ariaSend"
              @click="send"
            >
              <svg viewBox="0 0 20 20" fill="currentColor">
                <path
                  d="M3.105 2.288a.75.75 0 00-.826.95l1.337 4.01a.75.75 0 00.593.518l5.662.944-5.662.944a.75.75 0 00-.593.519l-1.337 4.01a.75.75 0 00.826.95 19.955 19.955 0 0016.233-8.568.75.75 0 000-.904A19.955 19.955 0 003.105 2.288z"
                />
              </svg>
            </button>
          </div>
          <p class="footer-hint">{{ currentI18n().footerHint }}</p>
        </div>
      </div>
    </Transition>

    <!-- FAB -->
    <ChatFab
      :open="open"
      :has-new-message="hasNewMessage"
      :label="currentI18n().ariaOpenChat"
      @click="toggle"
    />
  </div>
</template>

<style src="@/components/chat/FloatingChat.css"></style>
