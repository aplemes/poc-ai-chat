<script setup lang="ts">
import { ref, watch, nextTick, onBeforeUnmount } from 'vue'

const props = defineProps<{
  open: boolean
  loading: boolean
  text: string
  error: string
  renderMd: (text: string) => string
}>()

const emit = defineEmits<{
  close: []
  submit: []
}>()

const panelRef = ref<HTMLElement | null>(null)
let previouslyFocused: HTMLElement | null = null

const FOCUSABLE =
  'a[href], button:not([disabled]), input, select, textarea, [tabindex]:not([tabindex="-1"])'

function getFocusable(): HTMLElement[] {
  return panelRef.value ? Array.from(panelRef.value.querySelectorAll<HTMLElement>(FOCUSABLE)) : []
}

function handleKeydown(e: KeyboardEvent) {
  if (e.key === 'Escape') {
    emit('close')
    return
  }
  if (e.key !== 'Tab') return
  const focusable = getFocusable()
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

watch(
  () => props.open,
  async (isOpen) => {
    if (isOpen) {
      previouslyFocused = document.activeElement as HTMLElement
      await nextTick()
      panelRef.value?.focus()
      document.addEventListener('keydown', handleKeydown)
    } else {
      document.removeEventListener('keydown', handleKeydown)
      previouslyFocused?.focus()
    }
  },
)

onBeforeUnmount(() => document.removeEventListener('keydown', handleKeydown))
</script>

<template>
  <Teleport to="body">
    <Transition name="analysis-modal">
      <div v-if="open" class="analysis-backdrop" @click.self="$emit('close')">
        <div
          ref="panelRef"
          class="analysis-panel"
          role="dialog"
          aria-modal="true"
          aria-label="AI Form Review"
          tabindex="-1"
        >
          <div class="analysis-header">
            <div class="analysis-title">
              <div class="analysis-icon-wrap">
                <svg viewBox="0 0 20 20" fill="currentColor" class="analysis-icon-svg">
                  <path
                    d="M10 1a.75.75 0 0 1 .75.75v1.5a.75.75 0 0 1-1.5 0v-1.5A.75.75 0 0 1 10 1ZM5.05 3.05a.75.75 0 0 1 1.06 0l1.062 1.06A.75.75 0 1 1 6.11 5.173L5.05 4.11a.75.75 0 0 1 0-1.06ZM14.95 3.05a.75.75 0 0 1 0 1.06l-1.06 1.062a.75.75 0 0 1-1.062-1.061l1.061-1.061a.75.75 0 0 1 1.061 0ZM3 9.25a.75.75 0 0 1 .75-.75h1.5a.75.75 0 0 1 0 1.5h-1.5A.75.75 0 0 1 3 9.25ZM14.75 9.25a.75.75 0 0 1 .75-.75h1.5a.75.75 0 0 1 0 1.5h-1.5a.75.75 0 0 1-.75-.75ZM10 14a4 4 0 1 0 0-8 4 4 0 0 0 0 8Z"
                  />
                </svg>
              </div>
              <div>
                <p class="analysis-label">AI Review</p>
                <p class="analysis-subtitle">Checking your demand before submission</p>
              </div>
            </div>
            <button
              class="analysis-close"
              aria-label="Close AI review dialog"
              @click="$emit('close')"
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

          <div class="analysis-body">
            <div
              v-if="loading && !text && !error"
              class="analysis-loading"
              aria-live="polite"
              aria-busy="true"
            >
              <div class="analysis-loading-dots" aria-hidden="true">
                <span></span><span></span><span></span>
              </div>
              <p>Analysing your demand…</p>
            </div>
            <p v-else-if="error" class="analysis-error" role="alert">{{ error }}</p>
            <div v-else class="analysis-md" v-html="renderMd(text)" />
          </div>

          <div v-if="!loading" class="analysis-footer">
            <template v-if="error">
              <button class="analysis-btn analysis-btn--primary" @click="$emit('close')">
                <svg viewBox="0 0 20 20" fill="currentColor">
                  <path
                    fill-rule="evenodd"
                    d="M17 10a.75.75 0 0 1-.75.75H5.612l4.158 3.96a.75.75 0 1 1-1.04 1.08l-5.5-5.25a.75.75 0 0 1 0-1.08l5.5-5.25a.75.75 0 1 1 1.04 1.08L5.612 9.25H16.25A.75.75 0 0 1 17 10Z"
                    clip-rule="evenodd"
                  />
                </svg>
                Return to form
              </button>
            </template>
            <template v-else>
              <button class="analysis-btn analysis-btn--primary" @click="$emit('close')">
                <svg viewBox="0 0 20 20" fill="currentColor">
                  <path
                    fill-rule="evenodd"
                    d="M17 10a.75.75 0 0 1-.75.75H5.612l4.158 3.96a.75.75 0 1 1-1.04 1.08l-5.5-5.25a.75.75 0 0 1 0-1.08l5.5-5.25a.75.75 0 1 1 1.04 1.08L5.612 9.25H16.25A.75.75 0 0 1 17 10Z"
                    clip-rule="evenodd"
                  />
                </svg>
                Return to form
              </button>
              <button
                class="analysis-btn analysis-btn--ghost"
                aria-describedby="submit-anyway-desc"
                @click="$emit('submit')"
              >
                Submit anyway
              </button>
              <span id="submit-anyway-desc" class="visually-hidden">
                Proceed with submission ignoring the AI review suggestions
              </span>
            </template>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style src="@/components/AnalysisModal.css"></style>
