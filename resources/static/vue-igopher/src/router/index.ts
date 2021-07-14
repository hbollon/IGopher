import { createRouter, createWebHashHistory, RouteRecordRaw } from 'vue-router'
import DmAutomation from '../views/DmAutomation.vue'
import Settings from '../views/Settings.vue'
import Logs from '../views/Logs.vue'
import About from '../views/About.vue'
import NotFound from '../views/NotFound.vue'

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
  },
  {
    path: "/:pathMatch(.*)*",
    name: '404 Not Found',
    component: NotFound
  }
]

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

export default router
