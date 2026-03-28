import { createApp } from 'vue'
import { createVuetify } from 'vuetify'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'
import { aliases, mdi } from 'vuetify/iconsets/mdi'
import 'vuetify/styles'
import '@mdi/font/css/materialdesignicons.css'
import App from './App.vue'

const vuetify = createVuetify({
  components,
  directives,
  icons: {
    defaultSet: 'mdi',
    aliases,
    sets: { mdi },
  },
  theme: {
    defaultTheme: 'hacker',
    themes: {
      hacker: {
        dark: true,
        colors: {
          background: '#0a0a0f',
          surface: '#0d0d14',
          'surface-bright': '#12121c',
          primary: '#00ff88',
          secondary: '#0088ff',
          error: '#ff4466',
          warning: '#ffaa00',
          info: '#00aaff',
          success: '#00ff88',
        },
      },
    },
  },
})

createApp(App).use(vuetify).mount('#app')
