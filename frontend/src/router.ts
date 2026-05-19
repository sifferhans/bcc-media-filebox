import { createRouter, createWebHistory } from 'vue-router'

const Home = () => import('./views/Home.vue')
const Admin = () => import('./views/Admin.vue')

export const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', name: 'home', component: Home },
    { path: '/admin', name: 'admin', component: Admin },
  ],
})
