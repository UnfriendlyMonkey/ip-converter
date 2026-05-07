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
      button(@click="setMyIp") {{ t('myIp') }}
      button(@click="bulk = !bulk") {{ bulk ? t('singleMode') : t('bulkMode') }}
      button(@click="input = ''") {{ t('clear') }}
      button(@click="pasteInput") {{ t('paste') }}

  div.input-hint
    span
      | try:
      a.sample(href="#" @click.prevent="input = '192.168.1.1'") 192.168.1.1
      span.sep ·
      a.sample(href="#" @click.prevent="input = '3232235777'") 3232235777
      span.sep ·
      a.sample(href="#" @click.prevent="input = '0xC0A80101'") 0xC0A80101
      span.sep ·
      a.sample(href="#" @click.prevent="input = '2001:db8::1'") 2001:db8::1
    span.hint-keys / focus · ⌘K focus · ⌘↵ MY IP · B bulk · 1/2 tabs

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
    singlePlaceholder: 'Paste an IP · 192.168.1.1 · 3232235777 · 0xC0A80101 · 2001:db8::1',
    bulkPlaceholder: 'One value per line...',
    myIp: '⌘ MY IP',
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
    singlePlaceholder: 'IP · 192.168.1.1 · 3232235777 · 0xC0A80101 · 2001:db8::1',
    bulkPlaceholder: 'По одному значению на строку...',
    myIp: '⌘ МОЙ IP',
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
  JSON.parse(localStorage.getItem('ummy.ip.history') || '[]'),
)

const detectLabels: Record<string, string> = {
  'ipv4-dotted': 'IPv4 · dotted',
  'ipv4-int': 'IPv4 · integer',
  'ipv4-hex': 'IPv4 · hex',
  'ipv4-bin': 'IPv4 · binary',
  ipv6: 'IPv6 address',
}

function expandIPv6(s: string): string | null {
  const parts = s.split('::')
  if (parts.length > 2) return null
  const head = parts[0] ? parts[0].split(':') : []
  const tail = parts.length === 2 && parts[1] ? parts[1].split(':') : []
  const fill = 8 - head.length - tail.length
  if (fill < 0) return null
  const full = [...head, ...Array(fill).fill('0'), ...tail]
  if (full.length !== 8) return null
  if (!full.every(g => /^[0-9a-f]{0,4}$/i.test(g))) return null
  return full.map(g => g.padStart(4, '0').toLowerCase()).join(':')
}

function detectIp(raw: string): { kind: string | null; octets?: number[]; value?: number; expanded?: string } {
  const s = raw.trim()
  if (!s) return { kind: null }
  if (/^(\d{1,3}\.){3}\d{1,3}$/.test(s)) {
    const octets = s.split('.').map(Number)
    if (octets.every(o => o >= 0 && o <= 255)) return { kind: 'ipv4-dotted', octets }
  }
  if (/^0x[0-9a-f]+$/i.test(s)) {
    const n = parseInt(s, 16)
    if (n >= 0 && n <= 0xffffffff) return { kind: 'ipv4-hex', value: n }
  }
  if (/^[01]+$/.test(s) && s.length <= 32) {
    const n = parseInt(s, 2)
    if (!Number.isNaN(n)) return { kind: 'ipv4-bin', value: n }
  }
  if (/^\d+$/.test(s)) {
    const n = Number(s)
    if (n >= 0 && n <= 0xffffffff) return { kind: 'ipv4-int', value: n }
  }
  if (/^[0-9a-f:]+$/i.test(s) && s.includes(':')) {
    const expanded = expandIPv6(s)
    if (expanded) return { kind: 'ipv6', expanded }
  }
  return { kind: null }
}

function intToIpv4(n: number): number[] {
  return [(n >>> 24) & 255, (n >>> 16) & 255, (n >>> 8) & 255, n & 255]
}

function ipv4ToInt(o: number[]): number {
  return ((o[0] << 24) >>> 0) + (o[1] << 16) + (o[2] << 8) + o[3]
}

function rowsFor(value: string): { detected: string | null; rows: Row[] } {
  const det = detectIp(value)
  if (!det.kind) {
    const empty = '—'
    return {
      detected: null,
      rows: [
        { id: 'dotted', name: 'Dotted decimal', desc: 'IPv4 canonical', value: empty, invalid: true },
        { id: 'int', name: 'Decimal integer', desc: '32-bit unsigned', value: empty, invalid: true },
        { id: 'hex', name: 'Hexadecimal', desc: '32-bit · 0x...', value: empty, invalid: true },
        { id: 'bin', name: 'Binary', desc: '32 bits · grouped', value: empty, invalid: true },
        { id: 'ptr', name: 'Reverse DNS', desc: '.in-addr.arpa', value: empty, invalid: true },
        { id: 'mapped', name: 'IPv6 mapped', desc: '::ffff:0:0/96', value: empty, invalid: true },
      ],
    }
  }
  if (det.kind === 'ipv6' && det.expanded) {
    return {
      detected: detectLabels[det.kind],
      rows: [
        { id: 'ipv6-comp', name: 'IPv6 compressed', desc: 'Canonical form', value: value, featured: true },
        { id: 'ipv6-exp', name: 'IPv6 expanded', desc: 'Full 8-group', value: det.expanded },
        { id: 'ipv6-bin', name: 'Binary', desc: '128 bits', value: det.expanded.split(':').map(g => parseInt(g, 16).toString(2).padStart(16, '0')).join(' ') },
      ],
    }
  }

  const n = det.kind === 'ipv4-dotted' && det.octets ? ipv4ToInt(det.octets) : (det.value || 0)
  const octets = intToIpv4(n)
  const dotted = octets.join('.')
  const bin = n.toString(2).padStart(32, '0').match(/.{8}/g)?.join('.') || ''
  return {
    detected: detectLabels[det.kind],
    rows: [
      { id: 'dotted', name: 'Dotted decimal', desc: 'IPv4 canonical', value: dotted, featured: det.kind === 'ipv4-dotted' },
      { id: 'int', name: 'Decimal integer', desc: '32-bit unsigned', value: String(n), featured: det.kind === 'ipv4-int' },
      { id: 'hex', name: 'Hexadecimal', desc: '32-bit · 0x...', value: `0x${n.toString(16).toUpperCase().padStart(8, '0')}`, featured: det.kind === 'ipv4-hex' },
      { id: 'bin', name: 'Binary', desc: '32 bits · grouped', value: bin, featured: det.kind === 'ipv4-bin' },
      { id: 'ptr', name: 'Reverse DNS', desc: '.in-addr.arpa', value: `${octets.slice().reverse().join('.')}.in-addr.arpa` },
      { id: 'mapped', name: 'IPv6 mapped', desc: '::ffff:0:0/96', value: `::ffff:${((octets[0] << 8) | octets[1]).toString(16)}:${((octets[2] << 8) | octets[3]).toString(16)}` },
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
  localStorage.setItem('ummy.ip.history', JSON.stringify(history.value))
})

function setMyIp() {
  input.value = '192.168.1.1'
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
  localStorage.removeItem('ummy.ip.history')
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
