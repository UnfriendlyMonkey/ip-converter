<template lang="pug">
v-app(theme="hacker")
  v-main.app-main
    .app-root(:class="{ 'app-root--light': !darkMode }")
      header.app-header
        div.app-brand
          span.app-brand-icon ⬡
          span.app-brand-name CONVERTER

        nav.app-nav
          router-link.app-nav-tab(
            to="/ip"
            active-class=""
            exact-active-class="app-nav-tab--active"
          ) IP CONVERTER
          router-link.app-nav-tab(
            to="/timestamp"
            active-class=""
            exact-active-class="app-nav-tab--active"
          ) TIMESTAMP

        div.app-header-right
          button.app-ctrl-btn(
            @click="toggleLang"
            :title="lang === 'en' ? 'Switch to Russian' : 'Switch to English'"
          ) {{ lang === 'en' ? 'RU' : 'EN' }}
          button.app-ctrl-btn(
            @click="toggleTheme"
            :title="darkMode ? 'Switch to light mode' : 'Switch to dark mode'"
          ) {{ darkMode ? '☀' : '☾' }}

      router-view
</template>

<script setup lang="ts">
import { useSettings } from './composables/useSettings'

const { darkMode, lang, toggleTheme, toggleLang } = useSettings()
</script>

<style>
.app-main {
  background: transparent;
  min-height: 100vh;
}
</style>

<style scoped>
/* ── Design tokens ──────────────────────────────── */
.app-root {
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

/* ── Light theme ─────────────────────────────────── */
.app-root--light {
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

/* ── Header ──────────────────────────────────────── */
.app-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px 0 16px;
  border-bottom: 1px solid var(--border);
  margin-bottom: 24px;
}

.app-brand {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-shrink: 0;
}

.app-brand-icon {
  font-size: 22px;
  color: var(--green);
  filter: drop-shadow(0 0 6px var(--green));
  line-height: 1;
}

.app-brand-name {
  font-size: 18px;
  font-weight: 700;
  letter-spacing: 0.15em;
  color: var(--green);
  text-shadow: 0 0 20px rgba(0,255,136,0.4);
}

/* ── Nav tabs ────────────────────────────────────── */
.app-nav {
  display: flex;
  align-items: center;
  gap: 3px;
  background: var(--surface);
  border: 1px solid var(--border-hi);
  border-radius: 6px;
  padding: 3px;
}

.app-nav-tab {
  padding: 5px 16px;
  font-family: var(--mono);
  font-size: 11px;
  font-weight: 600;
  letter-spacing: 0.1em;
  border: 1px solid transparent;
  border-radius: 4px;
  background: transparent;
  color: var(--text-dim);
  cursor: pointer;
  text-decoration: none;
  transition: all 0.18s;
  white-space: nowrap;
}

.app-nav-tab:hover {
  color: var(--text);
}

.app-nav-tab--active {
  background: rgba(0, 255, 136, 0.12);
  border-color: var(--border-hi);
  color: var(--green);
  text-shadow: 0 0 12px rgba(0,255,136,0.5);
}

/* ── Controls ────────────────────────────────────── */
.app-header-right {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-shrink: 0;
}

.app-ctrl-btn {
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
  line-height: 1;
}

.app-ctrl-btn:hover {
  background: rgba(0,255,136,0.08);
  box-shadow: 0 0 12px rgba(0,255,136,0.2);
}

/* ── Responsive ──────────────────────────────────── */
@media (max-width: 768px) {
  .app-header {
    flex-wrap: wrap;
    gap: 12px;
  }

  .app-brand-name {
    font-size: 15px;
  }

  .app-nav {
    order: 3;
    flex-basis: 100%;
    justify-content: center;
  }
}
</style>
