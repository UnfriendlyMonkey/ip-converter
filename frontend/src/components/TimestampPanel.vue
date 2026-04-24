<template lang="pug">
div.cp-content

  //- ── Hero single input ────────────────────────────────────
  section.cp-hero
    div.cp-hero-input-wrap
      input.cp-hero-input.cp-mono(
        v-model="heroInput"
        :placeholder="singlePlaceholder"
        type="text"
        spellcheck="false"
        autocomplete="off"
        autocorrect="off"
        @input="onHeroInput"
      )
      transition(name="toast")
        span.cp-hero-copied(v-if="heroJustCopied") {{ t('copied') }}
      button.cp-hero-history-btn(
        v-if="history.length"
        @click="openHistory"
        :title="t('historyBtn')"
      ) ⏱
    div.cp-hero-output-wrap
      pre.cp-hero-output.cp-mono(v-if="heroOutputText") {{ heroOutputText }}
      span.cp-hint.cp-hero-hint(v-else-if="heroError" style="color: var(--error)") {{ heroError }}
      span.cp-hint.cp-hero-hint(v-else) {{ t('hintSingle') }}

  //- ── Direction switcher ──────────────────────────────────
  div.cp-dir-row
    div.cp-dir-toggle
      button(
        :class="['cp-dir-opt', direction === 'to_human' && 'cp-dir-opt--active']"
        @click="setDirection('to_human')"
      ) {{ t('unixToHuman') }}
      button.cp-dir-swap(@click="toggleDirection" :title="t('swapDir')") ⇄
      button(
        :class="['cp-dir-opt', direction === 'to_unix' && 'cp-dir-opt--active']"
        @click="setDirection('to_unix')"
      ) {{ t('humanToUnix') }}

  //- ── Multi / batch section ────────────────────────────────
  div.cp-multi-section
    div.cp-panels

      //- INPUT PANEL
      section.cp-panel.cp-panel--input
        div.cp-panel-header
          span.cp-panel-title {{ t('multi') }}

        div.cp-panel-body
          textarea.cp-textarea.cp-mono(
            v-model="multiInput"
            :placeholder="multiPlaceholder"
            spellcheck="false"
            @keydown.ctrl.enter="convertMulti"
          )

          div.cp-panel-footer
            span.cp-hint {{ t('hintMulti') }}
            button.cp-action-btn(
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
    div.cp-footer-links
      a.cp-footer-link(href="/about.html" target="_blank") {{ t('about') }}
      span.cp-footer-sep ·
      a.cp-footer-link(href="/license.txt" target="_blank") {{ t('license') }}

  //- Copy toast
  transition(name="toast")
    div.cp-toast(v-if="justCopied") {{ t('toastCopied') }}

  //- ── History overlay ──────────────────────────────────────
  transition(name="overlay")
    .cp-history-overlay(v-if="historyOpen" @click.self="historyOpen = false")
      .cp-history-panel
        .cp-history-header
          .cp-history-tabs
            button(
              :class="['cp-history-tab', historyTab === 'to_human' && 'cp-history-tab--active']"
              @click="historyTab = 'to_human'"
            ) {{ t('unixToHuman') }}
            button(
              :class="['cp-history-tab', historyTab === 'to_unix' && 'cp-history-tab--active']"
              @click="historyTab = 'to_unix'"
            ) {{ t('humanToUnix') }}
          button.cp-history-close(@click="historyOpen = false") ✕
        .cp-history-list
          .cp-history-empty(v-if="!historyByTab.length")
            span.cp-hint {{ t('historyEmpty') }}
          .cp-history-item(
            v-else
            v-for="(entry, i) in historyByTab"
            :key="i"
            @click="applyHistory(entry)"
          )
            span.cp-history-in {{ entry.input }}
            span.cp-history-arrow →
            span.cp-history-out {{ entry.output }}
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { convertTimestamps, type TsDirection, type TsConversionResult } from '../api/timestamp'
import { useSettings } from '../composables/useSettings'

// ── i18n ────────────────────────────────────────────────────
type Lang = 'en' | 'ru'

const translations = {
  en: {
    unixToHuman: 'UNIX → HUMAN',
    humanToUnix: 'HUMAN → UNIX',
    swapDir:     'Swap direction',
    output:      'OUTPUT',
    multi:       'MULTI',
    list:        'LIST',
    dict:        'DICT',
    copy:        'COPY',
    copied:      '✓ COPIED',
    convert:     'CONVERT',
    processing:  'processing...',
    hintSingle:  'converts live · unix timestamp or DD.MM.YYYY HH:mm:ss',
    hintMulti:   'one value per line (or comma-separated) · Ctrl+Enter to convert',
    emptyOutput: '// output will appear here',
    footer:      'unix seconds · milliseconds auto-detected · DD.MM.YYYY HH:mm:ss UTC',
    toastCopied: 'Copied to clipboard!',
    about:       'About',
    license:     'MIT License',
    historyBtn:  'Recent conversions',
    historyEmpty:'No conversions yet',
  },
  ru: {
    unixToHuman: 'UNIX → ДАТА',
    humanToUnix: 'ДАТА → UNIX',
    swapDir:     'Сменить направление',
    output:      'ВЫВОД',
    multi:       'МНОГО',
    list:        'СПИСОК',
    dict:        'СЛОВАРЬ',
    copy:        'КОПИРОВАТЬ',
    copied:      '✓ СКОПИРОВАНО',
    convert:     'КОНВЕРТИРОВАТЬ',
    processing:  'обработка...',
    hintSingle:  'конвертирует автоматически · unix timestamp или ДД.ММ.ГГГГ ЧЧ:мм:сс',
    hintMulti:   'по одному значению в строке, или разделенные запятой · Ctrl+Enter',
    emptyOutput: '// результат появится здесь',
    footer:      'unix секунды · миллисекунды определяются автоматически · ДД.ММ.ГГГГ ЧЧ:мм:сс UTC',
    toastCopied: 'Скопировано!',
    about:       'О сервисе',
    license:     'Лицензия MIT',
    historyBtn:  'Недавние конвертации',
    historyEmpty:'Конвертаций пока нет',
  },
} as const

type TKey = keyof typeof translations.en

const { lang } = useSettings()

function t(key: TKey): string {
  return translations[lang.value as Lang][key]
}

// ── State ──────────────────────────────────────────────────
const direction = ref<TsDirection>('to_human')

// Hero (single)
const heroInput = ref('')
const heroResult = ref<TsConversionResult | null>(null)
const heroError = ref('')
const heroJustCopied = ref(false)
const historyOpen = ref(false)
const historyTab = ref<TsDirection>('to_human')
const history = ref<Array<{ input: string; output: string; direction: TsDirection }>>(
  JSON.parse(localStorage.getItem('ts_history') || '[]')
)
const historyByTab = computed(() => history.value.filter(h => h.direction === historyTab.value))

function openHistory() {
  historyTab.value = direction.value
  historyOpen.value = true
}

// Multi / batch
const multiInput = ref('')
const outputFormat = ref<'list' | 'dict'>('list')
const results = ref<TsConversionResult[]>([])
const loading = ref(false)
const apiError = ref('')
const justCopied = ref(false)

// ── Placeholders ───────────────────────────────────────────
const singlePlaceholder = computed(() =>
  direction.value === 'to_human'
    ? '1700000000  or  1700000000000 (ms)'
    : '15.11.2023 21:53:20  or  15.11.2023 21:53:20 UTC'
)

const multiPlaceholder = computed(() =>
  direction.value === 'to_human'
    ? '1700000000\n1700000000000\n1699920000'
    : '15.11.2023 21:53:20 UTC\n01.01.2024 00:00:00\n31.12.2023 23:59:59'
)

// ── Hero output ────────────────────────────────────────────
const heroOutputText = computed(() => {
  if (!heroResult.value || heroResult.value.type === 'error') return ''
  return heroResult.value.output
})

// ── Multi output formatting ────────────────────────────────
const hasErrors = computed(() => results.value.some(r => r.type === 'error'))

const outputText = computed(() => {
  if (!results.value.length) return ''

  if (outputFormat.value === 'list') {
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

// ── Hero conversion (auto-copy on success) ─────────────────
async function heroRunConversion(value: string) {
  if (!value) {
    heroResult.value = null
    heroError.value = ''
    return
  }
  try {
    const resp = await convertTimestamps({ direction: direction.value, values: [value] })
    const r = resp.results[0]
    heroResult.value = r
    heroError.value = ''
    if (r && r.type !== 'error') {
      try {
        await navigator.clipboard.writeText(r.output)
      } catch {
        const ta = document.createElement('textarea')
        ta.value = r.output
        document.body.appendChild(ta)
        ta.select()
        document.execCommand('copy')
        document.body.removeChild(ta)
      }
      heroJustCopied.value = true
      setTimeout(() => { heroJustCopied.value = false }, 2000)
      const dir = direction.value
      const sameDir = [
        { input: value, output: r.output, direction: dir },
        ...history.value.filter(h => h.direction === dir && h.input !== value),
      ].slice(0, 5)
      const otherDir = history.value.filter(h => h.direction !== dir)
      history.value = [...sameDir, ...otherDir]
      localStorage.setItem('ts_history', JSON.stringify(history.value))
    }
  } catch (e: unknown) {
    heroError.value = e instanceof Error ? e.message : 'Unknown error'
    heroResult.value = null
  }
}

const debouncedHero = debounce(() => {
  heroRunConversion(heroInput.value.trim())
}, 350)

function onHeroInput() {
  debouncedHero()
}

function applyHistory(entry: { input: string; output: string; direction: TsDirection }) {
  historyOpen.value = false
  direction.value = entry.direction
  heroInput.value = entry.input
  heroRunConversion(entry.input)
}

// ── Multi conversion ───────────────────────────────────────
async function runConversion(values: string[]) {
  if (!values.length) {
    results.value = []
    apiError.value = ''
    return
  }
  loading.value = true
  apiError.value = ''
  try {
    const resp = await convertTimestamps({ direction: direction.value, values })
    results.value = resp.results
  } catch (e: unknown) {
    apiError.value = e instanceof Error ? e.message : 'Unknown error'
    results.value = []
  } finally {
    loading.value = false
  }
}

function splitMultiValues(raw: string): string[] {
  return raw.split(/[\n,]/).map(s => s.trim()).filter(s => s.length > 0)
}

function convertMulti() {
  runConversion(splitMultiValues(multiInput.value))
}

// ── Direction switching ────────────────────────────────────
function toggleDirection() {
  setDirection(direction.value === 'to_human' ? 'to_unix' : 'to_human')
}

function setDirection(dir: TsDirection) {
  direction.value = dir
  heroInput.value = ''
  heroResult.value = null
  heroError.value = ''
  multiInput.value = ''
  results.value = []
  apiError.value = ''
}

// ── Clipboard (multi panel) ────────────────────────────────
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
/* ── Dir row ─────────────────────────────────────── */
.cp-dir-row {
  display: flex;
  justify-content: center;
  margin: 0 0 24px;
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

/* ── Hero section ────────────────────────────────── */
.cp-hero {
  margin-bottom: 20px;
  padding: 24px 28px;
  background: var(--surface);
  border: 1px solid var(--border-hi);
  border-radius: 8px;
  box-shadow: 0 0 40px rgba(0,255,136,0.04), inset 0 0 60px rgba(0,255,136,0.015);
}

.cp-hero-input-wrap {
  display: flex;
  align-items: center;
  gap: 14px;
}

.cp-hero-input {
  flex: 1;
  padding: 14px 18px;
  font-size: 20px;
  background: var(--surface-hi);
  border: 1px solid var(--border-hi);
  border-radius: 6px;
  color: var(--text-bright);
  outline: none;
  transition: border-color 0.2s, box-shadow 0.2s;
  caret-color: var(--green);
}

.cp-hero-input::placeholder {
  color: var(--text-dim);
  font-size: 13px;
}

.cp-hero-input:focus {
  border-color: var(--green);
  box-shadow: 0 0 0 2px rgba(0,255,136,0.12), inset 0 0 30px rgba(0,255,136,0.03);
}

.cp-hero-copied {
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 0.12em;
  color: var(--green);
  text-shadow: 0 0 10px rgba(0,255,136,0.5);
  white-space: nowrap;
  flex-shrink: 0;
}

.cp-hero-output-wrap {
  margin-top: 14px;
  min-height: 28px;
  display: flex;
  align-items: baseline;
}

.cp-hero-output {
  margin: 0;
  padding: 0;
  font-size: 18px;
  line-height: 1.5;
  color: var(--green);
  text-shadow: 0 0 20px rgba(0,255,136,0.35);
  white-space: pre-wrap;
  word-break: break-all;
}

.cp-hero-hint {
  font-size: 12px;
}

/* ── Multi section ───────────────────────────────── */
.cp-multi-section {
  border-top: 1px solid var(--border);
  padding-top: 24px;
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
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
}

.cp-footer-links {
  display: flex;
  align-items: center;
  gap: 8px;
}

.cp-footer-link {
  font-size: 11px;
  color: var(--text-dim);
  text-decoration: none;
  letter-spacing: 0.05em;
  transition: color 0.15s;
}

.cp-footer-link:hover {
  color: var(--green);
}

.cp-footer-sep {
  color: var(--text-dim);
  font-size: 11px;
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

/* ── History button ──────────────────────────────────────── */
.cp-hero-history-btn {
  flex-shrink: 0;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
  background: transparent;
  border: 1px solid var(--border-hi);
  border-radius: 6px;
  color: var(--text-dim);
  cursor: pointer;
  transition: all 0.18s;
  line-height: 1;
}

.cp-hero-history-btn:hover {
  color: var(--green);
  border-color: var(--green);
  box-shadow: 0 0 10px rgba(0,255,136,0.2);
}

/* ── History overlay ─────────────────────────────────────── */
.cp-history-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.65);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 100;
}

.cp-history-panel {
  background: var(--surface);
  border: 1px solid var(--border-hi);
  border-radius: 8px;
  box-shadow: 0 0 60px rgba(0,255,136,0.08);
  width: min(560px, calc(100vw - 32px));
  max-height: 80vh;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.cp-history-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  border-bottom: 1px solid var(--border);
  background: rgba(0,255,136,0.03);
  flex-shrink: 0;
}

.cp-history-tabs {
  display: flex;
  gap: 4px;
}

.cp-history-tab {
  padding: 4px 12px;
  font-family: var(--mono);
  font-size: 10px;
  font-weight: 600;
  letter-spacing: 0.08em;
  border: 1px solid var(--border);
  border-radius: 4px;
  background: transparent;
  color: var(--text-dim);
  cursor: pointer;
  transition: all 0.15s;
  white-space: nowrap;
}

.cp-history-tab:hover {
  border-color: var(--border-hi);
  color: var(--text);
}

.cp-history-tab--active {
  background: rgba(0, 255, 136, 0.1);
  border-color: var(--border-hi);
  color: var(--green);
}

.cp-history-close {
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  background: transparent;
  border: 1px solid var(--border);
  border-radius: 4px;
  color: var(--text-dim);
  cursor: pointer;
  transition: all 0.15s;
  line-height: 1;
}

.cp-history-close:hover {
  color: var(--green);
  border-color: var(--border-hi);
}

.cp-history-list {
  overflow-y: auto;
  padding: 8px;
  display: flex;
  flex-direction: column;
  gap: 4px;
  min-height: 60px;
}

.cp-history-empty {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 16px;
}

.cp-history-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 12px;
  border: 1px solid var(--border);
  border-radius: 5px;
  cursor: pointer;
  transition: all 0.15s;
  font-family: var(--mono);
  font-size: 13px;
  overflow: hidden;
}

.cp-history-item:hover {
  border-color: var(--border-hi);
  background: rgba(0,255,136,0.05);
}

.cp-history-in {
  color: var(--text-bright);
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  flex: 1;
}

.cp-history-arrow {
  color: var(--text-dim);
  flex-shrink: 0;
}

.cp-history-out {
  color: var(--green);
  text-shadow: 0 0 10px rgba(0,255,136,0.25);
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  flex: 1;
  text-align: right;
}

/* ── Overlay transition ──────────────────────────────────── */
.overlay-enter-active,
.overlay-leave-active {
  transition: opacity 0.2s;
}

.overlay-enter-from,
.overlay-leave-to {
  opacity: 0;
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

  .cp-dir-opt {
    padding: 5px 10px;
    font-size: 10px;
  }

  .cp-hero-input {
    font-size: 15px;
    padding: 12px 14px;
  }

  .cp-hero-output {
    font-size: 14px;
  }
}
</style>
