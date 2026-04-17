import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { FormFillData } from '@/services/chatService'

export const useChatStore = defineStore('chat', () => {
  const pendingFormFill = ref<FormFillData | null>(null)

  function setFormFill(data: FormFillData) {
    pendingFormFill.value = data
  }

  function clearFormFill() {
    pendingFormFill.value = null
  }

  return { pendingFormFill, setFormFill, clearFormFill }
})
