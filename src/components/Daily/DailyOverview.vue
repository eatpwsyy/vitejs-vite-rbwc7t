<template>
  <div class="app-container">
    <!-- Sidebar -->
    <nav class="sidebar" :class="{ 'open': isSidebarOpen }">
      <!-- Close button inside sidebar -->
      <button @click="toggleSidebar" class="close-btn">×</button>
      <br /> <br /> <br /> <br /> <br />
      <ul class="sidebar-menu">

        <!-- DAILY OVERVIEW + nested links -->
        <li>
          <!-- Parent link (toggles chilc links) -->
          <div class="sidebar-link" @click="toggleDailyOverview">
            Daily Analytics
            <span class="arrow" :class="{ 'rotated': dailyOverviewOpen }">▸</span>
          </div>

          <!-- Child links (shown if dailyOverviewOpen is True) -->
          <ul v-show="dailyOverviewOpen" class="nested-links">
            <li>
              <router-link to="/">Daily Overview</router-link>
            </li>            
            <li>
              <router-link to="/daily-trending-topic">Daily Trending Topic</router-link>
            </li>
          </ul>
        </li>

        <!-- WEEKLY OVERVIEW + nested links -->
        <li>
          <!-- Parent link (toggles chilc links) -->
          <div class="sidebar-link" @click="toggleWeeklyOverview">
            Weekly Analytics
            <span class="arrow" :class="{ 'rotated': weeklyOverviewOpen }">▸</span>
          </div>

          <!-- Child links (shown if weeklyOverviewOpen is True) 
          <ul v-show="weeklyOverviewOpen" class="nested-links">
            <li>
              <router-link to="/weekly-overview">Weekly Overview</router-link>
            </li>
            <li>
              <router-link to="/weekly-trending-topic">Weekly Trending Topic</router-link>
            </li>
          </ul>-->
        </li>

        <li>
          <router-link to="/fw-data-request">Custom Query News</router-link>
        </li>
      </ul>
    </nav>

    <!-- Hamburger Button (always visible in the top-left) -->
    <button @click="toggleSidebar" class="hamburger-btn">☰</button>

    <!-- Main Content -->
    <div class="main-container">
      <div class="container-fluid">
        <div class="row">
          <header class="col-12 mb-4">
            <h1>News Analytics Dashboard</h1>
            <p>Analytics summary and insights for news articles and sources.</p>
          </header>
        </div>
      </div>

      <div class="analytics-container">
        <div class="current-date">
          Report for <br />
          {{ yesterdayDisplayFormat }}
        </div>

        <!-- Articles and Sources Cards -->
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

        <!-- Pie Chart Card -->
        <router-link to="/daily-trending-topic" class="card pie-chart-card">
          <h3 class="card-title">Trending Topic</h3>
          <div class="card-content">
            <div class="chart-container">
              <canvas ref="pieChart"></canvas>
            </div>
          </div>
        </router-link>
      </div>
    </div>
  </div>
</template>

<script>
import Chart from 'chart.js/auto';
import ChartDataLabels from 'chartjs-plugin-datalabels';

export default {
  name: 'DailyOverview',
  data() {
    return {
      isSidebarOpen: false,
      dailyOverviewOpen: false,
      weeklyOverviewOpen: false,
      totalArticles: 0,
      totalSources: 0,
      yesterdayAPIFormat: this.getYesterdayDate('YYYYMMDD'),
      yesterdayDisplayFormat: this.getYesterdayDate('readable'),
      trendingData: [],
      nerData: {},
      pieChartInstance: null
    };
  },
  methods: {
    toggleSidebar() {
      this.isSidebarOpen = !this.isSidebarOpen;
    },

    toggleDailyOverview(){
      this.dailyOverviewOpen = !this.dailyOverviewOpen;
    },

    toggleWeeklyOverview(){
      this.weeklyOverviewOpen = !this.weeklyOverviewOpen;
    },

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
        const options = {
          weekday: 'long',
          year: 'numeric',
          month: 'long',
          day: 'numeric'
        };
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
      const dataValues = this.trendingData.map((item) =>
        parseFloat(item.percentage.toFixed(2))
      );

      this.pieChartInstance = new Chart(ctx, {
        type: 'pie',
        data: {
          labels: topics,
          datasets: [
            {
              data: dataValues,
              backgroundColor: [
                '#FF6384',
                '#36A2EB',
                '#FFCE56',
                '#4BC0C0',
                '#9966FF',
                '#FF9F40',
                '#61FF84',
                '#3642EB',
                '#CE5645',
                '#40A0C0',
                '#2296FF',
                '#9F40FF'
              ]
            }
          ]
        },
        options: {
          responsive: true,
          maintainAspectRatio: false,
          plugins: {
            legend: {
              position: 'right',
              labels: {
                generateLabels: (chart) => {
                  return chart.data.labels.map((label, index) => ({
                    text: label,
                    fillStyle: chart.data.datasets[0].backgroundColor[index]
                  }));
                }
              }
            },
            tooltip: {
              callbacks: {
                label: (context) => {
                  const keyword = keywords[context.dataIndex];
                  return `${keyword}: ${context.raw}%`;
                }
              }
            },
            datalabels: {
              color: '#fff',
              formatter: (value, context) => {
                return value > 5
                  ? `${context.chart.data.labels[context.dataIndex]}: ${value}%`
                  : '';
              }
            }
          }
        }
      });
    }
  },
  mounted() {
    this.fetchData();
  }
};
</script>

<style scoped>
/* Parent container for the entire app */
.app-container {
  position: relative; /* so we can absolutely-position the hamburger if we want */
  min-height: 100vh;
  background-color: #fff;
}

/* SIDEBAR */
.sidebar {
  position: fixed;
  top: 0;
  left: 0;
  width: 240px;
  height: 100%;
  background-color: #343a40;
  color: #fff;
  padding: 2rem 1rem;
  transform: translateX(-240px); /* hidden by default */
  transition: transform 0.3s ease-in-out;
  z-index: 999; /* on top of main content */
  /*box-sizing: border-box;*/
}

.sidebar.open {
  transform: translateX(0);
}

.sidebar-menu {
  list-style: none;
  padding: 0;
  margin: 0;
}

/* Each top-level <li> */
.sidebar-menu > li {
  margin-bottom: 1rem;
  /*padding-top: 10px;*/
}

/* Parent clickable item (Daily / Weekly) */
.sidebar-link {
  display: flex;
  align-items: center;
  cursor: pointer;
  font-weight: 600; /* Bold, optional */
  color: #fff;
}

/* The arrow icon next to "Daily Overview" */
.arrow {
  margin-left: auto; /* pushes arrow to the right end */
  transition: transform 0.3s;
}

/* Rotate arrow 90 degrees when submenu is open */
.arrow.rotated {
  transform: rotate(90deg);
}

/* Submenu list */
.nested-links {
  list-style: none;
  padding-left: 1.5rem; /* indent child items */
  margin-top: 0.5rem;
}

/* Style nested link items */
.nested-links li {
  margin-bottom: 0.5rem;
}

/* Example styling for nested links */
.nested-links a {
  color: #ddd; /* slightly lighter than parent link color */
  text-decoration: none;
}

.nested-links a:hover {
  text-decoration: underline;
}

/* Close button inside sidebar */
.close-btn {
  background: none;
  border: none;
  color: #fff;
  font-size: 2rem;
  position: absolute;
  top: 10px;
  right: 10px;
  cursor: pointer;
}

/* HAMBURGER BUTTON (top-left of screen) */
.hamburger-btn {
  position: fixed;
  top: 20px;
  left: 20px;
  font-size: 26px;
  background: none;
  border: none;
  color: #343a40;
  cursor: pointer;
  z-index: 1000;
}

/* MAIN CONTENT - centered with margin: 0 auto */
.main-container {
  max-width: 1200px;
  margin: 0 auto; /* keeps the main content centered */
  /*padding-top: 80px; /* so content doesn't run under the hamburger button */
  padding-top: 0;
}

/* 
  2) EXISTING STYLES 
     (tweak them to ensure your analytics cards remain centered within .main-container)
*/

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

/* Analytics Container */
.analytics-container {
  display: flex;
  flex-direction: column;
  gap: 30px;
  width: 100%;
  background-color: #f5f5f5;
  padding: 20px;
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
  justify-content: space-between; /* ensure the cards row is centered */
}

@media (max-width: 768px) {
  .cards-row {
    flex-direction: column;
    align-items: center; /* center them in the column layout */
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
  min-width: 200px;
}

.card:hover {
  transform: scale(1.02);
}

.articles-card,
.sources-card {
  flex: 1;
  max-width: none; /* limit card width so they don't stretch too far */
}

.pie-chart-card {
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