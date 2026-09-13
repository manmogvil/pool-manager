import { createApp } from 'vue'
import { createVuetify } from 'vuetify'
import 'vuetify/styles'
import '@mdi/font/css/materialdesignicons.css'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'

import App from './App.vue'
import router from './router'
import './style.css'

const vuetify = createVuetify({
  components,
  directives,
  theme: {
    defaultTheme: 'light',
    themes: {
      light: {
        colors: {
          background: '#f8f9fa',
          surface: '#ffffff',
          primary: '#3b82f6',
          'primary-darken-1': '#2563eb',
          secondary: '#666666',
          error: '#ef4444',
          info: '#3b82f6',
          success: '#22c55e',
          warning: '#f59e0b'
        }
      }
    }
  }
})

const app = createApp(App)

app.use(vuetify)
app.use(router)

app.mount('#app')
