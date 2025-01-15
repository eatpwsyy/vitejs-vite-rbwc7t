<!-- JSONFileMaker.vue -->
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
        <FormField
          label="News Type"
          id="news_type"
          required
          tooltipId="news_type-tooltip"
          tooltip="Latest news provides access to the latest and breaking news while news archive provides access to the old news data based on the date you provide."
        >
          <template #input>
            <select v-model="form.news_type" id="news_type" class="form-control" required>
              <option value="" disabled>Select News Type</option>
              <option value="archive">Archive</option>
              <option value="latest">Latest</option>
            </select>
          </template>
        </FormField>

        <!-- Query (q) -->
        <FormField
          label="Query (q) (Optional)"
          id="q"
          tooltipId="q-tooltip"
          tooltip="Search news articles for specific keywords or phrases present in the news title, content, URL, meta keywords and meta description. Max 512 characters."
        >
          <template #input>
            <input type="text" v-model="form.q" id="q" class="form-control" />
          </template>
        </FormField>

        <!-- Conditional Fields Based on News Type -->
        <div v-if="form.news_type === 'latest'">
          <FormField
            label="Timeframe"
            id="timeframe"
            required
            tooltipId="timeframe-tooltip"
            tooltip="Search the news articles for a specific timeframe. Either 24 or 48 hours ago."
          >
            <template #input>
              <select v-model="form.timeframe" id="timeframe" class="form-control" required>
                <option value="" disabled>Select Timeframe</option>
                <option value="24">24 Jam Lalu</option>
                <option value="48">48 Jam Lalu</option>
              </select>
            </template>
          </FormField>
        </div>

        <div v-if="form.news_type === 'archive'">
          <FormField
            label="From Date (YYYY-MM-DD)"
            id="from_date"
            required
            tooltipId="from_date-tooltip"
            tooltip="Get news data from a particular start date in the past. ex: 2024-12-18"
          >
            <template #input>
              <input type="date" v-model="form.from_date" id="from_date" class="form-control" required />
            </template>
          </FormField>

          <FormField
            label="To Date (YYYY-MM-DD)"
            id="to_date"
            required
            tooltipId="to_date-tooltip"
            tooltip="Get news data from a particular end date in the past. ex: 2024-12-28"
          >
            <template #input>
              <input type="date" v-model="form.to_date" id="to_date" class="form-control" required />
            </template>
          </FormField>
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
            <!-- Query in Title (qInTitle) -->
            <FormField
              label="Query in Title (qInTitle) (Optional)"
              id="qInTitle"
              tooltipId="qInTitle-tooltip"
              tooltip="Search news articles for specific keywords or phrases present in the news titles only. Max 512 characters. Note: qInTitle can't be used with q or qInMeta parameter in the same query."
            >
              <template #input>
                <input type="text" v-model="form.qInTitle" id="qInTitle" class="form-control" />
              </template>
            </FormField>

            <!-- Query in Meta (qInMeta) -->
            <FormField
              label="Query in Meta (qInMeta) (Optional)"
              id="qInMeta"
              tooltipId="qInMeta-tooltip"
              tooltip="Search news articles for specific keywords or phrases present in the news titles, URL, meta keywords and meta description only. Max 512 characters. Note: qInMeta can't be used with q or qInTitle parameter in the same query."
            >
              <template #input>
                <input type="text" v-model="form.qInMeta" id="qInMeta" class="form-control" />
              </template>
            </FormField>

            <!-- Country -->
            <FormField
              label="Country (Optional)"
              id="country"
              tooltipId="country-tooltip"
              tooltip="Search the news articles from specific countries. Max 5 countries. Use country codes. ex: id,it,sg"
            >
              <template #input>
                <input type="text" v-model="form.country" id="country" class="form-control" />
              </template>
            </FormField>

            <!-- Category Selector (Using MultiSelector.vue) -->
            <MultiSelector
              :items="form.category"
              :available-items="availableCategories"
              label="Category (Optional)"
              id="category"
              tooltip="Search the news articles for specific categories. Max 5 categories."
              tooltipId="category-tooltip"
              placeholder="Select Categories (Max 5)"
              @add-item="addCategory"
              @remove-item="removeCategory"
              :item-label-map="categoryLabelMap"
            />

            <!-- Language -->
            <FormField
              label="Language (Optional)"
              id="language"
              tooltipId="language-tooltip"
              tooltip="Search the news articles for a specific language (English or Indonesia)."
            >
              <template #input>
                <select v-model="form.language" id="language" class="form-control">
                  <option value="">Select Language</option>
                  <option value="en">English</option>
                  <option value="id">Indonesia</option>
                </select>
              </template>
            </FormField>

            <!-- Domain -->
            <FormField
              label="Domain (Optional)"
              id="domain"
              tooltipId="domain-tooltip"
              tooltip="Search the news articles for specific news sources. Max 5 domains. ex: nytimes,bbc"
            >
              <template #input>
                <input type="text" v-model="form.domain" id="domain" class="form-control" />
              </template>
            </FormField>

            <!-- Domain URL -->
            <FormField
              label="Domain URL (Optional)"
              id="domainurl"
              tooltipId="domainurl-tooltip"
              tooltip="Search the news articles for specific domains or news sources URL. Max 5 sources. ex: nytimes.com,bbc.com,bbc.co.uk"
            >
              <template #input>
                <input type="text" v-model="form.domainurl" id="domainurl" class="form-control" />
              </template>
            </FormField>

            <!-- Exclude Domain -->
            <FormField
              label="Exclude Domain (Optional)"
              id="excludedomain"
              tooltipId="excludedomain-tooltip"
              tooltip="You can exclude specific domains or news sources to search the news articles. Max 5 domains. ex: nytimes.com,bbc.com,bbc.co.uk"
            >
              <template #input>
                <input type="text" v-model="form.excludedomain" id="excludedomain" class="form-control" />
              </template>
            </FormField>

            <!-- Exclude Field -->
            <FormField
              label="Exclude Field (Optional)"
              id="excludefield"
              tooltipId="excludefield-tooltip"
              tooltip="Limit the response field. You can exclude multiple response fields in a single query. ex: source_icon,pubdate,link"
            >
              <template #input>
                <input type="text" v-model="form.excludefield" id="excludefield" class="form-control" />
              </template>
            </FormField>

            <!-- Priority Domain -->
            <FormField
              label="Priority Domain (Optional)"
              id="prioritydomain"
              tooltipId="prioritydomain-tooltip"
              tooltip="Search the news articles only from specific category of news domains. Top: Fetches news articles from the top 10%. Medium: Fetches news articles from the top 30%. Low: Fetches news articles from the top 50%."
            >
              <template #input>
                <select v-model="form.prioritydomain" id="prioritydomain" class="form-control">
                  <option value="">Select Priority Domain</option>
                  <option value="top">Top</option>
                  <option value="medium">Medium</option>
                  <option value="low">Low</option>
                </select>
              </template>
            </FormField>

            <!-- Timezone -->
            <FormField
              label="Timezone (Optional)"
              id="timezone"
              tooltipId="timezone-tooltip"
              tooltip="Search the news articles for a specific timezone. ex: Asia/Jakarta"
            >
              <template #input>
                <input type="text" v-model="form.timezone" id="timezone" class="form-control" />
              </template>
            </FormField>

            <!-- Image -->
            <FormField
              label="Image (0 or 1)"
              id="image"
              tooltipId="image-tooltip"
              tooltip="Search the news articles with featured image or without featured image. Use '1' for articles with featured image and '0' for articles without it."
            >
              <template #input>
                <input type="number" v-model.number="form.image" id="image" class="form-control" min="0" max="1" />
              </template>
            </FormField>

            <!-- Video -->
            <FormField
              label="Video (0 or 1)"
              id="video"
              tooltipId="video-tooltip"
              tooltip="Search the news articles with featured video or without featured video. Use '1' for articles with featured video and '0' for articles without it."
            >
              <template #input>
                <input type="number" v-model.number="form.video" id="video" class="form-control" min="0" max="1" />
              </template>
            </FormField>

            <!-- Size -->
            <FormField
              label="Size (Max 50)"
              id="size"
              tooltipId="size-tooltip"
              tooltip="You can customize the number of articles you get per request."
            >
              <template #input>
                <input type="number" v-model.number="form.size" id="size" class="form-control" min="1" max="50" />
              </template>
            </FormField>

            <!-- Full Content -->
            <FormField
              label="Full Content (0 or 1)"
              id="full_content"
              tooltipId="full_content-tooltip"
              tooltip="Search the news articles with full content or without full content. Use '1' for news articles with full_content field and '0' for news articles without it."
            >
              <template #input>
                <input type="number" v-model.number="form.full_content" id="full_content" class="form-control" min="0" max="1" />
              </template>
            </FormField>

            <div v-if="form.news_type === 'latest'">
              <!-- Remove Duplicate -->
              <FormField
                label="Remove Duplicate (0 or 1)"
                id="removeduplicate"
                tooltipId="removeduplicate-tooltip"
                tooltip="The 'removeduplicate' parameter will allow users to filter out duplicate articles. Use '1' to remove duplicate articles."
              >
                <template #input>
                  <input type="number" v-model.number="form.removeduplicate" id="removeduplicate" class="form-control" min="0" max="1" />
                </template>
              </FormField>

              <!-- Sentiment -->
              <FormField
                label="Sentiment"
                id="sentiment"
                tooltipId="sentiment-tooltip"
                tooltip="Search the news articles based on the sentiment of the news article (positive, negative, neutral)."
              >
                <template #input>
                  <select v-model="form.sentiment" id="sentiment" class="form-control">
                    <option value="" disabled>Select Sentiment</option>
                    <option value="positive">Positive</option>
                    <option value="negative">Negative</option>
                    <option value="neutral">Neutral</option>
                  </select>
                </template>
              </FormField>

              <!-- Tag Selector (Using MultiSelector.vue) -->
              <MultiSelector
                :items="form.tag"
                :available-items="availableTags"
                label="Tag (Optional)"
                id="tag"
                tooltip="Search the news articles for specific AI-classified tags. Max 5 tags."
                tooltipId="tag-tooltip"
                placeholder="Select Tags (Max 5)"
                @add-item="addTag"
                @remove-item="removeTag"
                :item-label-map="tagLabelMap"
              />
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
import FormField from './FormField.vue';
import MultiSelector from './MultiSelector.vue';
import {
  categoryLabelMap,
  tagLabelMap,
  availableCategories,
  availableTags,
} from '../../config/newsConfig.js';

export default {
  name: "JSONFileMaker",
  components: {
    FormField,
    MultiSelector,
  },
  data() {
    return {
      categoryLabelMap,
      tagLabelMap,
      availableCategories,
      availableTags,
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
        excludefield: "ai_region,ai_org,",
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
        tag: [],
      },
      isLoading: false,
      isSubmitting: false,
      error: null,
      response: null,
      generatedLink: "",
      isUrlGenerated: false,
      showAdvanced: false,
    };
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

      // Validate 'category' (max 5)
      if (this.form.tag.length > 5) {
        this.error = "You can select up to 5 tags.";
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
     addCategory(category) {
      if (category && this.form.category.length < 5) {
        this.form.category.push(category);
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
     * Adds the selected category to the selectedCategories array.
     */
    addTag(tag){
      if(tag && this.form.tag.length < 5){
        this.form.tag.push(tag);
        this.error = null;
      }
    },

    /**
     * Removes a tag from the selectedTags array by index.
     * @param {number} index - The index of the tag to remove.
     */
     removeTag(index) {
      this.form.tag.splice(index, 1);
      this.error = null;
    },
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

/* Main Container */
.fw-data-request {
  display: flex;
  flex-direction: column;
  align-items: center;
}

/* Form Section */
.form-section {
  width: 80%;
  max-width: 800px;
  margin-bottom: 30px;
}

/* Button Styles */
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

.btn-primary {
  padding: 10px 20px;
  background-color: #007bff; /* Blue button */
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 1rem;
  transition: background-color 0.3s;
}

.btn-primary:hover {
  background-color: #0056b3;
}

.btn-success {
  padding: 10px 20px;
  background-color: #28a745; /* Green button */
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 1rem;
  transition: background-color 0.3s;
}

.btn-success:hover {
  background-color: #218838;
}

/* Loading, Error, and Response Sections */
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

/* Toggle Advanced Options */
.toggle-advanced {
  margin-bottom: 15px;
}

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

/* Fade Transition */
.fade-enter-active, .fade-leave-active {
  transition: opacity 0.5s;
}
.fade-enter, .fade-leave-to {
  opacity: 0;
}

/* Responsive Design */
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