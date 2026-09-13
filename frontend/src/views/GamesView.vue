<script setup>
import GamesTable from '../components/games/GamesTable.vue'
</script>

<template>
  <v-container fluid class="view-container">
    <h1 class="text-h5 font-weight-bold view-title">Games</h1>
    <GamesTable />
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
