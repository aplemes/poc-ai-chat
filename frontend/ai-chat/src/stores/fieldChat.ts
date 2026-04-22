import { defineStore } from 'pinia'
import { ref } from 'vue'

export interface FieldFillPayload {
  fieldName: string
  value: string | string[]
}

export const useFieldChatStore = defineStore('fieldChat', () => {
  const activeField = ref<string | null>(null)
  // Persists session IDs per field so conversation is remembered while panel reopens
  const fieldSessionIds = ref<Map<string, string>>(new Map())
  const pendingFieldFill = ref<FieldFillPayload | null>(null)

  function openPanel(fieldName: string) {
    activeField.value = fieldName
  }

  function closePanel() {
    activeField.value = null
  }

  function setSessionId(fieldName: string, sessionId: string) {
    fieldSessionIds.value.set(fieldName, sessionId)
  }

  function getSessionId(fieldName: string): string | null {
    return fieldSessionIds.value.get(fieldName) ?? null
  }

  function setFieldFill(payload: FieldFillPayload) {
    pendingFieldFill.value = payload
  }

  function clearFieldFill() {
    pendingFieldFill.value = null
  }

  return {
    activeField,
    pendingFieldFill,
    openPanel,
    closePanel,
    setSessionId,
    getSessionId,
    setFieldFill,
    clearFieldFill,
  }
})
