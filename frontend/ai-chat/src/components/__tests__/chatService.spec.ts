import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { sendMessage } from '@/services/chatService'
import type { ChatEvent, FormFillData } from '@/services/chatService'

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

/**
 * Builds a ReadableStream that emits the given raw text as a single chunk.
 * This simulates the SSE response body.
 */
function makeStream(text: string): ReadableStream<Uint8Array> {
  const encoder = new TextEncoder()
  return new ReadableStream<Uint8Array>({
    start(controller) {
      controller.enqueue(encoder.encode(text))
      controller.close()
    },
  })
}

/**
 * Formats a single SSE event line.
 */
function sseEvent(event: ChatEvent): string {
  return `data: ${JSON.stringify(event)}\n`
}

/**
 * Creates a mock fetch response for a successful SSE stream.
 */
function makeResponse(body: string, headers: Record<string, string> = {}): Response {
  const defaultHeaders: Record<string, string> = {
    'Content-Type': 'text/event-stream',
    'X-Session-ID': 'mock-session-001',
    ...headers,
  }

  return {
    ok: true,
    status: 200,
    headers: {
      get: (name: string) => defaultHeaders[name] ?? defaultHeaders[name.toLowerCase()] ?? null,
    },
    body: makeStream(body),
  } as unknown as Response
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------
describe('sendMessage', () => {
  const mockOnEvent = vi.fn<(event: ChatEvent) => void>()

  beforeEach(() => {
    mockOnEvent.mockClear()
    vi.stubGlobal('fetch', vi.fn())
  })

  afterEach(() => {
    vi.unstubAllGlobals()
  })

  // ---- SSE parsing ---------------------------------------------------------

  describe('SSE parsing', () => {
    it('calls onEvent for each valid data: line', async () => {
      const body =
        sseEvent({ type: 'token', content: 'Hello' }) +
        sseEvent({ type: 'token', content: ' world' }) +
        sseEvent({ type: 'done' })

      vi.mocked(fetch).mockResolvedValueOnce(makeResponse(body))

      await sendMessage(null, 'Hi', 'en', mockOnEvent)

      expect(mockOnEvent).toHaveBeenCalledTimes(3)
      expect(mockOnEvent).toHaveBeenNthCalledWith(1, { type: 'token', content: 'Hello' })
      expect(mockOnEvent).toHaveBeenNthCalledWith(2, { type: 'token', content: ' world' })
      expect(mockOnEvent).toHaveBeenNthCalledWith(3, { type: 'done' })
    })

    it('parses form_fill events and passes the full data object', async () => {
      const formFillData: FormFillData = {
        title: 'Add feature X',
        businessLine: '18518',
        requesterBU: 'ADEO-8095',
        busInterested: ['20059'],
        timeSensitive: 'No',
        whyDemand: 'Because',
        whoIsImpacted: 'Users',
        benefitCategory: 'Innovation',
        benefitHypothesis: 'It will work',
        measureBenefits: 'Track monthly',
      }

      const body = sseEvent({ type: 'form_fill', data: formFillData }) + sseEvent({ type: 'done' })

      vi.mocked(fetch).mockResolvedValueOnce(makeResponse(body))

      await sendMessage(null, 'Fill form', 'en', mockOnEvent)

      expect(mockOnEvent).toHaveBeenCalledWith({ type: 'form_fill', data: formFillData })
    })

    it('handles multi-chunk streams by buffering correctly', async () => {
      // Simulate chunk split in the middle of a SSE line
      const chunk1 = 'data: {"type":"tok'
      const chunk2 = 'en","content":"Hi"}\n'
      const chunk3 = sseEvent({ type: 'done' })

      const encoder = new TextEncoder()
      const response: Response = {
        ok: true,
        status: 200,
        headers: {
          get: (name: string) => (name === 'X-Session-ID' ? 'session-x' : null),
        },
        body: new ReadableStream<Uint8Array>({
          start(controller) {
            controller.enqueue(encoder.encode(chunk1))
            controller.enqueue(encoder.encode(chunk2))
            controller.enqueue(encoder.encode(chunk3))
            controller.close()
          },
        }),
      } as unknown as Response

      vi.mocked(fetch).mockResolvedValueOnce(response)

      await sendMessage(null, 'Hi', 'en', mockOnEvent)

      expect(mockOnEvent).toHaveBeenCalledWith({ type: 'token', content: 'Hi' })
      expect(mockOnEvent).toHaveBeenCalledWith({ type: 'done' })
    })

    it('handles lines separated by double newlines (blank SSE comment lines)', async () => {
      // SSE spec allows blank lines between events; our parser splits on \n
      const body =
        sseEvent({ type: 'token', content: 'A' }) +
        '\n' + // blank line
        sseEvent({ type: 'done' })

      vi.mocked(fetch).mockResolvedValueOnce(makeResponse(body))

      await sendMessage(null, 'Hi', 'en', mockOnEvent)

      expect(mockOnEvent).toHaveBeenCalledWith({ type: 'token', content: 'A' })
      expect(mockOnEvent).toHaveBeenCalledWith({ type: 'done' })
    })

    it('ignores lines that do not start with "data: "', async () => {
      const body =
        ': comment line\n' +
        'event: message\n' +
        sseEvent({ type: 'token', content: 'Valid' }) +
        sseEvent({ type: 'done' })

      vi.mocked(fetch).mockResolvedValueOnce(makeResponse(body))

      await sendMessage(null, 'Hi', 'en', mockOnEvent)

      // Only token and done events should fire, not the comment/event lines
      expect(mockOnEvent).toHaveBeenCalledTimes(2)
      expect(mockOnEvent).toHaveBeenCalledWith({ type: 'token', content: 'Valid' })
    })
  })

  // ---- Malformed JSON -------------------------------------------------------

  describe('malformed JSON handling', () => {
    it('skips malformed JSON lines without throwing', async () => {
      const body =
        'data: {this is not json}\n' +
        sseEvent({ type: 'token', content: 'OK' }) +
        sseEvent({ type: 'done' })

      vi.mocked(fetch).mockResolvedValueOnce(makeResponse(body))

      // Should resolve without throwing
      await expect(sendMessage(null, 'Hi', 'en', mockOnEvent)).resolves.not.toThrow()

      // The valid events should still fire
      expect(mockOnEvent).toHaveBeenCalledWith({ type: 'token', content: 'OK' })
      expect(mockOnEvent).toHaveBeenCalledWith({ type: 'done' })
    })

    it('skips truncated JSON without throwing', async () => {
      const body =
        'data: {"type":"token","content":"inc\n' + // truncated mid-string
        sseEvent({ type: 'done' })

      vi.mocked(fetch).mockResolvedValueOnce(makeResponse(body))

      await expect(sendMessage(null, 'Hi', 'en', mockOnEvent)).resolves.not.toThrow()
      expect(mockOnEvent).toHaveBeenCalledWith({ type: 'done' })
    })

    it('skips empty data: lines without throwing', async () => {
      const body = 'data: \n' + sseEvent({ type: 'done' })

      vi.mocked(fetch).mockResolvedValueOnce(makeResponse(body))

      await expect(sendMessage(null, 'Hi', 'en', mockOnEvent)).resolves.not.toThrow()
      expect(mockOnEvent).toHaveBeenCalledTimes(1)
      expect(mockOnEvent).toHaveBeenCalledWith({ type: 'done' })
    })

    it('handles a stream with only malformed lines and returns the session id', async () => {
      const body = 'data: not-json\ndata: also-not-json\n'

      vi.mocked(fetch).mockResolvedValueOnce(makeResponse(body))

      const result = await sendMessage(null, 'Hi', 'en', mockOnEvent)
      expect(result).toBe('mock-session-001')
      expect(mockOnEvent).not.toHaveBeenCalled()
    })
  })

  // ---- Session ID handling -------------------------------------------------

  describe('session ID handling', () => {
    it('returns X-Session-ID header value when present', async () => {
      vi.mocked(fetch).mockResolvedValueOnce(
        makeResponse(sseEvent({ type: 'done' }), { 'X-Session-ID': 'server-generated-id' }),
      )

      const result = await sendMessage(null, 'Hi', 'en', mockOnEvent)
      expect(result).toBe('server-generated-id')
    })

    it('falls back to passed sessionId when X-Session-ID header is absent', async () => {
      const response: Response = {
        ok: true,
        status: 200,
        headers: { get: (_name: string) => null },
        body: makeStream(sseEvent({ type: 'done' })),
      } as unknown as Response

      vi.mocked(fetch).mockResolvedValueOnce(response)

      const result = await sendMessage('existing-session', 'Hi', 'en', mockOnEvent)
      expect(result).toBe('existing-session')
    })

    it('returns empty string when both header and sessionId are absent', async () => {
      const response: Response = {
        ok: true,
        status: 200,
        headers: { get: (_name: string) => null },
        body: makeStream(sseEvent({ type: 'done' })),
      } as unknown as Response

      vi.mocked(fetch).mockResolvedValueOnce(response)

      const result = await sendMessage(null, 'Hi', 'en', mockOnEvent)
      expect(result).toBe('')
    })

    it('sends correct request body including sessionId, message, and language', async () => {
      vi.mocked(fetch).mockResolvedValueOnce(makeResponse(sseEvent({ type: 'done' })))

      await sendMessage('my-session', 'Hello backend', 'fr', mockOnEvent)

      expect(fetch).toHaveBeenCalledOnce()
      const [url, options] = vi.mocked(fetch).mock.calls[0]!
      expect(url).toContain('/api/chat/message')
      const body = JSON.parse(options!.body as string) as Record<string, unknown>
      expect(body.sessionId).toBe('my-session')
      expect(body.message).toBe('Hello backend')
      expect(body.language).toBe('fr')
    })

    it('sends null sessionId when none is provided', async () => {
      vi.mocked(fetch).mockResolvedValueOnce(makeResponse(sseEvent({ type: 'done' })))

      await sendMessage(null, 'Hello', 'en', mockOnEvent)

      const [, options] = vi.mocked(fetch).mock.calls[0]!
      const body = JSON.parse(options!.body as string) as Record<string, unknown>
      expect(body.sessionId).toBeNull()
    })
  })

  // ---- HTTP error handling -------------------------------------------------

  describe('HTTP error handling', () => {
    it('throws an error when response.ok is false (e.g. 500)', async () => {
      const errorResponse: Response = {
        ok: false,
        status: 500,
        headers: { get: () => null },
        body: null,
      } as unknown as Response

      vi.mocked(fetch).mockResolvedValueOnce(errorResponse)

      await expect(sendMessage(null, 'Hi', 'en', mockOnEvent)).rejects.toThrow('HTTP 500')
    })

    it('throws an error when response is 404', async () => {
      const errorResponse: Response = {
        ok: false,
        status: 404,
        headers: { get: () => null },
        body: null,
      } as unknown as Response

      vi.mocked(fetch).mockResolvedValueOnce(errorResponse)

      await expect(sendMessage(null, 'Hi', 'en', mockOnEvent)).rejects.toThrow('HTTP 404')
    })

    it('rethrows network-level errors (e.g. fetch rejection)', async () => {
      vi.mocked(fetch).mockRejectedValueOnce(new Error('Network failure'))

      await expect(sendMessage(null, 'Hi', 'en', mockOnEvent)).rejects.toThrow('Network failure')
    })
  })

  // ---- AbortSignal ---------------------------------------------------------

  describe('AbortSignal cancellation', () => {
    it('throws an AbortError when signal is already aborted before fetch', async () => {
      const controller = new AbortController()
      controller.abort()

      // fetch should throw DOMException with name=AbortError when signal is already aborted
      const abortError = new DOMException('The operation was aborted.', 'AbortError')
      vi.mocked(fetch).mockRejectedValueOnce(abortError)

      const err = await sendMessage(null, 'Hi', 'en', mockOnEvent, controller.signal).catch(
        (e) => e,
      )

      expect(err).toBeInstanceOf(DOMException)
      expect((err as DOMException).name).toBe('AbortError')
    })

    it('throws an AbortError when signal is aborted mid-stream', async () => {
      const controller = new AbortController()

      // Simulate stream that hangs then aborts
      const encoder = new TextEncoder()
      const response: Response = {
        ok: true,
        status: 200,
        headers: { get: (name: string) => (name === 'X-Session-ID' ? 'session-x' : null) },
        body: new ReadableStream<Uint8Array>({
          start(streamController) {
            streamController.enqueue(encoder.encode(sseEvent({ type: 'token', content: 'part' })))
            // Abort the controller to simulate user aborting mid-stream
            controller.abort()
            // Signal the abort by erroring the stream
            streamController.error(new DOMException('AbortError', 'AbortError'))
          },
        }),
      } as unknown as Response

      vi.mocked(fetch).mockResolvedValueOnce(response)

      const err = await sendMessage(null, 'Hi', 'en', mockOnEvent, controller.signal).catch(
        (e) => e,
      )

      // The function should propagate the error; it may be DOMException or Error
      expect(err).toBeDefined()
      expect(err instanceof Error || err instanceof DOMException).toBe(true)
    })

    it('passes AbortSignal to fetch options', async () => {
      const controller = new AbortController()
      vi.mocked(fetch).mockResolvedValueOnce(makeResponse(sseEvent({ type: 'done' })))

      await sendMessage(null, 'Hi', 'en', mockOnEvent, controller.signal)

      const [, options] = vi.mocked(fetch).mock.calls[0]!
      expect((options as RequestInit).signal).toBe(controller.signal)
    })
  })

  // ---- Edge cases ----------------------------------------------------------

  describe('edge cases', () => {
    it('handles a completely empty stream body', async () => {
      vi.mocked(fetch).mockResolvedValueOnce(makeResponse(''))

      await expect(sendMessage(null, 'Hi', 'en', mockOnEvent)).resolves.not.toThrow()
      expect(mockOnEvent).not.toHaveBeenCalled()
    })

    it('handles a stream with only newlines', async () => {
      vi.mocked(fetch).mockResolvedValueOnce(makeResponse('\n\n\n\n'))

      await expect(sendMessage(null, 'Hi', 'en', mockOnEvent)).resolves.not.toThrow()
      expect(mockOnEvent).not.toHaveBeenCalled()
    })

    it('handles multiple consecutive tokens without intermediate done events', async () => {
      const tokens = Array.from({ length: 10 }, (_, i) =>
        sseEvent({ type: 'token', content: `word${i}` }),
      ).join('')

      vi.mocked(fetch).mockResolvedValueOnce(makeResponse(tokens + sseEvent({ type: 'done' })))

      await sendMessage(null, 'Hi', 'en', mockOnEvent)

      // 10 tokens + 1 done = 11 calls
      expect(mockOnEvent).toHaveBeenCalledTimes(11)
    })

    it('handles a form_fill event with busInterested as an empty array', async () => {
      const formData: FormFillData = {
        title: 'Test',
        businessLine: '18518',
        requesterBU: 'ADEO-8095',
        busInterested: [],
        timeSensitive: 'No',
        whyDemand: 'Test',
        whoIsImpacted: 'Test',
        benefitCategory: 'Innovation',
        benefitHypothesis: 'Test',
        measureBenefits: 'Test',
      }

      const body = sseEvent({ type: 'form_fill', data: formData }) + sseEvent({ type: 'done' })
      vi.mocked(fetch).mockResolvedValueOnce(makeResponse(body))

      await sendMessage(null, 'Hi', 'en', mockOnEvent)

      expect(mockOnEvent).toHaveBeenCalledWith({ type: 'form_fill', data: formData })
    })

    it('uses POST method with Content-Type application/json', async () => {
      vi.mocked(fetch).mockResolvedValueOnce(makeResponse(sseEvent({ type: 'done' })))

      await sendMessage(null, 'Test', 'en', mockOnEvent)

      const [, options] = vi.mocked(fetch).mock.calls[0]!
      expect((options as RequestInit).method).toBe('POST')
      expect((options as RequestInit).headers).toMatchObject({
        'Content-Type': 'application/json',
      })
    })
  })
})
