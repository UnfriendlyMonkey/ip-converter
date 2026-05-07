<template lang="pug">
v-app(theme="hacker")
  v-main.app-main
    .app-root(data-density="comfy" data-radius="sharp")
      header.app-header
        div.app-brand
          span.app-brand-mark
            svg(width="34" height="34" viewBox="0 0 34 34" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="square")
              path(d="M6 11h17l-4 -4")
              path(d="M28 23h-17l4 4")
          span.app-brand-text
            | ummy/converter
            small brutal little tools for devs

        div.app-header-right
          span.app-pill.app-pill-live {{ clockLabel }}
          button.app-pill(@click="copyShareLink") ↗ SHARE
          button.app-pill(@click="cycleTheme" :title="`Theme: ${theme}`") {{ theme.toUpperCase() }}
          button.app-pill(@click="toggleLang") {{ lang === 'en' ? 'RU' : 'EN' }}

      div.app-header-divider

      nav.app-tabs
        router-link.app-nav-tab(
          to="/ip"
          active-class=""
          exact-active-class="app-nav-tab--active"
        )
          span.app-tab-num 01
          | IP ADDRESS
        router-link.app-nav-tab(
          to="/timestamp"
          active-class=""
          exact-active-class="app-nav-tab--active"
        )
          span.app-tab-num 02
          | TIMESTAMP

      router-view
</template>

<script setup lang="ts">
import { onMounted, onUnmounted, ref, watchEffect } from 'vue'
import { useSettings } from './composables/useSettings'

const { lang, theme, cycleTheme, toggleLang } = useSettings()
const clockLabel = ref('00:00:00')

function updateClock() {
  const d = new Date()
  const p = (n: number) => String(n).padStart(2, '0')
  clockLabel.value = `${p(d.getHours())}:${p(d.getMinutes())}:${p(d.getSeconds())}`
}

async function copyShareLink() {
  try { await navigator.clipboard.writeText(window.location.href) } catch {}
}

let timer: ReturnType<typeof setInterval> | null = null
onMounted(() => {
  updateClock()
  timer = setInterval(updateClock, 1000)
})
onUnmounted(() => {
  if (timer) clearInterval(timer)
})

watchEffect(() => {
  document.documentElement.setAttribute('data-theme', theme.value)
})
</script>

<style>
:root {
  --border-w: 2px;
  --radius: 0px;
  --pad: 16px;
  --gap: 12px;
  --shadow-offset: 6px;
  --bg: #f4f1ea;
  --fg: #0a0a0a;
  --muted: #6b6660;
  --line: #0a0a0a;
  --card: #ffffff;
  --accent: #ff5b1f;
  --accent-fg: #0a0a0a;
  --warn: #d4341a;
  --good: #1f7a3a;
  --shadow-color: #0a0a0a;
  --error: var(--warn);
  --surface: var(--card);
  --surface-hi: color-mix(in srgb, var(--card) 90%, var(--bg));
  --border: color-mix(in srgb, var(--line) 25%, transparent);
  --border-hi: var(--line);
  --green: var(--accent);
  --green-dim: color-mix(in srgb, var(--accent) 70%, transparent);
  --text: var(--fg);
  --text-bright: var(--fg);
  --text-dim: var(--muted);
  --mono: 'JetBrains Mono', ui-monospace, Menlo, monospace;
}

:root[data-theme='coal'] {
  --bg: #0e0e0e;
  --fg: #f0ece4;
  --muted: #8a8680;
  --line: #f0ece4;
  --card: #181818;
  --accent: #ffe14a;
  --accent-fg: #0a0a0a;
  --warn: #f59e0b;
  --good: #86efac;
  --shadow-color: #ffe14a;
}

:root[data-theme='terminal'] {
  --bg: #050805;
  --fg: #b8ffb8;
  --muted: #4a8a4a;
  --line: #4ade80;
  --card: #0a120a;
  --accent: #4ade80;
  --accent-fg: #050805;
  --warn: #f59e0b;
  --good: #86efac;
  --shadow-color: #4ade80;
}

:root[data-theme='pulp'] {
  --bg: #fff8e7;
  --fg: #1a0e08;
  --muted: #6e4a2a;
  --line: #1a0e08;
  --card: #ffefc4;
  --accent: #e63946;
  --accent-fg: #fff8e7;
  --warn: #dc2626;
  --good: #1f7a3a;
  --shadow-color: #1a0e08;
}

:root[data-theme='riot'] {
  --bg: #1f2330;
  --fg: #ece6d6;
  --muted: #8a8676;
  --line: #ece6d6;
  --card: #272c3b;
  --accent: #e84a2a;
  --accent-fg: #ece6d6;
  --warn: #f97316;
  --good: #86efac;
  --shadow-color: #f4c542;
}

.app-main {
  background: var(--bg);
  min-height: 100vh;
}

html,
body {
  margin: 0;
  padding: 0;
  background: var(--bg);
  color: var(--fg);
  font-family: 'JetBrains Mono', ui-monospace, Menlo, monospace;
}

body::before {
  content: '';
  position: fixed;
  inset: 0;
  background-image:
    linear-gradient(to right, color-mix(in srgb, var(--fg) 6%, transparent) 1px, transparent 1px),
    linear-gradient(to bottom, color-mix(in srgb, var(--fg) 6%, transparent) 1px, transparent 1px);
  background-size: 32px 32px;
  pointer-events: none;
  z-index: 0;
}

#app,
.v-application,
.v-application__wrap {
  background: transparent !important;
}
</style>

<style scoped>
.app-root {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
  color: var(--fg);
  font-family: var(--mono);
  font-size: 14px;
  padding: 32px 24px 96px;
  max-width: 1100px;
  margin: 0 auto;
  position: relative;
  z-index: 1;
}

.app-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 24px;
  margin-bottom: 18px;
}

.app-brand {
  display: flex;
  align-items: center;
  gap: 14px;
  flex-shrink: 0;
}

.app-brand-mark {
  width: 56px;
  height: 56px;
  display: grid;
  place-items: center;
  background: var(--fg);
  color: var(--bg);
  border: var(--border-w) solid var(--line);
}

.app-brand-text {
  display: inline-flex;
  flex-direction: column;
  font-size: 24px;
  font-weight: 700;
  letter-spacing: -0.02em;
  line-height: 1;
}

.app-brand-text small {
  font-size: 11px;
  font-weight: 500;
  letter-spacing: 0.12em;
  text-transform: uppercase;
  color: var(--muted);
  margin-top: 4px;
}

.app-tabs {
  display: flex;
  width: min(760px, 100%);
  border: var(--border-w) solid var(--line);
  background: var(--card);
  margin-bottom: 22px;
}

.app-nav-tab {
  flex: 1;
  padding: 14px 20px;
  border-right: var(--border-w) solid var(--line);
  text-decoration: none;
  font-family: var(--mono);
  font-size: 13px;
  font-weight: 600;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  background: transparent;
  color: var(--fg);
  cursor: pointer;
  transition: background 80ms;
  position: relative;
  text-align: left;
}

.app-nav-tab:last-child {
  border-right: none;
}

.app-tab-num {
  display: block;
  font-size: 10px;
  letter-spacing: 0.15em;
  opacity: 0.65;
  margin-bottom: 2px;
}

.app-nav-tab:hover {
  background: color-mix(in srgb, var(--accent) 15%, transparent);
}

.app-nav-tab--active {
  background: var(--accent);
  color: var(--accent-fg);
}

.app-header-right {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-shrink: 0;
}

.app-pill {
  min-width: 56px;
  height: 30px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-family: var(--mono);
  font-size: 11px;
  font-weight: 600;
  letter-spacing: 0.1em;
  text-transform: uppercase;
  border: var(--border-w) solid var(--line);
  background: var(--card);
  color: var(--fg);
  cursor: pointer;
  transition: background 80ms;
  line-height: 1;
  padding: 0 10px;
}

.app-pill-live::before {
  content: '';
  width: 6px;
  height: 6px;
  margin-right: 6px;
  border-radius: 50%;
  background: var(--accent);
  display: inline-block;
}

.app-pill:hover {
  background: var(--accent);
  color: var(--accent-fg);
}

.app-header-divider {
  height: var(--border-w);
  width: 100%;
  background: var(--line);
  margin-bottom: 18px;
}

@media (max-width: 768px) {
  .app-root {
    padding: 20px 14px 80px;
  }

  .app-header {
    flex-wrap: wrap;
    gap: 16px;
  }

  .app-brand-text {
    font-size: 20px;
  }

  .app-tabs {
    width: 100%;
  }

  .app-nav-tab {
    padding: 12px 10px;
    font-size: 11px;
  }
}

@media (max-width: 640px) {
  .app-header-right {
    width: 100%;
    flex-wrap: wrap;
  }
}
</style>
