<!-- FormField.vue -->
<template>
    <div class="form-group">
      <label :for="id">
        {{ label }}
        <span v-if="required" class="required">*</span>
        <Tooltip v-if="tooltip" :tooltip-id="tooltipId">{{ tooltip }}</Tooltip>
      </label>
      <slot name="input"></slot>
      <small v-if="hint">{{ hint }}</small>
    </div>
  </template>
  
  <script>
  import Tooltip from './Tooltip.vue';
  
  export default {
    name: 'FormField',
    components: { Tooltip },
    props: {
      label: {
        type: String,
        required: true,
      },
      id: {
        type: String,
        required: true,
      },
      required: {
        type: Boolean,
        default: false,
      },
      tooltip: {
        type: String,
        default: '',
      },
      tooltipId: {
        type: String,
        default: '',
      },
      hint: {
        type: String,
        default: '',
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
  
  .required {
    color: red;
    margin-left: 4px;
  }
  
  .form-group small {
    color: #6c757d;
    margin-top: 5px;
  }
  </style>  