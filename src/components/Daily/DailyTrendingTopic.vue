<template>
  <div class="trending-topic">
    <!-- Header Section with Back Link -->
    <header class="header col-12 mb-4">
      <router-link to="/" class="back-link" aria-label="Back to Daily Overview">
        &#8592; Back to Daily Overview
      </router-link>
      <h2>Daily Trending Topics</h2>
      <h3>{{ yesterdayDisplayFormat }}</h3>
    </header>

    <!-- Loading State -->
    <div v-if="isLoading" class="loading">
      Loading trending topics...
    </div>

    <!-- Error State -->
    <div v-if="error" class="error">
      {{ error }}
    </div>

    <!-- Instructions Section -->
    <div class="info-section">
      <p>
        Click on a slice of the pie chart to view related articles for that keyword. <br />
        Percentages are displayed only for topics with more than 5% coverage.
      </p>
    </div>

    <!-- Content Section -->
    <div v-if="!isLoading && !error" class="content-section">
      <!-- Pie Chart Section -->
      <div class="pie-chart-container">
        <PieChart :data="pieChartDataComputed" @sliceClick="handleSliceClick" />
      </div>

      <!-- Articles Section -->
      <div v-if="currentTopic" class="articles-container">
        <h3 class="keywords-title">
          Articles for "{{ currentTopic }}:
          {{ getKeywordsForCurrentTopic().join(', ') }}"
        </h3>
        <div v-if="selectedArticles.length === 0" class="no-articles">
          No articles available for this topic.
        </div>
        <transition-group name="fade" tag="div" class="article-grid">
          <div
            v-for="article in paginatedArticles"
            :key="article.article_id"
            class="article-card"
          >
            <div class="image-container">
              <img
                :src="article.image_url"
                :alt="`Image for ${article.title}`"
                class="article-image"
              />
            </div>
            <div class="article-details">
              <h4>{{ article.title }}</h4>
              <p>{{ truncateDescription(article.description) }}</p>
              <a :href="article.link" target="_blank" rel="noopener noreferrer">
                Read more
              </a>
            </div>
          </div>
        </transition-group>

        <!-- Pagination Controls -->
        <div class="pagination" v-if="totalPages > 1">
          <!-- First Page Button -->
          <button
            @click="setPage(1)"
            :disabled="currentPage === 1"
            class="pagination-arrow"
            aria-label="First Page"
          >
            &#171; First
          </button>

          <!-- Previous Window Button -->
          <button
            @click="prevWindow"
            :disabled="currentWindowStart === 1"
            class="pagination-arrow"
            aria-label="Previous Window"
          >
            &laquo;
          </button>

          <!-- Page Number Buttons -->
          <button
            v-for="page in visiblePages"
            :key="page"
            @click="setPage(page)"
            :class="{ active: currentPage === page }"
            class="pagination-page"
            :aria-current="currentPage === page ? 'page' : undefined"
          >
            {{ page }}
          </button>

          <!-- Last Page Button -->
          <button
            @click="setPage(totalPages)"
            :disabled="currentPage === totalPages"
            class="last-page-button"
            aria-label="Last Page"
          >
            Last &#187;
          </button>

          <!-- Next Window Button -->
          <button
            @click="nextWindow"
            :disabled="currentWindowEnd === totalPages"
            class="pagination-arrow"
            aria-label="Next Window"
          >
            &raquo;
          </button>

          <!-- Page Indicator -->
          <span class="pagination-info">
            Page {{ currentPage }} of {{ totalPages }}
          </span>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import PieChart from '../PieChart.vue'; // Ensure the PieChart component is correctly imported

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
      articlesPerPage: 6, // Updated from 3 to 6 articles per page
      isLoading: false, // Loading state
      error: null, // Error message

      // Pagination window variables
      paginationWindowSize: 5, // Number of page buttons to show at a time
      currentWindowStart: 1,
      currentWindowEnd: 5,
    };
  },
  computed: {
    /**
     * Retrieves the list of articles for the selected topic.
     */
    selectedArticles() {
      if (!this.currentTopic) return [];
      const topicData = this.trendingTopics.find(
        (topic) => topic.topic === this.currentTopic
      );
      if (!topicData) {
        console.warn(`No articles found for topic: ${this.currentTopic}`);
        return [];
      }
      return topicData.articles || [];
    },

    /**
     * Returns the articles for the current page based on pagination.
     */
    paginatedArticles() {
      const start = (this.currentPage - 1) * this.articlesPerPage;
      const end = start + this.articlesPerPage;
      return this.selectedArticles.slice(start, end);
    },

    /**
     * Calculates the total number of pages based on articles per page.
     */
    totalPages() {
      return Math.ceil(this.selectedArticles.length / this.articlesPerPage);
    },

    /**
     * Prepares the data for the Pie Chart as a computed property.
     */
    pieChartDataComputed() {
      if (!Array.isArray(this.trendingTopics) || this.trendingTopics.length === 0) {
        console.warn('No trending topics available to generate pie chart data.');
        return [];
      }

      const pieData = this.trendingTopics.map((topic, index) => ({
        label: topic.topic,
        value: topic.percentage, // Ensure 'percentage' exists and is a number
        keywords: topic.keywords,
        color: this.getRandomColor(index),
      }));

      return pieData;
    },

    /**
     * Determines the visible page buttons based on the current window.
     */
    visiblePages() {
      const pages = [];
      for (let i = this.currentWindowStart; i <= this.currentWindowEnd; i++) {
        if (i > this.totalPages) break;
        pages.push(i);
      }
      return pages;
    },
  },
  methods: {
    /**
     * Formats yesterday's date based on the specified format.
     * @param {string} format - The desired date format ('YYYYMMDD' or 'readable').
     * @returns {string} - Formatted date string.
     */
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
          day: 'numeric',
        };
        return yesterday.toLocaleDateString(undefined, options);
      }

      return '';
    },

    /**
     * Fetches trending topics data from the API.
     */
    async fetchData() {
      this.isLoading = true;
      this.error = null;

      try {
        const apiUrl = `/api/get-data?req_id=xxx&req_type=trendingtopic&freq=daily&date=${encodeURIComponent(
          this.yesterdayAPIFormat
        )}`;

        const response = await fetch(apiUrl);

        if (!response.ok) {
          const errorText = await response.text();
          throw new Error(`HTTP error! Status: ${response.status}, Response: ${errorText}`);
        }

        const contentType = response.headers.get('content-type');

        if (!contentType || !contentType.includes('application/json')) {
          const text = await response.text();
          throw new Error('Received non-JSON response');
        }

        const data = await response.json();

        // Runtime checks for data.trending
        if (data.hasOwnProperty('trending')) {
          if (Array.isArray(data.trending)) {
            const isValid = data.trending.every(
              (topic) =>
                typeof topic === 'object' &&
                topic !== null &&
                typeof topic.topic === 'string' &&
                typeof topic.percentage === 'number' &&
                Array.isArray(topic.keywords) &&
                Array.isArray(topic.articles)
            );

            if (isValid) {
              this.trendingTopics = data.trending;

              // If there's at least one topic, find the one with the highest percentage
              if (this.trendingTopics.length > 0) {
                const maxTopic = this.trendingTopics.reduce((prev, curr) =>
                  curr.percentage > prev.percentage ? curr : prev
                );
                this.currentTopic = maxTopic.topic; // Set currentTopic to the topic with max percentage
              }
            } else {
              throw new Error('Invalid data format received from API.');
            }
          } else {
            throw new Error(
              `Expected 'data.trending' to be an array, but received type '${typeof data.trending}'.`
            );
          }
        } else {
          throw new Error("'data.trending' property is missing in the API response.");
        }
      } catch (error) {
        console.error('An error occurred while fetching or processing API data:', error);
        this.error = 'Failed to load trending topics. Please try again later.';
        this.trendingTopics = []; // Fallback to empty array
      } finally {
        this.isLoading = false;
      }
    },

    /**
     * Handles slice clicks from the PieChart component.
     * @param {string} topic - The topic label that was clicked.
     */
    handleSliceClick(topic) {
      this.currentTopic = topic;
      this.currentPage = 1; // Reset to the first page when changing topics
      this.updatePaginationWindow(this.currentPage);
    },

    /**
     * Retrieves keywords for the currently selected topic.
     * @returns {Array<string>} - Array of keywords.
     */
    getKeywordsForCurrentTopic() {
      if (!this.currentTopic) return [];
      const topicData = this.trendingTopics.find((item) => item.topic === this.currentTopic);
      if (!topicData) {
        console.warn(`No keywords found for topic: ${this.currentTopic}`);
        return [];
      }
      return topicData.keywords || [];
    },

    /**
     * Sets the current page for pagination and updates the pagination window.
     * @param {number} page - The page number to set.
     */
    setPage(page) {
      this.currentPage = page;
      this.updatePaginationWindow(page);
    },

    /**
     * Updates the pagination window based on the current page.
     * Ensures that only a subset of page buttons are visible at a time.
     * @param {number} page - The current page number.
     */
    updatePaginationWindow(page) {
      const halfWindow = Math.floor(this.paginationWindowSize / 2);
      let start = page - halfWindow;
      let end = page + halfWindow;

      if (start < 1) {
        start = 1;
        end = Math.min(this.paginationWindowSize, this.totalPages);
      }

      if (end > this.totalPages) {
        end = this.totalPages;
        start = Math.max(end - this.paginationWindowSize + 1, 1);
      }

      this.currentWindowStart = start;
      this.currentWindowEnd = end;
    },

    /**
     * Navigates to the previous window of pages.
     */
    prevWindow() {
      if (this.currentWindowStart > 1) {
        this.currentWindowStart -= this.paginationWindowSize;
        this.currentWindowEnd = Math.max(
          this.currentWindowEnd - this.paginationWindowSize,
          this.paginationWindowSize
        );
      }
    },

    /**
     * Navigates to the next window of pages.
     */
    nextWindow() {
      if (this.currentWindowEnd < this.totalPages) {
        this.currentWindowStart += this.paginationWindowSize;
        this.currentWindowEnd = Math.min(
          this.currentWindowEnd + this.paginationWindowSize,
          this.totalPages
        );
      }
    },

    /**
     * Generates a random HEX color from a predefined palette for consistency.
     * @param {number} index - The index of the topic for color assignment.
     * @returns {string} - A HEX color string.
     */
    getRandomColor(index) {
      const colors = ['#FF6384', '#36A2EB', '#FFCE56', '#4BC0C0', '#9966FF', '#FF9F40'];
      return colors[index % colors.length];
    },

    /**
     * Truncates the description to a specified number of sentences.
     * @param {string} description - The full description text.
     * @param {number} maxSentences - Maximum number of sentences to display.
     * @returns {string} - Truncated description.
     */
    truncateDescription(description, maxSentences = 2) {
      if (!description) return '';
      const sentences = description.match(/[^\.!\?]+[\.!\?]+/g);
      if (!sentences) return description;
      return sentences.slice(0, maxSentences).join(' ');
    },
  },
  mounted() {
    this.fetchData(); // Fetch data when the component is mounted
  },
};
</script>

<style scoped>
/* Header Styles */
.header {
  position: relative;
  text-align: center;
  margin-bottom: 30px;
  padding: 20px;
  border-radius: 8px;
}

.back-link {
  position: absolute;
  left: 20px;
  top: 50%;
  transform: translateY(-50%);
  font-size: 1rem;
  color: #007bff;
  text-decoration: none;
  font-weight: bold;
  transition: color 0.3s;
}

.back-link:hover {
  color: #0056b3;
}

h2 {
  font-size: 2rem;
  color: #343a40; /* Darker text color */
}

h3 {
  font-size: 1.2rem;
  color: #6c757d; /* Light gray text color */
}

.trending-topic {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.content-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  width: 100%;
}

.info-section {
  text-align: center;
  margin-bottom: 10px;
  font-size: 1rem;
  color: #333;
}

.pie-chart-container {
  width: 50%; /* Adjust as needed */
  max-width: 400px; /* Adjust as needed */
  height: auto; /* Allow height to adjust based on width */
  margin-top: 20px;
}

.articles-container {
  width: 80%;
  margin-top: 20px;
}

.article-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr); /* Always three columns */
  gap: 20px;
}

.article-card {
  border: 1px solid #ccc;
  border-radius: 8px;
  padding: 15px;
  display: flex;
  flex-direction: column;
  height: 100%; /* Ensure uniform height */
  transition: transform 0.3s, box-shadow 0.3s; /* For hover effects */
}

.article-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
}

.image-container {
  width: 100%;
  height: 200px; /* Fixed height for images */
  overflow: hidden;
  border-radius: 8px;
  margin-bottom: 10px;
}

.article-image {
  width: 100%;
  height: 100%;
  object-fit: cover; /* Ensures the image covers the container */
}

.article-details {
  flex-grow: 1;
  display: flex;
  flex-direction: column;
}

.article-details h4 {
  margin: 0 0 10px 0;
  font-size: 1.1rem;
}

.article-details p {
  flex-grow: 1;
  margin-bottom: 10px;
  color: #555;
}

.article-details a {
  align-self: flex-start;
  color: #007bff;
  text-decoration: none;
  font-weight: bold;
}

.pagination {
  margin-top: 20px;
  display: flex;
  align-items: center;
  justify-content: center; /* Centers the pagination buttons */
}

.pagination button {
  margin: 0 3px;
  padding: 5px 10px;
  border: none;
  background-color: #f0f0f0;
  cursor: pointer;
  border-radius: 4px;
  transition: background-color 0.3s, color 0.3s;
}

.pagination button:hover:not(:disabled) {
  background-color: #e0e0e0;
}

.pagination button.active {
  background-color: #007bff;
  color: white;
}

.pagination button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.pagination-arrow {
  font-size: 1rem;
}

.last-page-button {
  /* Style similar to other buttons */
  background-color: #f0f0f0;
  cursor: pointer;
}

.last-page-button:hover:not(:disabled) {
  background-color: #e0e0e0;
}

.keywords-title {
  text-align: center; /* Center the text */
  width: 100%; /* Make sure it takes up the full width */
  margin: 0 auto 20px auto; /* Center it within its container and add bottom margin */
  font-size: 1.5rem; /* Adjust the font size */
  font-weight: bold; /* Make the text bold */
}

.loading,
.error,
.no-articles {
  text-align: center;
  font-size: 1.2rem;
  margin-top: 20px;
}

.loading {
  color: #6c757d; /* Light gray text color */
}

.error {
  color: red;
}

.no-articles {
  color: #6c757d; /* Light gray text color */
}

.pagination-info {
  margin-left: 10px;
  font-size: 1rem;
  color: #333;
}

/* Transition Styles */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.5s;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

/* Responsive Adjustments */
@media (max-width: 768px) {
  .pie-chart-container {
    width: 80%;
  }

  .articles-container {
    width: 95%;
  }

  .image-container {
    height: 150px;
  }

  .pagination-info {
    display: none; /* Hide on smaller screens */
  }

  .article-grid {
    grid-template-columns: repeat(2, 1fr); /* Two columns on small screens */
  }

  .header {
    padding: 15px;
  }

  .back-link {
    left: 10px;
    font-size: 0.9rem;
  }
}

@media (max-width: 480px) {
  .article-grid {
    grid-template-columns: 1fr; /* Single column on very small screens */
  }

  .pie-chart-container {
    width: 100%;
  }

  .article-card {
    padding: 10px;
  }

  .image-container {
    height: 120px;
  }
}
</style>