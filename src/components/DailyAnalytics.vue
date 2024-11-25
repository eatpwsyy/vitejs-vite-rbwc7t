<template>
  <div class="analytics-container">
    <!-- Yesterday's Date Display -->
    <div class="current-date">
      Report for <br />
      {{ yesterdayDisplayFormat }}
    </div>
    <!-- Cards: Articles and Sources -->
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
      <h3 class="card-title">Organization, Place, and People</h3>
      <div class="card-content">
        <ol class="numbered-list">
          <li v-for="(entity, key) in nerData" :key="key">
            {{ key }} <!-- Entity Type -->
            <ol type="a">
              <li v-for="(details, name) in entity" :key="name">
                {{ name }} - Count: {{ details.count }},
                Positive: {{ details.sentiment.positive }},
                Negative: {{ details.sentiment.negative }}
              </li>
            </ol>
          </li>
        </ol>
      </div>
    </div>
  </div>
</template>

<script>
import Chart from 'chart.js/auto';
import ChartDataLabels from 'chartjs-plugin-datalabels';

export default {
  name: 'DailyAnalytics',
  data() {
    return {
      totalArticles: 0, // Initialize to 0
      totalSources: 0, // Initialize to 0
      yesterdayAPIFormat: this.getYesterdayDate('YYYYMMDD'), // Date for API request
      yesterdayDisplayFormat: this.getYesterdayDate('readable'), // Date for display
      trendingData: [], // Placeholder for trending topics
      nerData: {}, // Placeholder for NER data
    };
  },
  methods: {
    getYesterdayDate(format) {
      const yesterday = new Date();
      yesterday.setDate(yesterday.getDate() - 1); // Subtract one day

      if (format === 'YYYYMMDD') {
        // API format: YYYYMMDD
        const year = yesterday.getFullYear();
        const month = String(yesterday.getMonth() + 1).padStart(2, '0'); // Leading zero for month
        const day = String(yesterday.getDate()).padStart(2, '0'); // Leading zero for day
        return `${year}${month}${day}`;
      }

      if (format === 'readable') {
        // Readable format: e.g., "Wednesday, 20 November 2024"
        const options = { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' };
        return yesterday.toLocaleDateString(undefined, options);
      }

      return '';
    },
    async fetchData() {
      try {
        console.log('Fetching API data...');
        const response = await fetch(`/api/get-data?req_id=xxx&req_type=oview&freq=daily&date=${this.yesterdayAPIFormat}`);
        console.log('API Response Status:', response.status);

        if (!response.ok) {
          throw new Error(`HTTP error! Status: ${response.status}`);
        }

        const data = await response.json();
        console.log('Parsed API Data:', data);

        this.totalArticles = data.articles;
        this.totalSources = data.sources;
        this.trendingData = data.trending;
        this.nerData = data.OPP;

        this.generatePieChart();
      } catch (error) {
        console.error('Error in fetchData:', error);
      }
    },
    generatePieChart() {
      const ctx = document.getElementById('pieChart');
      if (!ctx) {
        console.warn('Canvas not found for pie chart');
        return;
      }
      const chartCtx = ctx.getContext('2d');
      Chart.register(ChartDataLabels);

      const labels = this.trendingData.map((item) => item.topic);
      const dataValues = this.trendingData.map((item) => item.percentage);

      new Chart(chartCtx, {
        type: 'pie',
        data: {
          labels,
          datasets: [{ data: dataValues, backgroundColor: ['#FF6384', '#36A2EB'] }],
        },
        options: {
          responsive: true,
          maintainAspectRatio: true,
          plugins: {
            legend: { position: 'bottom' },
            datalabels: { color: '#fff', formatter: (value) => `${value}%` },
          },
        },
      });
    },
  },
  mounted() {
    console.log('Component mounted');
    this.fetchData(); // Fetch data when the component is mounted
    this.totalArticles = 920;
    this.totalSources = 69;
  },
};
</script>

<style>
/* General body styling */
body {
  font-family: 'Roboto', sans-serif;
  color: #333;
  background-color: #ffffff;
  margin: 0;
  padding: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  box-sizing: border-box;
}

/* Analytics container */
.analytics-container {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 30px;
  width: 100%;
  max-width: 1200px; /* Prevent stretching too wide */
  min-width: 60%; /* Ensure it takes more space on large screens */
  margin: 0 auto;
  padding: 20px;
  box-sizing: border-box;
  background-color: #f5f5f5;
  border-radius: 8px;
  box-shadow: 0px 4px 6px rgba(0, 0, 0, 0.1);
  align-items: center;
  justify-content: center;
}

/* Center content on desktop screens */
@media (min-width: 1024px) {
  .analytics-container {
    grid-template-columns: repeat(3, 1fr); /* Use three equal columns */
    justify-items: center; /* Center items within their grid cells */
    gap: 40px;
    max-width: 85%; /* Make sure the content uses more horizontal space */
  }
}

/* Current date styling */
.current-date {
  grid-column: span 3; /* Center across all columns */
  font-size: 2.2em;
  font-weight: bold;
  color: #444;
  text-align: center;
  margin-bottom: 30px;
  line-height: 1.5;
}

/* Card styling */
.card {
  background: linear-gradient(to bottom, #ffffff, #f9fff0);
  text-align: center;
  padding: 30px;
  border-radius: 16px;
  box-shadow: 0px 4px 6px rgba(0, 0, 0, 0.1);
  transition: transform 0.3s ease-in-out, box-shadow 0.3s ease-in-out;
}

.card:hover {
  transform: scale(1.05); /* Slight scale-up on hover */
  box-shadow: 0px 8px 14px rgba(0, 0, 0, 0.2); /* Enhanced shadow on hover */
}

/* Card content styling */
.number {
  font-size: 4em;
  font-weight: bold;
  color: #333;
  margin-bottom: 10px;
}

.label {
  font-size: 1.2em;
  color: #555;
}

/* Pie chart card and nested list adjustments */
.pie-chart-card,
.nested-list-card {
  grid-column: span 2; /* Spans two columns for larger screens */
}

.chart-container {
  width: 100%;
  max-width: 800px; /* Wider chart on larger screens */
  height: 400px; /* Slightly taller for balance */
  margin: 0 auto;
}

/* Responsive adjustments for smaller screens */
@media (max-width: 768px) {
  .analytics-container {
    grid-template-columns: 1fr; /* Stack content vertically */
  }

  .current-date {
    grid-column: span 1;
  }

  .pie-chart-card,
  .nested-list-card {
    grid-column: span 1;
  }
}
</style>
