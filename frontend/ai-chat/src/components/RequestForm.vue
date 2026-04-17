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
  busInterested: '',
  timeSensitive: 'No' as string,
  whyDemand: '',
  whoIsImpacted: '',
  benefitCategory: '',
  benefitHypothesis: '',
  measureBenefits: '',
})

const aiFilledFields = ref<Set<string>>(new Set())

watch(
  () => chatStore.pendingFormFill,
  (data) => {
    if (!data) return

    const filled = new Set<string>()

    if (data.title) { form.value.title = data.title; filled.add('title') }
    if (data.businessLine) { form.value.businessLine = data.businessLine; filled.add('businessLine') }
    if (data.requesterBU) { form.value.requesterBU = data.requesterBU; filled.add('requesterBU') }
    if (data.busInterested?.length) { form.value.busInterested = data.busInterested[0] ?? ''; filled.add('busInterested') }
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

function handleSubmit() {
  console.log(form.value)
}
</script>

<template>
  <form class="request-form" @submit.prevent="handleSubmit">
    <h2 class="form-title">New Request</h2>

    <!-- Tell us more about this demand -->
    <div class="field field--full">
      <label for="title">
        Title of the demand <span class="required">- required</span>
        <span v-if="aiFilledFields.has('title')" class="ai-badge">IA</span>
      </label>
      <p class="field-hint">Start with an infinitive verb + scope. e.g. "Add the new Payment Method 'XXX' on the website only"</p>
      <input
        id="title"
        v-model="form.title"
        type="text"
        placeholder="Add the new... / Improve the... / Enable..."
        :class="{ 'ai-filled': aiFilledFields.has('title') }"
        @input="aiFilledFields.delete('title')"
      />
    </div>

    <div class="field">
      <label for="businessLine">
        Business Line <span class="required">- required</span>
        <span v-if="aiFilledFields.has('businessLine')" class="ai-badge">IA</span>
      </label>
      <select
        id="businessLine"
        v-model="form.businessLine"
        :class="{ 'ai-filled': aiFilledFields.has('businessLine') }"
        @change="aiFilledFields.delete('businessLine')"
      >
        <option value="" disabled>-- Choose an option --</option>
        <option v-for="org in organizations" :key="org.id" :value="org.id">
          {{ org.value }}
        </option>
      </select>
    </div>

    <div class="field">
      <label for="requesterBU">
        Requester BU <span class="required">- required</span>
        <span v-if="aiFilledFields.has('requesterBU')" class="ai-badge">IA</span>
      </label>
      <select
        id="requesterBU"
        v-model="form.requesterBU"
        :class="{ 'ai-filled': aiFilledFields.has('requesterBU') }"
        @change="aiFilledFields.delete('requesterBU')"
      >
        <option value="" disabled>-- Choose an option --</option>
        <option v-for="req in requesters" :key="req.id" :value="req.id">
          {{ req.value }}
        </option>
      </select>
    </div>

    <div class="field field--full">
      <label for="busInterested">
        BUs interested <span class="required">- required</span>
        <span v-if="aiFilledFields.has('busInterested')" class="ai-badge">IA</span>
      </label>
      <select
        id="busInterested"
        v-model="form.busInterested"
        :class="{ 'ai-filled': aiFilledFields.has('busInterested') }"
        @change="aiFilledFields.delete('busInterested')"
      >
        <option value="" disabled>-- Choose an option --</option>
        <option v-for="bu in concerned" :key="bu.id" :value="bu.id">
          {{ bu.value }}
        </option>
      </select>
    </div>

    <!-- Is your demand time sensitive? -->
    <div class="field field--full">
      <label>
        Is your demand time sensitive?
        <span v-if="aiFilledFields.has('timeSensitive')" class="ai-badge">IA</span>
      </label>
      <div class="toggle-group" :class="{ 'ai-filled': aiFilledFields.has('timeSensitive') }">
        <button
          v-for="opt in timeSensitiveOptions"
          :key="opt"
          type="button"
          class="toggle-btn"
          :class="{ 'toggle-btn--active': form.timeSensitive === opt }"
          @click="form.timeSensitive = opt; aiFilledFields.delete('timeSensitive')"
        >
          {{ opt }}
        </button>
      </div>
    </div>

    <!-- Why do you make this demand? -->
    <div class="field field--full">
      <label for="whyDemand">
        Describe current situation, the pain points, comparison with competitors...
        <span class="required">- required</span>
        <span v-if="aiFilledFields.has('whyDemand')" class="ai-badge">IA</span>
      </label>
      <textarea
        id="whyDemand"
        v-model="form.whyDemand"
        rows="5"
        placeholder="Insert your text here"
        :class="{ 'ai-filled': aiFilledFields.has('whyDemand') }"
        @input="aiFilledFields.delete('whyDemand')"
      />
    </div>

    <!-- Who is impacted -->
    <div class="field field--full">
      <label for="whoIsImpacted">
        Describe who is impacted <span class="required">- required</span>
        <span v-if="aiFilledFields.has('whoIsImpacted')" class="ai-badge">IA</span>
      </label>
      <textarea
        id="whoIsImpacted"
        v-model="form.whoIsImpacted"
        rows="4"
        placeholder="Insert your text here"
        :class="{ 'ai-filled': aiFilledFields.has('whoIsImpacted') }"
        @input="aiFilledFields.delete('whoIsImpacted')"
      />
    </div>

    <!-- What are the benefits? -->
    <div class="field">
      <label for="benefitCategory">
        Select a category <span class="required">- required</span>
        <span v-if="aiFilledFields.has('benefitCategory')" class="ai-badge">IA</span>
      </label>
      <select
        id="benefitCategory"
        v-model="form.benefitCategory"
        :class="{ 'ai-filled': aiFilledFields.has('benefitCategory') }"
        @change="aiFilledFields.delete('benefitCategory')"
      >
        <option value="" disabled>-- Choose an option --</option>
        <option v-for="cat in benefitCategories" :key="cat" :value="cat">
          {{ cat }}
        </option>
      </select>
    </div>

    <div class="field field--full">
      <label for="benefitHypothesis">
        What are your hypothesis to meet those benefits? <span class="required">- required</span>
        <span v-if="aiFilledFields.has('benefitHypothesis')" class="ai-badge">IA</span>
      </label>
      <textarea
        id="benefitHypothesis"
        v-model="form.benefitHypothesis"
        rows="4"
        placeholder="Insert your text here"
        :class="{ 'ai-filled': aiFilledFields.has('benefitHypothesis') }"
        @input="aiFilledFields.delete('benefitHypothesis')"
      />
    </div>

    <!-- How will you measure those benefits? -->
    <div class="field field--full">
      <label for="measureBenefits">
        What are the KPI to measure? <span class="required">- required</span>
        <span v-if="aiFilledFields.has('measureBenefits')" class="ai-badge">IA</span>
      </label>
      <textarea
        id="measureBenefits"
        v-model="form.measureBenefits"
        rows="4"
        placeholder="Insert your text here"
        :class="{ 'ai-filled': aiFilledFields.has('measureBenefits') }"
        @input="aiFilledFields.delete('measureBenefits')"
      />
    </div>

    <div class="form-actions field--full">
      <button type="submit">Submit</button>
    </div>
  </form>
</template>

<style scoped>
.request-form {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1.25rem;
  max-width: 860px;
  width: 100%;
  padding: 2rem;
  background: #ffffff;
  border-radius: 10px;
  border: 1px solid #e2e8f0;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.06);
}

.form-title {
  grid-column: 1 / -1;
  font-size: 1.25rem;
  font-weight: 600;
  color: #1e293b;
  margin: 0;
}

.field {
  display: flex;
  flex-direction: column;
  gap: 0.4rem;
}

.field--full {
  grid-column: 1 / -1;
}

label {
  display: flex;
  align-items: center;
  gap: 0.4rem;
  font-size: 0.875rem;
  font-weight: 500;
  color: #374151;
}

.ai-badge {
  font-size: 0.65rem;
  font-weight: 700;
  padding: 0.1rem 0.35rem;
  border-radius: 4px;
  background: hsla(160, 100%, 37%, 0.12);
  color: hsla(160, 100%, 28%, 1);
  letter-spacing: 0.03em;
}

input,
select,
textarea {
  padding: 0.5rem 0.75rem;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  background: #f9fafb;
  color: #1e293b;
  font-size: 0.9rem;
  font-family: inherit;
  outline: none;
  transition: border-color 0.2s, box-shadow 0.2s, background 0.2s;
}

input:focus,
select:focus,
textarea:focus {
  border-color: #3b82f6;
  background: #ffffff;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.ai-filled {
  border-color: hsla(160, 100%, 37%, 0.7);
  background: hsla(160, 100%, 37%, 0.04);
  box-shadow: 0 0 0 3px hsla(160, 100%, 37%, 0.08);
}

textarea {
  resize: vertical;
  min-height: 80px;
  line-height: 1.5;
}

.required {
  font-weight: 400;
  color: #e67e22;
  font-size: 0.8rem;
}

.field-hint {
  font-size: 0.78rem;
  color: #6b7280;
  margin: 0;
  line-height: 1.4;
}

.toggle-group {
  display: flex;
  gap: 0.5rem;
  flex-wrap: wrap;
}

.toggle-btn {
  padding: 0.45rem 1.25rem;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  background: #f9fafb;
  color: #374151;
  font-size: 0.9rem;
  font-family: inherit;
  cursor: pointer;
  transition: border-color 0.15s, background 0.15s, color 0.15s;
}

.toggle-btn:hover {
  border-color: #9ca3af;
}

.toggle-btn--active {
  border-color: #0d9488;
  color: #0d9488;
  background: #f0fdfa;
  font-weight: 500;
}

.toggle-group.ai-filled .toggle-btn--active {
  border-color: hsla(160, 100%, 37%, 0.7);
  background: hsla(160, 100%, 37%, 0.04);
  box-shadow: 0 0 0 3px hsla(160, 100%, 37%, 0.08);
}

.form-actions {
  display: flex;
  justify-content: flex-end;
}

button[type='submit'] {
  padding: 0.55rem 1.75rem;
  background: hsla(160, 100%, 37%, 1);
  color: #fff;
  border: none;
  border-radius: 6px;
  font-size: 0.9rem;
  font-weight: 600;
  cursor: pointer;
  transition: background 0.2s;
}

button[type='submit']:hover {
  background: hsla(160, 100%, 30%, 1);
}

@media (max-width: 700px) {
  .request-form {
    grid-template-columns: 1fr;
    padding: 1.25rem;
    border-radius: 0;
    border-left: none;
    border-right: none;
    box-shadow: none;
  }

  .field--full {
    grid-column: 1;
  }

  .form-actions {
    justify-content: stretch;
  }

  button[type='submit'] {
    width: 100%;
  }
}
</style>
