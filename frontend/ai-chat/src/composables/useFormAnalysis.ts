import { ref } from 'vue'
import { analyzeForm, type FormFillData } from '@/services/chatService'
import { renderMarkdown } from '@/utils/markdown'

export function useFormAnalysis(
  getFormData: () => FormFillData,
  validate: () => boolean,
  onConfirmedSubmit?: () => void,
) {
  const analysisOpen = ref(false)
  const analysisLoading = ref(false)
  const analysisText = ref('')
  const analysisError = ref('')
  const analysisAbort = ref<AbortController | null>(null)

  async function handleSubmit() {
    if (!validate()) return
    // Cancel any in-flight analysis before starting a new one (FE-07)
    analysisAbort.value?.abort()
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

  // Called when user confirms submission from the analysis modal (UX-11)
  function confirmSubmit() {
    closeAnalysis()
    onConfirmedSubmit?.()
  }

  return {
    analysisOpen,
    analysisLoading,
    analysisText,
    analysisError,
    handleSubmit,
    closeAnalysis,
    confirmSubmit,
    renderMd: renderMarkdown,
  }
}
