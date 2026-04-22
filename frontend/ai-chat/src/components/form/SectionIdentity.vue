<script setup lang="ts">
import { useFormContext } from '@/composables/useFormContext'
import { useFieldChatStore } from '@/stores/fieldChat'
import { requesters, concerned, organizations, timeSensitiveOptions } from '@/data/formOptions'

const { form, errors, aiFilledFields, clearField, badgeClass } = useFormContext()
const fieldChatStore = useFieldChatStore()

function selectTimeSensitive(opt: string) {
  form.value.timeSensitive = opt
  clearField('timeSensitive')
}

function toggleBusInterested(id: string) {
  const idx = form.value.busInterested.indexOf(id)
  form.value.busInterested =
    idx === -1
      ? [...form.value.busInterested, id]
      : form.value.busInterested.filter((v) => v !== id)
  clearField('busInterested')
}

function removeBusInterested(id: string) {
  form.value.busInterested = form.value.busInterested.filter((v) => v !== id)
  clearField('busInterested')
}

function handleBusInterestedSelect(e: Event) {
  const select = e.target as HTMLSelectElement
  const value = select.value
  if (value) toggleBusInterested(value)
  select.value = ''
}

function getBuLabel(id: string): string {
  return concerned.find((c) => c.id === id)?.value ?? id
}
</script>

<template>
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
          <div class="label-actions">
            <button
              type="button"
              class="field-ai-btn"
              :class="{ active: fieldChatStore.activeField === 'title' }"
              aria-label="Open AI assistant for Title"
              @click="fieldChatStore.openPanel('title')"
            >
              ✦
            </button>
            <Transition name="badge"
              ><span
                v-if="aiFilledFields.has('title')"
                :class="badgeClass('title')"
                aria-live="polite"
                aria-label="Field filled by AI"
                ><span class="ai-dot"></span>AI</span
              ></Transition
            >
          </div>
        </div>
        <p class="field-hint">
          Start with an infinitive verb + scope. e.g. "Add the new Payment Method X on website"
        </p>
        <input
          id="title"
          v-model="form.title"
          type="text"
          placeholder="Add the new... / Improve the... / Enable..."
          :class="{ 'ai-filled': aiFilledFields.has('title'), 'has-error': errors['title'] }"
          @input="clearField('title')"
        />
        <Transition name="err"
          ><p v-if="errors['title']" class="field-error" role="alert">
            {{ errors['title'] }}
          </p></Transition
        >
      </div>

      <!-- Business Line -->
      <div class="field">
        <div class="label-row">
          <label for="businessLine">Business Line <span class="required">*</span></label>
          <div class="label-actions">
            <button
              type="button"
              class="field-ai-btn"
              :class="{ active: fieldChatStore.activeField === 'businessLine' }"
              aria-label="Open AI assistant for Business Line"
              @click="fieldChatStore.openPanel('businessLine')"
            >
              ✦
            </button>
            <Transition name="badge"
              ><span
                v-if="aiFilledFields.has('businessLine')"
                :class="badgeClass('businessLine')"
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
            'ai-filled': aiFilledFields.has('businessLine'),
            'has-error': errors['businessLine'],
          }"
        >
          <select
            id="businessLine"
            v-model="form.businessLine"
            @change="clearField('businessLine')"
          >
            <option value="" disabled>Select business line</option>
            <option v-for="org in organizations" :key="org.id" :value="org.id">
              {{ org.value }}
            </option>
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
          ><p v-if="errors['businessLine']" class="field-error" role="alert">
            {{ errors['businessLine'] }}
          </p></Transition
        >
      </div>

      <!-- Requester BU -->
      <div class="field">
        <div class="label-row">
          <label for="requesterBU">Requester BU <span class="required">*</span></label>
          <div class="label-actions">
            <button
              type="button"
              class="field-ai-btn"
              :class="{ active: fieldChatStore.activeField === 'requesterBU' }"
              aria-label="Open AI assistant for Requester BU"
              @click="fieldChatStore.openPanel('requesterBU')"
            >
              ✦
            </button>
            <Transition name="badge"
              ><span
                v-if="aiFilledFields.has('requesterBU')"
                :class="badgeClass('requesterBU')"
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
            'ai-filled': aiFilledFields.has('requesterBU'),
            'has-error': errors['requesterBU'],
          }"
        >
          <select id="requesterBU" v-model="form.requesterBU" @change="clearField('requesterBU')">
            <option value="" disabled>Select your BU</option>
            <option v-for="req in requesters" :key="req.id" :value="req.id">{{ req.value }}</option>
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
          ><p v-if="errors['requesterBU']" class="field-error" role="alert">
            {{ errors['requesterBU'] }}
          </p></Transition
        >
      </div>

      <!-- BUs Interested -->
      <div class="field field--full">
        <div class="label-row">
          <label id="busInterested-label">BUs Interested</label>
          <div class="label-actions">
            <button
              type="button"
              class="field-ai-btn"
              :class="{ active: fieldChatStore.activeField === 'busInterested' }"
              aria-label="Open AI assistant for BUs Interested"
              @click="fieldChatStore.openPanel('busInterested')"
            >
              ✦
            </button>
            <Transition name="badge"
              ><span
                v-if="aiFilledFields.has('busInterested')"
                :class="badgeClass('busInterested')"
                aria-live="polite"
                aria-label="Field filled by AI"
                ><span class="ai-dot"></span>AI</span
              ></Transition
            >
          </div>
        </div>
        <div class="chips-field" :class="{ 'ai-filled': aiFilledFields.has('busInterested') }">
          <div v-if="form.busInterested.length" class="chips-list">
            <span
              v-for="id in form.busInterested"
              :key="id"
              class="chip"
              :class="{ 'chip--ai': aiFilledFields.has('busInterested') }"
            >
              {{ getBuLabel(id) }}
              <button
                type="button"
                class="chip-x"
                :aria-label="`Remove ${getBuLabel(id)}`"
                @click="removeBusInterested(id)"
              >
                <svg viewBox="0 0 12 12" fill="currentColor">
                  <path
                    d="M1 1l10 10M11 1L1 11"
                    stroke="currentColor"
                    stroke-width="2"
                    stroke-linecap="round"
                  />
                </svg>
              </button>
            </span>
          </div>
          <div class="select-wrap chips-select">
            <select aria-labelledby="busInterested-label" @change="handleBusInterestedSelect">
              <option value="" disabled>Add a BU...</option>
              <option
                v-for="bu in concerned"
                :key="bu.id"
                :value="bu.id"
                :disabled="form.busInterested.includes(bu.id)"
              >
                {{ bu.value }}
              </option>
            </select>
            <svg class="select-chevron" viewBox="0 0 20 20" fill="currentColor">
              <path
                fill-rule="evenodd"
                d="M5.23 7.21a.75.75 0 011.06.02L10 11.168l3.71-3.938a.75.75 0 111.08 1.04l-4.25 4.5a.75.75 0 01-1.08 0l-4.25-4.5a.75.75 0 01.02-1.06z"
                clip-rule="evenodd"
              />
            </svg>
          </div>
        </div>
      </div>

      <!-- Time Sensitive -->
      <div class="field field--full">
        <div class="label-row">
          <label id="timeSensitive-label">Is this demand time sensitive?</label>
          <div class="label-actions">
            <button
              type="button"
              class="field-ai-btn"
              :class="{ active: fieldChatStore.activeField === 'timeSensitive' }"
              aria-label="Open AI assistant for Time Sensitive"
              @click="fieldChatStore.openPanel('timeSensitive')"
            >
              ✦
            </button>
            <Transition name="badge"
              ><span
                v-if="aiFilledFields.has('timeSensitive')"
                :class="badgeClass('timeSensitive')"
                aria-live="polite"
                aria-label="Field filled by AI"
                ><span class="ai-dot"></span>AI</span
              ></Transition
            >
          </div>
        </div>
        <div
          class="segmented"
          :class="{ 'ai-filled': aiFilledFields.has('timeSensitive') }"
          role="group"
          aria-labelledby="timeSensitive-label"
        >
          <button
            v-for="opt in timeSensitiveOptions"
            :key="opt"
            type="button"
            class="seg-btn"
            :class="{ active: form.timeSensitive === opt }"
            :aria-pressed="form.timeSensitive === opt"
            @click="selectTimeSensitive(opt)"
          >
            {{ opt }}
          </button>
        </div>
      </div>
    </div>
  </section>
</template>
