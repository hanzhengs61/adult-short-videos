import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { createHead } from '@vueuse/head'
import router from './router'
import App from './App.vue'
import './assets/main.css'

createApp(App).use(createPinia()).use(router).use(createHead()).mount('#app')
