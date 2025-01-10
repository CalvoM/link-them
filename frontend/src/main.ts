import './assets/main.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import router from './router'
import PrimeVue from 'primevue/config'
import { InputText, Button, FloatLabel, AutoComplete, Message } from 'primevue'

const app = createApp(App)

app.component('InputText', InputText)
app.component('Button', Button)
app.component('FloatLabel', FloatLabel)
app.component('AutoComplete', AutoComplete)
app.component('Message', Message)

app.use(createPinia())
app.use(router)
app.use(PrimeVue, { theme: 'none' })

app.mount('#app')
