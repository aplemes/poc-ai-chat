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

const chatStore = useChatStore()

const form = ref({
  title: '',
  demandScope: '' as 'Intra-BU' | 'Adeo Platform' | '',
  businessLine: '',
  requesterBU: '',
  busInterested: [] as string[],
  demandContext: '',
  currentSituation: '',
  problemsToSolve: '',
  whoIsImpacted: '',
  measureBenefits: '',
})

const aiFilledFields = ref<Set<string>>(new Set())

watch(
  () => chatStore.pendingFormFill,
  (data) => {
    if (!data) return

    const filled = new Set<string>()

    if (data.title) { form.value.title = data.title; filled.add('title') }
    if (data.demandScope) { form.value.demandScope = data.demandScope as 'Intra-BU' | 'Adeo Platform'; filled.add('demandScope') }
    if (data.businessLine) { form.value.businessLine = data.businessLine; filled.add('businessLine') }
    if (data.requesterBU) { form.value.requesterBU = data.requesterBU; filled.add('requesterBU') }
    if (data.busInterested?.length) { form.value.busInterested = data.busInterested; filled.add('busInterested') }
    if (data.demandContext) { form.value.demandContext = data.demandContext; filled.add('demandContext') }
    if (data.currentSituation) { form.value.currentSituation = data.currentSituation; filled.add('currentSituation') }
    if (data.problemsToSolve) { form.value.problemsToSolve = data.problemsToSolve; filled.add('problemsToSolve') }
    if (data.whoIsImpacted) { form.value.whoIsImpacted = data.whoIsImpacted; filled.add('whoIsImpacted') }
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

    <div class="field field--full">
      <label for="title">
        Title
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
      <label for="demandScope">
        For whom is this demand?
        <span v-if="aiFilledFields.has('demandScope')" class="ai-badge">IA</span>
      </label>
      <p class="field-hint">Intra-BU: stays within your BU. Adeo Platform: sent to Adeo Platform for analysis.</p>
      <select
        id="demandScope"
        v-model="form.demandScope"
        :class="{ 'ai-filled': aiFilledFields.has('demandScope') }"
        @change="aiFilledFields.delete('demandScope')"
      >
        <option value="" disabled>Select scope</option>
        <option value="Intra-BU">Intra-BU</option>
        <option value="Adeo Platform">Adeo Platform</option>
      </select>
    </div>

    <div class="field">
      <label for="businessLine">
        Business Line
        <span v-if="aiFilledFields.has('businessLine')" class="ai-badge">IA</span>
      </label>
      <select
        id="businessLine"
        v-model="form.businessLine"
        :class="{ 'ai-filled': aiFilledFields.has('businessLine') }"
        @change="aiFilledFields.delete('businessLine')"
      >
        <option value="" disabled>Select a business line</option>
        <option v-for="org in organizations" :key="org.id" :value="org.id">
          {{ org.value }}
        </option>
      </select>
    </div>

    <div class="field">
      <label for="requesterBU">
        Requester BU
        <span v-if="aiFilledFields.has('requesterBU')" class="ai-badge">IA</span>
      </label>
      <select
        id="requesterBU"
        v-model="form.requesterBU"
        :class="{ 'ai-filled': aiFilledFields.has('requesterBU') }"
        @change="aiFilledFields.delete('requesterBU')"
      >
        <option value="" disabled>Select a requester BU</option>
        <option v-for="req in requesters" :key="req.id" :value="req.id">
          {{ req.value }}
        </option>
      </select>
    </div>

    <div class="field">
      <label for="busInterested">
        BUs Interested
        <span v-if="aiFilledFields.has('busInterested')" class="ai-badge">IA</span>
      </label>
      <select
        id="busInterested"
        v-model="form.busInterested"
        multiple
        :class="{ 'ai-filled': aiFilledFields.has('busInterested') }"
        @change="aiFilledFields.delete('busInterested')"
      >
        <option v-for="bu in concerned" :key="bu.id" :value="bu.id">
          {{ bu.value }}
        </option>
      </select>
    </div>

    <div class="field-group field--full">
      <h3 class="group-title">Why do you make this demand?</h3>

      <div class="field">
        <label for="demandContext">
          1. The precise context
          <span v-if="aiFilledFields.has('demandContext')" class="ai-badge">IA</span>
        </label>
        <p class="field-hint">What is the event or situation that motivates it?</p>
        <textarea
          id="demandContext"
          v-model="form.demandContext"
          rows="3"
          placeholder="Describe the event or situation that triggers this demand..."
          :class="{ 'ai-filled': aiFilledFields.has('demandContext') }"
          @input="aiFilledFields.delete('demandContext')"
        />
      </div>

      <div class="field">
        <label for="currentSituation">
          2. The current situation
          <span v-if="aiFilledFields.has('currentSituation')" class="ai-badge">IA</span>
        </label>
        <p class="field-hint">How are things going today? What processes or tools are being used?</p>
        <textarea
          id="currentSituation"
          v-model="form.currentSituation"
          rows="3"
          placeholder="Describe the current state, tools and processes..."
          :class="{ 'ai-filled': aiFilledFields.has('currentSituation') }"
          @input="aiFilledFields.delete('currentSituation')"
        />
      </div>

      <div class="field">
        <label for="problemsToSolve">
          3. The problems to solve
          <span v-if="aiFilledFields.has('problemsToSolve')" class="ai-badge">IA</span>
        </label>
        <p class="field-hint">What are the sticking points, frustrations or inefficiencies you encounter?</p>
        <textarea
          id="problemsToSolve"
          v-model="form.problemsToSolve"
          rows="3"
          placeholder="List the main pain points and inefficiencies..."
          :class="{ 'ai-filled': aiFilledFields.has('problemsToSolve') }"
          @input="aiFilledFields.delete('problemsToSolve')"
        />
      </div>
    </div>

    <div class="field">
      <label for="whoIsImpacted">
        Who is impacted by this demand?
        <span v-if="aiFilledFields.has('whoIsImpacted')" class="ai-badge">IA</span>
      </label>
      <p class="field-hint">List the personas and estimated numbers. e.g. "Customers online: 5% of total = 50,000. Coworkers: not concerned."</p>
      <textarea
        id="whoIsImpacted"
        v-model="form.whoIsImpacted"
        rows="3"
        placeholder="Customers: ... / Coworkers: ... / Partners: ..."
        :class="{ 'ai-filled': aiFilledFields.has('whoIsImpacted') }"
        @input="aiFilledFields.delete('whoIsImpacted')"
      />
    </div>

    <div class="field">
      <label for="measureBenefits">
        How will you measure those benefits?
        <span v-if="aiFilledFields.has('measureBenefits')" class="ai-badge">IA</span>
      </label>
      <p class="field-hint">Provide the metrics and timing. e.g. "GMV per payment method during the first 3 months after activation."</p>
      <textarea
        id="measureBenefits"
        v-model="form.measureBenefits"
        rows="3"
        placeholder="Metric: ... / Timing: ..."
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

select[multiple] {
  min-height: 120px;
  padding: 0.25rem 0;
}

select[multiple] option {
  padding: 0.35rem 0.75rem;
}

textarea {
  resize: vertical;
  min-height: 80px;
  line-height: 1.5;
}

.field-hint {
  font-size: 0.78rem;
  color: #6b7280;
  margin: 0;
  line-height: 1.4;
}

.field-group {
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
  padding: 1.25rem;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  background: #f8fafc;
}

.group-title {
  font-size: 0.95rem;
  font-weight: 600;
  color: #374151;
  margin: 0 0 0.25rem;
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
