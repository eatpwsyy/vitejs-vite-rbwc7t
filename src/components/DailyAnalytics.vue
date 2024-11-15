<template>
  <div class="analytics-container">
    <!-- Current Date Display -->
    <div class="current-date">
      {{ currentDate }}
    </div>
    <!-- Cards Container for Articles and Sources -->
    <div class="cards-container">
      <div class="card articles-card">
        <div class="card-content">
          <span class="number">{{ totalArticles }}</span>
          <span class="label">Articles</span>
        </div>
      </div>
      <div class="card sources-card">
        <div class="card-content">
          <span class="number">{{ totalSources }}</span>
          <span class="label">Sources</span>
        </div>
      </div>
    </div>
    <!-- New Cards Container for Pie Chart and Nested List -->
    <div class="cards-container">
      <!-- Pie Chart Card -->
      <div class="card pie-chart-card">
        <h3 class="card-title">Trending Topic</h3>
        <div class="card-content">
          <div class="chart-container">
            <canvas id="pieChart"></canvas>
          </div>
        </div>
      </div>
      <!-- Nested List Card -->
      <div class="card nested-list-card">
        <h3 class="card-title">Organization, Place and People</h3>
        <div class="card-content">
          <ol class="numbered-list">
            <li>List item one
              <ol type="a">
                <li>Sub-item one</li>
                <li>Sub-item two</li>
                <li>Sub-item three</li>
              </ol>
            </li>
            <li>List item two
              <ol type="a">
                <li>Sub-item one</li>
                <li>Sub-item two</li>
                <li>Sub-item three</li>
              </ol>
            </li>
            <li>List item three
              <ol type="a">
                <li>Sub-item one</li>
                <li>Sub-item two</li>
                <li>Sub-item three</li>
              </ol>
            </li>
          </ol>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { onMounted } from 'vue';
import Chart from 'chart.js/auto';
import ChartDataLabels from 'chartjs-plugin-datalabels'; // Import Data Labels Plugin

export default {
  name: 'DailyAnalytics',
  data() {
    return {
      totalArticles: 2078,
      totalSources: 164,
      currentDate: this.formatDate(new Date()),
    };
  },
  methods: {
    formatDate(date) {
      // Format date as "Thursday, 14 November 2024"
      const options = { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' };
      return date.toLocaleDateString(undefined, options);
    },
    generatePieChart() {
      const ctx = document.getElementById('pieChart').getContext('2d');
      Chart.register(ChartDataLabels); // Register Data Labels Plugin

      new Chart(ctx, {
        type: 'pie',
        data: {
          labels: ['Category A', 'Category B', 'Category C'],
          datasets: [
            {
              data: [40, 30, 30], // Sample data percentages
              backgroundColor: ['#FF6384', '#36A2EB', '#FFCE56'],
            },
          ],
        },
        options: {
          responsive: true,
          maintainAspectRatio: true,
          plugins: {
            legend: {
              position: 'bottom', // Ensures legend stays within the card
              labels: {
                boxWidth: 15, // Adjust legend box size for better fit
              },
            },
            datalabels: {
              color: '#fff', // Color of the label text
              formatter: (value, context) => {
                const total = context.chart.data.datasets[0].data.reduce((sum, val) => sum + val, 0);
                const percentage = ((value / total) * 100).toFixed(1) + '%';
                return percentage; // Show percentage
              },
              font: {
                weight: 'bold',
                size: 14,
              },
            },
          },
        },
      });
    },
  },
  mounted() {
    this.generatePieChart();
  },
};
</script>

<style scoped>
.analytics-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 20px;
  margin: 20px;
  background-color: #e0e0e0; /* Light grey background to enhance card visibility */
  padding: 20px;
  border-radius: 8px;
}

.current-date {
  font-size: 1.5em;
  font-weight: bold;
  color: #333; /* Dark color for the date text */
  margin-bottom: 20px;
}

.cards-container {
  display: flex;
  justify-content: center; /* Ensure the cards are centered */
  gap: 20px; /* Space between the cards */
}

.card {
  flex: 1;
  background-color: #f0ffd9; /* Light green background */
  text-align: center;
  padding: 20px; /* Adjusted padding for better spacing */
  border-radius: 10px;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
  min-width: 200px; /* Ensure card width is enough to fit the content */
  max-width: 350px; /* Consistent max-width for all cards */
}

.card-title {
  font-size: 1.5em;
  font-weight: bold;
  color: #333;
  margin-bottom: 15px;
}

.card-content {
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center; /* Ensures content is centered */
  gap: 10px; /* Added gap to prevent overlap between number and label */
}

.chart-container {
  width: 100%;
  height: 300px; /* Adjust height to ensure space for chart and legend */
  position: relative;
}

.pie-chart-card {
  height: 400px; /* Height to accommodate both the pie chart and the legend */
}

.nested-list-card ol {
  text-align: left; /* Align the list to the left for better readability */
  padding-left: 20px;
}

.nested-list-card .numbered-list {
  list-style-type: decimal; /* Remove bullet points for the top-level list */
}

.number {
  font-size: 6.5em; /* Adjusted font size to balance the card */
  font-weight: bold;
  color: #333; /* Dark color for good contrast */
  line-height: 1; /* Tighter line-height to keep content compact */
}

.label {
  font-size: 1.5em; /* Increased the size slightly for better readability */
  font-weight: normal;
  color: #333; /* Dark color for the label as well */
}
</style>