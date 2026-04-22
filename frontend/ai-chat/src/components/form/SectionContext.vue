<script setup lang="ts">
import { useFormContext } from '@/composables/useFormContext'
import { useFieldChatStore } from '@/stores/fieldChat'

const { form, errors, aiFilledFields, clearField, badgeClass } = useFormContext()
const fieldChatStore = useFieldChatStore()
</script>

<template>
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
          <div class="label-actions">
            <button type="button" class="field-ai-btn" :class="{ active: fieldChatStore.activeField === 'whyDemand' }" @click="fieldChatStore.openPanel('whyDemand')">✦</button>
            <Transition name="badge"><span v-if="aiFilledFields.has('whyDemand')" :class="badgeClass('whyDemand')"><span class="ai-dot"></span>AI</span></Transition>
          </div>
        </div>
        <p class="field-hint">Describe the current situation, pain points, and comparison with competitors.</p>
        <textarea
          id="whyDemand"
          v-model="form.whyDemand"
          rows="4"
          placeholder="The current situation is..."
          :class="{ 'ai-filled': aiFilledFields.has('whyDemand'), 'has-error': errors['whyDemand'] }"
          @input="clearField('whyDemand')"
        />
        <Transition name="err"><p v-if="errors['whyDemand']" class="field-error" role="alert">{{ errors['whyDemand'] }}</p></Transition>
      </div>

      <div class="field field--full">
        <div class="label-row">
          <label for="whoIsImpacted">Who is impacted? <span class="required">*</span></label>
          <div class="label-actions">
            <button type="button" class="field-ai-btn" :class="{ active: fieldChatStore.activeField === 'whoIsImpacted' }" @click="fieldChatStore.openPanel('whoIsImpacted')">✦</button>
            <Transition name="badge"><span v-if="aiFilledFields.has('whoIsImpacted')" :class="badgeClass('whoIsImpacted')"><span class="ai-dot"></span>AI</span></Transition>
          </div>
        </div>
        <p class="field-hint">Describe the personas and estimated number of users affected.</p>
        <textarea
          id="whoIsImpacted"
          v-model="form.whoIsImpacted"
          rows="3"
          placeholder="This affects approximately..."
          :class="{ 'ai-filled': aiFilledFields.has('whoIsImpacted'), 'has-error': errors['whoIsImpacted'] }"
          @input="clearField('whoIsImpacted')"
        />
        <Transition name="err"><p v-if="errors['whoIsImpacted']" class="field-error" role="alert">{{ errors['whoIsImpacted'] }}</p></Transition>
      </div>
    </div>
  </section>
</template>
