import { createApp } from 'vue'
import './asset/css/main.css'

import Sub from './view/sub/Sub.vue'

const app = createApp(Sub)

app.mount('#sub')