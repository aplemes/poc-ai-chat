export async function readSSEStream<T>(
  response: Response,
  onEvent: (event: T) => void,
): Promise<void> {
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

      for (const rawLine of lines) {
        const line = rawLine.replace(/\r$/, '')
        if (!line.startsWith('data: ')) continue
        try {
          const event = JSON.parse(line.slice(6)) as T
          onEvent(event)
        } catch {
          // skip malformed lines
        }
      }
    }
  } finally {
    reader.releaseLock()
  }
}
