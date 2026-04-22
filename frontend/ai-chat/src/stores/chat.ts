import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { FormFillData } from '@/services/chatService'

export const useChatStore = defineStore('chat', () => {
  const pendingFormFill = ref<FormFillData | null>(null)
  const formFilled = ref(false)

  function setFormFill(data: FormFillData) {
    pendingFormFill.value = data
    formFilled.value = true
  }

  function clearFormFill() {
    pendingFormFill.value = null
    formFilled.value = false
  }

  return { pendingFormFill, formFilled, setFormFill, clearFormFill }
})
