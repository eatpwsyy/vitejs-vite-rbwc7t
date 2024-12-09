<template>
  <div>
    <canvas ref="pieChart"></canvas>
    <!-- Optionally, display a message if there's no data -->
    <div v-if="!hasValidData" class="no-data">
      No data available to display the pie chart.
    </div>
  </div>
</template>

<script>
import { Chart, registerables } from 'chart.js';
import ChartDataLabels from 'chartjs-plugin-datalabels';

// Register Chart.js components and the datalabels plugin
Chart.register(...registerables, ChartDataLabels);

export default {
  props: {
    data: {
      type: Array,
      required: true,
      default: () => [],
    },
  },
  data() {
    return {
      chart: null,
    };
  },
  computed: {
    /**
     * Checks if the provided data is a non-empty array.
     */
    hasValidData() {
      return Array.isArray(this.data) && this.data.length > 0;
    },
  },
  methods: {
    /**
     * Renders the pie chart using Chart.js.
     */
    renderChart() {
      if (!this.hasValidData) {
        console.warn('PieChart expects a non-empty array as data prop.');
        return;
      }

      const ctx = this.$refs.pieChart.getContext('2d');

      // Prepare labels and datasets
      const labels = this.data.map((item) => item.label);
      const values = this.data.map((item) => item.value);
      const backgroundColors = this.data.map((item) => item.color);

      // Extract keywords as strings for each slice (used in tooltips)
      const keywords = this.data.map((item) => item.keywords.join(', '));

      // Create the chart
      this.chart = new Chart(ctx, {
        type: 'pie',
        data: {
          labels: labels,
          datasets: [
            {
              data: values,
              backgroundColor: backgroundColors,
              borderWidth: 1,
              borderColor: '#fff', // Default border color
            },
          ],
        },
        options: {
          responsive: true,
          plugins: {
            legend: {
              display: false, // Hide legend if not needed
            },
            tooltip: {
              callbacks: {
                label: function (context) {
                  const keyword = keywords[context.dataIndex];
                  const value = context.parsed;
                  return `${keyword}: ${value.toFixed(2)}%`;
                },
              },
            },
            datalabels: {
              color: '#fff',
              formatter: function(value, context) {
                if (value > 5) {
                  return `${context.chart.data.labels[context.dataIndex]}: ${value.toFixed(2)}%`;
                }
                return '';
              },
              font: {
                weight: 'bold',
                size: 14,
              },
              align: 'center',
            },
          },
          onClick: (event, elements) => {
            if (elements.length > 0) {
              const chartElement = elements[0];
              const index = chartElement.index;
              const clickedLabel = labels[index];
              this.$emit('sliceClick', clickedLabel);
            }
          },
        },
      });
    },
  },
  mounted() {
    this.renderChart();
  },
  watch: {
    /**
     * Watches for changes in the data prop and updates the chart accordingly.
     */
    data: {
      handler(newData) {
        if (this.chart) {
          // Update the chart data
          const labels = newData.map((item) => item.label);
          const values = newData.map((item) => item.value);
          const backgroundColors = newData.map((item) => item.color);

          this.chart.data.labels = labels;
          this.chart.data.datasets[0].data = values;
          this.chart.data.datasets[0].backgroundColor = backgroundColors;

          // Update the chart
          this.chart.update();
        } else {
          this.renderChart();
        }
      },
      deep: true,
      immediate: true,
    },
  },
};
</script>

<style scoped>
.no-data {
  text-align: center;
  font-size: 1.2rem;
  color: #6c757d;
  margin-top: 20px;
}
</style>