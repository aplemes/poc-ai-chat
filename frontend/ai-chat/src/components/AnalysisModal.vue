<script setup lang="ts">
defineProps<{
  open: boolean
  loading: boolean
  text: string
  error: string
  renderMd: (text: string) => string
}>()

defineEmits<{
  close: []
  submit: []
}>()
</script>

<template>
  <Teleport to="body">
    <Transition name="analysis-modal">
      <div v-if="open" class="analysis-backdrop" @click.self="$emit('close')">
        <div class="analysis-panel" role="dialog" aria-label="AI Form Review">
          <div class="analysis-header">
            <div class="analysis-title">
              <span class="analysis-sparkle">✦</span>
              <div>
                <p class="analysis-label">AI Review</p>
                <p class="analysis-subtitle">Checking your demand before submission</p>
              </div>
            </div>
            <button class="analysis-close" aria-label="Close" @click="$emit('close')">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round">
                <path d="M18 6L6 18M6 6l12 12" />
              </svg>
            </button>
          </div>

          <div class="analysis-body">
            <div v-if="loading && !text && !error" class="analysis-loading">
              <span class="typing"><span></span><span></span><span></span></span>
              <p>Analysing your demand...</p>
            </div>
            <p v-else-if="error" class="analysis-error">{{ error }}</p>
            <div v-else class="analysis-md" v-html="renderMd(text)" />
          </div>

          <div v-if="!loading" class="analysis-footer">
            <button class="analysis-btn analysis-btn--secondary" @click="$emit('close')">Fix issues</button>
            <button class="analysis-btn analysis-btn--primary" @click="$emit('submit')">Submit anyway</button>
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
  background: rgba(0, 0, 0, 0.35);
  z-index: 9000;
  display: flex;
  align-items: flex-end;
  justify-content: center;
  padding: 0 1rem;
}

.analysis-panel {
  width: 100%;
  max-width: 640px;
  max-height: 80vh;
  background: var(--color-neutral-0);
  border-radius: var(--radius-xl) var(--radius-xl) 0 0;
  box-shadow: 0 -8px 40px rgba(0, 0, 0, 0.18);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.analysis-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 1.1rem 1.25rem;
  background: var(--gradient-primary);
  flex-shrink: 0;
}

.analysis-title {
  display: flex;
  align-items: center;
  gap: 0.65rem;
}

.analysis-sparkle {
  font-size: 1.3rem;
  color: rgba(255, 255, 255, 0.9);
}

.analysis-label {
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-semibold);
  color: #fff;
  margin: 0;
  line-height: 1.2;
}

.analysis-subtitle {
  font-size: 0.65rem;
  color: rgba(255, 255, 255, 0.7);
  margin: 0;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.analysis-close {
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

.analysis-close:hover {
  background: rgba(0, 0, 0, 0.22);
  color: #fff;
}

.analysis-close svg {
  width: 0.9rem;
  height: 0.9rem;
}

.analysis-body {
  flex: 1;
  overflow-y: auto;
  padding: 1.5rem 1.5rem 1rem;
  min-height: 0;
}

.analysis-loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
  padding: 2rem 0;
  color: var(--color-neutral-500);
  font-size: var(--font-size-sm);
}

.analysis-error {
  color: var(--color-error);
  font-size: var(--font-size-sm);
}

.analysis-md {
  font-size: var(--font-size-sm);
  color: var(--color-neutral-800);
  line-height: var(--line-height-relaxed);
}

.analysis-md :deep(p) {
  margin: 0 0 0.6rem;
}
.analysis-md :deep(p:last-child) {
  margin-bottom: 0;
}
.analysis-md :deep(ul),
.analysis-md :deep(ol) {
  margin: 0.25rem 0 0.6rem;
  padding-left: 1.3rem;
}
.analysis-md :deep(li) {
  margin-bottom: 0.35rem;
}
.analysis-md :deep(strong) {
  font-weight: var(--font-weight-semibold);
}

.analysis-footer {
  display: flex;
  gap: 0.75rem;
  padding: 1rem 1.5rem 1.5rem;
  border-top: 1px solid var(--color-neutral-100);
  flex-shrink: 0;
}

.analysis-btn {
  flex: 1;
  padding: 0.65rem 1rem;
  border-radius: var(--radius-md);
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-medium);
  cursor: pointer;
  transition: all var(--transition-fast);
  border: none;
}

.analysis-btn--secondary {
  background: var(--color-neutral-100);
  color: var(--color-neutral-700);
}

.analysis-btn--secondary:hover {
  background: var(--color-neutral-200);
}

.analysis-btn--primary {
  background: var(--color-primary);
  color: #fff;
  box-shadow: var(--shadow-primary);
}

.analysis-btn--primary:hover {
  background: var(--color-primary-hover);
}

.analysis-modal-enter-active {
  transition: opacity 0.2s ease, transform 0.25s ease;
}
.analysis-modal-leave-active {
  transition: opacity 0.15s ease, transform 0.18s ease;
}
.analysis-modal-enter-from,
.analysis-modal-leave-to {
  opacity: 0;
  transform: translateY(20px);
}
.analysis-modal-enter-from .analysis-panel,
.analysis-modal-leave-to .analysis-panel {
  transform: translateY(40px);
}
</style>
