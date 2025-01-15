<!-- MultiSelector.vue -->
<template>
    <div class="form-group">
      <label :for="id">
        {{ label }}
        <Tooltip v-if="tooltip" :tooltip-id="tooltipId">
          {{ tooltip }}
        </Tooltip>
      </label>
      <select 
        v-model="selected" 
        :id="id" 
        class="form-control" 
        @change="addItem"
        :disabled="isDisabled"
        aria-describedby="tooltipId"
      >
        <option disabled value="">{{ placeholder }}</option>
        <option
          v-for="option in availableOptions"
          :key="option.value"
          :value="option.value"
        >
          {{ option.text }}
        </option>
      </select>
  
      <div class="selected-items" v-if="items.length">
        <span
          v-for="(item, index) in items"
          :key="index"
          class="item-tag"
        >
          {{ getItemLabel(item) }}
          <button type="button" @click="removeItem(index)" aria-label="Remove Item">&times;</button>
        </span>
      </div>
    </div>
  </template>
  
  <script>
  import Tooltip from './Tooltip.vue';
  
  export default {
    name: 'MultiSelector',
    components: { Tooltip },
    props: {
      /**
       * The list of currently selected items (categories or tags).
       * Example: ['business', 'crime']
       */
      items: {
        type: Array,
        required: true,
      },
      /**
       * The maximum number of items a user can select.
       */
      maxItems: {
        type: Number,
        default: 5,
      },
      /**
       * The complete list of available items for selection.
       * Example: [{ value: 'business', text: 'Business' }, ...]
       */
      availableItems: {
        type: Array,
        required: true,
      },
      /**
       * The label for the selector (e.g., "Category (Optional)", "Tag (Optional)").
       */
      label: {
        type: String,
        required: true,
      },
      /**
       * The ID for the selector (used for accessibility and binding).
       */
      id: {
        type: String,
        required: true,
      },
      /**
       * The tooltip text providing additional information about the selector.
       */
      tooltip: {
        type: String,
        default: '',
      },
      /**
       * The ID for the tooltip (used for accessibility).
       */
      tooltipId: {
        type: String,
        default: '',
      },
      /**
       * The placeholder text for the dropdown.
       */
      placeholder: {
        type: String,
        default: 'Select an option',
      },
      /**
       * Optional mapping for item labels. If not provided, the component uses the `getItemLabel` method.
       */
      itemLabelMap: {
        type: Object,
        default: () => ({}),
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
       * Computes the list of items available for selection by excluding
       * the ones already selected.
       */
      availableOptions() {
        return this.availableItems.filter(
          (option) => !this.items.includes(option.value)
        );
      },
      /**
       * Determines if the dropdown should be disabled based on the maximum
       * number of selected items.
       */
      isDisabled() {
        return this.items.length >= this.maxItems;
      },
    },
    methods: {
      /**
       * Emits an event to add a selected item and resets the dropdown.
       */
      addItem() {
        if (this.selected && !this.isDisabled) {
          this.$emit('add-item', this.selected);
          this.selected = ''; // Reset the dropdown
        }
      },
      /**
       * Emits an event to remove an item based on its index.
       * @param {number} index - The index of the item to remove.
       */
      removeItem(index) {
        this.$emit('remove-item', index);
      },
      /**
       * Retrieves the display label for a given item value.
       * @param {string} value - The value of the item.
       * @returns {string} - The display label of the item.
       */
      getItemLabel(value) {
        return this.itemLabelMap[value] || value;
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
  
  .selected-items {
    margin-top: 10px;
    margin-bottom: 15px;
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
  }
  
  .item-tag {
    background-color: #17a2b8; /* Teal color for tags */
    color: white;
    padding: 5px 10px;
    border-radius: 15px;
    display: flex;
    align-items: center;
    font-size: 0.9em;
  }
  
  .item-tag button {
    background: none;
    border: none;
    color: white;
    margin-left: 8px;
    cursor: pointer;
    font-size: 1em;
    line-height: 1;
  }
  
  .item-tag button:hover {
    color: #ffdddd;
  }
  </style>  