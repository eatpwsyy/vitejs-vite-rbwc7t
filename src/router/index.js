// src/router/index.js
import { createRouter, createWebHistory } from 'vue-router';
import DailyOverview from '@/components/Daily/DailyOverview.vue';
import DailyTrendingTopic from '@/components/Daily/DailyTrendingTopic.vue';
import JSONFileMaker from '@/components/Maker/JSONFileMaker.vue';

const routes = [
  {
    path: '/',
    name: 'DailyOverview',
    component: DailyOverview,
  },
  {
    path: '/daily-trending-topic',
    name: 'DailyTrendingTopic',
    component: DailyTrendingTopic,
  },
  {
    path: '/fw-data-request',
    name: 'JSONFileMaker',
    component: JSONFileMaker,
  }
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;