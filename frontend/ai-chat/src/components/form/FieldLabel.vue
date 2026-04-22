<script setup lang="ts">
import { useFormContext } from '@/composables/useFormContext'
import { useFieldChatStore } from '@/stores/fieldChat'

const props = defineProps<{
  fieldName: string
  label: string
  labelFor?: string
  labelId?: string
  required?: boolean
}>()

const { aiFilledFields, badgeClass } = useFormContext()
const fieldChatStore = useFieldChatStore()

const SPARKLE = '✦'
</script>

<template>
  <div class="label-row">
    <label
      :for="props.labelId ? props.labelFor : (props.labelFor ?? props.fieldName)"
      :id="props.labelId"
    >
      {{ props.label }}<span v-if="props.required" class="required"> *</span>
    </label>
    <div class="label-actions">
      <button
        type="button"
        class="field-ai-btn"
        :class="{ active: fieldChatStore.activeField === props.fieldName }"
        :aria-label="`Open AI assistant for ${props.label}`"
        @click="fieldChatStore.openPanel(props.fieldName)"
      >
        {{ SPARKLE }}
      </button>
      <Transition name="badge">
        <span
          v-if="aiFilledFields.has(props.fieldName)"
          :class="badgeClass(props.fieldName)"
          aria-live="polite"
          aria-label="Field filled by AI"
          ><span class="ai-dot"></span>AI</span
        >
      </Transition>
    </div>
  </div>
</template>
