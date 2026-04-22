import { ref, watch, type Ref } from 'vue'
import { useChatStore } from '@/stores/chat'
import { useFieldChatStore } from '@/stores/fieldChat'
import type { FormState } from './useFormContext'

export function useAiFormFill(form: Ref<FormState>, errors: Ref<Record<string, string>>) {
  const chatStore = useChatStore()
  const fieldChatStore = useFieldChatStore()

  const aiFilledFields = ref<Set<string>>(new Set())
  const aiUncertainFields = ref<Set<string>>(new Set())

  watch(
    () => chatStore.pendingFormFill,
    (data) => {
      if (!data) return
      const filled = new Set<string>()
      const fieldMap: [keyof FormState, string | undefined][] = [
        ['title', data.title],
        ['businessLine', data.businessLine],
        ['requesterBU', data.requesterBU],
        ['timeSensitive', data.timeSensitive],
        ['whyDemand', data.whyDemand],
        ['whoIsImpacted', data.whoIsImpacted],
        ['benefitCategory', data.benefitCategory],
        ['benefitHypothesis', data.benefitHypothesis],
        ['measureBenefits', data.measureBenefits],
      ]
      for (const [field, value] of fieldMap) {
        if (value) {
          ;(form.value as Record<string, unknown>)[field] = value
          filled.add(field)
        }
      }
      const busIds = Array.isArray(data.busInterested)
        ? data.busInterested
        : data.busInterested
          ? [data.busInterested as unknown as string]
          : []
      if (busIds.length) {
        form.value.busInterested = busIds
        filled.add('busInterested')
      }
      aiFilledFields.value = filled
      aiUncertainFields.value = new Set(data.lowConfidenceFields ?? [])
      chatStore.clearFormFill()
    },
  )

  watch(
    () => fieldChatStore.pendingFieldFill,
    (payload) => {
      if (!payload) return
      const { fieldName, value } = payload
      if (fieldName === 'busInterested') {
        form.value.busInterested = Array.isArray(value) ? value : [value]
      } else if (fieldName in form.value) {
        ;(form.value as Record<string, unknown>)[fieldName] = value
      }
      aiFilledFields.value.add(fieldName)
      aiUncertainFields.value.delete(fieldName)
      fieldChatStore.clearFieldFill()
    },
  )

  function clearField(field: string) {
    aiFilledFields.value.delete(field)
    delete errors.value[field]
  }

  function badgeClass(field: string) {
    return ['ai-badge', { 'ai-badge--uncertain': aiUncertainFields.value.has(field) }]
  }

  return { aiFilledFields, clearField, badgeClass }
}
