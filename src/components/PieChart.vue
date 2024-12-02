<template>
  <canvas ref="pieChart" @click="handleClick"></canvas>
</template>

<script>
import { Chart } from 'chart.js'; // Import the full Chart.js
import ChartDataLabels from 'chartjs-plugin-datalabels'; // Import the data labels plugin

export default {
  props: {
    data: {
      type: Array,
      required: true,
    },
  },
  watch: {
    // Watch for changes in the `data` prop and re-render the chart when it changes
    data(newData) {
      this.renderChart(newData);
    },
  },
  mounted() {
    // Initially render the chart
    this.renderChart(this.data);
  },
  beforeDestroy() {
    // Clean up the chart instance when the component is destroyed
    if (this.chartInstance) {
      this.chartInstance.destroy();
    }
  },
  methods: {
    renderChart(data) {
      const ctx = this.$refs.pieChart.getContext('2d');
      Chart.register(ChartDataLabels);  // Register the data labels plugin

      // Destroy any existing chart instance before creating a new one
      if (this.chartInstance) {
        this.chartInstance.destroy();
      }

      const chartData = {
        labels: data.map(item => item.label),
        datasets: [{
          data: data.map(item => item.value),
          backgroundColor: data.map(item => item.color),
        }],
      };

      // Create the pie chart (use 'Chart' instead of 'Pie')
      this.chartInstance = new Chart(ctx, {
        type: 'pie',  // Use the 'pie' chart type here
        data: chartData,
        options: {
          responsive: true,
          maintainAspectRatio: false,
          plugins: {
            legend: {
              display: false,
            },
            tooltip: {
              callbacks: {
                label: (context) => {
                  return `${context.label}: ${context.raw}%`;
                },
              },
            },
            datalabels: {
              color: '#fff',
              formatter: (value, context) => {
                return `${context.chart.data.labels[context.dataIndex]}: ${value}%`;
              },
            },
          },
        },
      });
    },

    // Handle the click event on the pie chart
    handleClick(event) {
      // Get the chart instance's context to process the click event
      const activePoints = this.chartInstance.getElementsAtEventForMode(event, 'nearest', { intersect: true }, true);

      if (activePoints.length > 0) {
        // Get the clicked slice index (the index of the slice)
        const clickedIndex = activePoints[0].index;
        
        // Get the label of the clicked slice
        const clickedTopic = this.data[clickedIndex].label;

        // Emit the custom 'sliceClick' event with the clicked topic
        this.$emit('sliceClick', clickedTopic);
      }
    },
  },
};
</script>