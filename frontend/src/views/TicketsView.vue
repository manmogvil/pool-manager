<script setup>
import TicketsTable from '../components/tickets/TicketsTable.vue'
</script>

<template>
  <v-container fluid class="view-container">
    <h1 class="text-h5 font-weight-bold view-title">Tickets</h1>
    <TicketsTable />
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
