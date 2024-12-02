<template>
  <div class="trending-topic">
    <header class="col-12 mb-4">
        <h2>Daily Trending Topics <br /></h2>
        <h3>{{ yesterdayDisplayFormat }}</h3>
    </header>
    <!-- Pie Chart Section -->
    <div class="pie-chart-container">
      <!-- Ensure PieChart component exists and receives 'data' and emits 'sliceClick' correctly -->
      <PieChart :data="pieChartData" @sliceClick="handleSliceClick" />
    </div>

    <!-- Articles Section -->
    <div v-if="currentTopic" class="articles-container">
      <h3 class="keywords-title">Articles for "{{ getKeywordsForCurrentTopic().join(', ') }}"</h3>
      <div class="article-grid">
        <div
          v-for="article in 2"
          :key="article.article_id"
          class="article-card"
        >
          <img :src="article.image_url" alt="article image" class="article-image" />
          <div class="article-details">
            <h4>{{ article.title }}</h4>
            <p>{{ article.description }}</p>
            <a :href="article.link" target="_blank">Read more</a>
          </div>
        </div>
      </div>

      <!-- Pagination -->
      <div class="pagination">
        <button
          v-for="page in totalPages"
          :key="page"
          @click="setPage(page)"
          :class="{ active: currentPage === page }"
        >
          {{ page }}
        </button>
      </div>
    </div>
  </div>
</template>

<script>
import PieChart from './PieChart.vue'; // Assuming the PieChart component is created separately

export default {
  components: {
    PieChart,
  },
  data() {
    return {
      yesterdayAPIFormat: this.getYesterdayDate('YYYYMMDD'), // Date for API request
      yesterdayDisplayFormat: this.getYesterdayDate('readable'),
      trendingTopics: [],
      currentTopic: null, // Topic selected by the user
      currentPage: 1, // Current page for pagination
      articlesPerPage: 3, // Number of articles per page
    };
  },
  computed: {
    // Get the list of articles for the selected topic
    selectedArticles() {
      console.log('selectedArticles computed', this.currentTopic); // Log currentTopic to see if it's correct
      if (!this.currentTopic) return [];
      return this.trendingTopics.find((article) => article.topic === this.currentTopic)?.articles || [];
    },

    // Paginate the selected articles
    paginatedArticles() {
      console.log('paginatedArticles computed');
      const start = (this.currentPage - 1) * this.articlesPerPage;
      const end = start + this.articlesPerPage;
      return this.selectedArticles.slice(start, end);
    },

    // Calculate the total pages
    totalPages() {
      console.log('totalPages computed');
      return Math.ceil(this.selectedArticles.length / this.articlesPerPage);
    },
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
        const response = await fetch(`/api/get-data?req_id=xxx&req_type=trendingtopic&freq=daily&date=${this.yesterdayAPIFormat}`);
        console.log('API Response Status:', response.status);

        if (!response.ok) {
          throw new Error(`HTTP error! Status: ${response.status}`);
        }

        const data = await response.json();
        console.log('Parsed API Data:', data);
        this.trendingTopics = data.trending;

        this.pieChartData();
      } catch (error) {
        console.error('Error in fetchData:', error);
      }
    },

    // Prepare the data for the Pie Chart
    pieChartData() {
      console.log('pieChartData computed');
      return this.trendingTopics.map((article) => ({
        label: article.topic,
        value: article.percentage,
        keywords: article.keywords,
        articles: article.articles,
        color: this.getRandomColor(),
      }));
    },

    // Handle the slice click from the Pie Chart
    handleSliceClick(topic) {
      console.log('handleSliceClick', topic); // Make sure this logs the correct topic
      this.currentTopic = topic;
      this.currentPage = 1; // Reset to the first page when changing topics
    },

    // Get the keywords for the selected topic
    getKeywordsForCurrentTopic() {
      if (!this.currentTopic) return [];
      const topicData = this.trendingTopics.find((item) => item.topic === this.currentTopic);
      return topicData ? topicData.keywords : [];
    },

    // Set the current page for pagination
    setPage(page) {
      console.log('setPage', page);
      this.currentPage = page;
    },

    // Generate a random color for the pie chart slice
    getRandomColor() {
      const letters = '0123456789ABCDEF';
      let color = '#';
      for (let i = 0; i < 6; i++) {
        color += letters[Math.floor(Math.random() * 16)];
      }
      return color;
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

.trending-topic {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.pie-chart-container {
  width: 80%;
  height: 300px;
  margin-top: 20px;
}

.articles-container {
  width: 80%;
  margin-top: 20px;
}

.article-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 20px;
}

.article-card {
  border: 1px solid #ccc;
  border-radius: 8px;
  padding: 15px;
  min-height: 250px; /* Adjust this as needed */
}

.article-image {
  width: 100%;
  height: auto;
  border-radius: 8px;
}

.article-details {
  margin-top: 10px;
  overflow: hidden; /* This ensures no content is cut off */
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: center;
}

.pagination button {
  margin: 0 5px;
  padding: 5px 10px;
  border: none;
  background-color: #f0f0f0;
  cursor: pointer;
  border-radius: 4px;
}

.pagination button.active {
  background-color: #007bff;
  color: white;
}

.keywords-title {
  text-align: center;  /* Center the text */
  width: 100%;         /* Make sure it takes up the full width */
  margin: 0 auto;      /* Center it within its container */
  font-size: 1.5rem;   /* Optional: adjust the font size */
  font-weight: bold;   /* Optional: make the text bold */
}
</style>