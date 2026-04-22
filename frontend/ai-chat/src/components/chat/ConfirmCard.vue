<script setup lang="ts">
import { computed } from 'vue'
import type { FormFillData } from '@/services/chatService'

interface Props {
  confirmData: FormFillData
  confirmed: boolean
  loading: boolean
  confirmTitle: string
  confirmBtn: string
  correctBtn: string
  formFilled: string
}

const props = defineProps<Props>()

const emit = defineEmits<{
  confirm: []
  correct: []
}>()

const FIELD_LABELS: Record<string, string> = {
  title: 'Title',
  businessLine: 'Business Line',
  requesterBU: 'Requester BU',
  busInterested: 'BUs Interested',
  timeSensitive: 'Time Sensitive',
  whyDemand: 'Why demand',
  whoIsImpacted: 'Who is impacted',
  benefitCategory: 'Benefit category',
  benefitHypothesis: 'Benefit hypothesis',
  measureBenefits: 'Measure benefits',
}

const visibleFields = computed(() =>
  (Object.entries(props.confirmData) as [string, unknown][]).filter(
    ([k, val]) => k !== 'lowConfidenceFields' && Boolean(val),
  ),
)
</script>

<template>
  <div class="confirm-card">
    <p class="confirm-title">{{ props.confirmTitle }}</p>
    <ul class="confirm-fields">
      <li v-for="[key, val] in visibleFields" :key="key">
        <span class="confirm-key">{{ FIELD_LABELS[key] ?? key }}</span>
        <span
          class="confirm-val"
          :class="{ uncertain: props.confirmData.lowConfidenceFields?.includes(key) }"
        >
          {{ val }}
          <span
            v-if="props.confirmData.lowConfidenceFields?.includes(key)"
            class="uncertain-badge"
            title="Inferred from context"
            >?</span
          >
        </span>
      </li>
    </ul>
    <div class="confirm-actions">
      <button
        class="confirm-btn primary"
        :disabled="props.confirmed || props.loading"
        @click="emit('confirm')"
      >
        <svg v-if="props.confirmed" viewBox="0 0 20 20" fill="currentColor" class="check-icon">
          <path
            fill-rule="evenodd"
            d="M16.707 5.293a1 1 0 010 1.414L8.414 15l-4.121-4.121a1 1 0 011.414-1.414L8.414 12.172l6.879-6.879a1 1 0 011.414 0z"
            clip-rule="evenodd"
          />
        </svg>
        {{ props.confirmed ? props.formFilled : props.confirmBtn }}
      </button>
      <button
        class="confirm-btn secondary"
        :disabled="props.confirmed || props.loading"
        @click="emit('correct')"
      >
        {{ props.correctBtn }}
      </button>
    </div>
  </div>
</template>

<style scoped>
.confirm-card {
  max-width: 88%;
  background: var(--color-neutral-0);
  border: 1px solid var(--color-primary);
  border-radius: var(--radius-xl);
  border-bottom-left-radius: var(--radius-xs);
  padding: 0.9rem 1rem;
  box-shadow: var(--shadow-sm);
  display: flex;
  flex-direction: column;
  gap: 0.65rem;
}

.confirm-title {
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-semibold);
  color: var(--color-primary);
  margin: 0;
}

.confirm-fields {
  list-style: none;
  margin: 0;
  padding: 0;
  display: flex;
  flex-direction: column;
  gap: 0.35rem;
}

.confirm-fields li {
  display: flex;
  gap: 0.5rem;
  font-size: var(--font-size-xs);
  line-height: 1.4;
}

.confirm-key {
  color: var(--color-neutral-500);
  flex-shrink: 0;
  min-width: 7rem;
  font-weight: var(--font-weight-medium);
}

.confirm-val {
  color: var(--color-neutral-800);
  word-break: break-word;
  display: flex;
  align-items: center;
  gap: 0.3rem;
}

.confirm-val.uncertain {
  color: var(--color-warning);
}

.uncertain-badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 1rem;
  height: 1rem;
  border-radius: 50%;
  background: var(--color-warning);
  color: #fff;
  font-size: 0.6rem;
  font-weight: var(--font-weight-bold);
  flex-shrink: 0;
  cursor: help;
}

.confirm-actions {
  display: flex;
  gap: 0.5rem;
  margin-top: 0.25rem;
}

.confirm-btn {
  flex: 1;
  padding: 0.45rem 0.75rem;
  border-radius: var(--radius-md);
  font-size: var(--font-size-xs);
  font-weight: var(--font-weight-semibold);
  cursor: pointer;
  border: none;
  transition: all var(--transition-fast);
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.3rem;
}

.confirm-btn.primary {
  background: var(--color-primary);
  color: #fff;
  box-shadow: var(--shadow-primary);
}

.confirm-btn.primary:hover:not(:disabled) {
  background: var(--color-primary-hover);
}

.confirm-btn.secondary {
  background: var(--color-neutral-100);
  color: var(--color-neutral-700);
  border: 1px solid var(--color-neutral-200);
}

.confirm-btn.secondary:hover:not(:disabled) {
  background: var(--color-neutral-200);
}

.confirm-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.check-icon {
  width: 0.85rem;
  height: 0.85rem;
}
</style>
