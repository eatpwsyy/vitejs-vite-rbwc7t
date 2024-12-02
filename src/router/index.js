// src/router/index.js
import { createRouter, createWebHistory } from 'vue-router';
import Overview from '@/components/Overview.vue';
import TrendingTopic from '@/components/TrendingTopic.vue';

const routes = [
  {
    path: '/',
    name: 'Overview',
    component: Overview,
  },
  {
    path: '/trending-topic',
    name: 'TrendingTopic',
    component: TrendingTopic,
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;