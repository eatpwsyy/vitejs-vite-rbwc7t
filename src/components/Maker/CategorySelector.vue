<!-- CategorySelector.vue -->
<template>
    <div class="form-group">
      <label for="category">
        Category (Optional)
        <Tooltip tooltip-id="category-tooltip">
          Search the news articles for specific categories. Max 5 categories.
        </Tooltip>
      </label>
      <select 
        v-model="selected" 
        id="category" 
        class="form-control" 
        @change="addCategory"
        :disabled="isDisabled"
      >
        <option disabled value="">Select Categories (Max 5)</option>
        <option
          v-for="option in availableOptions"
          :key="option.value"
          :value="option.value"
        >{{ option.text }}</option>
      </select>
  
      <div class="selected-categories" v-if="categories.length">
        <span
          v-for="(category, index) in categories"
          :key="index"
          class="category-tag"
        >
          {{ getCategoryLabel(category) }}
          <button type="button" @click="removeCategory(index)" aria-label="Remove Category">&times;</button>
        </span>
      </div>
    </div>
  </template>
  
  <script>
  import Tooltip from './Tooltip.vue';
  
  export default {
    name: 'CategorySelector',
    components: { Tooltip },
    props: {
      /**
       * The list of currently selected categories.
       * Example: ['business', 'crime']
       */
      categories: {
        type: Array,
        required: true,
      },
      /**
       * The maximum number of categories a user can select.
       */
      maxCategories: {
        type: Number,
        default: 5,
      },
      /**
       * The complete list of available categories.
       * Example: [{ value: 'business', text: 'Business' }, ...]
       */
      availableCategories: {
        type: Array,
        required: true,
      },
    },
    data() {
      return {
        /**
         * Holds the currently selected value in the dropdown.
         */
        selected: '',
      };
    },
    computed: {
      /**
       * Computes the list of categories available for selection by excluding
       * the ones already selected.
       */
      availableOptions() {
        return this.availableCategories.filter(
          (option) => !this.categories.includes(option.value)
        );
      },
      /**
       * Determines if the dropdown should be disabled based on the maximum
       * number of selected categories.
       */
      isDisabled() {
        return this.categories.length >= this.maxCategories;
      },
    },
    methods: {
      /**
       * Emits an event to add a selected category and resets the dropdown.
       */
      addCategory() {
        if (this.selected && !this.isDisabled) {
          this.$emit('add-category', this.selected);
          this.selected = ''; // Reset the dropdown
        }
      },
      /**
       * Emits an event to remove a category based on its index.
       * @param {number} index - The index of the category to remove.
       */
      removeCategory(index) {
        this.$emit('remove-category', index);
      },
      /**
       * Retrieves the display label for a given category value.
       * @param {string} value - The value of the category.
       * @returns {string} - The display label of the category.
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
  .form-group {
    margin-bottom: 15px;
    display: flex;
    flex-direction: column;
  }
  
  .selected-categories {
    margin-top: 10px;
    margin-bottom: 15px;
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
  
  .category-tag button:hover {
    color: #ffdddd;
  }
  </style>  