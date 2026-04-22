<script setup lang="ts">
import { ref, provide, computed } from 'vue'
import { FORM_CONTEXT_KEY } from '@/composables/useFormContext'
import { useAiFormFill } from '@/composables/useAiFormFill'
import { useFormAnalysis } from '@/composables/useFormAnalysis'
import FieldChatPanel from '@/components/FieldChatPanel.vue'
import AnalysisModal from '@/components/AnalysisModal.vue'
import SectionIdentity from '@/components/form/SectionIdentity.vue'
import SectionContext from '@/components/form/SectionContext.vue'
import SectionBenefits from '@/components/form/SectionBenefits.vue'

const form = ref({
  title: '',
  businessLine: '',
  requesterBU: '',
  busInterested: [] as string[],
  timeSensitive: 'No' as string,
  whyDemand: '',
  whoIsImpacted: '',
  benefitCategory: '',
  benefitHypothesis: '',
  measureBenefits: '',
})

const errors = ref<Record<string, string>>({})

const { aiFilledFields, clearField, badgeClass } = useAiFormFill(form, errors)

function validate(): boolean {
  errors.value = {}
  if (!form.value.title.trim()) errors.value['title'] = 'Required field'
  if (!form.value.businessLine) errors.value['businessLine'] = 'Please select an option'
  if (!form.value.requesterBU) errors.value['requesterBU'] = 'Please select an option'
  if (!form.value.whyDemand.trim()) errors.value['whyDemand'] = 'Required field'
  if (!form.value.whoIsImpacted.trim()) errors.value['whoIsImpacted'] = 'Required field'
  if (!form.value.benefitCategory) errors.value['benefitCategory'] = 'Please select an option'
  if (!form.value.benefitHypothesis.trim()) errors.value['benefitHypothesis'] = 'Required field'
  if (!form.value.measureBenefits.trim()) errors.value['measureBenefits'] = 'Required field'
  return Object.keys(errors.value).length === 0
}

function submitForm() {
  // TODO: wire to backend submission API
  console.log('Form submitted:', form.value)
}

const {
  analysisOpen,
  analysisLoading,
  analysisText,
  analysisError,
  handleSubmit,
  closeAnalysis,
  confirmSubmit,
  renderMd,
} = useFormAnalysis(() => ({ ...form.value, lowConfidenceFields: [] }), validate, submitForm)

provide(FORM_CONTEXT_KEY, { form, errors, aiFilledFields, clearField, badgeClass })

// UX-05: announce AI fills to screen readers via live region
const aiAnnouncement = computed(() =>
  aiFilledFields.value.size > 0
    ? `${aiFilledFields.value.size} field${aiFilledFields.value.size > 1 ? 's' : ''} filled by AI assistant`
    : '',
)
</script>

<template>
  <div aria-live="polite" aria-atomic="true" class="sr-only">{{ aiAnnouncement }}</div>

  <form class="request-form" novalidate @submit.prevent="handleSubmit">
    <FieldChatPanel />

    <div class="form-header">
      <div class="form-header-icon">
        <svg viewBox="0 0 20 20" fill="currentColor">
          <path d="M9 2a1 1 0 000 2h2a1 1 0 100-2H9z" />
          <path
            fill-rule="evenodd"
            d="M4 5a2 2 0 012-2 3 3 0 003 3h2a3 3 0 003-3 2 2 0 012 2v11a2 2 0 01-2 2H6a2 2 0 01-2-2V5zm3 4a1 1 0 000 2h.01a1 1 0 100-2H7zm3 0a1 1 0 000 2h3a1 1 0 100-2h-3zm-3 4a1 1 0 100 2h.01a1 1 0 100-2H7zm3 0a1 1 0 100 2h3a1 1 0 100-2h-3z"
            clip-rule="evenodd"
          />
        </svg>
      </div>
      <div>
        <h1 class="form-title">New Demand Request</h1>
        <p class="form-subtitle">
          Fill in the details below or use the AI assistant to auto-complete the form.
        </p>
      </div>
    </div>

    <SectionIdentity />
    <SectionContext />
    <SectionBenefits />

    <div class="form-actions">
      <p class="actions-note"><span class="required">*</span> Required fields</p>
      <button type="submit" class="submit-btn" :disabled="analysisLoading">
        <svg v-if="analysisLoading" class="spinner" viewBox="0 0 24 24" fill="none">
          <circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="3" opacity=".25" />
          <path
            d="M12 2a10 10 0 0 1 10 10"
            stroke="currentColor"
            stroke-width="3"
            stroke-linecap="round"
          />
        </svg>
        <svg v-else viewBox="0 0 20 20" fill="currentColor">
          <path
            d="M10.894 2.553a1 1 0 00-1.788 0l-7 14a1 1 0 001.169 1.409l5-1.429A1 1 0 009 15.571V11a1 1 0 112 0v4.571a1 1 0 00.725.962l5 1.428a1 1 0 001.17-1.408l-7-14z"
          />
        </svg>
        {{ analysisLoading ? 'Reviewing...' : 'Submit Request' }}
      </button>
    </div>
  </form>

  <AnalysisModal
    :open="analysisOpen"
    :loading="analysisLoading"
    :text="analysisText"
    :error="analysisError"
    :render-md="renderMd"
    @close="closeAnalysis"
    @submit="confirmSubmit"
  />
</template>
