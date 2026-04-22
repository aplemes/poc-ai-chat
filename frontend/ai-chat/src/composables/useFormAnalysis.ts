import { ref } from 'vue'
import { marked } from 'marked'
import DOMPurify from 'dompurify'
import { analyzeForm, type FormFillData } from '@/services/chatService'

export function useFormAnalysis(getFormData: () => FormFillData, validate: () => boolean) {
  const analysisOpen = ref(false)
  const analysisLoading = ref(false)
  const analysisText = ref('')
  const analysisError = ref('')
  const analysisAbort = ref<AbortController | null>(null)

  function renderMd(text: string) {
    return DOMPurify.sanitize(marked.parse(text) as string)
  }

  async function handleSubmit() {
    if (!validate()) return
    const language = localStorage.getItem('chat_language') ?? 'en'
    const ctrl = new AbortController()
    analysisAbort.value = ctrl
    analysisOpen.value = true
    analysisLoading.value = true
    analysisText.value = ''
    analysisError.value = ''
    try {
      await analyzeForm(
        getFormData(),
        language,
        (event) => {
          if (event.type === 'token' && event.content) analysisText.value += event.content
          else if (event.type === 'error') analysisError.value = event.content ?? 'Error'
        },
        ctrl.signal,
      )
    } catch (e) {
      if (e instanceof Error && e.name !== 'AbortError') analysisError.value = 'Could not connect.'
    } finally {
      analysisLoading.value = false
    }
  }

  function closeAnalysis() {
    analysisAbort.value?.abort()
    analysisOpen.value = false
    analysisLoading.value = false
    analysisText.value = ''
    analysisError.value = ''
  }

  return {
    analysisOpen,
    analysisLoading,
    analysisText,
    analysisError,
    handleSubmit,
    closeAnalysis,
    renderMd,
  }
}
