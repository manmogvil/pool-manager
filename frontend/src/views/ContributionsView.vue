<script setup>
import ContributionsTable from '../components/contributions/ContributionsTable.vue'
</script>

<template>
  <v-container fluid class="view-container">
    <h1 class="text-h5 font-weight-bold view-title">Contributions</h1>
    <ContributionsTable />
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
