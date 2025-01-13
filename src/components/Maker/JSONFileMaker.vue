<template>
  <div class="fw-data-request">
    <!-- Header Section -->
    <header class="header col-12 mb-4">
      <router-link to="/" class="back-link" aria-label="Back to Daily Overview">
        &#8592; Back to Daily Overview
      </router-link>
      <h2>JSON Data Request</h2>
      <h3>Generate JSON and Link for NewsData.io</h3>
    </header>

    <!-- Form Section -->
    <div class="form-section">
      <form @submit.prevent="generateUrl">
        <!-- Required Fields -->
        <div class="form-group">
          <label for="news_type">
            News Type <span class="required">*</span>
            <div class="tooltip-container">
              <span 
                class="tooltip-icon"
                aria-describedby="news_type-tooltip"
                >ⓘ</span>
              <div class="tooltip-content" id="news_type-tooltip" role="tooltip">
                Latest news provides access to the latest and breaking news while news archive provides access to the old news data based on the date you provide.
              </div>
            </div>
          </label>
          <select v-model="form.news_type" id="news_type" class="form-control" required>
            <option value="" disabled>Select News Type</option>
            <option value="archive">Archive</option>
            <option value="latest">Latest</option>
          </select>
        </div>

        <!-- Toggle Button for Optional Fields -->
        <div class="form-group toggle-advanced">
          <a 
            href="#" 
            @click.prevent="toggleAdvanced" 
            class="toggle-link"
            :aria-expanded="showAdvanced"
            aria-controls="optional-fields"
          >
            {{ showAdvanced ? 'Hide Advanced Options' : 'Show Advanced Options' }}
          </a>
        </div>

        <!-- Optional Fields -->
        <transition name="fade">
          <div v-if="showAdvanced" id="optional-fields">
            <!-- Query (q) -->
            <div class="form-group">
              <label for="q">
                Query (q) (Optional)
                <div class="tooltip-container">
                  <span 
                    class="tooltip-icon"
                    aria-describedby="q-tooltip"
                    >ⓘ</span>
                  <div class="tooltip-content" id="q-tooltip" role="tooltip">
                    Search news articles for specific keywords or phrases present in the news title, content, URL, meta keywords and meta description. Max 512 characters.
                  </div>
                </div>
              </label>
              <input type="text" v-model="form.q" id="q" class="form-control" />
            </div>

            <!-- Query in Title (qInTitle) -->
            <div class="form-group">
              <label for="qInTitle">
                Query in Title (qInTitle) (Optional)
                <div class="tooltip-container">
                  <span 
                    class="tooltip-icon"
                    aria-describedby="qInTitle-tooltip"
                    >ⓘ</span>
                  <div class="tooltip-content" id="qInTitle-tooltip" role="tooltip">
                    Search news articles for specific keywords or phrases present in the news titles only. Max 512 characters. Note: qInTitle can't be used with q or qInMeta parameter in the same query.
                  </div>
                </div>
              </label>
              <input type="text" v-model="form.qInTitle" id="qInTitle" class="form-control" />
            </div>

            <!-- Query in Meta (qInMeta) -->
            <div class="form-group">
              <label for="qInMeta">
                Query in Meta (qInMeta) (Optional)
                <div class="tooltip-container">
                  <span 
                    class="tooltip-icon"
                    aria-describedby="qInMeta-tooltip"
                    >ⓘ</span>
                  <div class="tooltip-content" id="qInMeta-tooltip" role="tooltip">
                    Search news articles for specific keywords or phrases present in the news titles, URL, meta keywords and meta description only. Max 512 characters. Note: qInMeta can't be used with q or qInTitle parameter in the same query.
                  </div>
                </div>
              </label>
              <input type="text" v-model="form.qInMeta" id="qInMeta" class="form-control" />
            </div>

            <!-- Country -->
            <div class="form-group">
              <label for="country">
                Country (Optional)
                <div class="tooltip-container">
                  <span 
                    class="tooltip-icon"
                    aria-describedby="country-tooltip"
                    >ⓘ</span>
                  <div class="tooltip-content" id="country-tooltip" role="tooltip">
                    Search the news articles from specific countries. Max 5 countries. Use country codes. ex: id,it,sg
                  </div>
                </div>
              </label>
              <input type="text" v-model="form.country" id="country" class="form-control" />
            </div>

            <!-- Category -->
            <div class="form-group">
              <label for="category">
                Category (Optional)
                <div class="tooltip-container">
                  <span class="tooltip-icon" aria-describedby="category-tooltip">ⓘ</span>
                  <div class="tooltip-content" id="category-tooltip" role="tooltip">
                    Search the news articles for specific categories. Max 5 categories.
                  </div>
                </div>
              </label>

              <select 
                v-model="selectedCategory" 
                id="category" 
                class="form-control" 
                @change="addCategory"
                :disabled="form.category.length >= 5"
              >
                <option disabled value="">Select Categories (Max 5)</option>
                <option
                  v-for="option in availableCategories"
                  :key="option.value"
                  :value="option.value"
                >{{ option.text }}</option>
              </select>
              <small>
                Selected {{ form.category.length }}/5 categories.
              </small>

              <!-- Display Selected Categories as Tags -->
              <div class="selected-categories" v-if="form.category.length">
                <span
                  v-for="(category, index) in form.category"
                  :key="index"
                  class="category-tag"
                >
                  {{ getCategoryLabel(category) }}
                  <button type="button" @click="removeCategory(index)" aria-label="Remove Category">&times;</button>
                </span>
              </div>
            </div>

            <!-- Language -->
            <div class="form-group">
              <label for="language">
                Language (Optional)
                <div class="tooltip-container">
                  <span 
                    class="tooltip-icon"
                    aria-describedby="language-tooltip"
                    >ⓘ</span>
                  <div class="tooltip-content" id="language-tooltip" role="tooltip">
                    Search the news articles for a specific language (English or Indonesia).
                  </div>
                </div>
              </label>
              <select v-model="form.language" id="language" class="form-control">
                <option value="">Select Language</option>
                <option value="en">English</option>
                <option value="id">Indonesia</option>
              </select>
            </div>

            <!-- Domain -->
            <div class="form-group">
              <label for="domain">
                Domain (Optional)
                <div class="tooltip-container">
                  <span 
                    class="tooltip-icon"
                    aria-describedby="domain-tooltip"
                    >ⓘ</span>
                  <div class="tooltip-content" id="domain-tooltip" role="tooltip">
                    Search the news articles for specific news sources. Max 5 domains. ex: nytimes,bbc
                  </div>
                </div>
              </label>
              <input type="text" v-model="form.domain" id="domain" class="form-control" />
            </div>

            <!-- Domain URL -->
            <div class="form-group">
              <label for="domainurl">
                Domain URL (Optional)
                <div class="tooltip-container">
                  <span 
                    class="tooltip-icon"
                    aria-describedby="domainurl-tooltip"
                    >ⓘ</span>
                  <div class="tooltip-content" id="domainurl-tooltip" role="tooltip">
                    Search the news articles for specific domains or news sources URL. Max 5 sources. ex: nytimes.com,bbc.com,bbc.co.uk
                  </div>
                </div>
              </label>
              <input type="text" v-model="form.domainurl" id="domainurl" class="form-control" />
            </div>

            <!-- Exclude Domain -->
            <div class="form-group">
              <label for="excludedomain">
                Exclude Domain (Optional)
                <div class="tooltip-container">
                  <span 
                    class="tooltip-icon"
                    aria-describedby="excludedomain-tooltip"
                    >ⓘ</span>
                  <div class="tooltip-content" id="excludedomain-tooltip" role="tooltip">
                    You can exclude specific domains or news sources to search the news articles. Max 5 domains. ex: nytimes.com,bbc.com,bbc.co.uk
                  </div>
                </div>
              </label>
              <input type="text" v-model="form.excludedomain" id="excludedomain" class="form-control" />
            </div>

            <!-- Exclude Field -->
            <div class="form-group">
              <label for="excludefield">
                Exclude Field (Optional)
                <div class="tooltip-container">
                  <span 
                    class="tooltip-icon"
                    aria-describedby="excludefield-tooltip"
                    >ⓘ</span>
                  <div class="tooltip-content" id="excludefield-tooltip" role="tooltip">
                    Limit the response field. You can exclude multiple response fields in a single query. ex: source_icon,pubdate,link
                  </div>
                </div>
              </label>
              <input type="text" v-model="form.excludefield" id="excludefield" class="form-control" />
            </div>

            <!-- Priority Domain -->
            <div class="form-group">
              <label for="prioritydomain">
                Priority Domain (Optional)
                <div class="tooltip-container">
                  <span 
                    class="tooltip-icon"
                    aria-describedby="prioritydomain-tooltip"
                    >ⓘ</span>
                  <div class="tooltip-content" id="prioritydomain-tooltip" role="tooltip">
                    Search the news articles only from specific category of news domains. Top: Fetches news articles from the top 10%. Medium: Fetches news articles from the top 30%. Low: Fetches news articles from the top 50%.
                  </div>
                </div>
              </label>
              <select v-model="form.prioritydomain" id="prioritydomain" class="form-control">
                <option value="">Select Priority Domain</option>
                <option value="top">Top</option>
                <option value="medium">Medium</option>
                <option value="low">Low</option>
              </select>
            </div>

            <!-- Timezone -->
            <div class="form-group">
              <label for="timezone">
                Timezone (Optional)
                <div class="tooltip-container">
                  <span 
                    class="tooltip-icon"
                    aria-describedby="timezone-tooltip"
                    >ⓘ</span>
                  <div class="tooltip-content" id="timezone-tooltip" role="tooltip">
                    Search the news articles for a specific timezone. ex: Asia/Jakarta
                  </div>
                </div>
              </label>
              <input type="text" v-model="form.timezone" id="timezone" class="form-control" />
            </div>

            <!-- Image -->
            <div class="form-group">
              <label for="image">
                Image (0 or 1)
                <div class="tooltip-container">
                  <span 
                    class="tooltip-icon"
                    aria-describedby="image-tooltip"
                    >ⓘ</span>
                  <div class="tooltip-content" id="image-tooltip" role="tooltip">
                    Search the news articles with featured image or without featured image. Use '1' for articles with featured image and '0' for articles without it.
                  </div>
                </div>
              </label>
              <input type="number" v-model.number="form.image" id="image" class="form-control" min="0" max="1" />
            </div>

            <!-- Video -->
            <div class="form-group">
              <label for="video">
                Video (0 or 1)
                <div class="tooltip-container">
                  <span 
                    class="tooltip-icon"
                    aria-describedby="video-tooltip"
                    >ⓘ</span>
                  <div class="tooltip-content" id="video-tooltip" role="tooltip">
                    Search the news articles with featured video or without featured video. Use '1' for articles with featured video and '0' for articles without it.
                  </div>
                </div>
              </label>
              <input type="number" v-model.number="form.video" id="video" class="form-control" min="0" max="1" />
            </div>

            <!-- Size -->
            <div class="form-group">
              <label for="size">
                Size (Max 50)
                <div class="tooltip-container">
                  <span 
                    class="tooltip-icon"
                    aria-describedby="size-tooltip"
                    >ⓘ</span>
                  <div class="tooltip-content" id="size-tooltip" role="tooltip">
                    You can customize the number of articles you get per request.
                  </div>
                </div>
              </label>
              <input type="number" v-model.number="form.size" id="size" class="form-control" min="1" max="50" />
            </div>

            <!-- Full Content -->
            <div class="form-group">
              <label for="full_content">
                Full Content (0 or 1)
                <div class="tooltip-container">
                  <span 
                    class="tooltip-icon"
                    aria-describedby="full_content-tooltip"
                    >ⓘ</span>
                  <div class="tooltip-content" id="full_content-tooltip" role="tooltip">
                    Search the news articles with full content or without full content. Use '1' for news articles with full_content field and '0' for news articles without it.
                  </div>
                </div>
              </label>
              <input type="number" v-model.number="form.full_content" id="full_content" class="form-control" min="0" max="1" />
            </div>

            <!-- Conditional Fields Based on News Type -->
            <div v-if="form.news_type === 'archive'">
              <!-- From Date -->
              <div class="form-group">
                <label for="from_date">
                  From Date (YYYY-MM-DD) <span class="required">*</span>
                  <div class="tooltip-container">
                    <span 
                      class="tooltip-icon"
                      aria-describedby="from_date-tooltip"
                      >ⓘ</span>
                    <div class="tooltip-content" id="from_date-tooltip" role="tooltip">
                      Get news data from a particular start date in the past. ex: 2024-12-18
                    </div>
                  </div>
                </label>
                <input type="date" v-model="form.from_date" id="from_date" class="form-control" required />
              </div>

              <!-- To Date -->
              <div class="form-group">
                <label for="to_date">
                  To Date (YYYY-MM-DD) <span class="required">*</span>
                  <div class="tooltip-container">
                    <span 
                      class="tooltip-icon"
                      aria-describedby="to_date-tooltip"
                      >ⓘ</span>
                    <div class="tooltip-content" id="to_date-tooltip" role="tooltip">
                      Get news data from a particular end date in the past. ex: 2024-12-28
                    </div>
                  </div>
                </label>
                <input type="date" v-model="form.to_date" id="to_date" class="form-control" required />
              </div>
            </div>

            <div v-if="form.news_type === 'latest'">
              <!-- Remove Duplicate -->
              <div class="form-group">
                <label for="removeduplicate">
                  Remove Duplicate (0 or 1)
                  <div class="tooltip-container">
                    <span 
                      class="tooltip-icon"
                      aria-describedby="removeduplicate-tooltip"
                      >ⓘ</span>
                    <div class="tooltip-content" id="removeduplicate-tooltip" role="tooltip">
                      The 'removeduplicate' parameter will allow users to filter out duplicate articles. Use '1' to remove duplicate articles.
                    </div>
                  </div>
                </label>
                <input type="number" v-model.number="form.removeduplicate" id="removeduplicate" class="form-control" min="0" max="1" />
              </div>

              <!-- Sentiment -->
              <div class="form-group">
                <label for="sentiment">
                  Sentiment <span class="required">*</span>
                  <div class="tooltip-container">
                    <span 
                      class="tooltip-icon"
                      aria-describedby="sentiment-tooltip"
                      >ⓘ</span>
                    <div class="tooltip-content" id="sentiment-tooltip" role="tooltip">
                      Search the news articles based on the sentiment of the news article (positive, negative, neutral).
                    </div>
                  </div>
                </label>
                <select v-model="form.sentiment" id="sentiment" class="form-control" required>
                  <option value="" disabled>Select Sentiment</option>
                  <option value="positive">Positive</option>
                  <option value="negative">Negative</option>
                  <option value="neutral">Neutral</option>
                </select>
              </div>

              <!-- Timeframe -->
              <div class="form-group">
                <label for="timeframe">
                  Timeframe <span class="required">*</span>
                  <div class="tooltip-container">
                    <span 
                      class="tooltip-icon"
                      aria-describedby="timeframe-tooltip"
                      >ⓘ</span>
                    <div class="tooltip-content" id="timeframe-tooltip" role="tooltip">
                      Search the news articles for a specific timeframe. Either 24 or 48 hours ago.
                    </div>
                  </div>
                </label>
                <select v-model="form.timeframe" id="timeframe" class="form-control" required>
                  <option value="" disabled>Select Timeframe</option>
                  <option value="24">24 Jam Lalu</option>
                  <option value="48">48 Jam Lalu</option>
                </select>
              </div>

              <!-- Tag -->
              <div class="form-group">
                <label for="tag">
                  Tag (Optional)
                  <div class="tooltip-container">
                    <span 
                      class="tooltip-icon"
                      aria-describedby="tag-tooltip"
                      >ⓘ</span>
                    <div class="tooltip-content" id="tag-tooltip" role="tooltip">
                      Search the news articles for specific AI-classified tags.
                    </div>
                  </div>
                </label>
                <input type="text" v-model="form.tag" id="tag" class="form-control" />
              </div>
            </div>
          </div>
          </transition>
            <!-- Submit Button for Generating URL -->
            <div class="form-group">
              <button type="submit" class="btn-submit" :disabled="isSubmitting">
                {{ isSubmitting ? 'Generating URL...' : 'Generate URL' }}
              </button>
            </div>
          </form>
        </div>

        <!-- Generated URL Section -->
        <div v-if="isUrlGenerated" class="generated-link">
          <p><strong>Generated API URL:</strong></p>
          <a :href="generatedLink" target="_blank">{{ generatedLink }}</a>
        </div>

        <!-- Button to Fetch API Response -->
        <div v-if="isUrlGenerated" class="fetch-api-button">
          <button @click="fetchApiResponse" class="btn btn-primary" :disabled="isLoading">
            {{ isLoading ? 'Fetching Data...' : 'Fetch API Response' }}
          </button>
        </div>

        <!-- Loading State -->
        <div v-if="isLoading" class="loading">Sending request...</div>

        <!-- Error State -->
        <div v-if="error" class="error">{{ error }}</div>

        <!-- API Response Section -->
        <div v-if="response && !isLoading && !error" class="response-section">
          <h3>API Response</h3>
          
          <button @click="downloadJSON" class="btn btn-success">
            Download JSON
          </button>
          <br />
          <pre>{{ JSON.stringify(response, null, 2) }}</pre>
          <br />
        </div>
      </div>
    </template>

<script>
export default {
  name: "FwDataRequest",
  data() {
    return {
      form: {
        news_type: "archive",
        q: "",
        qInTitle: "",
        qInMeta: "",
        country: "",
        category: [],
        language: "",
        domain: "",
        domainurl: "",
        excludedomain: "",
        excludefield: "",
        prioritydomain: "",
        timezone: "",
        image: null,
        video: null,
        size: null,
        full_content: null,
        from_date: "",
        to_date: "",
        removeduplicate: null,
        sentiment: "",
        timeframe: "",
        tag: "",
      },
      selectedCategory: "",
      isLoading: false,
      isSubmitting: false,
      error: null,
      response: null,
      generatedLink: "",
      isUrlGenerated: false,
      showAdvanced: false, // New state variable
    };
  },
  computed: {
    /**
     * Computes the list of available categories by excluding already selected ones.
     */
    availableCategories() {
      const allCategories = [
        { value: "business", text: "Business" },
        { value: "crime", text: "Crime" },
        { value: "domestic", text: "Domestic" },
        { value: "education", text: "Education" },
        { value: "entertainment", text: "Entertainment" },
        { value: "environment", text: "Environment" },
        { value: "food", text: "Food" },
        { value: "health", text: "Health" },
        { value: "lifestyle", text: "Lifestyle" },
        { value: "politics", text: "Politics" },
        { value: "science", text: "Science" },
        { value: "sports", text: "Sports" },
        { value: "technology", text: "Technology" },
        { value: "top", text: "Top" },
        { value: "tourism", text: "Tourism" },
        { value: "world", text: "World" },
        { value: "other", text: "Other" },
      ];
      return allCategories.filter(
        (cat) => !this.form.category.includes(cat.value)
      );
    },
  },
  methods: {
    /**
     * Toggles the visibility of advanced (optional) fields.
     */
    toggleAdvanced() {
      this.showAdvanced = !this.showAdvanced;
    },

    /**
     * Validates the form before submission.
     * Returns true if valid, otherwise sets the error message.
     */
    validateForm() {
      // Reset error
      this.error = null;

      // Validate 'news_type' (already required via HTML)
      if (!this.form.news_type) {
        this.error = "News Type is required.";
        return false;
      }

      // Conditional Validation based on 'news_type'
      if (this.form.news_type === "archive") {
        if (!this.form.from_date) {
          this.error = "From Date is required when News Type is 'Archive'.";
          return false;
        }
        if (!this.form.to_date) {
          this.error = "To Date is required when News Type is 'Archive'.";
          return false;
        }
      }

      if (this.form.news_type === "latest") {
        if (!this.form.timeframe) {
          this.error = "Timeframe is required when News Type is 'Latest'.";
          return false;
        }
      }

      // Additional Validations
      // Validate 'size' if provided
      if (this.form.size !== null && (isNaN(this.form.size) || this.form.size > 50)) {
        this.error = "Size must be an integer not greater than 50.";
        return false;
      }

      // Validate 'full_content' if provided
      if (
        this.form.full_content !== null &&
        ![0, 1].includes(this.form.full_content)
      ) {
        this.error = "full_content must be 0 or 1.";
        return false;
      }

      // Validate 'removeduplicate' if provided
      if (
        this.form.news_type === "latest" &&
        this.form.removeduplicate !== null &&
        ![0, 1].includes(this.form.removeduplicate)
      ) {
        this.error = "removeduplicate must be 0 or 1.";
        return false;
      }

      // Validate 'category' (max 5)
      if (this.form.category.length > 5) {
        this.error = "You can select up to 5 categories.";
        return false;
      }

      // All validations passed
      this.error = null;
      return true;
    },
    
    /**
     * Constructs the query parameters object based on the form data.
     */
    constructParams() {
      const params = {};
      for (const key in this.form) {
        const value = this.form[key];
        if (
          value !== null &&
          value !== undefined &&
          value !== "" &&
          !(typeof value === "number" && isNaN(value))
        ) {
          if (Array.isArray(value)) {
            // Join array values with commas
            params[key] = value.join(",");
          } else {
            params[key] = value;
          }
        }
      }
      return params;
    },
  
    /**
     * Generates the API URL based on form data.
     */
    generateUrl() {
      if (!this.validateForm()) {
        return;
      }
  
      // Construct the query parameters
      const params = this.constructParams();
  
      // Construct the query string
      const queryString = Object.keys(params)
        .map(
          (key) =>
            `${encodeURIComponent(key)}=${encodeURIComponent(params[key])}`
        )
        .join("&");
  
      // Construct the generated URL
      this.generatedLink = `http://192.168.166.42:5000/api/fw-data?${queryString}`;
  
      // Show the URL immediately
      this.isUrlGenerated = true;
    },
  
    /**
     * Sends the API request to fetch the JSON data.
     */
    async fetchApiResponse() {
      const params = this.constructParams();
      const queryString = Object.keys(params)
        .map(
          (key) =>
            `${encodeURIComponent(key)}=${encodeURIComponent(params[key])}`
        )
        .join("&");
  
      this.isLoading = true;
      this.isSubmitting = true;
      this.error = null;
      this.response = null;
  
      try {
        const response = await fetch(`/api/fw-data?${queryString}`, {
          method: "GET",
          headers: {
            "Content-Type": "application/json",
          },
        });
  
        if (!response.ok) {
          const errorText = await response.text();
          throw new Error(
            `HTTP error! Status: ${response.status}, Response: ${errorText}`
          );
        }
  
        const data = await response.json();
        this.response = data;
      } catch (err) {
        console.error("Error fetching data:", err);
        this.error =
          err.message || "An unexpected error occurred while fetching data.";
      } finally {
        this.isLoading = false;
        this.isSubmitting = false;
      }
    },
  
    /**
     * Downloads the JSON data as a file.
     */
    downloadJSON() {
      const blob = new Blob([JSON.stringify(this.response, null, 2)], {
        type: "application/json",
      });
  
      const link = document.createElement("a");
      link.href = URL.createObjectURL(blob);
      link.download = "response.json";
      link.click();
    },
    
    /**
     * Adds the selected category to the selectedCategories array.
     */
    addCategory() {
      if (this.selectedCategory && this.form.category.length < 5) {
        this.form.category.push(this.selectedCategory);
        this.selectedCategory = ""; // Reset the select
        this.error = null;
      }
    },

    /**
     * Removes a category from the selectedCategories array by index.
     * @param {number} index - The index of the category to remove.
     */
    removeCategory(index) {
      this.form.category.splice(index, 1);
      this.error = null;
    },

    /**
     * Helper method to get the display label for a category value.
     * @param {string} value - The category value.
     * @returns {string} - The display text for the category.
     */
    getCategoryLabel(value) {
      const categoryMap = {
        business: "Business",
        crime: "Crime",
        domestic: "Domestic",
        education: "Education",
        entertainment: "Entertainment",
        environment: "Environment",
        food: "Food",
        health: "Health",
        lifestyle: "Lifestyle",
        politics: "Politics",
        science: "Science",
        sports: "Sports",
        technology: "Technology",
        top: "Top",
        tourism: "Tourism",
        world: "World",
        other: "Other",
      };
      return categoryMap[value] || value;
    },
  },
};
</script>

<style scoped>
/* Existing Styles */

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

.fw-data-request {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.form-section {
  width: 80%;
  max-width: 800px;
  margin-bottom: 30px;
}

.form-group {
  margin-bottom: 15px;
  display: flex;
  flex-direction: column;
}

/* Indicate required fields dynamically */
.form-group label .required {
  color: red;
  margin-left: 4px;
}

/* Disabled select and input styles */
.form-group input[disabled],
.form-group select[disabled] {
  background-color: #e9ecef;
  cursor: not-allowed;
}

.form-control {
  padding: 8px 12px;
  border: 1px solid #ccc;
  border-radius: 4px;
}

.btn-submit {
  padding: 10px 20px;
  background-color: #28a745; /* Green button */
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 1rem;
  transition: background-color 0.3s;
}

.btn-submit:hover {
  background-color: #218838;
}

.loading,
.error,
.response-section,
.generated-link,
.fetch-api-button {
  width: 80%;
  max-width: 800px;
  margin-top: 20px;
  text-align: center;
}

.loading {
  color: #6c757d; /* Light gray text color */
}

.error {
  color: red;
}

.response-section pre {
  background-color: #f8f9fa;
  padding: 15px;
  border-radius: 4px;
  overflow-x: auto;
  text-align: left;
}

/* Selected Categories Tags */
.selected-categories {
  margin-top: 10px;
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.category-tag {
  background-color: #007bff;
  color: white;
  padding: 5px 10px;
  border-radius: 15px;
  display: flex;
  align-items: center;
  font-size: 0.9em;
}

.category-tag button {
  background: none;
  border: none;
  color: white;
  margin-left: 8px;
  cursor: pointer;
  font-size: 1em;
  line-height: 1;
}

/* Optional: Hover effect for remove button */
.category-tag button:hover {
  color: #ffdddd;
}

/* Tooltip Container */
.tooltip-container {
  position: relative;
  display: inline-block;
  margin-left: 4px; /* Space between label and icon */
}

/* Tooltip Icon */
.tooltip-icon {
  color: #007bff;
  cursor: help;
  font-size: 0.9em;
}

/* Tooltip Text */
.tooltip-content {
  visibility: hidden;
  width: 300px;
  background-color: #333;
  color: #fff;
  text-align: left;
  border-radius: 4px;
  padding: 8px;
  position: absolute;
  z-index: 1;
  bottom: 100%; /* Position above the icon */
  left: 50%;
  transform: translateX(-50%);
  opacity: 0;
  transition: opacity 0.3s;
}

/* Tooltip Arrow */
.tooltip-content::after {
  content: "";
  position: absolute;
  top: 100%; /* At the bottom edge of the tooltip */
  left: 50%;
  transform: translateX(-50%);
  border-width: 2px;
  border-style: solid;
  border-color: #333 transparent transparent transparent;
}

/* Show Tooltip on Hover */
.tooltip-container:hover .tooltip-content {
  visibility: visible;
  opacity: 1;
}

.form-group small {
  color: #6c757d; /* Light gray */
  margin-top: 5px;
}

/* Highlight required fields */
.required {
  color: red;
  margin-left: 4px;
}

/* Style for the toggle button */
.btn-toggle {
  padding: 8px 12px;
  background-color: #007bff; /* Blue color */
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.9rem;
  transition: background-color 0.3s;
}

.btn-toggle:hover {
  background-color: #0056b3;
}

/* Fade Transition */
.fade-enter-active, .fade-leave-active {
  transition: opacity 0.5s;
}
.fade-enter, .fade-leave-to {
  opacity: 0;
}

/* Style for the toggle link */
.toggle-link {
  color: #007bff; /* Bootstrap primary color */
  text-decoration: none;
  cursor: pointer;
  font-size: 0.9rem;
}

.toggle-link:hover {
  text-decoration: underline;
  color: #0056b3; /* Darker shade on hover */
}

/* Remove default button styles if still using a button */
.btn-toggle {
  background: none;
  border: none;
  padding: 0;
  color: #007bff;
  text-decoration: underline;
  cursor: pointer;
  font-size: 0.9rem;
}

.btn-toggle:hover {
  color: #0056b3;
}

@media (max-width: 768px) {
  .form-section,
  .loading,
  .error,
  .response-section,
  .generated-link,
  .fetch-api-button {
    width: 95%;
  }
}

@media (max-width: 480px) {
  .form-section,
  .loading,
  .error,
  .response-section,
  .generated-link,
  .fetch-api-button {
    width: 100%;
  }

}
</style>