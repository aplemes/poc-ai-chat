<script setup lang="ts">
import { renderMarkdown } from '@/utils/markdown'
import ConfirmCard from '@/components/chat/ConfirmCard.vue'
import type { FormFillData } from '@/services/chatService'

export interface ChatMessage {
  id: number
  role: 'user' | 'assistant'
  text: string
  type?: 'text' | 'confirm'
  confirmData?: FormFillData
  confirmed?: boolean
}

interface Props {
  messages: ChatMessage[]
  loading: boolean
  confirmTitle: string
  confirmBtn: string
  correctBtn: string
  formFilled: string
  conversationLabel: string
}

defineProps<Props>()

const emit = defineEmits<{
  confirm: [msg: ChatMessage]
  correct: []
}>()
</script>

<template>
  <div role="log" aria-live="polite" :aria-label="conversationLabel">
    <TransitionGroup name="msg">
      <div
        v-for="msg in messages"
        :key="msg.id"
        class="msg-row"
        :class="[msg.role, msg.type === 'confirm' ? 'confirm-row' : '']"
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

        <!-- Confirmation card -->
        <ConfirmCard
          v-if="msg.type === 'confirm' && msg.confirmData"
          :confirm-data="msg.confirmData"
          :confirmed="msg.confirmed ?? false"
          :loading="loading"
          :confirm-title="confirmTitle"
          :confirm-btn="confirmBtn"
          :correct-btn="correctBtn"
          :form-filled="formFilled"
          @confirm="emit('confirm', msg)"
          @correct="emit('correct')"
        />

        <!-- Normal message bubble -->
        <div v-else class="msg-bubble">
          <span v-if="msg.role === 'assistant' && !msg.text && loading" class="typing">
            <span></span><span></span><span></span>
          </span>
          <span v-else-if="msg.role === 'assistant'" class="md" v-html="renderMarkdown(msg.text)" />
          <template v-else>{{ msg.text }}</template>
        </div>
      </div>
    </TransitionGroup>
  </div>
</template>

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
</style>
