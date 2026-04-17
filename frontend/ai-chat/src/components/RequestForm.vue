<script setup lang="ts">
import { ref, watch } from 'vue'
import { useChatStore } from '@/stores/chat'

const requesters = [
  { id: 'ADEO-8052', value: 'Adeo Marketplace Services' },
  { id: 'ADEO-35430', value: 'Adeo Production' },
  { id: 'ADEO-35424', value: 'Adeo Services Chine' },
  { id: 'ADEO-8078', value: 'Adeo Services France' },
  { id: 'ADEO-35426', value: 'Adeo Services Poland' },
  { id: 'ADEO-35427', value: 'Adeo Services Vietnam' },
  { id: 'ADEO-8062', value: 'Bricocenter Italy' },
  { id: 'ADEO-8087', value: 'Bricoman Poland' },
  { id: 'ADEO-35431', value: 'Enki Home' },
  { id: 'ADEO-36214', value: 'GO XL' },
  { id: 'ADEO-8089', value: 'Golilla' },
  { id: 'ADEO-8054', value: 'Kbane France' },
  { id: 'ADEO-8095', value: 'Leroy Merlin Brazil' },
  { id: 'ADEO-8064', value: 'Leroy Merlin France' },
  { id: 'ADEO-8075', value: 'Leroy Merlin Greece/Cyprus' },
  { id: 'ADEO-8100', value: 'Leroy Merlin Italy' },
  { id: 'ADEO-8067', value: 'Leroy Merlin Poland' },
  { id: 'ADEO-8057', value: 'Leroy Merlin Portugal' },
  { id: 'ADEO-40345', value: 'Leroy Merlin Renovation' },
  { id: 'ADEO-8091', value: 'Leroy Merlin Romania' },
  { id: 'ADEO-8061', value: 'Leroy Merlin South Africa' },
  { id: 'ADEO-8084', value: 'Leroy Merlin Spain' },
  { id: 'ADEO-8071', value: 'Leroy Merlin Ukraine' },
  { id: 'ADEO-8051', value: 'Obramat Portugal' },
  { id: 'ADEO-8053', value: 'Obramat Spain' },
  { id: 'ADEO-8092', value: 'Obramax Brazil' },
  { id: 'ADEO-8070', value: 'Quotatis' },
  { id: 'ADEO-23566', value: 'Saint Maclou France' },
  { id: 'ADEO-8074', value: 'Tecnomat France' },
  { id: 'ADEO-8055', value: 'Tecnomat Italy' },
  { id: 'ADEO-8056', value: 'Terra Incognita' },
  { id: 'ADEO-8060', value: 'Weldom France' },
]

const concerned = [
  { id: '20047', value: 'Adeo Marketplace Services' },
  { id: '20048', value: 'Adeo Production' },
  { id: '20049', value: 'Adeo Services Chine' },
  { id: '20050', value: 'Adeo Services France' },
  { id: '20051', value: 'Adeo Services Poland' },
  { id: '20052', value: 'Adeo Services Vietnam' },
  { id: '20053', value: 'Bricocenter Italy' },
  { id: '20054', value: 'Bricoman Poland' },
  { id: '20055', value: 'Enki Home' },
  { id: '20056', value: 'GO XL' },
  { id: '20057', value: 'Golilla' },
  { id: '20058', value: 'Kbane France' },
  { id: '20059', value: 'Leroy Merlin Brazil' },
  { id: '20060', value: 'Leroy Merlin France' },
  { id: '20061', value: 'Leroy Merlin Greece/Cyprus' },
  { id: '20062', value: 'Leroy Merlin Italy' },
  { id: '20063', value: 'Leroy Merlin Poland' },
  { id: '20064', value: 'Leroy Merlin Portugal' },
  { id: '20065', value: 'Leroy Merlin Renovation' },
  { id: '20066', value: 'Leroy Merlin Romania' },
  { id: '20067', value: 'Leroy Merlin South Africa' },
  { id: '20068', value: 'Leroy Merlin Spain' },
  { id: '20069', value: 'Leroy Merlin Ukraine' },
  { id: '20070', value: 'Obramat Portugal' },
  { id: '20071', value: 'Obramat Spain' },
  { id: '20072', value: 'Obramax Brazil' },
  { id: '20073', value: 'Quotatis' },
  { id: '20074', value: 'Saint Maclou France' },
  { id: '20075', value: 'Tecnomat France' },
  { id: '20076', value: 'Tecnomat Italy' },
  { id: '20077', value: 'Terra Incognita' },
  { id: '20078', value: 'Weldom France' },
]

const organizations = [
  { id: '18518', value: 'Omnicommerce Experience' },
  { id: '18519', value: 'Services & Renovation' },
  { id: '18520', value: 'Supply Chain & Delivery' },
  { id: '18521', value: 'Offer & Industry' },
  { id: '18522', value: 'Finance' },
  { id: '18523', value: 'Positive Impacts' },
  { id: '18524', value: 'Human & Sharing' },
  { id: '19033', value: 'Executive Succession Plan' },
  { id: '18525', value: 'Digital Data Tech' },
]

const benefitCategories = [
  'Cost efficiency',
  'Environmental & social sustainability',
  'Service quality & security risk',
  'Customer satisfaction & revenue',
  'Innovation',
  'Other',
]

const timeSensitiveOptions = ['No', 'Legal', 'Security'] as const

const chatStore = useChatStore()

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

const aiFilledFields = ref<Set<string>>(new Set())
const errors = ref<Record<string, string>>({})
const submitting = ref(false)

watch(
  () => chatStore.pendingFormFill,
  (data) => {
    if (!data) return
    const filled = new Set<string>()
    if (data.title) { form.value.title = data.title; filled.add('title') }
    if (data.businessLine) { form.value.businessLine = data.businessLine; filled.add('businessLine') }
    if (data.requesterBU) { form.value.requesterBU = data.requesterBU; filled.add('requesterBU') }
    const busIds = Array.isArray(data.busInterested)
      ? data.busInterested
      : data.busInterested
        ? [data.busInterested as unknown as string]
        : []
    if (busIds.length) { form.value.busInterested = busIds; filled.add('busInterested') }
    if (data.timeSensitive) { form.value.timeSensitive = data.timeSensitive; filled.add('timeSensitive') }
    if (data.whyDemand) { form.value.whyDemand = data.whyDemand; filled.add('whyDemand') }
    if (data.whoIsImpacted) { form.value.whoIsImpacted = data.whoIsImpacted; filled.add('whoIsImpacted') }
    if (data.benefitCategory) { form.value.benefitCategory = data.benefitCategory; filled.add('benefitCategory') }
    if (data.benefitHypothesis) { form.value.benefitHypothesis = data.benefitHypothesis; filled.add('benefitHypothesis') }
    if (data.measureBenefits) { form.value.measureBenefits = data.measureBenefits; filled.add('measureBenefits') }
    aiFilledFields.value = filled
    chatStore.clearFormFill()
  },
)

function clearError(field: string) { delete errors.value[field] }

function toggleBusInterested(id: string) {
  const idx = form.value.busInterested.indexOf(id)
  form.value.busInterested = idx === -1
    ? [...form.value.busInterested, id]
    : form.value.busInterested.filter((v) => v !== id)
  aiFilledFields.value.delete('busInterested')
}

function removeBusInterested(id: string) {
  form.value.busInterested = form.value.busInterested.filter((v) => v !== id)
  aiFilledFields.value.delete('busInterested')
}

function getBuLabel(id: string): string {
  return concerned.find((c) => c.id === id)?.value ?? id
}

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

async function handleSubmit() {
  if (!validate()) return
  submitting.value = true
  try {
    console.log(form.value)
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <form class="request-form" @submit.prevent="handleSubmit" novalidate>

    <!-- Form Header -->
    <div class="form-header">
      <div class="form-header-icon">
        <svg viewBox="0 0 20 20" fill="currentColor">
          <path d="M9 2a1 1 0 000 2h2a1 1 0 100-2H9z" />
          <path fill-rule="evenodd" d="M4 5a2 2 0 012-2 3 3 0 003 3h2a3 3 0 003-3 2 2 0 012 2v11a2 2 0 01-2 2H6a2 2 0 01-2-2V5zm3 4a1 1 0 000 2h.01a1 1 0 100-2H7zm3 0a1 1 0 000 2h3a1 1 0 100-2h-3zm-3 4a1 1 0 100 2h.01a1 1 0 100-2H7zm3 0a1 1 0 100 2h3a1 1 0 100-2h-3z" clip-rule="evenodd" />
        </svg>
      </div>
      <div>
        <h1 class="form-title">New Demand Request</h1>
        <p class="form-subtitle">Fill in the details below or use the AI assistant to auto-complete the form.</p>
      </div>
    </div>

    <!-- Section 1: Demand identity -->
    <section class="form-section">
      <div class="section-header">
        <span class="section-number">01</span>
        <div>
          <h2 class="section-title">Demand Identity</h2>
          <p class="section-desc">Basic information about your request</p>
        </div>
      </div>

      <div class="fields-grid">

        <!-- Title -->
        <div class="field field--full">
          <div class="label-row">
            <label for="title">Title <span class="required">*</span></label>
            <Transition name="badge"><span v-if="aiFilledFields.has('title')" class="ai-badge"><span class="ai-dot"></span>AI</span></Transition>
          </div>
          <p class="field-hint">Start with an infinitive verb + scope. e.g. "Add the new Payment Method X on website"</p>
          <input
            id="title"
            v-model="form.title"
            type="text"
            placeholder="Add the new... / Improve the... / Enable..."
            :class="{ 'ai-filled': aiFilledFields.has('title'), 'has-error': errors['title'] }"
            @input="aiFilledFields.delete('title'); clearError('title')"
          />
          <Transition name="err"><p v-if="errors['title']" class="field-error" role="alert">{{ errors['title'] }}</p></Transition>
        </div>

        <!-- Business Line -->
        <div class="field">
          <div class="label-row">
            <label for="businessLine">Business Line <span class="required">*</span></label>
            <Transition name="badge"><span v-if="aiFilledFields.has('businessLine')" class="ai-badge"><span class="ai-dot"></span>AI</span></Transition>
          </div>
          <div class="select-wrap" :class="{ 'ai-filled': aiFilledFields.has('businessLine'), 'has-error': errors['businessLine'] }">
            <select id="businessLine" v-model="form.businessLine" @change="aiFilledFields.delete('businessLine'); clearError('businessLine')">
              <option value="" disabled>Select business line</option>
              <option v-for="org in organizations" :key="org.id" :value="org.id">{{ org.value }}</option>
            </select>
            <svg class="select-chevron" viewBox="0 0 20 20" fill="currentColor"><path fill-rule="evenodd" d="M5.23 7.21a.75.75 0 011.06.02L10 11.168l3.71-3.938a.75.75 0 111.08 1.04l-4.25 4.5a.75.75 0 01-1.08 0l-4.25-4.5a.75.75 0 01.02-1.06z" clip-rule="evenodd" /></svg>
          </div>
          <Transition name="err"><p v-if="errors['businessLine']" class="field-error" role="alert">{{ errors['businessLine'] }}</p></Transition>
        </div>

        <!-- Requester BU -->
        <div class="field">
          <div class="label-row">
            <label for="requesterBU">Requester BU <span class="required">*</span></label>
            <Transition name="badge"><span v-if="aiFilledFields.has('requesterBU')" class="ai-badge"><span class="ai-dot"></span>AI</span></Transition>
          </div>
          <div class="select-wrap" :class="{ 'ai-filled': aiFilledFields.has('requesterBU'), 'has-error': errors['requesterBU'] }">
            <select id="requesterBU" v-model="form.requesterBU" @change="aiFilledFields.delete('requesterBU'); clearError('requesterBU')">
              <option value="" disabled>Select your BU</option>
              <option v-for="req in requesters" :key="req.id" :value="req.id">{{ req.value }}</option>
            </select>
            <svg class="select-chevron" viewBox="0 0 20 20" fill="currentColor"><path fill-rule="evenodd" d="M5.23 7.21a.75.75 0 011.06.02L10 11.168l3.71-3.938a.75.75 0 111.08 1.04l-4.25 4.5a.75.75 0 01-1.08 0l-4.25-4.5a.75.75 0 01.02-1.06z" clip-rule="evenodd" /></svg>
          </div>
          <Transition name="err"><p v-if="errors['requesterBU']" class="field-error" role="alert">{{ errors['requesterBU'] }}</p></Transition>
        </div>

        <!-- BUs Interested -->
        <div class="field field--full">
          <div class="label-row">
            <label>BUs Interested</label>
            <Transition name="badge"><span v-if="aiFilledFields.has('busInterested')" class="ai-badge"><span class="ai-dot"></span>AI</span></Transition>
          </div>
          <div class="chips-field" :class="{ 'ai-filled': aiFilledFields.has('busInterested') }">
            <div v-if="form.busInterested.length" class="chips-list">
              <span v-for="id in form.busInterested" :key="id" class="chip" :class="{ 'chip--ai': aiFilledFields.has('busInterested') }">
                {{ getBuLabel(id) }}
                <button type="button" class="chip-x" :aria-label="`Remove ${getBuLabel(id)}`" @click="removeBusInterested(id)">
                  <svg viewBox="0 0 12 12" fill="currentColor"><path d="M1 1l10 10M11 1L1 11" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg>
                </button>
              </span>
            </div>
            <div class="select-wrap chips-select">
              <select
                value=""
                @change="(e) => { const v = (e.target as HTMLSelectElement).value; if (v) toggleBusInterested(v); (e.target as HTMLSelectElement).value = '' }"
              >
                <option value="" disabled>Add a BU...</option>
                <option v-for="bu in concerned" :key="bu.id" :value="bu.id" :disabled="form.busInterested.includes(bu.id)">{{ bu.value }}</option>
              </select>
              <svg class="select-chevron" viewBox="0 0 20 20" fill="currentColor"><path fill-rule="evenodd" d="M5.23 7.21a.75.75 0 011.06.02L10 11.168l3.71-3.938a.75.75 0 111.08 1.04l-4.25 4.5a.75.75 0 01-1.08 0l-4.25-4.5a.75.75 0 01.02-1.06z" clip-rule="evenodd" /></svg>
            </div>
          </div>
        </div>

        <!-- Time Sensitive -->
        <div class="field field--full">
          <div class="label-row">
            <label>Is this demand time sensitive?</label>
            <Transition name="badge"><span v-if="aiFilledFields.has('timeSensitive')" class="ai-badge"><span class="ai-dot"></span>AI</span></Transition>
          </div>
          <div class="segmented" :class="{ 'ai-filled': aiFilledFields.has('timeSensitive') }">
            <button
              v-for="opt in timeSensitiveOptions"
              :key="opt"
              type="button"
              class="seg-btn"
              :class="{ active: form.timeSensitive === opt }"
              @click="form.timeSensitive = opt; aiFilledFields.delete('timeSensitive')"
            >{{ opt }}</button>
          </div>
        </div>

      </div>
    </section>

    <!-- Section 2: Context -->
    <section class="form-section">
      <div class="section-header">
        <span class="section-number">02</span>
        <div>
          <h2 class="section-title">Context &amp; Impact</h2>
          <p class="section-desc">Help us understand the problem you're solving</p>
        </div>
      </div>

      <div class="fields-grid">

        <div class="field field--full">
          <div class="label-row">
            <label for="whyDemand">Why are you making this demand? <span class="required">*</span></label>
            <Transition name="badge"><span v-if="aiFilledFields.has('whyDemand')" class="ai-badge"><span class="ai-dot"></span>AI</span></Transition>
          </div>
          <p class="field-hint">Describe the current situation, pain points, and comparison with competitors.</p>
          <textarea
            id="whyDemand"
            v-model="form.whyDemand"
            rows="4"
            placeholder="The current situation is..."
            :class="{ 'ai-filled': aiFilledFields.has('whyDemand'), 'has-error': errors['whyDemand'] }"
            @input="aiFilledFields.delete('whyDemand'); clearError('whyDemand')"
          />
          <Transition name="err"><p v-if="errors['whyDemand']" class="field-error" role="alert">{{ errors['whyDemand'] }}</p></Transition>
        </div>

        <div class="field field--full">
          <div class="label-row">
            <label for="whoIsImpacted">Who is impacted? <span class="required">*</span></label>
            <Transition name="badge"><span v-if="aiFilledFields.has('whoIsImpacted')" class="ai-badge"><span class="ai-dot"></span>AI</span></Transition>
          </div>
          <p class="field-hint">Describe the personas and estimated number of users affected.</p>
          <textarea
            id="whoIsImpacted"
            v-model="form.whoIsImpacted"
            rows="3"
            placeholder="This affects approximately..."
            :class="{ 'ai-filled': aiFilledFields.has('whoIsImpacted'), 'has-error': errors['whoIsImpacted'] }"
            @input="aiFilledFields.delete('whoIsImpacted'); clearError('whoIsImpacted')"
          />
          <Transition name="err"><p v-if="errors['whoIsImpacted']" class="field-error" role="alert">{{ errors['whoIsImpacted'] }}</p></Transition>
        </div>

      </div>
    </section>

    <!-- Section 3: Benefits -->
    <section class="form-section">
      <div class="section-header">
        <span class="section-number">03</span>
        <div>
          <h2 class="section-title">Expected Benefits</h2>
          <p class="section-desc">Define success and how you'll measure it</p>
        </div>
      </div>

      <div class="fields-grid">

        <div class="field">
          <div class="label-row">
            <label for="benefitCategory">Benefit Category <span class="required">*</span></label>
            <Transition name="badge"><span v-if="aiFilledFields.has('benefitCategory')" class="ai-badge"><span class="ai-dot"></span>AI</span></Transition>
          </div>
          <div class="select-wrap" :class="{ 'ai-filled': aiFilledFields.has('benefitCategory'), 'has-error': errors['benefitCategory'] }">
            <select id="benefitCategory" v-model="form.benefitCategory" @change="aiFilledFields.delete('benefitCategory'); clearError('benefitCategory')">
              <option value="" disabled>Select a category</option>
              <option v-for="cat in benefitCategories" :key="cat" :value="cat">{{ cat }}</option>
            </select>
            <svg class="select-chevron" viewBox="0 0 20 20" fill="currentColor"><path fill-rule="evenodd" d="M5.23 7.21a.75.75 0 011.06.02L10 11.168l3.71-3.938a.75.75 0 111.08 1.04l-4.25 4.5a.75.75 0 01-1.08 0l-4.25-4.5a.75.75 0 01.02-1.06z" clip-rule="evenodd" /></svg>
          </div>
          <Transition name="err"><p v-if="errors['benefitCategory']" class="field-error" role="alert">{{ errors['benefitCategory'] }}</p></Transition>
        </div>

        <div class="field field--placeholder"></div>

        <div class="field field--full">
          <div class="label-row">
            <label for="benefitHypothesis">What's your hypothesis to achieve those benefits? <span class="required">*</span></label>
            <Transition name="badge"><span v-if="aiFilledFields.has('benefitHypothesis')" class="ai-badge"><span class="ai-dot"></span>AI</span></Transition>
          </div>
          <textarea
            id="benefitHypothesis"
            v-model="form.benefitHypothesis"
            rows="3"
            placeholder="We believe that by doing X we will achieve Y because..."
            :class="{ 'ai-filled': aiFilledFields.has('benefitHypothesis'), 'has-error': errors['benefitHypothesis'] }"
            @input="aiFilledFields.delete('benefitHypothesis'); clearError('benefitHypothesis')"
          />
          <Transition name="err"><p v-if="errors['benefitHypothesis']" class="field-error" role="alert">{{ errors['benefitHypothesis'] }}</p></Transition>
        </div>

        <div class="field field--full">
          <div class="label-row">
            <label for="measureBenefits">How will you measure success? <span class="required">*</span></label>
            <Transition name="badge"><span v-if="aiFilledFields.has('measureBenefits')" class="ai-badge"><span class="ai-dot"></span>AI</span></Transition>
          </div>
          <p class="field-hint">Define KPIs and the timeframe to evaluate results.</p>
          <textarea
            id="measureBenefits"
            v-model="form.measureBenefits"
            rows="3"
            placeholder="We will measure success by tracking..."
            :class="{ 'ai-filled': aiFilledFields.has('measureBenefits'), 'has-error': errors['measureBenefits'] }"
            @input="aiFilledFields.delete('measureBenefits'); clearError('measureBenefits')"
          />
          <Transition name="err"><p v-if="errors['measureBenefits']" class="field-error" role="alert">{{ errors['measureBenefits'] }}</p></Transition>
        </div>

      </div>
    </section>

    <!-- Actions -->
    <div class="form-actions">
      <p class="actions-note"><span class="required">*</span> Required fields</p>
      <button type="submit" class="submit-btn" :disabled="submitting">
        <svg v-if="submitting" class="spinner" viewBox="0 0 24 24" fill="none">
          <circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="3" opacity=".25"/>
          <path d="M12 2a10 10 0 0 1 10 10" stroke="currentColor" stroke-width="3" stroke-linecap="round"/>
        </svg>
        <svg v-else viewBox="0 0 20 20" fill="currentColor">
          <path d="M10.894 2.553a1 1 0 00-1.788 0l-7 14a1 1 0 001.169 1.409l5-1.429A1 1 0 009 15.571V11a1 1 0 112 0v4.571a1 1 0 00.725.962l5 1.428a1 1 0 001.17-1.408l-7-14z" />
        </svg>
        {{ submitting ? 'Submitting...' : 'Submit Request' }}
      </button>
    </div>

  </form>
</template>

<style scoped>
/* ── Layout ── */
.request-form {
  width: 100%;
  max-width: 780px;
  display: flex;
  flex-direction: column;
  gap: 0;
  background: var(--color-neutral-0);
  border-radius: var(--radius-xl);
  border: 1px solid var(--color-neutral-200);
  box-shadow: var(--shadow-lg);
  overflow: hidden;
}

/* ── Form Header ── */
.form-header {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1.75rem 2rem;
  background: var(--gradient-surface);
  border-bottom: 1px solid var(--color-neutral-200);
}

.form-header-icon {
  width: 2.75rem;
  height: 2.75rem;
  border-radius: var(--radius-md);
  background: var(--color-primary-subtle);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  border: 1px solid var(--color-primary-muted);
}

.form-header-icon svg {
  width: 1.25rem;
  height: 1.25rem;
  color: var(--color-primary);
}

.form-title {
  font-size: var(--font-size-xl);
  font-weight: var(--font-weight-bold);
  color: var(--color-neutral-900);
  margin: 0;
  letter-spacing: -0.02em;
  line-height: 1.2;
}

.form-subtitle {
  font-size: var(--font-size-sm);
  color: var(--color-neutral-500);
  margin: 0.2rem 0 0;
  line-height: 1.4;
}

/* ── Sections ── */
.form-section {
  padding: 1.75rem 2rem;
  border-bottom: 1px solid var(--color-neutral-100);
}

.form-section:last-of-type {
  border-bottom: none;
}

.section-header {
  display: flex;
  align-items: flex-start;
  gap: 0.875rem;
  margin-bottom: 1.5rem;
}

.section-number {
  font-size: var(--font-size-xs);
  font-weight: var(--font-weight-bold);
  color: var(--color-primary);
  background: var(--color-primary-subtle);
  border: 1px solid var(--color-primary-muted);
  padding: 0.2rem 0.45rem;
  border-radius: var(--radius-xs);
  letter-spacing: 0.08em;
  flex-shrink: 0;
  margin-top: 3px;
}

.section-title {
  font-size: var(--font-size-md);
  font-weight: var(--font-weight-semibold);
  color: var(--color-neutral-800);
  margin: 0;
  letter-spacing: -0.01em;
}

.section-desc {
  font-size: var(--font-size-sm);
  color: var(--color-neutral-400);
  margin: 0.15rem 0 0;
}

/* ── Fields Grid ── */
.fields-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1.25rem;
}

.field {
  display: flex;
  flex-direction: column;
  gap: 0.35rem;
}

.field--full { grid-column: 1 / -1; }
.field--placeholder { visibility: hidden; }

/* ── Labels ── */
.label-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 0.5rem;
}

label {
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-medium);
  color: var(--color-neutral-700);
  line-height: 1;
}

.required { color: var(--color-error); margin-left: 1px; }

.field-hint {
  font-size: var(--font-size-xs);
  color: var(--color-neutral-400);
  margin: 0;
  line-height: 1.5;
}

/* ── AI Badge ── */
.ai-badge {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  font-size: 0.65rem;
  font-weight: var(--font-weight-bold);
  padding: 0.15rem 0.45rem 0.15rem 0.3rem;
  border-radius: var(--radius-full);
  background: var(--color-ai-bg);
  color: var(--color-ai-text);
  border: 1px solid var(--color-ai-border);
  letter-spacing: 0.04em;
  white-space: nowrap;
}

.ai-dot {
  width: 5px;
  height: 5px;
  border-radius: 50%;
  background: var(--color-ai-dot);
  flex-shrink: 0;
}

/* ── Inputs ── */
input, textarea, select {
  padding: 0.6rem 0.875rem;
  border: 1.5px solid var(--color-neutral-200);
  border-radius: var(--radius-md);
  background: var(--color-neutral-50);
  color: var(--color-neutral-900);
  font-size: var(--font-size-sm);
  font-family: var(--font-family);
  outline: none;
  transition: border-color var(--transition-fast), box-shadow var(--transition-fast), background var(--transition-fast);
  width: 100%;
  box-sizing: border-box;
}

input::placeholder, textarea::placeholder { color: var(--color-neutral-400); }

input:hover:not(:disabled), textarea:hover:not(:disabled), select:hover:not(:disabled) {
  border-color: var(--color-neutral-300);
}

input:focus, textarea:focus {
  border-color: var(--color-primary);
  background: var(--color-neutral-0);
  box-shadow: 0 0 0 3px rgba(0, 135, 74, 0.1);
}

input.has-error, textarea.has-error {
  border-color: var(--color-error);
  background: var(--color-error-subtle);
}

input.has-error:focus, textarea.has-error:focus {
  box-shadow: 0 0 0 3px rgba(230, 57, 70, 0.1);
}

input.ai-filled, textarea.ai-filled {
  border-color: var(--color-ai-border);
  background: rgba(237, 233, 254, 0.3);
  box-shadow: 0 0 0 3px rgba(196, 181, 253, 0.2);
}

textarea {
  resize: vertical;
  min-height: 90px;
  line-height: var(--line-height-relaxed);
}

/* ── Select ── */
.select-wrap {
  position: relative;
  display: flex;
  align-items: center;
}

.select-wrap select {
  appearance: none;
  -webkit-appearance: none;
  padding-right: 2.5rem;
  cursor: pointer;
}

.select-wrap select:focus { border-color: var(--color-primary); box-shadow: 0 0 0 3px rgba(0, 135, 74, 0.1); background: var(--color-neutral-0); }
.select-wrap.has-error select { border-color: var(--color-error); background: var(--color-error-subtle); }
.select-wrap.ai-filled select { border-color: var(--color-ai-border); background: rgba(237, 233, 254, 0.3); }

.select-chevron {
  position: absolute;
  right: 0.7rem;
  width: 1rem;
  height: 1rem;
  color: var(--color-neutral-400);
  pointer-events: none;
  flex-shrink: 0;
}

/* ── Chips ── */
.chips-field {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  padding: 0.5rem;
  border: 1.5px solid var(--color-neutral-200);
  border-radius: var(--radius-md);
  background: var(--color-neutral-50);
  transition: border-color var(--transition-fast), box-shadow var(--transition-fast);
}

.chips-field:focus-within {
  border-color: var(--color-primary);
  box-shadow: 0 0 0 3px rgba(0, 135, 74, 0.1);
  background: var(--color-neutral-0);
}

.chips-field.ai-filled {
  border-color: var(--color-ai-border);
  background: rgba(237, 233, 254, 0.3);
}

.chips-list {
  display: flex;
  flex-wrap: wrap;
  gap: 0.35rem;
}

.chip {
  display: inline-flex;
  align-items: center;
  gap: 0.3rem;
  padding: 0.2rem 0.35rem 0.2rem 0.6rem;
  border-radius: var(--radius-full);
  background: var(--color-neutral-0);
  border: 1px solid var(--color-neutral-300);
  color: var(--color-neutral-700);
  font-size: var(--font-size-xs);
  font-weight: var(--font-weight-medium);
  box-shadow: var(--shadow-xs);
}

.chip--ai {
  background: var(--color-ai-bg);
  border-color: var(--color-ai-border);
  color: var(--color-ai-text);
}

.chip-x {
  width: 1.1rem;
  height: 1.1rem;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(0,0,0,0.06);
  border: none;
  border-radius: 50%;
  cursor: pointer;
  color: inherit;
  padding: 0;
  transition: background var(--transition-fast);
  flex-shrink: 0;
}

.chip-x:hover { background: rgba(0,0,0,0.14); }
.chip-x svg { width: 7px; height: 7px; }

.chips-select { border: none; padding: 0; }
.chips-select select { background: transparent; border: none; box-shadow: none; font-size: var(--font-size-xs); color: var(--color-neutral-500); padding: 0.25rem 1.75rem 0.25rem 0.25rem; }
.chips-select select:focus { box-shadow: none; border: none; }

/* ── Segmented Control ── */
.segmented {
  display: inline-flex;
  background: var(--color-neutral-100);
  border: 1px solid var(--color-neutral-200);
  border-radius: var(--radius-md);
  padding: 3px;
  gap: 2px;
}

.segmented.ai-filled {
  border-color: var(--color-ai-border);
  background: rgba(237, 233, 254, 0.3);
}

.seg-btn {
  padding: 0.4rem 1.1rem;
  border: none;
  border-radius: calc(var(--radius-md) - 3px);
  background: transparent;
  color: var(--color-neutral-500);
  font-size: var(--font-size-sm);
  font-family: var(--font-family);
  font-weight: var(--font-weight-medium);
  cursor: pointer;
  transition: all var(--transition-fast);
}

.seg-btn:hover { color: var(--color-neutral-700); background: rgba(0,0,0,0.04); }

.seg-btn.active {
  background: var(--color-neutral-0);
  color: var(--color-neutral-900);
  box-shadow: var(--shadow-sm);
  font-weight: var(--font-weight-semibold);
}

.seg-btn:focus-visible {
  outline: 2px solid var(--color-primary);
  outline-offset: 1px;
}

/* ── Errors ── */
.field-error {
  font-size: var(--font-size-xs);
  color: var(--color-error);
  margin: 0;
  display: flex;
  align-items: center;
  gap: 0.25rem;
}

/* ── Actions ── */
.form-actions {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 1.25rem 2rem;
  background: var(--color-neutral-50);
  border-top: 1px solid var(--color-neutral-200);
}

.actions-note {
  font-size: var(--font-size-xs);
  color: var(--color-neutral-400);
  margin: 0;
}

.submit-btn {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.65rem 1.5rem;
  background: var(--gradient-primary);
  color: #fff;
  border: none;
  border-radius: var(--radius-md);
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-semibold);
  font-family: var(--font-family);
  cursor: pointer;
  box-shadow: var(--shadow-primary);
  transition: all var(--transition-fast);
  letter-spacing: -0.01em;
}

.submit-btn:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 6px 20px rgba(0, 135, 74, 0.45);
}

.submit-btn:active:not(:disabled) { transform: translateY(0); }

.submit-btn:disabled {
  opacity: 0.65;
  cursor: not-allowed;
  transform: none;
  box-shadow: none;
}

.submit-btn:focus-visible {
  outline: 2px solid var(--color-primary);
  outline-offset: 2px;
}

.submit-btn svg { width: 0.875rem; height: 0.875rem; }

@keyframes spin { to { transform: rotate(360deg); } }
.spinner { animation: spin 0.8s linear infinite; }

/* ── Transitions ── */
.badge-enter-active, .badge-leave-active { transition: opacity var(--transition-fast), transform var(--transition-fast); }
.badge-enter-from, .badge-leave-to { opacity: 0; transform: scale(0.75); }

.err-enter-active, .err-leave-active { transition: opacity var(--transition-fast), transform var(--transition-fast); max-height: 2rem; overflow: hidden; }
.err-enter-from, .err-leave-to { opacity: 0; transform: translateY(-4px); }

@media (prefers-reduced-motion: reduce) {
  *, *::before, *::after { transition: none !important; animation: none !important; }
}

@media (max-width: 640px) {
  .form-header { padding: 1.25rem; }
  .form-section { padding: 1.25rem; }
  .fields-grid { grid-template-columns: 1fr; }
  .field--full { grid-column: 1; }
  .field--placeholder { display: none; }
  .form-actions { padding: 1rem 1.25rem; flex-direction: column; gap: 0.75rem; align-items: stretch; }
  .submit-btn { justify-content: center; }
}
</style>
