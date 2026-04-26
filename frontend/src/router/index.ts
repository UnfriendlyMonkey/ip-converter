import { createRouter, createWebHistory } from 'vue-router'
import IpConverterPanel from '../components/IpConverterPanel.vue'
import TimestampPanel from '../components/TimestampPanel.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', redirect: '/ip' },
    { path: '/ip', component: IpConverterPanel },
    { path: '/timestamp', component: TimestampPanel },
  ],
})

export default router
