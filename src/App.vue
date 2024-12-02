<template>
  <div class="container-fluid">
    <div class="row">
      <!-- This is where routed components will be injected -->
      <router-view />
    </div>
  </div>
</template>

<script>
export default {
  name: 'App',
};
</script>