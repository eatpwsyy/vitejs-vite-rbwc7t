<!-- Tooltip.vue -->
<template>
    <div class="tooltip-container">
      <span class="tooltip-icon" :aria-describedby="tooltipId">ⓘ</span>
      <div class="tooltip-content" :id="tooltipId" role="tooltip">
        <slot></slot>
      </div>
    </div>
  </template>
  
  <script>
  export default {
    name: 'Tooltip',
    props: {
      tooltipId: {
        type: String,
        required: true,
      },
    },
  };
  </script>
  
  <style scoped>
  .tooltip-container {
    position: relative;
    display: inline-block;
    margin-left: 4px;
  }
  
  .tooltip-icon {
    color: #007bff;
    cursor: help;
    font-size: 0.9em;
  }
  
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
    bottom: 100%;
    left: 50%;
    transform: translateX(-50%);
    opacity: 0;
    transition: opacity 0.3s;
  }
  
  .tooltip-content::after {
    content: "";
    position: absolute;
    top: 100%;
    left: 50%;
    transform: translateX(-50%);
    border-width: 2px;
    border-style: solid;
    border-color: #333 transparent transparent transparent;
  }
  
  .tooltip-container:hover .tooltip-content {
    visibility: visible;
    opacity: 1;
  }
  </style>  