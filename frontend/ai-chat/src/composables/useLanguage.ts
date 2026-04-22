import { ref } from 'vue'

export interface Language {
  code: string
  label: string
}

export const LANGUAGES: Language[] = [
  { code: 'pt', label: 'PT' },
  { code: 'en', label: 'EN' },
  { code: 'es', label: 'ES' },
  { code: 'fr', label: 'FR' },
]

export function useLanguage() {
  const stored = localStorage.getItem('chat_language')
  const validCodes = LANGUAGES.map((l) => l.code)
  const language = ref<string>(stored && validCodes.includes(stored) ? stored : 'en')

  function setLanguage(code: string) {
    language.value = code
    localStorage.setItem('chat_language', code)
  }

  return { language, setLanguage, languages: LANGUAGES }
}
