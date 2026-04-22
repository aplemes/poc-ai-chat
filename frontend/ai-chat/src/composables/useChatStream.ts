import { ref, nextTick, onUnmounted } from 'vue'

export interface StreamMessage {
  id: number
  role: 'user' | 'assistant'
  text: string
  type?: 'text' | 'confirm'
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  confirmData?: any
  confirmed?: boolean
}

export interface UseChatStreamOptions<TEvent> {
  /**
   * Called for each SSE event received during streaming.
   * Receives the current assistant placeholder message so the handler can
   * mutate its text in-place, and the addMessage helper for injecting new
   * messages (e.g. correction hints).
   */
  onEvent: (
    event: TEvent,
    assistantMsg: StreamMessage,
    helpers: {
      addMessage: (msg: Omit<StreamMessage, 'id'>) => void
      scrollToBottom: () => Promise<void>
      nextId: () => number
    },
  ) => void
  /**
   * Performs the actual API call. Receives sessionId, text, language, the
   * per-event callback, and the AbortSignal. Returns the new session ID or
   * null.
   */
  sendFn: (
    text: string,
    language: string,
    onEvent: (event: TEvent) => void,
    signal: AbortSignal,
  ) => Promise<string | null>
  messagesEl: () => HTMLElement | null
  textareaEl: () => HTMLTextAreaElement | null
}

export function useChatStream<TEvent>(options: UseChatStreamOptions<TEvent>) {
  const messages = ref<StreamMessage[]>([])
  const input = ref('')
  const loading = ref(false)
  let _nextId = 1
  let abortController: AbortController | null = null

  function nextId() {
    return _nextId++
  }

  async function scrollToBottom() {
    await nextTick()
    const el = options.messagesEl()
    if (el) el.scrollTop = el.scrollHeight
  }

  function autoResize() {
    const el = options.textareaEl()
    if (!el) return
    el.style.height = 'auto'
    el.style.height = Math.min(el.scrollHeight, 140) + 'px'
  }

  function addMessage(msg: Omit<StreamMessage, 'id'>) {
    messages.value.push({ ...msg, id: nextId() })
  }

  async function send(language: string): Promise<string | null> {
    const text = input.value.trim()
    if (!text || loading.value) return null

    addMessage({ role: 'user', text })
    input.value = ''
    loading.value = true
    nextTick(() => {
      const el = options.textareaEl()
      if (el) el.style.height = 'auto'
    })
    await scrollToBottom()

    const assistantMsg: StreamMessage = { id: nextId(), role: 'assistant', text: '' }
    messages.value.push(assistantMsg)

    abortController?.abort()
    abortController = new AbortController()

    try {
      const newSessionId = await options.sendFn(
        text,
        language,
        (event) => {
          options.onEvent(event, assistantMsg, { addMessage, scrollToBottom, nextId })
        },
        abortController.signal,
      )
      return newSessionId
    } finally {
      loading.value = false
      abortController = null
    }
  }

  function abort() {
    abortController?.abort()
  }

  onUnmounted(() => abortController?.abort())

  return {
    messages,
    input,
    loading,
    nextId,
    scrollToBottom,
    autoResize,
    addMessage,
    send,
    abort,
  }
}
