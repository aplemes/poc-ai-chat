<script setup lang="ts">
import { useFormContext } from '@/composables/useFormContext'
import { useFieldChatStore } from '@/stores/fieldChat'
import { benefitCategories } from '@/data/formOptions'

const { form, errors, aiFilledFields, clearField, badgeClass } = useFormContext()
const fieldChatStore = useFieldChatStore()
</script>

<template>
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
          <div class="label-actions">
            <button
              type="button"
              class="field-ai-btn"
              :class="{ active: fieldChatStore.activeField === 'benefitCategory' }"
              aria-label="Open AI assistant for Benefit Category"
              @click="fieldChatStore.openPanel('benefitCategory')"
            >
              ✦
            </button>
            <Transition name="badge"
              ><span
                v-if="aiFilledFields.has('benefitCategory')"
                :class="badgeClass('benefitCategory')"
                aria-live="polite"
                aria-label="Field filled by AI"
                ><span class="ai-dot"></span>AI</span
              ></Transition
            >
          </div>
        </div>
        <div
          class="select-wrap"
          :class="{
            'ai-filled': aiFilledFields.has('benefitCategory'),
            'has-error': errors['benefitCategory'],
          }"
        >
          <select
            id="benefitCategory"
            v-model="form.benefitCategory"
            @change="clearField('benefitCategory')"
          >
            <option value="" disabled>Select a category</option>
            <option v-for="cat in benefitCategories" :key="cat" :value="cat">{{ cat }}</option>
          </select>
          <svg class="select-chevron" viewBox="0 0 20 20" fill="currentColor">
            <path
              fill-rule="evenodd"
              d="M5.23 7.21a.75.75 0 011.06.02L10 11.168l3.71-3.938a.75.75 0 111.08 1.04l-4.25 4.5a.75.75 0 01-1.08 0l-4.25-4.5a.75.75 0 01.02-1.06z"
              clip-rule="evenodd"
            />
          </svg>
        </div>
        <Transition name="err"
          ><p v-if="errors['benefitCategory']" class="field-error" role="alert">
            {{ errors['benefitCategory'] }}
          </p></Transition
        >
      </div>

      <div class="field field--placeholder"></div>

      <div class="field field--full">
        <div class="label-row">
          <label for="benefitHypothesis"
            >What's your hypothesis to achieve those benefits?
            <span class="required">*</span></label
          >
          <div class="label-actions">
            <button
              type="button"
              class="field-ai-btn"
              :class="{ active: fieldChatStore.activeField === 'benefitHypothesis' }"
              aria-label="Open AI assistant for Benefit Hypothesis"
              @click="fieldChatStore.openPanel('benefitHypothesis')"
            >
              ✦
            </button>
            <Transition name="badge"
              ><span
                v-if="aiFilledFields.has('benefitHypothesis')"
                :class="badgeClass('benefitHypothesis')"
                aria-live="polite"
                aria-label="Field filled by AI"
                ><span class="ai-dot"></span>AI</span
              ></Transition
            >
          </div>
        </div>
        <textarea
          id="benefitHypothesis"
          v-model="form.benefitHypothesis"
          rows="3"
          placeholder="We believe that by doing X we will achieve Y because..."
          :class="{
            'ai-filled': aiFilledFields.has('benefitHypothesis'),
            'has-error': errors['benefitHypothesis'],
          }"
          @input="clearField('benefitHypothesis')"
        />
        <Transition name="err"
          ><p v-if="errors['benefitHypothesis']" class="field-error" role="alert">
            {{ errors['benefitHypothesis'] }}
          </p></Transition
        >
      </div>

      <div class="field field--full">
        <div class="label-row">
          <label for="measureBenefits"
            >How will you measure success? <span class="required">*</span></label
          >
          <div class="label-actions">
            <button
              type="button"
              class="field-ai-btn"
              :class="{ active: fieldChatStore.activeField === 'measureBenefits' }"
              aria-label="Open AI assistant for Measure Benefits"
              @click="fieldChatStore.openPanel('measureBenefits')"
            >
              ✦
            </button>
            <Transition name="badge"
              ><span
                v-if="aiFilledFields.has('measureBenefits')"
                :class="badgeClass('measureBenefits')"
                aria-live="polite"
                aria-label="Field filled by AI"
                ><span class="ai-dot"></span>AI</span
              ></Transition
            >
          </div>
        </div>
        <p class="field-hint">Define KPIs and the timeframe to evaluate results.</p>
        <textarea
          id="measureBenefits"
          v-model="form.measureBenefits"
          rows="3"
          placeholder="We will measure success by tracking..."
          :class="{
            'ai-filled': aiFilledFields.has('measureBenefits'),
            'has-error': errors['measureBenefits'],
          }"
          @input="clearField('measureBenefits')"
        />
        <Transition name="err"
          ><p v-if="errors['measureBenefits']" class="field-error" role="alert">
            {{ errors['measureBenefits'] }}
          </p></Transition
        >
      </div>
    </div>
  </section>
</template>
