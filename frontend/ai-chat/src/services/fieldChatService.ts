import { readSSEStream } from '@/utils/sse'
import { API_BASE } from '@/services/config'

export interface FieldChatEvent {
  type: 'token' | 'field_fill' | 'done' | 'error'
  content?: string
  data?: { fieldName: string; value: string | string[] }
}

export async function sendFieldMessage(
  sessionId: string | null,
  fieldName: string,
  message: string,
  language: string,
  onEvent: (event: FieldChatEvent) => void,
  signal?: AbortSignal,
): Promise<string | null> {
  const response = await fetch(`${API_BASE}/chat/field-message`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ sessionId, fieldName, message, language }),
    signal,
  })

  if (!response.ok) {
    throw new Error(`HTTP ${response.status}`)
  }

  const newSessionId = response.headers.get('X-Session-ID') ?? sessionId
  await readSSEStream(response, onEvent)
  return newSessionId
}
