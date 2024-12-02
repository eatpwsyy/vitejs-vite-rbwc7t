<template>
  <div class="container-fluid">
      <div class="row">
      <header class="col-12 mb-4">
        <h1>News Analytics Dashboard</h1>
        <p>Analytics summary and insights for news articles and sources.</p>
      </header>
    </div>
  </div>
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
    <router-link to="/trending-topic" class="card pie-chart-card">
      <h3 class="card-title">Trending Topic</h3>
       <div class="card-content">
         <div class="chart-container">
           <canvas id="pieChart"></canvas>
         </div>
       </div>
    </router-link>
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
  name: 'Overview',
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
    
        const topic = this.trendingData.map((item) => item.topic);
        const keywords = this.trendingData.map((item) => item.keywords);
        const dataValues = this.trendingData.map((item) => parseFloat(item.percentage.toFixed(2)));
    
        new Chart(chartCtx, {
            type: 'pie',
            data: {
                labels: topic, // Use topic for the default chart labels
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
                            // Customize the legend labels to display `topic`
                            generateLabels: (chart) => {
                                return chart.data.labels.map((label, index) => {
                                    return {
                                        text: label, // Use topic as legend text
                                        fillStyle: chart.data.datasets[0].backgroundColor[index], // Set the corresponding color
                                    };
                                });
                            },
                        },
                    },
                    tooltip: {
                        callbacks: {
                            label: (context) => {
                                // Show the `keywords` in the tooltip instead of `topic`
                                const keyword = keywords[context.dataIndex];
                                return `${keyword}: ${context.raw}%`; // Use `keywords` here
                            },
                        },
                    },
                    datalabels: {
                        color: '#fff',
                        formatter: (value, context) => {
                            if (value > 5) { // Only display data labels for larger slices
                                return `${context.chart.data.labels[context.dataIndex]}: ${value}%`;
                            } else {
                                return ''; // Avoid displaying labels for very small segments
                            }
                        },
                    },
                },
            },
        });
    },
  },
  mounted() {
    console.log('Component mounted');
    this.fetchData(); // Fetch data when the component is mounted
  },
};
</script>

<style scoped>
header {
  text-align: center;
  margin-bottom: 30px;
  background-color: #f8f9fa;  /* Light background */
  padding: 20px;
  border-radius: 8px;
}

h1 {
  font-size: 2.5rem;
  color: #343a40;  /* Darker text color */
}

p {
  font-size: 1.2rem;
  color: #6c757d;  /* Light gray text color */
}

.analytics-container {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 30px;
  width: 100%;
  max-width: 1200px;
  margin: 0 auto;
  min-height: 100vh;
  height: 400px;
  padding: 20px;
  box-sizing: border-box;
  background-color: #f5f5f5;
  border-radius: 8px;
  box-shadow: 0px 4px 6px rgba(0, 0, 0, 0.1);
}

.current-date {
  grid-column: span 2; /* Span two columns for larger screens */
  font-size: 1.8em;
  font-weight: bold;
  color: #444;
  text-align: center;
  margin-bottom: 20px;
}

.cards-container {
  display: contents; /* Remove redundant grid inside the container */
}

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

.articles-card,
.sources-card {
  flex: 1 1 300px;
}

.pie-chart-card {
  grid-column: span 2; /* Make the pie chart span two columns */
}

.nested-list-card {
  grid-column: span 2; /* Make the nested list span two columns */
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
  max-width: 600px; /* Allow a larger pie chart */
  height: 350px;
  margin: 0 auto;
}

@media (max-width: 768px) {
  .analytics-container {
    grid-template-columns: 1fr; /* Stack everything vertically on smaller screens */
  }

  .current-date {
    grid-column: span 1; /* Adjust date display */
  }

  .pie-chart-card,
  .nested-list-card {
    grid-column: span 1; /* Make large cards fit single-column */
  }
}
</style>