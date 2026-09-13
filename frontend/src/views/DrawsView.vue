<script setup>
import DrawsTable from '../components/draws/DrawsTable.vue'
</script>

<template>
  <v-container fluid class="view-container">
    <h1 class="text-h5 font-weight-bold view-title">Draws</h1>
    <DrawsTable />
  </v-container>
</template>

<style scoped>
.view-container {
  padding: 24px;
}

.view-title {
  margin-bottom: 24px;
}
</style>
