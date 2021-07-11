import { createRouter, createWebHashHistory, RouteRecordRaw } from 'vue-router'
import DmAutomation from '../views/DmAutomation.vue'
import Settings from '../views/Settings.vue'
import Logs from '../views/Logs.vue'
import About from '../views/About.vue'

const routes: Array<RouteRecordRaw> = [
  {
    path: '/',
    name: 'DmAutomation',
    component: DmAutomation
  },
  {
    path: '/settings',
    name: 'Settings',
    component: Settings
  },
  {
    path: '/logs',
    name: 'Logs',
    component: Logs
  },
  {
    path: '/about',
    name: 'About',
    component: About
  }
]

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

export default router
