const API_BASE = 'http://localhost:8080/api'

export interface ChatEvent {
  type: 'token' | 'form_fill' | 'done' | 'error'
  content?: string
  data?: FormFillData
}

export interface FormFillData {
  title: string
  demandScope: string
  businessLine: string
  requesterBU: string
  busInterested: string[]
  demandContext: string
  currentSituation: string
  problemsToSolve: string
  whoIsImpacted: string
  measureBenefits: string
}

export async function sendMessage(
  sessionId: string | null,
  message: string,
  onEvent: (event: ChatEvent) => void,
  signal?: AbortSignal,
): Promise<string> {
  const response = await fetch(`${API_BASE}/chat/message`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ sessionId, message }),
    signal,
  })

  if (!response.ok) {
    throw new Error(`HTTP ${response.status}`)
  }

  const newSessionId = response.headers.get('X-Session-ID') ?? sessionId ?? ''

  const reader = response.body!.getReader()
  const decoder = new TextDecoder()
  let buffer = ''

  try {
    while (true) {
      const { done, value } = await reader.read()
      if (done) break

      buffer += decoder.decode(value, { stream: true })
      const lines = buffer.split('\n')
      buffer = lines.pop() ?? ''

      for (const line of lines) {
        if (!line.startsWith('data: ')) continue
        try {
          const event = JSON.parse(line.slice(6)) as ChatEvent
          onEvent(event)
        } catch {
          // skip malformed lines
        }
      }
    }
  } finally {
    reader.releaseLock()
  }

  return newSessionId
}
