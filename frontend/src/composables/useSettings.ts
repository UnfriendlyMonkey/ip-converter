import { ref } from 'vue'

export type Lang = 'en' | 'ru'

const darkMode = ref(localStorage.getItem('theme') !== 'light')
const lang = ref<Lang>((localStorage.getItem('lang') as Lang) || 'en')

export function useSettings() {
  function toggleTheme() {
    darkMode.value = !darkMode.value
    localStorage.setItem('theme', darkMode.value ? 'dark' : 'light')
  }

  function toggleLang() {
    lang.value = lang.value === 'en' ? 'ru' : 'en'
    localStorage.setItem('lang', lang.value)
  }

  return { darkMode, lang, toggleTheme, toggleLang }
}
