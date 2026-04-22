import { readSSEStream } from '@/utils/sse'

const API_BASE = import.meta.env.VITE_API_BASE ?? 'http://localhost:8080/api'

export interface ChatEvent {
  type: 'token' | 'form_fill' | 'form_confirm' | 'done' | 'error'
  content?: string
  data?: FormFillData
}

export interface FormFillData {
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
  lowConfidenceFields?: string[]
}

export async function sendMessage(
  sessionId: string | null,
  message: string,
  language: string,
  onEvent: (event: ChatEvent) => void,
  signal?: AbortSignal,
): Promise<string | null> {
  const response = await fetch(`${API_BASE}/chat/message`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ sessionId, message, language }),
    signal,
  })

  if (!response.ok) {
    throw new Error(`HTTP ${response.status}`)
  }

  const newSessionId = response.headers.get('X-Session-ID') ?? sessionId
  await readSSEStream(response, onEvent)
  return newSessionId
}

export async function analyzeForm(
  formData: FormFillData,
  language: string,
  onEvent: (event: ChatEvent) => void,
  signal?: AbortSignal,
): Promise<void> {
  const response = await fetch(`${API_BASE}/chat/analyze-form`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ formData, language }),
    signal,
  })

  if (!response.ok) {
    throw new Error(`HTTP ${response.status}`)
  }

  await readSSEStream(response, onEvent)
}

export async function confirmForm(
  sessionId: string,
  onEvent: (event: ChatEvent) => void,
  signal?: AbortSignal,
): Promise<void> {
  const response = await fetch(`${API_BASE}/chat/confirm`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ sessionId }),
    signal,
  })

  if (!response.ok) {
    throw new Error(`HTTP ${response.status}`)
  }

  await readSSEStream(response, onEvent)
}
