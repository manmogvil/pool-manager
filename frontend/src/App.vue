<script setup>
import { useAuth } from './composables/useAuth'
import AppSidebar from './components/layout/AppSidebar.vue'

const { isAuthenticated } = useAuth()
</script>

<template>
  <v-app>
    <AppSidebar v-if="isAuthenticated" />
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
