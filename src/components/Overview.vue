<template>
  <div class="main-container">
    <!-- Header Section -->
    <div class="container-fluid">
      <div class="row">
        <header class="col-12 mb-4">
          <h1>News Analytics Dashboard</h1>
          <p>Analytics summary and insights for news articles and sources.</p>
        </header>
      </div>
    </div>

    <!-- Analytics Sections -->
    <div class="analytics-container">
      <!-- 1. Report for Yesterday's Date -->
      <div class="current-date">
        Report for <br />
        {{ yesterdayDisplayFormat }}
      </div>

      <!-- 2. Articles and Sources Cards -->
      <div class="cards-row">
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

      <!-- 3. Pie Chart Card -->
      <router-link to="/trending-topic" class="card pie-chart-card">
        <h3 class="card-title">Trending Topic</h3>
        <div class="card-content">
          <div class="chart-container">
            <canvas ref="pieChart"></canvas>
          </div>
        </div>
      </router-link>

      <!-- 4. OPP Card -->
      <div class="card nested-list-card">
        <h3 class="card-title">Organization, Place, and People</h3>
        <div class="card-content">
          <ol class="numbered-list">
            <li v-for="(entity, key) in nerData" :key="key">
              {{ key }}
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
  </div>
</template>

<script>
import Chart from 'chart.js/auto';
import ChartDataLabels from 'chartjs-plugin-datalabels';

export default {
  name: 'Overview',
  data() {
    return {
      totalArticles: 0,
      totalSources: 0,
      yesterdayAPIFormat: this.getYesterdayDate('YYYYMMDD'),
      yesterdayDisplayFormat: this.getYesterdayDate('readable'),
      trendingData: [],
      nerData: {},
      pieChartInstance: null, // To keep track of the chart instance
    };
  },
  methods: {
    getYesterdayDate(format) {
      const yesterday = new Date();
      yesterday.setDate(yesterday.getDate() - 1);

      if (format === 'YYYYMMDD') {
        const year = yesterday.getFullYear();
        const month = String(yesterday.getMonth() + 1).padStart(2, '0');
        const day = String(yesterday.getDate()).padStart(2, '0');
        return `${year}${month}${day}`;
      }

      if (format === 'readable') {
        const options = { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' };
        return yesterday.toLocaleDateString(undefined, options);
      }

      return '';
    },

    async fetchData() {
      try {
        console.log('Fetching API data...');
        const response = await fetch(
          `/api/get-data?req_id=xxx&req_type=oview&freq=daily&date=${this.yesterdayAPIFormat}`
        );
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
      if (this.pieChartInstance) {
        this.pieChartInstance.destroy();
      }

      const ctx = this.$refs.pieChart;
      if (!ctx) {
        console.warn('Canvas not found for pie chart');
        return;
      }

      Chart.register(ChartDataLabels);

      const topics = this.trendingData.map((item) => item.topic);
      const keywords = this.trendingData.map((item) => item.keywords);
      const dataValues = this.trendingData.map((item) => parseFloat(item.percentage.toFixed(2)));

      this.pieChartInstance = new Chart(ctx, {
        type: 'pie',
        data: {
          labels: topics,
          datasets: [
            {
              data: dataValues,
              backgroundColor: [
                '#FF6384', '#36A2EB', '#FFCE56', '#4BC0C0', '#9966FF', '#FF9F40',
                '#61FF84', '#3642EB', '#CE5645', '#40A0C0', '#2296FF', '#9F40FF'
              ],
            },
          ],
        },
        options: {
          responsive: true,
          maintainAspectRatio: false,
          plugins: {
            legend: {
              position: 'bottom',
              labels: {
                generateLabels: (chart) => {
                  return chart.data.labels.map((label, index) => ({
                    text: label,
                    fillStyle: chart.data.datasets[0].backgroundColor[index],
                  }));
                },
              },
            },
            tooltip: {
              callbacks: {
                label: (context) => {
                  const keyword = keywords[context.dataIndex];
                  return `${keyword}: ${context.raw}%`;
                },
              },
            },
            datalabels: {
              color: '#fff',
              formatter: (value, context) => {
                return value > 5
                  ? `${context.chart.data.labels[context.dataIndex]}: ${value}%`
                  : '';
              },
            },
          },
        },
      });
    },
  },
  mounted() {
    console.log('Component mounted');
    this.fetchData();
  },
};
</script>

<style scoped>
/* Header Styles */
header {
  text-align: center;
  margin-bottom: 30px;
  padding: 20px;
  border-radius: 8px;
}

h1 {
  font-size: 2.5rem;
  color: #343a40; /* Darker text color */
}

p {
  font-size: 1.2rem;
  color: #6c757d; /* Light gray text color */
}

/* Main Container */
.main-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 20px;
}

/* Analytics Container */
.analytics-container {
  display: flex;
  flex-direction: column;
  gap: 30px;
  width: 100%;
  max-width: 1200px;
  padding: 20px;
  box-sizing: border-box;
  background-color: #f5f5f5;
  border-radius: 8px;
  box-shadow: 0px 4px 6px rgba(0, 0, 0, 0.1);
}

/* 1. Current Date */
.current-date {
  font-size: 1.8em;
  font-weight: bold;
  color: #444;
  text-align: center;
}

/* 2. Articles and Sources Cards Row */
.cards-row {
  display: flex;
  gap: 20px;
}

@media (max-width: 768px) {
  .cards-row {
    flex-direction: column;
  }
}

/* Cards Styles */
.card {
  background-color: #f9fff0;
  text-align: center;
  padding: 20px;
  border-radius: 12px;
  box-shadow: 0px 4px 6px rgba(0, 0, 0, 0.1);
  transition: transform 0.2s ease-in-out;
}

.card:hover {
  transform: scale(1.02);
}

/* Specific Card Styles */
.articles-card,
.sources-card {
  flex: 1;
}

.pie-chart-card,
.nested-list-card {
  width: 100%;
}

.number {
  font-size: 3em;
  font-weight: bold;
  color: #333;
}

.label {
  font-size: 1.2em;
  color: #555;
}

.chart-container {
  width: 100%;
  max-width: 600px;
  height: 350px;
  margin: 0 auto;
}

/* 4. Nested List Card */
.nested-list-card .numbered-list {
  padding-left: 20px;
}

.nested-list-card .numbered-list li {
  margin-bottom: 10px;
}

/* Responsive Adjustments */
@media (max-width: 768px) {
  .analytics-container {
    padding: 10px;
  }

  .current-date {
    font-size: 1.5em;
  }

  .number {
    font-size: 2.5em;
  }

  .label {
    font-size: 1em;
  }

  .chart-container {
    height: 300px;
  }
}
</style>