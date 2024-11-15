<template>
    <div class="heatmap">
      <h5>Heatmap</h5>
      <div id="map" class="map-container"></div>
    </div>
  </template>
  
  <script>
  import { onMounted } from 'vue';
  import L from 'leaflet';
  import 'leaflet/dist/leaflet.css';
  import 'leaflet.heat'; // Importing the heatmap plugin
  
  export default {
    name: 'Heatmap',
    setup() {
      onMounted(() => {
        // Initialize the map
        const map = L.map('map').setView([20, 0], 2); // Coordinates for global view
  
        // Set up the tile layer from OpenStreetMap
        L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
          attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
        }).addTo(map);
  
        // Dummy data points (latitude, longitude, intensity)
        const heatData = [
          [37.7749, -122.4194, 0.5],  // San Francisco
          [40.7128, -74.0060, 0.6],   // New York
          [51.5074, -0.1278, 0.8],    // London
          [35.6895, 139.6917, 0.4],   // Tokyo
          [-33.8688, 151.2093, 0.7],  // Sydney
        ];
  
        // Adding heat layer to the map
        L.heatLayer(heatData, {
          radius: 25,
          blur: 15,
          maxZoom: 10,
        }).addTo(map);
      });
    }
  }
  </script>
  
  <style scoped>
  .map-container {
    height: 400px;
    width: 100%;
    margin-top: 20px;
  }
  .heatmap {
    padding: 20px;
    border: 1px solid #ccc;
    text-align: center;
  }
  </style>
  