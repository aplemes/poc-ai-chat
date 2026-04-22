<script setup lang="ts">
import { useFormContext } from '@/composables/useFormContext'
import FieldLabel from './FieldLabel.vue'

const { form, errors, aiFilledFields, clearField } = useFormContext()
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
        <FieldLabel field-name="whyDemand" label="Why are you making this demand?" required />
        <p class="field-hint">
          Describe the triggering event, the current situation and tools in use, and the specific pain points.
        </p>
        <textarea
          id="whyDemand"
          v-model="form.whyDemand"
          rows="4"
          placeholder="The current situation is..."
          :class="{
            'ai-filled': aiFilledFields.has('whyDemand'),
            'has-error': errors['whyDemand'],
          }"
          @input="clearField('whyDemand')"
        />
        <Transition name="err"
          ><p v-if="errors['whyDemand']" class="field-error" role="alert">
            {{ errors['whyDemand'] }}
          </p></Transition
        >
      </div>

      <div class="field field--full">
        <FieldLabel field-name="whoIsImpacted" label="Who is impacted?" required />
        <p class="field-hint">Describe the personas and estimated number of users affected.</p>
        <textarea
          id="whoIsImpacted"
          v-model="form.whoIsImpacted"
          rows="3"
          placeholder="This affects approximately..."
          :class="{
            'ai-filled': aiFilledFields.has('whoIsImpacted'),
            'has-error': errors['whoIsImpacted'],
          }"
          @input="clearField('whoIsImpacted')"
        />
        <Transition name="err"
          ><p v-if="errors['whoIsImpacted']" class="field-error" role="alert">
            {{ errors['whoIsImpacted'] }}
          </p></Transition
        >
      </div>
    </div>
  </section>
</template>
