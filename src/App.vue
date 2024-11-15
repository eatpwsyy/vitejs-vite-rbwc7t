<template>
  <div class="container-fluid">
    <div class="row">
      <header class="col-12 mb-4">
        <h1>News Analytics Dashboard</h1>
        <p>Analytics summary and insights for news articles and sources.</p>
      </header>
    </div>
    <div class="row">
      <!-- Summary Cards -->
      <div class="col-md-4 mb-4">
        <DailyAnalytics />
      </div>
      <div class="col-md-4 mb-4">
        <WeeklyAnalytics />
      </div>
      <div class="col-md-4 mb-4">
        <MonthlyAnalytics />
      </div>
    </div>
    <div class="row">
      <!-- Graphical Visualization Section -->
      <div class="col-12 mb-4">
        <Heatmap />
      </div>
    </div>
  </div>
</template>

<script>
import DailyAnalytics from './components/DailyAnalytics.vue';
import Heatmap from './components/Heatmap.vue';

export default {
  name: 'App',
  components: {
    DailyAnalytics,
    Heatmap
  }
}
</script>

<style scoped>
header {
  text-align: center;
  margin-bottom: 30px;
}
</style>
