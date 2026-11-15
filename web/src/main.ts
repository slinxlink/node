import { createApp } from 'vue'
import './asset/css/index.css'

import App from './App.vue'
import router from './router'

createApp(App).use(router).mount('#app')
