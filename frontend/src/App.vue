<template lang="pug">
v-app(theme="hacker")
  v-main.app-main
    .app-root(
      :data-theme="darkMode ? 'riot' : 'bone'"
      data-density="comfy"
      data-radius="sharp"
    )
      header.app-header
        router-link.app-brand(to="/timestamp")
          span.app-brand-mark(aria-hidden="true")
            svg(width="34" height="34" viewBox="0 0 34 34" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="square")
              path(d="M6 11h17l-4 -4")
              path(d="M28 23h-17l4 4")
          span.app-brand-text
            span ummy/converter
            small brutal little tools for devs

        nav.app-nav
          router-link.app-nav-tab(
            to="/timestamp"
            active-class=""
            exact-active-class="app-nav-tab--active"
          )
            span.app-nav-num 01
            span Timestamp
          router-link.app-nav-tab(
            to="/ip"
            active-class=""
            exact-active-class="app-nav-tab--active"
          )
            span.app-nav-num 02
            span IP Address

        div.app-header-right
          button.app-pill(
            @click="toggleLang"
            :title="lang === 'en' ? 'Switch to Russian' : 'Switch to English'"
          ) {{ lang === 'en' ? 'RU' : 'EN' }}
          button.app-pill(
            @click="toggleTheme"
            :title="darkMode ? 'Switch to light mode' : 'Switch to dark mode'"
          ) {{ darkMode ? 'BONE' : 'RIOT' }}

      router-view
</template>

<script setup lang="ts">
import { useSettings } from './composables/useSettings'

const { darkMode, lang, toggleTheme, toggleLang } = useSettings()
</script>

<style>
/* ummy converter theme tokens */
html,
body,
#app {
  min-height: 100%;
}

body {
  margin: 0;
  background: #1f2330;
}

.app-main {
  background: transparent;
  min-height: 100vh;
}

.app-root {
  --border-w: 2px;
  --radius: 0px;
  --pad: 16px;
  --gap: 12px;
  --shadow-offset: 6px;
  --grid-size: 32px;
  --grid-line: rgba(10, 10, 10, 0.08);
  --bg: #f4f1ea;
  --fg: #0a0a0a;
  --muted: #6b6660;
  --line: #0a0a0a;
  --card: #ffffff;
  --accent: #ff5b1f;
  --accent-fg: #0a0a0a;
  --error: #d4341a;
  --good: #1f7a3a;
  --shadow-color: #0a0a0a;
  --surface: var(--card);
  --surface-hi: var(--bg);
  --border: color-mix(in srgb, var(--line) 35%, transparent);
  --border-hi: var(--line);
  --green: var(--accent);
  --green-dim: color-mix(in srgb, var(--accent) 70%, transparent);
  --cyan: var(--accent);
  --text: var(--fg);
  --text-bright: var(--fg);
  --text-dim: var(--muted);
  --mono: "JetBrains Mono", "Fira Code", "Cascadia Code", ui-monospace, Menlo, monospace;
  --display: "Space Grotesk", -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;

  min-height: 100vh;
  background: var(--bg);
  background-image:
    linear-gradient(to right, var(--grid-line) 1px, transparent 1px),
    linear-gradient(to bottom, var(--grid-line) 1px, transparent 1px);
  background-size: var(--grid-size) var(--grid-size);
  color: var(--fg);
  font-family: var(--mono);
  font-size: 14px;
  line-height: 1.5;
  padding: 32px 24px 96px;
  position: relative;
  transition: background 120ms linear, color 120ms linear;
}

.app-root[data-theme="bone"] {
  --bg: #f4f1ea;
  --fg: #0a0a0a;
  --muted: #6b6660;
  --line: #0a0a0a;
  --card: #ffffff;
  --accent: #ff5b1f;
  --accent-fg: #0a0a0a;
  --error: #d4341a;
  --shadow-color: #0a0a0a;
  --grid-line: rgba(10, 10, 10, 0.08);
}

.app-root[data-theme="riot"] {
  --bg: #1f2330;
  --fg: #ece6d6;
  --muted: #8a8676;
  --line: #ece6d6;
  --card: #272c3b;
  --accent: #e84a2a;
  --accent-fg: #ece6d6;
  --error: #ff7a5f;
  --shadow-color: #f4c542;
  --grid-line: rgba(236, 230, 214, 0.08);
}

.app-root[data-density="compact"] { --pad: 10px; --gap: 8px; }
.app-root[data-density="comfy"] { --pad: 16px; --gap: 12px; }
.app-root[data-density="spacious"] { --pad: 24px; --gap: 18px; }
.app-root[data-radius="sharp"] { --radius: 0px; }
.app-root[data-radius="soft"] { --radius: 4px; }
.app-root[data-radius="round"] { --radius: 12px; }

.app-root > * {
  position: relative;
  max-width: 1100px;
  margin-left: auto;
  margin-right: auto;
}

.app-root ::selection {
  background: var(--accent);
  color: var(--accent-fg);
}

.app-root button,
.app-root input,
.app-root textarea,
.app-root pre {
  font-family: inherit;
}

.app-header {
  display: grid;
  grid-template-columns: 1fr minmax(320px, 520px) auto;
  align-items: start;
  gap: 24px;
  margin-bottom: 32px;
  padding-bottom: 24px;
  border-bottom: var(--border-w) solid var(--line);
}

.app-brand {
  display: flex;
  align-items: center;
  gap: 14px;
  color: inherit;
  text-decoration: none;
}

.app-brand-mark {
  width: 56px;
  height: 56px;
  display: grid;
  place-items: center;
  flex-shrink: 0;
  background: var(--fg);
  border: var(--border-w) solid var(--line);
  color: var(--bg);
}

.app-brand-text {
  display: grid;
  gap: 4px;
  font-family: var(--display);
  font-size: 24px;
  font-weight: 700;
  line-height: 1;
}

.app-brand-text small {
  color: var(--muted);
  font-family: var(--mono);
  font-size: 11px;
  font-weight: 500;
  letter-spacing: 0.12em;
  text-transform: uppercase;
}

.app-nav {
  display: flex;
  border: var(--border-w) solid var(--line);
  background: var(--card);
}

.app-nav-tab {
  flex: 1;
  position: relative;
  display: grid;
  gap: 2px;
  min-width: 0;
  padding: 14px 20px;
  border-right: var(--border-w) solid var(--line);
  color: var(--fg);
  font-size: 13px;
  font-weight: 600;
  letter-spacing: 0.08em;
  text-decoration: none;
  text-transform: uppercase;
  transition: background 80ms;
}

.app-nav-tab:last-child {
  border-right: 0;
}

.app-nav-tab:hover {
  background: color-mix(in srgb, var(--accent) 15%, transparent);
}

.app-nav-tab--active {
  background: var(--accent);
  color: var(--accent-fg);
}

.app-nav-num {
  font-size: 10px;
  letter-spacing: 0.15em;
  opacity: 0.65;
}

.app-header-right {
  display: flex;
  align-items: flex-start;
  gap: 8px;
}

.app-pill {
  appearance: none;
  padding: 4px 10px;
  border: var(--border-w) solid var(--line);
  border-radius: var(--radius);
  background: var(--card);
  color: var(--fg);
  cursor: pointer;
  font-size: 11px;
  font-weight: 600;
  letter-spacing: 0.1em;
  line-height: 1.5;
  text-transform: uppercase;
}

.app-pill:hover {
  background: var(--accent);
  color: var(--accent-fg);
}

/* Converter panels */
.app-root .cp-content {
  max-width: 1100px;
}

.app-root .cp-dir-row {
  display: flex;
  justify-content: center;
  margin: 0 0 24px;
}

.app-root .cp-dir-toggle,
.app-root .cp-output-controls,
.app-root .cp-history-tabs {
  display: flex;
  align-items: stretch;
  gap: 0;
  border: var(--border-w) solid var(--line);
  border-radius: var(--radius);
  background: var(--card);
}

.app-root .cp-dir-opt,
.app-root .cp-dir-swap,
.app-root .cp-tab,
.app-root .cp-copy-btn,
.app-root .cp-action-btn,
.app-root .cp-clock-btn,
.app-root .cp-hero-history-btn,
.app-root .cp-history-tab,
.app-root .cp-history-close {
  appearance: none;
  border: 0;
  border-right: var(--border-w) solid var(--line);
  border-radius: var(--radius);
  background: transparent;
  color: var(--fg);
  cursor: pointer;
  font-family: inherit;
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 0.1em;
  line-height: 1.4;
  min-height: 34px;
  padding: 7px 12px;
  text-transform: uppercase;
  transition: background 80ms, color 80ms;
}

.app-root .cp-dir-opt:last-child,
.app-root .cp-tab:last-child,
.app-root .cp-copy-btn:last-child,
.app-root .cp-history-tab:last-child {
  border-right: 0;
}

.app-root .cp-dir-swap {
  color: var(--muted);
  font-size: 18px;
  min-width: 42px;
  padding: 4px 10px;
}

.app-root .cp-dir-opt:hover,
.app-root .cp-dir-swap:hover,
.app-root .cp-tab:hover,
.app-root .cp-copy-btn:hover:not(:disabled),
.app-root .cp-action-btn:hover:not(:disabled),
.app-root .cp-clock-btn:hover,
.app-root .cp-hero-history-btn:hover,
.app-root .cp-history-tab:hover,
.app-root .cp-history-close:hover {
  background: color-mix(in srgb, var(--accent) 15%, transparent);
  color: var(--fg);
}

.app-root .cp-dir-opt--active,
.app-root .cp-tab--active,
.app-root .cp-history-tab--active,
.app-root .cp-copy-btn--ok {
  background: var(--accent);
  color: var(--accent-fg);
}

.app-root .cp-hero {
  margin-bottom: 20px;
  padding: 0;
  border: var(--border-w) solid var(--line);
  border-radius: var(--radius);
  background: var(--card);
  box-shadow: var(--shadow-offset) var(--shadow-offset) 0 var(--shadow-color);
}

.app-root .cp-clock-row {
  display: flex;
  align-items: center;
  gap: 0;
  border-bottom: var(--border-w) solid var(--line);
  background: var(--fg);
  color: var(--bg);
}

.app-root .cp-clock {
  flex: 1;
  min-width: 0;
  padding: 10px 14px;
  font-size: 12px;
  font-weight: 700;
  letter-spacing: 0.12em;
  text-transform: uppercase;
}

.app-root .cp-clock-btn {
  border-left: var(--border-w) solid var(--line);
  border-right: 0;
  color: var(--bg);
  min-width: 42px;
}

.app-root .cp-clock-btn:hover {
  background: var(--accent);
  color: var(--accent-fg);
}

.app-root .cp-hero-input-wrap {
  display: flex;
  align-items: stretch;
  gap: 0;
}

.app-root .cp-hero-input {
  flex: 1;
  min-width: 0;
  border: 0;
  border-radius: 0;
  background: transparent;
  color: var(--fg);
  caret-color: var(--accent);
  font-size: 22px;
  font-weight: 500;
  line-height: 1.4;
  outline: none;
  padding: 18px 16px;
}

.app-root .cp-hero-input::placeholder,
.app-root .cp-textarea::placeholder {
  color: var(--muted);
}

.app-root .cp-hero-copied {
  align-self: center;
  background: var(--accent);
  color: var(--accent-fg);
  font-size: 10px;
  font-weight: 700;
  letter-spacing: 0.08em;
  margin-right: 12px;
  padding: 3px 8px;
  text-transform: uppercase;
  white-space: nowrap;
}

.app-root .cp-hero-history-btn {
  border-left: var(--border-w) solid var(--line);
  border-right: 0;
  border-radius: 0;
  height: auto;
  min-width: 52px;
}

.app-root .cp-hero-output-wrap {
  display: flex;
  align-items: baseline;
  min-height: 58px;
  padding: 0 16px 16px;
}

.app-root .cp-hero-output {
  margin: 0;
  color: var(--accent);
  font-size: 18px;
  font-weight: 700;
  line-height: 1.5;
  white-space: pre-wrap;
  word-break: break-all;
}

.app-root .cp-hero-hint,
.app-root .cp-hint {
  color: var(--muted);
  font-size: 11px;
  letter-spacing: 0.05em;
}

.app-root .cp-multi-section {
  border-top: var(--border-w) solid var(--line);
  padding-top: 24px;
}

.app-root .cp-panels {
  display: grid;
  grid-template-columns: minmax(0, 1fr) auto minmax(0, 1fr);
  gap: 0;
  align-items: stretch;
}

.app-root .cp-panel {
  display: flex;
  min-width: 0;
  flex-direction: column;
  overflow: hidden;
  border: var(--border-w) solid var(--line);
  border-radius: var(--radius);
  background: var(--card);
}

.app-root .cp-arrow {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0 16px;
  color: var(--muted);
  font-size: 28px;
  user-select: none;
}

.app-root .cp-panel-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 10px 14px;
  border-bottom: var(--border-w) solid var(--line);
  background: var(--fg);
  color: var(--bg);
}

.app-root .cp-panel-title {
  color: inherit;
  font-size: 10px;
  font-weight: 700;
  letter-spacing: 0.2em;
  text-transform: uppercase;
}

.app-root .cp-panel-body {
  display: flex;
  flex: 1;
  flex-direction: column;
  padding: 14px;
}

.app-root .cp-panel-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-top: 10px;
}

.app-root .cp-textarea {
  flex: 1;
  width: 100%;
  min-height: 220px;
  padding: 10px 12px;
  resize: vertical;
  border: 1px solid var(--line);
  border-radius: var(--radius);
  background: var(--bg);
  color: var(--fg);
  caret-color: var(--accent);
  font-size: 13px;
  line-height: 1.7;
  outline: none;
}

.app-root .cp-hero-input:focus,
.app-root .cp-textarea:focus {
  box-shadow: inset 0 -4px 0 color-mix(in srgb, var(--accent) 35%, transparent);
}

.app-root .cp-output-state {
  display: flex;
  flex: 1;
  align-items: center;
  justify-content: center;
  min-height: 120px;
}

.app-root .cp-output {
  flex: 1;
  width: 100%;
  margin: 0;
  padding: 0;
  overflow-y: auto;
  color: var(--fg);
  font-size: 13px;
  line-height: 1.7;
  white-space: pre-wrap;
  word-break: break-all;
}

.app-root .cp-output--has-errors,
.app-root .cp-error-msg,
.app-root .cp-hero-hint[style] {
  color: var(--error) !important;
}

.app-root .cp-action-btn {
  flex-shrink: 0;
  border: var(--border-w) solid var(--line);
  background: var(--accent);
  color: var(--accent-fg);
}

.app-root .cp-copy-btn:disabled,
.app-root .cp-action-btn:disabled {
  cursor: not-allowed;
  opacity: 0.45;
}

.app-root .cp-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
  margin-top: 64px;
  padding-top: 16px;
  border-top: var(--border-w) solid var(--line);
}

.app-root .cp-footer-links {
  display: flex;
  align-items: center;
  gap: 8px;
}

.app-root .cp-footer-link {
  color: var(--muted);
  font-size: 11px;
  letter-spacing: 0.05em;
  text-decoration: none;
}

.app-root .cp-footer-link:hover {
  color: var(--fg);
  text-decoration: underline;
  text-underline-offset: 3px;
}

.app-root .cp-toast {
  position: fixed;
  right: auto;
  bottom: 24px;
  left: 50%;
  z-index: 100;
  transform: translateX(-50%);
  padding: 10px 18px;
  border: var(--border-w) solid var(--line);
  border-radius: var(--radius);
  background: var(--fg);
  box-shadow: 4px 4px 0 var(--accent);
  color: var(--bg);
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 0.12em;
  pointer-events: none;
  text-transform: uppercase;
}

.app-root .cp-history-overlay {
  position: fixed;
  inset: 0;
  z-index: 100;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 16px;
  background: color-mix(in srgb, var(--bg) 78%, #000000);
}

.app-root .cp-history-panel {
  display: flex;
  width: min(620px, calc(100vw - 32px));
  max-height: 80vh;
  flex-direction: column;
  overflow: hidden;
  border: var(--border-w) solid var(--line);
  border-radius: var(--radius);
  background: var(--card);
  box-shadow: var(--shadow-offset) var(--shadow-offset) 0 var(--shadow-color);
}

.app-root .cp-history-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 12px 16px;
  border-bottom: var(--border-w) solid var(--line);
  background: var(--fg);
}

.app-root .cp-history-close {
  width: 34px;
  padding: 0;
  border: var(--border-w) solid var(--line);
  background: var(--card);
}

.app-root .cp-history-list {
  display: grid;
  gap: 6px;
  min-height: 80px;
  overflow-y: auto;
  padding: 8px;
}

.app-root .cp-history-empty {
  display: grid;
  place-items: center;
  min-height: 80px;
  padding: 16px;
}

.app-root .cp-history-item {
  display: grid;
  grid-template-columns: minmax(0, 1fr) auto minmax(0, 1fr);
  align-items: center;
  gap: 10px;
  min-width: 0;
  padding: 10px 12px;
  border: 1px solid var(--line);
  border-radius: var(--radius);
  background: var(--card);
  cursor: pointer;
  font-size: 13px;
}

.app-root .cp-history-item:hover {
  background: color-mix(in srgb, var(--accent) 15%, var(--card));
}

.app-root .cp-history-in,
.app-root .cp-history-out {
  min-width: 0;
  overflow: hidden;
  color: var(--fg);
  text-overflow: ellipsis;
  white-space: nowrap;
}

.app-root .cp-history-out {
  color: var(--accent);
  text-align: right;
}

.app-root .cp-history-arrow {
  color: var(--muted);
}

.app-root .cp-blink {
  animation: cp-blink 1s step-end infinite;
}

@keyframes cp-blink {
  0%, 100% { opacity: 1; }
  50% { opacity: 0; }
}

.toast-enter-active,
.toast-leave-active,
.overlay-enter-active,
.overlay-leave-active {
  transition: opacity 0.2s, transform 0.2s;
}

.toast-enter-from,
.toast-leave-to,
.overlay-enter-from,
.overlay-leave-to {
  opacity: 0;
}

@media (max-width: 900px) {
  .app-header {
    grid-template-columns: 1fr auto;
  }

  .app-nav {
    grid-column: 1 / -1;
    grid-row: 2;
  }
}

@media (max-width: 720px) {
  .app-root {
    padding: 20px 14px 80px;
  }

  .app-header {
    grid-template-columns: 1fr;
    gap: 16px;
  }

  .app-header-right {
    align-items: stretch;
  }

  .app-brand-text {
    font-size: 21px;
  }

  .app-root .cp-panels {
    grid-template-columns: 1fr;
  }

  .app-root .cp-arrow {
    padding: 8px 0;
    rotate: 90deg;
  }

  .app-root .cp-hero-input-wrap {
    flex-wrap: wrap;
  }

  .app-root .cp-hero-input {
    flex-basis: 100%;
    font-size: 18px;
  }

  .app-root .cp-hero-copied {
    margin: 0 0 12px 16px;
  }

  .app-root .cp-hero-history-btn {
    flex: 1;
    border-top: var(--border-w) solid var(--line);
    border-left: 0;
    min-height: 40px;
  }

  .app-root .cp-panel-header,
  .app-root .cp-panel-footer,
  .app-root .cp-footer {
    align-items: stretch;
    flex-direction: column;
  }

  .app-root .cp-output-controls,
  .app-root .cp-history-tabs {
    width: 100%;
  }

  .app-root .cp-tab,
  .app-root .cp-copy-btn,
  .app-root .cp-history-tab {
    flex: 1;
  }

  .app-root .cp-history-item {
    grid-template-columns: 1fr;
    gap: 4px;
  }

  .app-root .cp-history-out {
    text-align: left;
  }
}
</style>
