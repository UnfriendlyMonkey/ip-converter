<template lang="pug">
div.ummy
  div.input-wrap
    div.input-label
      span ▸ {{ t('inputTitle') }}
      span.detected(v-if="detectedLabel") ⊕ {{ detectedLabel }}
      span.detected.miss(v-else-if="input.trim()") ? {{ t('unrecognized') }}
      span.detected.miss(v-else) — {{ t('awaiting') }}

    textarea.input-field(
      v-if="bulk"
      v-model="input"
      rows="5"
      spellcheck="false"
      :placeholder="t('bulkPlaceholder')"
    )
    input.input-field(
      v-else
      v-model="input"
      spellcheck="false"
      autocomplete="off"
      :placeholder="t('singlePlaceholder')"
    )

    div.input-actions
      button(@click="setNow") {{ t('now') }}
      button(@click="bulk = !bulk") {{ bulk ? t('singleMode') : t('bulkMode') }}
      button(@click="input = ''") {{ t('clear') }}
      button(@click="pasteInput") {{ t('paste') }}

  div.input-hint
    span
      | try:
      a.sample(href="#" @click.prevent="input = '1714000000'") 1714000000
      span.sep ·
      a.sample(href="#" @click.prevent="input = '2026-04-26T12:34:56Z'") 2026-04-26T12:34:56Z
      span.sep ·
      a.sample(href="#" @click.prevent="input = '1714000000000'") 1714000000000
    span.hint-keys / focus · ⌘K focus · ⌘↵ NOW · B bulk · 1/2 tabs

  table.bulk-table(v-if="bulk")
    thead
      tr
        th #
        th {{ t('inputCol') }}
        th {{ t('detectedCol') }}
        th {{ t('primaryCol') }}
        th
    tbody
      tr(v-for="(row, idx) in bulkRows" :key="idx")
        td {{ String(idx + 1).padStart(2, '0') }}
        td.input-cell {{ row.input }}
        td(:class="{ invalid: !row.detected }") {{ row.detected || '—' }}
        td(:class="{ invalid: !row.detected }") {{ row.primary }}
        td
          button.result-copy(v-if="row.detected" @click="copy(row.primary)") ⧉

  div.results(v-else)
    div.result(v-for="row in rows" :key="row.id" :class="{ featured: row.featured }")
      div.result-label
        span.name {{ row.name }}
        span.desc {{ row.desc }}
      div.result-value(:class="{ invalid: row.invalid }") {{ row.value }}
      button.result-copy(:disabled="row.invalid" @click="copy(row.value)") {{ copiedValue === row.value ? '✓' : '⧉ COPY' }}

  div.section-head(@click="historyOpen = !historyOpen" :data-open="historyOpen ? '1' : '0'")
    h2
      span.chev ▾
      span {{ t('recent') }}
      span.count {{ history.length }}
    div.section-actions
      button(v-if="history.length" @click.stop="clearHistory") {{ t('clear') }}

  div.history-wrap(:data-open="historyOpen ? '1' : '0'")
    div
      div.empty(v-if="!history.length") {{ t('historyEmpty') }}
      div.history(v-else)
        div.history-row(v-for="(entry, idx) in history.slice(0, 8)" :key="idx" @click="input = entry.input")
          span.type {{ entry.detected }}
          span.val {{ entry.input }}
          span.time {{ timeAgo(entry.ts) }}
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useSettings } from '../composables/useSettings'

const { lang } = useSettings()
type Lang = 'en' | 'ru'

const tr = {
  en: {
    inputTitle: 'INPUT · PASTE ANYTHING',
    awaiting: 'AWAITING INPUT',
    unrecognized: 'UNRECOGNIZED',
    singlePlaceholder: 'Paste a timestamp · 1714000000 · 2026-04-26T12:34:56Z · ...',
    bulkPlaceholder: 'One value per line...',
    now: '⌘ NOW',
    bulkMode: '≡ BULK MODE',
    singleMode: '× SINGLE MODE',
    clear: '⌫ CLEAR',
    paste: '⤓ PASTE',
    inputCol: 'Input',
    detectedCol: 'Detected',
    primaryCol: 'Primary output',
    recent: 'Recent',
    historyEmpty: 'No history yet. Conversions you make will appear here.',
  },
  ru: {
    inputTitle: 'ВВОД · ВСТАВЬТЕ ЧТО-НИБУДЬ',
    awaiting: 'ОЖИДАНИЕ ВВОДА',
    unrecognized: 'НЕ РАСПОЗНАНО',
    singlePlaceholder: 'Вставьте timestamp · 1714000000 · 2026-04-26T12:34:56Z · ...',
    bulkPlaceholder: 'По одному значению на строку...',
    now: '⌘ NOW',
    bulkMode: '≡ МНОГО',
    singleMode: '× ОДИН',
    clear: '⌫ ОЧИСТИТЬ',
    paste: '⤓ ВСТАВИТЬ',
    inputCol: 'Ввод',
    detectedCol: 'Тип',
    primaryCol: 'Основной результат',
    recent: 'Недавние',
    historyEmpty: 'История пока пуста.',
  },
} as const

function t(k: keyof typeof tr.en) {
  return tr[lang.value as Lang][k]
}

type Row = { id: string; name: string; desc: string; value: string; invalid?: boolean; featured?: boolean }
const input = ref('')
const bulk = ref(false)
const copiedValue = ref('')
const historyOpen = ref(true)
const history = ref<Array<{ input: string; detected: string; ts: number }>>(
  JSON.parse(localStorage.getItem('ummy.ts.history') || '[]'),
)

const detectLabels: Record<string, string> = {
  'unix-s': 'UNIX seconds',
  'unix-ms': 'UNIX millis',
  'unix-us': 'UNIX micros',
  'unix-ns': 'UNIX nanos',
  iso: 'ISO 8601 date',
}

function detectTimestamp(raw: string): { kind: string | null; value?: number } {
  const s = raw.trim()
  if (!s) return { kind: null }
  if (/^-?\d+$/.test(s)) {
    const n = Number(s)
    if (s.length <= 10) return { kind: 'unix-s', value: n }
    if (s.length <= 13) return { kind: 'unix-ms', value: n }
    if (s.length <= 16) return { kind: 'unix-us', value: n }
    return { kind: 'unix-ns', value: n }
  }
  const d = new Date(s)
  if (!Number.isNaN(d.getTime())) return { kind: 'iso', value: d.getTime() }
  return { kind: null }
}

function relativeTime(ms: number): string {
  const diff = ms - Date.now()
  const abs = Math.abs(diff)
  const sec = abs / 1000
  const past = diff < 0
  const fmt = (n: number, u: string) => `${Math.round(n)} ${u}${n >= 2 ? 's' : ''} ${past ? 'ago' : 'from now'}`
  if (sec < 60) return fmt(sec, 'second')
  if (sec < 3600) return fmt(sec / 60, 'minute')
  if (sec < 86400) return fmt(sec / 3600, 'hour')
  return fmt(sec / 86400, 'day')
}

function rowsFor(value: string): { detected: string | null; rows: Row[] } {
  const det = detectTimestamp(value)
  if (!det.kind || det.value == null) {
    const empty = '—'
    return {
      detected: null,
      rows: [
        { id: 'unix-s', name: 'UNIX seconds', desc: 'Epoch · 10 digits', value: empty, invalid: true },
        { id: 'unix-ms', name: 'UNIX millis', desc: 'Epoch · 13 digits', value: empty, invalid: true },
        { id: 'iso', name: 'ISO 8601', desc: 'UTC · RFC 3339', value: empty, invalid: true },
        { id: 'local', name: 'Local', desc: 'System TZ', value: empty, invalid: true },
        { id: 'human', name: 'Human', desc: 'Locale formatted', value: empty, invalid: true },
        { id: 'rel', name: 'Relative', desc: 'From now', value: empty, invalid: true },
        { id: 'rfc2822', name: 'RFC 2822', desc: 'Email / HTTP', value: empty, invalid: true },
        { id: 'hex', name: 'Hex (ms)', desc: 'Base 16', value: empty, invalid: true },
      ],
    }
  }

  let ms = det.value
  if (det.kind === 'unix-s') ms = det.value * 1000
  else if (det.kind === 'unix-us') ms = Math.floor(det.value / 1000)
  else if (det.kind === 'unix-ns') ms = Math.floor(det.value / 1e6)
  const d = new Date(ms)
  const tzOffsetMin = -d.getTimezoneOffset()
  const sign = tzOffsetMin >= 0 ? '+' : '-'
  const pad = (n: number) => String(n).padStart(2, '0')
  const tz = `UTC${sign}${pad(Math.floor(Math.abs(tzOffsetMin) / 60))}:${pad(Math.abs(tzOffsetMin) % 60)}`
  const local = `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}:${pad(d.getSeconds())}`
  const human = d.toLocaleString()

  return {
    detected: detectLabels[det.kind],
    rows: [
      { id: 'unix-s', name: 'UNIX seconds', desc: 'Epoch · 10 digits', value: String(Math.floor(ms / 1000)), featured: det.kind === 'unix-s' },
      { id: 'unix-ms', name: 'UNIX millis', desc: 'Epoch · 13 digits', value: String(ms), featured: det.kind === 'unix-ms' || det.kind === 'iso' },
      { id: 'iso', name: 'ISO 8601', desc: 'UTC · RFC 3339', value: d.toISOString(), featured: det.kind === 'iso' },
      { id: 'local', name: 'Local', desc: tz, value: local },
      { id: 'human', name: 'Human', desc: 'Locale formatted', value: human },
      { id: 'rel', name: 'Relative', desc: 'From now', value: relativeTime(ms) },
      { id: 'rfc2822', name: 'RFC 2822', desc: 'Email / HTTP', value: d.toUTCString() },
      { id: 'hex', name: 'Hex (ms)', desc: 'Base 16', value: `0x${ms.toString(16).toUpperCase()}` },
    ],
  }
}

const current = computed(() => rowsFor(input.value))
const rows = computed(() => current.value.rows)
const detectedLabel = computed(() => current.value.detected)
const bulkRows = computed(() =>
  input.value.split('\n').map(v => v.trim()).filter(Boolean).map((line) => {
    const r = rowsFor(line)
    const primary = r.rows.find(x => x.featured) || r.rows[0]
    return { input: line, detected: r.detected, primary: primary.value }
  }),
)

watch(input, (val) => {
  if (!val.trim()) return
  if (!detectedLabel.value) return
  history.value = [{ input: val, detected: detectedLabel.value, ts: Date.now() }, ...history.value.filter(h => h.input !== val)].slice(0, 30)
  localStorage.setItem('ummy.ts.history', JSON.stringify(history.value))
})

function setNow() {
  input.value = String(Date.now())
}

async function copy(value: string) {
  if (!value || value === '—') return
  await navigator.clipboard.writeText(value)
  copiedValue.value = value
  setTimeout(() => { copiedValue.value = '' }, 900)
}

async function pasteInput() {
  try { input.value = await navigator.clipboard.readText() } catch {}
}

function clearHistory() {
  history.value = []
  localStorage.removeItem('ummy.ts.history')
}

function timeAgo(ts: number): string {
  const s = (Date.now() - ts) / 1000
  if (s < 60) return 'just now'
  if (s < 3600) return `${Math.floor(s / 60)}m ago`
  if (s < 86400) return `${Math.floor(s / 3600)}h ago`
  return `${Math.floor(s / 86400)}d ago`
}
</script>

<style scoped>
.ummy { color: var(--fg); }
.input-wrap { border: var(--border-w) solid var(--line); background: var(--card); box-shadow: var(--shadow-offset) var(--shadow-offset) 0 var(--shadow-color); margin-bottom: 8px; }
.input-label { display: flex; justify-content: space-between; align-items: center; padding: 10px 14px; border-bottom: var(--border-w) solid var(--line); font-size: 11px; letter-spacing: .12em; text-transform: uppercase; font-weight: 600; background: var(--fg); color: var(--bg); }
.detected { background: var(--accent); color: var(--accent-fg); padding: 3px 8px; font-size: 10px; }
.detected.miss { background: transparent; color: var(--bg); border: 1px solid var(--bg); }
.input-field { width: 100%; border: none; background: transparent; color: var(--fg); font-family: var(--mono); font-size: 34px; padding: 16px 16px; outline: none; resize: none; line-height: 1.25; letter-spacing: .01em; }
.input-actions { display: flex; border-top: var(--border-w) solid var(--line); }
.input-actions button { flex: 1; background: transparent; border: none; border-right: var(--border-w) solid var(--line); padding: 9px 12px; font-family: var(--mono); font-size: 10px; letter-spacing: .1em; text-transform: uppercase; font-weight: 600; color: var(--fg); cursor: pointer; }
.input-actions button:last-child { border-right: none; }
.input-actions button:hover { background: var(--accent); color: var(--accent-fg); }
.input-hint { font-size: 10px; color: var(--muted); margin-bottom: 22px; display: flex; align-items: center; justify-content: space-between; gap: 8px; flex-wrap: wrap; letter-spacing: .04em; }
.sample { color: var(--fg); text-decoration: underline; text-underline-offset: 2px; margin-left: 6px; }
.sep { color: var(--muted); margin: 0 4px; }
.hint-keys { text-transform: lowercase; }
.results { display: grid; gap: var(--gap); }
.result { border: var(--border-w) solid var(--line); background: var(--card); display: grid; grid-template-columns: 166px 1fr auto; }
.result-label { background: var(--bg); border-right: var(--border-w) solid var(--line); padding: var(--pad); display: flex; flex-direction: column; gap: 4px; justify-content: center; }
.name { font-size: 12px; font-weight: 700; letter-spacing: .08em; text-transform: uppercase; }
.desc { font-size: 10px; color: var(--muted); }
.result-value { padding: var(--pad); font-size: 15px; display: flex; align-items: center; user-select: all; word-break: break-all; }
.result-value.invalid { color: var(--muted); font-style: italic; }
.result-copy { border: none; border-left: var(--border-w) solid var(--line); padding: 0 16px; font-size: 11px; font-weight: 600; letter-spacing: .1em; background: transparent; color: var(--fg); cursor: pointer; text-transform: uppercase; }
.result-copy:hover:not(:disabled) { background: var(--accent); color: var(--accent-fg); }
.result-copy:disabled { opacity: .45; cursor: not-allowed; }
.featured { border-color: var(--accent); position: relative; }
.featured::before { content: 'PRIMARY'; position: absolute; top: -2px; left: -2px; background: var(--accent); color: var(--accent-fg); font-size: 9px; letter-spacing: .15em; font-weight: 700; padding: 3px 8px; }
.section-head { display: flex; justify-content: space-between; align-items: baseline; margin: 34px 0 10px; padding-bottom: 8px; border-bottom: 1px dashed var(--line); cursor: pointer; }
.section-head h2 { margin: 0; display: flex; align-items: center; gap: 10px; font-size: 14px; letter-spacing: .1em; text-transform: uppercase; }
.chev { display: inline-block; width: 14px; font-size: 12px; color: var(--muted); transition: transform 160ms; }
.section-head[data-open='0'] .chev { transform: rotate(-90deg); }
.count { font-size: 11px; color: var(--muted); border: 1px solid var(--line); padding: 2px 8px; }
.section-actions button { border: 1px solid var(--line); background: transparent; padding: 4px 10px; font-size: 10px; text-transform: uppercase; cursor: pointer; color: var(--fg); }
.section-actions button:hover { background: var(--accent); color: var(--accent-fg); }
.history-wrap { display: grid; grid-template-rows: 1fr; overflow: hidden; transition: grid-template-rows .2s, opacity .2s; }
.history-wrap[data-open='0'] { grid-template-rows: 0fr; opacity: 0; }
.history-wrap > div { min-height: 0; }
.history { display: grid; gap: 6px; }
.history-row { display: grid; grid-template-columns: 90px 1fr auto; gap: 12px; padding: 7px 12px; border: 1px solid var(--line); background: var(--card); font-size: 12px; align-items: center; cursor: pointer; }
.history-row:hover { background: color-mix(in srgb, var(--accent) 15%, var(--card)); }
.type { font-size: 10px; text-transform: uppercase; color: var(--muted); }
.val { font-weight: 500; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.time { font-size: 10px; color: var(--muted); }
.empty { padding: 20px; border: 1px dashed var(--line); background: var(--card); text-align: center; color: var(--muted); font-size: 12px; }
.bulk-table { width: 100%; border-collapse: collapse; border: var(--border-w) solid var(--line); background: var(--card); font-size: 13px; }
.bulk-table th,.bulk-table td { border: 1px solid var(--line); padding: 8px 10px; text-align: left; }
.bulk-table th { background: var(--fg); color: var(--bg); font-size: 10px; letter-spacing: .12em; text-transform: uppercase; }
.input-cell { font-weight: 600; }
.invalid { color: var(--muted); font-style: italic; }
@media (max-width: 720px) {
  .input-field { font-size: 24px; }
  .result { grid-template-columns: 1fr; }
  .result-label { border-right: none; border-bottom: var(--border-w) solid var(--line); }
  .result-copy { border-left: none; border-top: var(--border-w) solid var(--line); padding: 10px; }
}
</style>
