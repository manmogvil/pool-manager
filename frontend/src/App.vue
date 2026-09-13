<script setup>
import AppSidebar from './components/layout/AppSidebar.vue'
</script>

<template>
  <v-app>
    <AppSidebar />
    <v-main class="main-content">
      <router-view />
    </v-main>
  </v-app>
</template>

<style scoped>
.main-content {
  background: #f8f9fa;
  margin-left: 56px;
}
</style>
