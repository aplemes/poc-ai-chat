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
            <!-- Error state: only allow returning to form, no submit -->
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
            <!-- Success state: allow returning to fix or submitting anyway -->
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

<style scoped>
.analysis-backdrop {
  position: fixed;
  inset: 0;
  background: rgba(10, 15, 26, 0.55);
  backdrop-filter: blur(3px);
  z-index: 9000;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 1.5rem;
}

.analysis-panel {
  width: 100%;
  max-width: 560px;
  max-height: 85vh;
  background: var(--color-neutral-0);
  border-radius: var(--radius-xl);
  box-shadow: var(--shadow-2xl);
  display: flex;
  flex-direction: column;
  overflow: hidden;
  outline: none;
}

/* ── Header ── */
.analysis-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 1.125rem 1.25rem;
  background: var(--gradient-primary);
  flex-shrink: 0;
}

.analysis-title {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.analysis-icon-wrap {
  width: 2rem;
  height: 2rem;
  border-radius: var(--radius-sm);
  background: rgba(255, 255, 255, 0.18);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.analysis-icon-svg {
  width: 1.1rem;
  height: 1.1rem;
  color: #fff;
}

.analysis-label {
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-semibold);
  color: #fff;
  margin: 0;
  line-height: 1.2;
}

.analysis-subtitle {
  font-size: 0.625rem;
  color: rgba(255, 255, 255, 0.65);
  margin: 0.1rem 0 0;
  text-transform: uppercase;
  letter-spacing: 0.07em;
}

.analysis-close {
  width: 1.75rem;
  height: 1.75rem;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(0, 0, 0, 0.15);
  border: none;
  border-radius: var(--radius-sm);
  color: rgba(255, 255, 255, 0.8);
  cursor: pointer;
  transition: background var(--transition-fast);
  flex-shrink: 0;
}

.analysis-close:hover {
  background: rgba(0, 0, 0, 0.25);
  color: #fff;
}
.analysis-close:focus-visible {
  outline: 2px solid rgba(255, 255, 255, 0.6);
  outline-offset: 2px;
}

.analysis-close svg {
  width: 0.85rem;
  height: 0.85rem;
}

/* ── Body ── */
.analysis-body {
  flex: 1;
  overflow-y: auto;
  padding: 1.25rem 1.25rem 0.75rem;
  min-height: 0;
  background: var(--color-neutral-50);
}

/* ── Loading ── */
.analysis-loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
  padding: 2.5rem 0;
  color: var(--color-neutral-500);
  font-size: var(--font-size-sm);
}

.analysis-loading-dots {
  display: flex;
  gap: 5px;
}

.analysis-loading-dots span {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: var(--color-primary);
  animation: dot-bounce 1.1s ease-in-out infinite;
}

.analysis-loading-dots span:nth-child(2) {
  animation-delay: 0.18s;
}
.analysis-loading-dots span:nth-child(3) {
  animation-delay: 0.36s;
}

@keyframes dot-bounce {
  0%,
  80%,
  100% {
    transform: scale(0.6);
    opacity: 0.4;
  }
  40% {
    transform: scale(1);
    opacity: 1;
  }
}

/* ── Error ── */
.analysis-error {
  color: var(--color-error);
  font-size: var(--font-size-sm);
  padding: 1rem;
  background: var(--color-error-subtle);
  border-radius: var(--radius-md);
  border-left: 3px solid var(--color-error);
}

/* ── Markdown content ── */
.analysis-md {
  font-size: var(--font-size-sm);
  color: var(--color-neutral-700);
  line-height: var(--line-height-normal);
}

.analysis-md :deep(p) {
  margin: 0 0 0.875rem;
  color: var(--color-neutral-600);
  font-size: var(--font-size-sm);
}

.analysis-md :deep(ul),
.analysis-md :deep(ol) {
  margin: 0;
  padding: 0;
  list-style: none;
  display: flex;
  flex-direction: column;
  gap: 0.625rem;
}

.analysis-md :deep(li) {
  background: var(--color-neutral-0);
  border: 1px solid var(--color-neutral-200);
  border-left: 3px solid var(--color-warning);
  border-radius: var(--radius-md);
  padding: 0.875rem 1rem;
  font-size: var(--font-size-sm);
  line-height: var(--line-height-relaxed);
  color: var(--color-neutral-700);
  box-shadow: var(--shadow-xs);
  position: relative;
}

.analysis-md :deep(code) {
  font-family: 'JetBrains Mono', 'Fira Code', 'Courier New', monospace;
  font-size: 0.75rem;
  background: var(--color-neutral-100);
  color: var(--color-neutral-800);
  border: 1px solid var(--color-neutral-200);
  border-radius: var(--radius-xs);
  padding: 0.1em 0.35em;
}

.analysis-md :deep(strong) {
  font-weight: var(--font-weight-semibold);
  color: var(--color-neutral-800);
}

/* ── Footer ── */
.analysis-footer {
  display: flex;
  gap: 0.625rem;
  padding: 0.875rem 1.25rem 1.125rem;
  border-top: 1px solid var(--color-neutral-100);
  background: var(--color-neutral-0);
  flex-shrink: 0;
}

.analysis-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 0.4rem;
  padding: 0.6rem 1.1rem;
  border-radius: var(--radius-md);
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-medium);
  cursor: pointer;
  transition: all var(--transition-fast);
  border: none;
  white-space: nowrap;
}

.analysis-btn svg {
  width: 0.9rem;
  height: 0.9rem;
  flex-shrink: 0;
}

.analysis-btn--primary {
  flex: 1;
  background: var(--color-primary);
  color: #fff;
  box-shadow: var(--shadow-primary);
}

.analysis-btn--primary:hover {
  background: var(--color-primary-hover);
}
.analysis-btn--primary:active {
  background: var(--color-primary-active);
}

.analysis-btn--ghost {
  background: transparent;
  color: var(--color-neutral-500);
  border: 1px solid var(--color-neutral-200);
  padding-inline: 1rem;
}

.analysis-btn--ghost:hover {
  background: var(--color-neutral-100);
  color: var(--color-neutral-700);
  border-color: var(--color-neutral-300);
}

.analysis-btn:focus-visible {
  outline: 2px solid var(--color-primary);
  outline-offset: 2px;
}

/* ── Accessibility ── */
.visually-hidden {
  position: absolute;
  width: 1px;
  height: 1px;
  padding: 0;
  margin: -1px;
  overflow: hidden;
  clip: rect(0, 0, 0, 0);
  white-space: nowrap;
  border: 0;
}

/* ── Transitions ── */
.analysis-modal-enter-active {
  transition:
    opacity 0.2s var(--ease-out),
    transform 0.25s var(--ease-out);
}
.analysis-modal-leave-active {
  transition:
    opacity 0.15s var(--ease-in),
    transform 0.18s var(--ease-in);
}
.analysis-modal-enter-from,
.analysis-modal-leave-to {
  opacity: 0;
}
.analysis-modal-enter-from .analysis-panel {
  transform: scale(0.96) translateY(8px);
}
.analysis-modal-leave-to .analysis-panel {
  transform: scale(0.97) translateY(4px);
}

@media (prefers-reduced-motion: reduce) {
  .analysis-modal-enter-active,
  .analysis-modal-leave-active {
    transition: none;
  }
  .analysis-loading-dots span {
    animation: none;
  }
}
</style>
