<script setup>
import UsersTable from '../components/users/UsersTable.vue'
</script>

<template>
  <v-container fluid class="view-container">
    <h1 class="text-h5 font-weight-bold view-title">Users</h1>
    <UsersTable />
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
