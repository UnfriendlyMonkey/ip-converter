import { ref } from 'vue'

export type Lang = 'en' | 'ru'
export type ThemeName = 'bone' | 'coal' | 'terminal' | 'pulp' | 'riot'

const THEME_ORDER: ThemeName[] = ['bone', 'coal', 'terminal', 'pulp', 'riot']
const savedTheme = (localStorage.getItem('theme') as ThemeName) || 'bone'
const theme = ref<ThemeName>(THEME_ORDER.includes(savedTheme) ? savedTheme : 'bone')
const darkMode = ref(theme.value !== 'bone')
const lang = ref<Lang>((localStorage.getItem('lang') as Lang) || 'en')

export function useSettings() {
  function syncThemeState(next: ThemeName) {
    theme.value = next
    darkMode.value = next !== 'bone'
    localStorage.setItem('theme', next)
  }

  function cycleTheme() {
    const i = THEME_ORDER.indexOf(theme.value)
    const next = THEME_ORDER[(i + 1) % THEME_ORDER.length]
    syncThemeState(next)
  }

  function setTheme(next: ThemeName) {
    syncThemeState(next)
  }

  function toggleTheme() {
    cycleTheme()
  }

  function toggleLang() {
    lang.value = lang.value === 'en' ? 'ru' : 'en'
    localStorage.setItem('lang', lang.value)
  }

  return { darkMode, lang, theme, cycleTheme, setTheme, toggleTheme, toggleLang, themes: THEME_ORDER }
}
