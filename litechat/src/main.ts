import { createApp } from 'vue'
import { createPinia } from 'pinia'
import './style.css'
import App from './App.vue'
import router from './router'
import { useThemeStore } from './stores/theme'

const pinia = createPinia()
const app = createApp(App)

app.use(pinia)
app.use(router)
app.mount('#app')

// Initialize theme store to apply theme on load
useThemeStore()
