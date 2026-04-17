# QA Bug Report

Generated: 2026-04-17  
Scope: `/home/apaneto/chat-ai/frontend/ai-chat/` — FloatingChat.vue, RequestForm.vue, stores/chat.ts, services/chatService.ts

---

## BUG-001

**Severity:** Major  
**File + line:** `frontend/ai-chat/src/components/FloatingChat.vue` — line 216  
**Description:** The footer hint text "Enter para enviar · Shift+Enter nova linha" is hard-coded in Portuguese regardless of the selected language. When the user switches to EN, ES, or FR the rest of the UI correctly localises, but this hint stays in Portuguese.  
**Steps to reproduce:**
1. Open the app in a browser.
2. Click the FAB to open the chat.
3. Select EN (or ES or FR) in the language switcher.
4. Observe the hint below the textarea.

**Expected:** Hint text matches the active language (e.g. "Press Enter to send · Shift+Enter for a new line" for EN).  
**Actual:** Hint always reads "Enter para enviar · Shift+Enter nova linha" (Portuguese).  
**Suggested fix:** Add a `footerHint` property to each entry in the `languages` array and bind it: `:text="currentLang().footerHint"`.

---

## BUG-002

**Severity:** Major  
**File + line:** `frontend/ai-chat/src/components/FloatingChat.vue` — lines 112–114  
**Description:** When a network-level error occurs (e.g. the backend is unreachable), the `catch` block maps any non-`AbortError` to `errorConnection`. However, when the server responds with a non-2xx HTTP status code (e.g. 500), `chatService.ts` throws `new Error('HTTP 500')` — this is also not an `AbortError`, so the UI correctly shows `errorConnection`. But if the error event type comes from the SSE stream body itself (the backend emits `{type:"error"}`), the handler sets `assistantMsg.text = currentLang().errorProcessing`. The two code paths produce different user-facing strings for what appears to be the same situation (server-side error), and only one path clears `loading`. The `event.type === 'error'` branch (line 106–108) does NOT set `loading = false` by itself — it relies on the `finally` block, which IS present. This is acceptable, but the semantic inconsistency between the two strings is confusing and undocumented.  
**Steps to reproduce:**
1. Configure the backend to return a streaming `{type:"error"}` event.
2. Observe the assistant message.
3. Now kill the backend and send another message.
4. Observe the different error message despite both being server errors.

**Expected:** A single, consistent error string for server-side problems.  
**Actual:** Two different strings (`errorProcessing` vs `errorConnection`) depending on which failure path was taken.  
**Suggested fix:** Unify the two error messages or clearly differentiate them in the UX (e.g. one for "processing error", one for "connection error") and document the distinction.

---

## BUG-003

**Severity:** Major  
**File + line:** `frontend/ai-chat/src/components/FloatingChat.vue` — lines 94–109  
**Description:** Race condition in the abort/send flow. When the user sends a new message while a previous stream is in-flight, the code creates a new `AbortController` and assigns it to the module-level `abortController` variable (line 94). The previous controller is never explicitly aborted — the old `send()` call's own `abortController` local reference is already unreachable because the variable was overwritten. The previous request therefore continues streaming in the background and its `onEvent` callback still mutates the old `assistantMsg` object, which may still be rendered in the `messages` array.

**Steps to reproduce:**
1. Open chat and send a message.
2. While the first stream is still producing tokens, type and send a second message.
3. Observe that the first assistant bubble may continue receiving tokens even after the second message is sent.

**Expected:** Sending a new message should cancel the previous in-flight request before starting a new one.  
**Actual:** Both streams can be active simultaneously; the first stream's tokens continue to update its assistant message bubble.  
**Suggested fix:** Before assigning the new controller, explicitly abort the existing one:
```ts
abortController?.abort()
abortController = new AbortController()
```

---

## BUG-004

**Severity:** Major  
**File + line:** `frontend/ai-chat/src/services/chatService.ts` — line 40  
**Description:** The fallback logic for `newSessionId` uses the nullish coalescing chain `response.headers.get('X-Session-ID') ?? sessionId ?? ''`. When the header is absent AND `sessionId` is `null`, the result is `''` (empty string). The caller (`FloatingChat.vue` line 110–111) then stores this empty string in `localStorage` and uses it as the session ID for subsequent requests. The backend will treat an empty-string session ID differently from a missing/null one, potentially creating a broken session.  
**Steps to reproduce:**
1. Clear localStorage.
2. Configure the mock server to NOT return `X-Session-ID`.
3. Send a message.
4. Observe `localStorage.getItem('chat_session_id')` — it will be `''`.
5. Send another message — the body will include `sessionId: ""` instead of `sessionId: null`.

**Expected:** When no session ID is available from either source, the value should remain `null` and not be stored.  
**Actual:** An empty string `''` is stored in localStorage and sent as the session ID on the next request.  
**Suggested fix:**
```ts
const newSessionId = response.headers.get('X-Session-ID') ?? sessionId
// In FloatingChat.vue, only persist when non-null/non-empty:
if (newId) {
  sessionId.value = newId
  localStorage.setItem('chat_session_id', newId)
}
```

---

## BUG-005

**Severity:** Minor  
**File + line:** `frontend/ai-chat/src/components/RequestForm.vue` — lines 117–135  
**Description:** The `watch` on `chatStore.pendingFormFill` does not use `{ immediate: false }` explicitly, but more importantly it does not guard against re-entrant triggers. If `setFormFill` is called while the watcher is still executing (e.g. a reactive side-effect triggers another store write), the watcher's `chatStore.clearFormFill()` call at line 133 will fire synchronously, nullifying `pendingFormFill` before any re-queued watcher invocation can run. In practice this is benign given the current flow, but creates a hidden coupling: `RequestForm` is the only consumer allowed to call `clearFormFill()`, and any second consumer (e.g. a future analytics hook) would never see the value.  
**Steps to reproduce:** Not directly reproducible with a single consumer; becomes a bug when a second watcher on `pendingFormFill` is added elsewhere.  
**Expected:** The store is cleared only after all consumers have had a chance to read the value, or the clearing is the explicit responsibility of the producer (`FloatingChat`) rather than the consumer (`RequestForm`).  
**Actual:** `RequestForm` silently clears the store immediately after reading, creating an implicit contract that only one consumer may ever exist.  
**Suggested fix:** Move `clearFormFill()` into `FloatingChat.vue` after `setFormFill()` is confirmed as delivered, or use a short `nextTick` delay before clearing, or implement an explicit acknowledgement pattern.

---

## BUG-006

**Severity:** Minor  
**File + line:** `frontend/ai-chat/src/components/RequestForm.vue` — lines 156–167 (validate function)  
**Description:** The `validate()` function does not validate `demandContext`, `currentSituation`, or `problemsToSolve` fields — these fields are listed in `CLAUDE.md` as required (`✓`) but do not appear at all in the form component or in the validation logic. The form silently omits them, so the submitted data will always be missing these fields.  
**Steps to reproduce:**
1. Fill all visible required fields and submit.
2. Check the `console.log(form.value)` output — `demandContext`, `currentSituation`, and `problemsToSolve` are absent.

**Expected:** Either the form includes all required fields from the domain model, or the architecture doc is updated to reflect that these three fields are deferred.  
**Actual:** Three required fields from the CLAUDE.md specification are completely missing from the form.  
**Suggested fix:** Add the three missing textarea fields to the form (likely in Section 2, "Context & Impact") and add them to the `FormFillData` interface in `chatService.ts` and to the backend tool schema.

---

## BUG-007

**Severity:** Minor  
**File + line:** `frontend/ai-chat/src/components/RequestForm.vue` — line 125  
**Description:** The `busInterested` spread `[...data.busInterested]` assumes `data.busInterested` is always an array. The `FormFillData` interface in `chatService.ts` types it as `string[]`, but the backend could theoretically emit `null` (if the model omits the field) or a JSON string `"20059"` (a single item without wrapping in an array). Both cases would cause a runtime error. The `data.busInterested?.length` guard (line 125) correctly handles `undefined` and empty-array cases but does not handle `null` or a raw string.  
**Steps to reproduce:**
1. Send a mock SSE event where `busInterested` is `null` instead of `[]`.
2. The `?.length` on `null` evaluates to `undefined` (falsy), so the guard correctly prevents the spread — this specific case is actually safe.
3. However if `busInterested` is a raw string `"20059"`, `"20059".length` is 5 (truthy) and `[..."20059"]` spreads individual characters `['2','0','0','5','9']` into `form.value.busInterested`, corrupting the data silently.

**Expected:** If the backend sends a single ID as a plain string, it should be normalised to a single-element array.  
**Actual:** Spreading a string value character-by-character into the chips array, rendering individual digits as chip IDs that won't match any `concerned` entry.  
**Suggested fix:** Normalise before spreading:
```ts
const ids = Array.isArray(data.busInterested) ? data.busInterested : [data.busInterested].filter(Boolean)
if (ids.length) { form.value.busInterested = ids; filled.add('busInterested') }
```

---

## BUG-008

**Severity:** Minor  
**File + line:** `frontend/ai-chat/src/services/chatService.ts` — lines 52–53  
**Description:** The SSE line buffer split uses `buffer.split('\n')` and `buffer = lines.pop() ?? ''`. This correctly handles `\n`-delimited streams but will silently misparse streams that use `\r\n` line endings (Windows-style), which some HTTP proxies and certain Nginx SSE configurations emit. A line like `data: {"type":"token"}\r` would fail `JSON.parse` because of the trailing `\r`, the error is swallowed, and the event is lost.  
**Steps to reproduce:**
1. Configure a proxy (e.g. Nginx) that rewrites line endings to CRLF.
2. Send a message.
3. No tokens appear in the chat, no error is visible.

**Expected:** SSE events are parsed regardless of `\n` vs `\r\n` line endings.  
**Actual:** CRLF-terminated lines cause silent JSON parse failures.  
**Suggested fix:**
```ts
const lines = buffer.split('\n')
// Strip carriage returns before parsing
for (const rawLine of lines) {
  const line = rawLine.replace(/\r$/, '')
  ...
}
```

---

## BUG-009

**Severity:** Minor  
**File + line:** `frontend/ai-chat/src/components/RequestForm.vue` — line 169–177 (handleSubmit)  
**Description:** The `handleSubmit` function calls `console.log(form.value)` instead of making an actual API submission. There is no TODO comment or feature-flag guard. This means the "Submit Request" button is a no-op in production.  
**Steps to reproduce:**
1. Fill all required fields.
2. Click "Submit Request".
3. Open the browser console — the form data is logged but no network request is made.

**Expected:** A real HTTP request to submit the demand, or at minimum a placeholder that throws a clear "not implemented" error.  
**Actual:** Silent `console.log` with no user feedback confirming or denying submission.  
**Suggested fix:** Either implement the submission endpoint call or replace the log with an explicit `throw new Error('Submission endpoint not yet implemented')` and surface a user-visible error state.

---

## BUG-010

**Severity:** Minor  
**File + line:** `frontend/ai-chat/src/components/FloatingChat.vue` — line 36  
**Description:** The initial `language` ref reads from `localStorage.getItem('chat_language') ?? 'en'`. If the stored value is a string that doesn't match any language code (e.g. `'de'` from a previous app version, or a corrupted value), `currentLang()` falls back to `languages[1]` (EN) via the `?? languages[1]!` guard. This is safe but silent — the localStorage value is not corrected, so every page load will go through the fallback path without self-healing.  
**Steps to reproduce:**
1. Set `localStorage.setItem('chat_language', 'de')` in the console.
2. Reload the page.
3. The UI shows EN (correct fallback) but the language button for EN is not highlighted as active because `language.value` is `'de'`, not `'en'`.

**Expected:** The active language button reflects the actual active language.  
**Actual:** No language button appears active; the `currentLang()` fallback works but the UI highlight is broken because `language.value` still holds the invalid `'de'` string.  
**Suggested fix:** Validate the stored value on init and reset if invalid:
```ts
const storedLang = localStorage.getItem('chat_language')
const validCodes = languages.map(l => l.code)
const language = ref<string>(
  storedLang && validCodes.includes(storedLang) ? storedLang : 'en'
)
```
