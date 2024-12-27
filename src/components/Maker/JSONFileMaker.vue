<template>
  <div class="fw-data-request">
    <!-- Header Section -->
    <header class="header col-12 mb-4">
      <router-link to="/" class="back-link" aria-label="Back to Dashboard">
        &#8592; Back to Dashboard
      </router-link>
      <h2>JSON Data Request</h2>
      <h3>Generate JSON and Link for NewsData.io</h3>
    </header>

    <!-- Form Section -->
    <div class="form-section">
      <form @submit.prevent="generateUrl">
        <!-- News Type Selection -->
        <div class="form-group">
          <label for="news_type">News Type</label>
          <select v-model="form.news_type" id="news_type" class="form-control">
            <option value="archive">Archive</option>
            <option value="latest">Latest</option>
          </select>
        </div>

        <!-- Common Parameters -->
        <div class="form-group">
          <label for="q">Query (q) (Optional)</label>
          <input type="text" v-model="form.q" id="q" class="form-control" />
        </div>

        <div class="form-group">
          <label for="qInTitle">Query in Title / qInTitle (Optional)</label>
          <input type="text" v-model="form.qInTitle" id="qInTitle" class="form-control" />
        </div>

        <div class="form-group">
          <label for="qInMeta">Query in Meta / qInMeta (Optional)</label>
          <input type="text" v-model="form.qInMeta" id="qInMeta" class="form-control" />
        </div>

        <div class="form-group">
          <label for="country">Country (Optional)</label>
          <input type="text" v-model="form.country" id="country" class="form-control" />
        </div>

        <div class="form-group">
          <label for="category">Category (Optional)</label>
          <select v-model="form.category" id="category" class="form-control">
            <option value="">Select Category</option>
            <option value="business">Business</option>           
            <option value="crime">Crime</option>
            <option value="domestic">Domestic</option>
            <option value="education">Education</option>
            <option value="entertainment">Entertainment</option>
            <option value="environment">Environment</option>
            <option value="food">Food</option>
            <option value="health">Health</option>
            <option value="lifestyle">Lifestyle</option>
            <option value="politics">Politics</option>
            <option value="science">Science</option>
            <option value="sports">Sports</option>
            <option value="technology">Technology</option>
            <option value="top">Top</option>
            <option value="tourism">Tourism</option>
            <option value="world">World</option>
            <option value="other">Other</option>
          </select>
        </div>

        <div class="form-group">
          <label for="language">Language (Optional)</label>
          <select v-model="form.language" id="language" class="form-control">
            <option value="">Select Language</option>
            <option value="en">English</option>
            <option value="id">Indonesia</option>
          </select>
        </div>

        <div class="form-group">
          <label for="domain">Domain (Optional)</label>
          <input type="text" v-model="form.domain" id="domain" class="form-control" />
        </div>

        <div class="form-group">
          <label for="domainurl">Domain URL (Optional)</label>
          <input type="text" v-model="form.domainurl" id="domainurl" class="form-control" />
        </div>

        <div class="form-group">
          <label for="excludedomain">Exclude Domain (Optional)</label>
          <input type="text" v-model="form.excludedomain" id="excludedomain" class="form-control" />
        </div>

        <div class="form-group">
          <label for="excludefield">Exclude Field (Optional)</label>
          <input type="text" v-model="form.excludefield" id="excludefield" class="form-control" />
        </div>

        <div class="form-group">
          <label for="prioritydomain">Priority Domain (Optional)</label>
          <select v-model="form.prioritydomain" id="prioritydomain" class="form-control">
            <option value="">Select Priority Domain</option>
            <option value="top">Top</option>
            <option value="medium">Medium</option>
            <option value="low">Low</option>
          </select>
        </div>

        <div class="form-group">
          <label for="timezone">Timezone (Optional)</label>
          <input type="text" v-model="form.timezone" id="timezone" class="form-control" />
        </div>

        <div class="form-group">
          <label for="image">Image (0 or 1)</label>
          <input type="number" v-model.number="form.image" id="image" class="form-control" min="0" max="1" />
        </div>

        <div class="form-group">
          <label for="video">Video (0 or 1)</label>
          <input type="number" v-model.number="form.video" id="video" class="form-control" min="0" max="1" />
        </div>

        <div class="form-group">
          <label for="size">Size (Max 50)</label>
          <input type="number" v-model.number="form.size" id="size" class="form-control" min="1" max="50" />
        </div>

        <div class="form-group">
          <label for="full_content">Full Content (0 or 1)</label>
          <input type="number" v-model.number="form.full_content" id="full_content" class="form-control" min="0" max="1" />
        </div>

        <!-- Conditional Fields Based on News Type -->
        <div v-if="form.news_type !== 'latest'">
          <div class="form-group">
            <label for="from_date">From Date (YYYY-MM-DD)</label>
            <input type="date" v-model="form.from_date" id="from_date" class="form-control" />
          </div>

          <div class="form-group">
            <label for="to_date">To Date (YYYY-MM-DD)</label>
            <input type="date" v-model="form.to_date" id="to_date" class="form-control" />
          </div>
        </div>

        <div v-else>
          <div class="form-group">
            <label for="removeduplicate">Remove Duplicate (0 or 1)</label>
            <input type="number" v-model.number="form.removeduplicate" id="removeduplicate" class="form-control" min="0" max="1" />
          </div>

          <div class="form-group">
            <label for="sentiment">Sentiment (Optional)</label>
            <select v-model="form.sentiment" id="sentiment" class="form-control">
              <option value="">Select Sentiment</option>
              <option value="positive">Positive</option>
              <option value="negative">Negative</option>
              <option value="neutral">Neutral</option>
            </select>
          </div>

          <div class="form-group">
            <label for="timeframe">Timeframe (Optional)</label>
            <select v-model="form.timeframe" id="timeframe" class="form-control">
              <option value="">Select Timeframe</option>
              <option value="24">24 Jam Lalu</option>
              <option value="48">48 Jam Lalu</option>
            </select>
          </div>

          <div class="form-group">
            <label for="tag">Tag (Optional)</label>
            <input type="text" v-model="form.tag" id="tag" class="form-control" />
          </div>
        </div>

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
        category: "",
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
      isLoading: false,
      isSubmitting: false,
      error: null,
      response: null,
      generatedLink: "",
      isUrlGenerated: false,  // Flag to track URL generation
    };
  },
  
  methods: {
    /**
     * Validates the form before submission.
     * Returns true if valid, otherwise sets the error message.
     */
    validateForm() {
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
  
      // If news_type is not 'latest', ensure from_date and to_date are provided
      if (
        this.form.news_type !== "latest" &&
        (!this.form.from_date || !this.form.to_date)
      ) {
        this.error =
          "Both 'from_date' and 'to_date' must be provided when news_type is 'archive'.";
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
          params[key] = value;
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
  },
};
</script>

<style scoped>
/* Reuse existing styles */

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

.form-group label {
  margin-bottom: 5px;
  font-weight: bold;
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
.response-section {
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

/* Responsive Adjustments */
@media (max-width: 768px) {
  .form-section,
  .loading,
  .error,
  .response-section {
    width: 95%;
  }
}

@media (max-width: 480px) {
  .form-section,
  .loading,
  .error,
  .response-section {
    width: 100%;
  }

  .btn-submit {
    width: 100%;
  }
}
</style>