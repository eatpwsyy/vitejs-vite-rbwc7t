import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import path from 'path';

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, 'src'), // Set the `@` alias to the `src` directory
    },
  },
  optimizeDeps: {
    include: ['chart.js', 'chartjs-plugin-datalabels'],
  },
});