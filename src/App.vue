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
    </div>
  </div>
</template>

<script>
import DailyAnalytics from './components/DailyAnalytics.vue';

export default {
  name: 'App',
  components: {
    DailyAnalytics
  }
}
</script>

<style scoped>
header {
  text-align: center;
  margin-bottom: 30px;
}
</style>
