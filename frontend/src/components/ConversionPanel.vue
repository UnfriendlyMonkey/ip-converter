<template lang="pug">
div.cp-root(:class="{ 'cp-root--light': !darkMode }")

  //- ── Top bar ──────────────────────────────────────────────
  header.cp-header
    div.cp-brand
      span.cp-brand-icon ⬡
      span.cp-brand-name IP CONVERTER

    div.cp-header-right
      button.cp-lang-btn(@click="toggleLang" :title="t('switchLang')") {{ lang === 'en' ? 'RU' : 'EN' }}
      button.cp-theme-btn(@click="toggleTheme" :title="darkMode ? t('toLight') : t('toDark')")
        | {{ darkMode ? '☀' : '☾' }}
      div.cp-dir-controls
        div.cp-dir-toggle
          button(
            :class="['cp-dir-opt', direction === 'to_numeric' && 'cp-dir-opt--active']"
            @click="setDirection('to_numeric')"
          ) {{ t('strToNum') }}
          button.cp-dir-swap(@click="toggleDirection" :title="t('swapDir')") ⇄
          button(
            :class="['cp-dir-opt', direction === 'to_string' && 'cp-dir-opt--active']"
            @click="setDirection('to_string')"
          ) {{ t('numToStr') }}

  //- ── Panels ───────────────────────────────────────────────
  div.cp-panels

    //- INPUT PANEL
    section.cp-panel.cp-panel--input
      div.cp-panel-header
        span.cp-panel-title {{ t('input') }}
        div.cp-mode-toggle
          button(
            :class="['cp-tab', !multiMode && 'cp-tab--active']"
            @click="setMode(false)"
          ) {{ t('single') }}
          button(
            :class="['cp-tab', multiMode && 'cp-tab--active']"
            @click="setMode(true)"
          ) {{ t('multi') }}

      div.cp-panel-body
        input.cp-input.cp-mono(
          v-if="!multiMode"
          v-model="singleInput"
          :placeholder="singlePlaceholder"
          type="text"
          spellcheck="false"
          autocomplete="off"
          autocorrect="off"
          @input="onSingleInput"
        )

        textarea.cp-textarea.cp-mono(
          v-else
          v-model="multiInput"
          :placeholder="multiPlaceholder"
          spellcheck="false"
          @keydown.ctrl.enter="convertMulti"
        )

        div.cp-panel-footer
          span.cp-hint {{ !multiMode ? t('hintSingle') : t('hintMulti') }}
          button.cp-action-btn(
            v-if="multiMode"
            :class="{ 'cp-action-btn--loading': loading }"
            :disabled="loading"
            @click="convertMulti"
          )
            span.cp-blink(v-if="loading") ···
            span(v-else) ⚡ {{ t('convert') }}

    //- Divider arrow
    div.cp-arrow
      span ›

    //- OUTPUT PANEL
    section.cp-panel.cp-panel--output
      div.cp-panel-header
        span.cp-panel-title {{ t('output') }}
        div.cp-output-controls
          template(v-if="multiMode")
            button(
              :class="['cp-tab', outputFormat === 'list' && 'cp-tab--active']"
              @click="outputFormat = 'list'"
            ) {{ t('list') }}
            button(
              :class="['cp-tab', outputFormat === 'dict' && 'cp-tab--active']"
              @click="outputFormat = 'dict'"
            ) {{ t('dict') }}
          button.cp-copy-btn(
            :class="{ 'cp-copy-btn--ok': justCopied }"
            :disabled="!outputText"
            @click="copyOutput"
          ) {{ justCopied ? t('copied') : t('copy') }}

      div.cp-panel-body
        div.cp-output-state(v-if="loading")
          span.cp-blink.cp-hint {{ t('processing') }}

        div.cp-output-state.cp-error-msg(v-else-if="apiError") {{ apiError }}

        pre.cp-output.cp-mono(
          v-else-if="outputText"
          :class="{ 'cp-output--has-errors': hasErrors }"
        ) {{ outputText }}

        div.cp-output-state(v-else)
          span.cp-hint {{ t('emptyOutput') }}

  //- ── Footer ───────────────────────────────────────────────
  footer.cp-footer
    span.cp-hint {{ t('footer') }}

  //- Copy toast
  transition(name="toast")
    div.cp-toast(v-if="justCopied") {{ t('toastCopied') }}
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { convertValues, type Direction, type ConversionResult } from '../api/converter'

// ── i18n ────────────────────────────────────────────────────
type Lang = 'en' | 'ru'

const translations = {
  en: {
    strToNum:   'STRING → NUM',
    numToStr:   'NUM → STRING',
    swapDir:    'Swap direction',
    toLight:    'Switch to light mode',
    toDark:     'Switch to dark mode',
    switchLang: 'Switch to Russian',
    input:      'INPUT',
    output:     'OUTPUT',
    single:     'SINGLE',
    multi:      'MULTI',
    list:       'LIST',
    dict:       'DICT',
    copy:       'COPY',
    copied:     '✓ COPIED',
    convert:    'CONVERT',
    processing: 'processing...',
    hintSingle: 'converts live · type a valid IP',
    hintMulti:  'one value per line · Ctrl+Enter to convert',
    emptyOutput:'// output will appear here',
    footer:     'ipv4 ints · ipv6 byte arrays · CIDR masks supported',
    toastCopied:'Copied to clipboard!',
  },
  ru: {
    strToNum:   'СТРОКА → ЧИСЛО',
    numToStr:   'ЧИСЛО → СТРОКА',
    swapDir:    'Сменить направление',
    toLight:    'Светлая тема',
    toDark:     'Тёмная тема',
    switchLang: 'Switch to English',
    input:      'ВВОД',
    output:     'ВЫВОД',
    single:     'ОДНО',
    multi:      'МНОГО',
    list:       'СПИСОК',
    dict:       'СЛОВАРЬ',
    copy:       'КОПИРОВАТЬ',
    copied:     '✓ СКОПИРОВАНО',
    convert:    'КОНВЕРТИРОВАТЬ',
    processing: 'обработка...',
    hintSingle: 'конвертирует автоматически · введите IP',
    hintMulti:  'по одному значению в строке · Ctrl+Enter',
    emptyOutput:'// результат появится здесь',
    footer:     'ipv4 целые числа · ipv6 байт-массивы · маски CIDR',
    toastCopied:'Скопировано!',
  },
} as const

type TKey = keyof typeof translations.en

const lang = ref<Lang>((localStorage.getItem('lang') as Lang) || 'en')

function t(key: TKey): string {
  return translations[lang.value][key]
}

function toggleLang() {
  lang.value = lang.value === 'en' ? 'ru' : 'en'
  localStorage.setItem('lang', lang.value)
}

// ── Theme ──────────────────────────────────────────────────
const darkMode = ref(localStorage.getItem('theme') !== 'light')

function toggleTheme() {
  darkMode.value = !darkMode.value
  localStorage.setItem('theme', darkMode.value ? 'dark' : 'light')
}

// ── State ──────────────────────────────────────────────────
const direction = ref<Direction>('to_numeric')
const multiMode = ref(false)
const singleInput = ref('')
const multiInput = ref('')
const outputFormat = ref<'list' | 'dict'>('list')
const results = ref<ConversionResult[]>([])
const loading = ref(false)
const apiError = ref('')
const justCopied = ref(false)

// ── Placeholders ───────────────────────────────────────────
const singlePlaceholder = computed(() =>
  direction.value === 'to_numeric'
    ? '192.168.1.1  or  192.168.1.0/24  or  2001:db8::1'
    : '3232235777  or  3232235776/24  or  [32,1,13,184,...,1]'
)

const multiPlaceholder = computed(() =>
  direction.value === 'to_numeric'
    ? '192.168.1.1\n192.168.1.0/24\n2001:db8::1\n2001:db8::/32'
    : '3232235777\n3232235776/24\n[32,1,13,184,0,0,0,0,0,0,0,0,0,0,0,1]\n[32,1,13,184,0,0,0,0,0,0,0,0,0,0,0,0]/32'
)

// ── Output formatting ──────────────────────────────────────
const hasErrors = computed(() => results.value.some(r => r.type === 'error'))

const outputText = computed(() => {
  if (!results.value.length) return ''

  if (!multiMode.value || outputFormat.value === 'list') {
    return results.value
      .map(r => r.type === 'error' ? `# ERROR: ${r.error}` : r.output)
      .join('\n')
  }

  const dict: Record<string, string> = {}
  for (const r of results.value) {
    dict[r.input] = r.type === 'error' ? `ERROR: ${r.error}` : r.output
  }
  return JSON.stringify(dict, null, 2)
})

// ── Debounce helper ────────────────────────────────────────
function debounce<T extends (...args: unknown[]) => void>(fn: T, ms: number): T {
  let timer: ReturnType<typeof setTimeout>
  return ((...args: unknown[]) => {
    clearTimeout(timer)
    timer = setTimeout(() => fn(...args), ms)
  }) as T
}

// ── Conversion ─────────────────────────────────────────────
async function runConversion(values: string[]) {
  if (!values.length) {
    results.value = []
    apiError.value = ''
    return
  }
  loading.value = true
  apiError.value = ''
  try {
    const resp = await convertValues({ direction: direction.value, values })
    results.value = resp.results
  } catch (e: unknown) {
    apiError.value = e instanceof Error ? e.message : 'Unknown error'
    results.value = []
  } finally {
    loading.value = false
  }
}

const debouncedSingle = debounce(() => {
  const v = singleInput.value.trim()
  runConversion(v ? [v] : [])
}, 350)

function onSingleInput() {
  debouncedSingle()
}

function convertMulti() {
  const lines = multiInput.value
    .split('\n')
    .map(l => l.trim())
    .filter(l => l.length > 0)
  runConversion(lines)
}

// ── Mode / direction switching ─────────────────────────────
function setMode(multi: boolean) {
  multiMode.value = multi
  results.value = []
  apiError.value = ''
}

function toggleDirection() {
  setDirection(direction.value === 'to_numeric' ? 'to_string' : 'to_numeric')
}

function setDirection(dir: Direction) {
  direction.value = dir
  singleInput.value = ''
  multiInput.value = ''
  results.value = []
  apiError.value = ''
}

// ── Clipboard ──────────────────────────────────────────────
async function copyOutput() {
  if (!outputText.value) return
  try {
    await navigator.clipboard.writeText(outputText.value)
    justCopied.value = true
    setTimeout(() => { justCopied.value = false }, 2000)
  } catch {
    const ta = document.createElement('textarea')
    ta.value = outputText.value
    document.body.appendChild(ta)
    ta.select()
    document.execCommand('copy')
    document.body.removeChild(ta)
    justCopied.value = true
    setTimeout(() => { justCopied.value = false }, 2000)
  }
}
</script>

<style scoped>
/* ── Design tokens ──────────────────────────────── */
.cp-root {
  --bg: #0a0a0f;
  --surface: #0d0d14;
  --surface-hi: #131320;
  --border: rgba(0, 255, 136, 0.12);
  --border-hi: rgba(0, 255, 136, 0.35);
  --green: #00ff88;
  --green-dim: rgba(0, 255, 136, 0.6);
  --cyan: #00e5ff;
  --text: #8aa88a;
  --text-bright: #ccffcc;
  --text-dim: #3a4a3a;
  --error: #ff4466;
  --mono: 'JetBrains Mono', 'Fira Code', 'Cascadia Code', 'Courier New', monospace;

  display: flex;
  flex-direction: column;
  min-height: 100vh;
  background: var(--bg);
  color: var(--text);
  font-family: var(--mono);
  font-size: 13px;
  padding: 0 16px 24px;
}

/* ── Header ──────────────────────────────────────── */
.cp-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px 0 16px;
  border-bottom: 1px solid var(--border);
  margin-bottom: 24px;
}

.cp-brand {
  display: flex;
  align-items: center;
  gap: 10px;
}

.cp-brand-icon {
  font-size: 22px;
  color: var(--green);
  filter: drop-shadow(0 0 6px var(--green));
  line-height: 1;
}

.cp-brand-name {
  font-size: 18px;
  font-weight: 700;
  letter-spacing: 0.15em;
  color: var(--green);
  text-shadow: 0 0 20px rgba(0,255,136,0.4);
}

.cp-header-right {
  display: flex;
  align-items: center;
  gap: 12px;
}

/* ── Lang + Theme buttons ────────────────────────── */
.cp-lang-btn,
.cp-theme-btn {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-family: var(--mono);
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 0.05em;
  border: 1px solid var(--border-hi);
  border-radius: 6px;
  background: transparent;
  color: var(--green);
  cursor: pointer;
  transition: all 0.18s;
  flex-shrink: 0;
  line-height: 1;
}

.cp-lang-btn:hover,
.cp-theme-btn:hover {
  background: rgba(0,255,136,0.08);
  box-shadow: 0 0 12px rgba(0,255,136,0.2);
}

/* ── Dir controls ────────────────────────────────── */
.cp-dir-controls {
  display: flex;
  align-items: center;
}

.cp-dir-toggle {
  display: flex;
  align-items: center;
  background: var(--surface);
  border: 1px solid var(--border-hi);
  border-radius: 6px;
  padding: 3px;
  gap: 2px;
}

.cp-dir-opt {
  padding: 6px 14px;
  font-family: var(--mono);
  font-size: 11px;
  font-weight: 600;
  letter-spacing: 0.08em;
  border: 1px solid transparent;
  border-radius: 4px;
  background: transparent;
  color: var(--text-dim);
  cursor: pointer;
  transition: all 0.18s;
  white-space: nowrap;
}

.cp-dir-opt:hover {
  color: var(--text);
}

.cp-dir-opt--active {
  background: rgba(0, 255, 136, 0.12);
  border-color: var(--border-hi);
  color: var(--green);
  text-shadow: 0 0 12px rgba(0,255,136,0.5);
}

.cp-dir-swap {
  padding: 4px 8px;
  font-size: 16px;
  font-family: var(--mono);
  border: none;
  border-radius: 3px;
  background: transparent;
  color: var(--text-dim);
  cursor: pointer;
  transition: all 0.15s;
  line-height: 1;
}

.cp-dir-swap:hover {
  color: var(--green);
  text-shadow: 0 0 10px rgba(0,255,136,0.6);
}

/* ── Panels layout ───────────────────────────────── */
.cp-panels {
  display: flex;
  gap: 0;
  flex: 1;
  align-items: stretch;
}

.cp-panel {
  flex: 1;
  display: flex;
  flex-direction: column;
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: 6px;
  overflow: hidden;
  transition: border-color 0.2s;
}

.cp-panel:focus-within {
  border-color: var(--border-hi);
}

.cp-arrow {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0 16px;
  font-size: 28px;
  color: var(--text-dim);
  flex-shrink: 0;
  user-select: none;
}

/* ── Panel sections ──────────────────────────────── */
.cp-panel-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 14px;
  border-bottom: 1px solid var(--border);
  background: rgba(0,255,136,0.03);
  flex-shrink: 0;
}

.cp-panel-title {
  font-size: 10px;
  font-weight: 700;
  letter-spacing: 0.2em;
  color: var(--green-dim);
}

.cp-panel-body {
  flex: 1;
  display: flex;
  flex-direction: column;
  padding: 14px;
}

.cp-panel-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: 10px;
}

/* ── Tabs ────────────────────────────────────────── */
.cp-mode-toggle,
.cp-output-controls {
  display: flex;
  align-items: center;
  gap: 4px;
}

.cp-tab {
  padding: 3px 10px;
  font-family: var(--mono);
  font-size: 10px;
  letter-spacing: 0.1em;
  border: 1px solid var(--border);
  border-radius: 3px;
  background: transparent;
  color: var(--text-dim);
  cursor: pointer;
  transition: all 0.15s;
}

.cp-tab:hover {
  border-color: var(--border-hi);
  color: var(--text);
}

.cp-tab--active {
  background: rgba(0,255,136,0.1);
  border-color: var(--border-hi);
  color: var(--green);
}

/* ── Input fields ────────────────────────────────── */
.cp-mono {
  font-family: var(--mono) !important;
}

.cp-input {
  width: 100%;
  padding: 10px 12px;
  font-size: 13px;
  background: var(--surface-hi);
  border: 1px solid var(--border);
  border-radius: 4px;
  color: var(--text-bright);
  outline: none;
  transition: border-color 0.2s, box-shadow 0.2s;
  caret-color: var(--green);
}

.cp-input::placeholder {
  color: var(--text-dim);
}

.cp-input:focus {
  border-color: var(--border-hi);
  box-shadow: 0 0 0 1px rgba(0,255,136,0.12), inset 0 0 20px rgba(0,255,136,0.03);
}

.cp-textarea {
  flex: 1;
  width: 100%;
  min-height: 200px;
  padding: 10px 12px;
  font-size: 13px;
  background: var(--surface-hi);
  border: 1px solid var(--border);
  border-radius: 4px;
  color: var(--text-bright);
  outline: none;
  resize: vertical;
  line-height: 1.7;
  transition: border-color 0.2s, box-shadow 0.2s;
  caret-color: var(--green);
}

.cp-textarea::placeholder {
  color: var(--text-dim);
}

.cp-textarea:focus {
  border-color: var(--border-hi);
  box-shadow: 0 0 0 1px rgba(0,255,136,0.12), inset 0 0 20px rgba(0,255,136,0.03);
}

/* ── Hints ───────────────────────────────────────── */
.cp-hint {
  font-size: 11px;
  color: var(--text-dim);
  letter-spacing: 0.05em;
}

/* ── Action button ───────────────────────────────── */
.cp-action-btn {
  padding: 5px 14px;
  font-family: var(--mono);
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 0.12em;
  background: rgba(0,255,136,0.08);
  border: 1px solid var(--green-dim);
  border-radius: 4px;
  color: var(--green);
  cursor: pointer;
  transition: all 0.15s;
  flex-shrink: 0;
}

.cp-action-btn:hover:not(:disabled) {
  background: rgba(0,255,136,0.15);
  box-shadow: 0 0 14px rgba(0,255,136,0.25);
}

.cp-action-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

/* ── Copy button ─────────────────────────────────── */
.cp-copy-btn {
  padding: 3px 10px;
  font-family: var(--mono);
  font-size: 10px;
  letter-spacing: 0.1em;
  border: 1px solid var(--border);
  border-radius: 3px;
  background: transparent;
  color: var(--text-dim);
  cursor: pointer;
  transition: all 0.15s;
}

.cp-copy-btn:hover:not(:disabled) {
  border-color: var(--border-hi);
  color: var(--text);
}

.cp-copy-btn--ok {
  border-color: var(--green) !important;
  color: var(--green) !important;
}

.cp-copy-btn:disabled {
  opacity: 0.3;
  cursor: not-allowed;
}

/* ── Output ──────────────────────────────────────── */
.cp-output-state {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 120px;
}

.cp-output {
  flex: 1;
  width: 100%;
  padding: 0;
  margin: 0;
  font-size: 13px;
  line-height: 1.7;
  color: var(--green);
  white-space: pre-wrap;
  word-break: break-all;
  text-shadow: 0 0 10px rgba(0,255,136,0.2);
  overflow-y: auto;
}

.cp-error-msg {
  color: var(--error);
  font-size: 12px;
  text-align: center;
}

/* ── Blink animation ─────────────────────────────── */
.cp-blink {
  animation: blink 1s step-end infinite;
}

@keyframes blink {
  0%, 100% { opacity: 1; }
  50% { opacity: 0; }
}

/* ── Footer ──────────────────────────────────────── */
.cp-footer {
  padding: 16px 0 0;
  text-align: center;
  border-top: 1px solid var(--border);
  margin-top: 24px;
}

/* ── Toast ───────────────────────────────────────── */
.cp-toast {
  position: fixed;
  bottom: 24px;
  right: 24px;
  background: rgba(0,255,136,0.12);
  border: 1px solid var(--green-dim);
  color: var(--green);
  padding: 8px 16px;
  border-radius: 4px;
  font-size: 12px;
  letter-spacing: 0.1em;
  pointer-events: none;
}

.toast-enter-active,
.toast-leave-active {
  transition: opacity 0.3s, transform 0.3s;
}

.toast-enter-from,
.toast-leave-to {
  opacity: 0;
  transform: translateY(8px);
}

/* ── Light theme ─────────────────────────────────── */
.cp-root--light {
  --bg: #f0f5f0;
  --surface: #ffffff;
  --surface-hi: #f8fbf8;
  --border: rgba(0, 100, 50, 0.15);
  --border-hi: rgba(0, 100, 50, 0.4);
  --green: #006633;
  --green-dim: rgba(0, 102, 51, 0.75);
  --cyan: #007788;
  --text: #2a4a2a;
  --text-bright: #0a1a0a;
  --text-dim: #7a9a7a;
  --error: #cc0033;
}

/* ── Responsive ──────────────────────────────────── */
@media (max-width: 768px) {
  .cp-panels {
    flex-direction: column;
  }

  .cp-arrow {
    padding: 8px 0;
    rotate: 90deg;
    font-size: 22px;
  }

  .cp-brand-name {
    font-size: 15px;
  }

  .cp-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
  }

  .cp-dir-opt {
    padding: 5px 10px;
    font-size: 10px;
  }
}
</style>
