import { inject, type InjectionKey, type Ref } from 'vue'

export interface FormState {
  title: string
  businessLine: string
  requesterBU: string
  busInterested: string[]
  timeSensitive: string
  whyDemand: string
  whoIsImpacted: string
  benefitCategory: string
  benefitHypothesis: string
  measureBenefits: string
}

export interface FormContext {
  form: Ref<FormState>
  errors: Ref<Record<string, string>>
  aiFilledFields: Ref<Set<string>>
  clearField: (field: string) => void
  badgeClass: (field: string) => (string | Record<string, boolean>)[]
}

export const FORM_CONTEXT_KEY: InjectionKey<FormContext> = Symbol('formContext')

export function useFormContext(): FormContext {
  const ctx = inject(FORM_CONTEXT_KEY)
  if (!ctx) throw new Error('useFormContext must be used inside RequestForm')
  return ctx
}
