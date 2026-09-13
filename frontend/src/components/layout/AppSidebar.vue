<script setup>
import { ref } from 'vue'

const items = [
  { title: 'Dashboard', icon: 'mdi-home', route: '/' },
  { title: 'Participants', icon: 'mdi-account-group', route: '/participants' },
  { title: 'Contributions', icon: 'mdi-wallet', route: '/contributions' },
  { title: 'Games', icon: 'mdi-ticket', route: '/games' },
  { title: 'Draws', icon: 'mdi-calendar', route: '/draws' },
  { title: 'Tickets', icon: 'mdi-book', route: '/tickets' }
]
const defaultPath = window.location.pathname;
const activeItem = ref(defaultPath)
</script>

<template>
  <v-navigation-drawer
    rail
    expand-on-hover
    width="250"
    class="sidebar"
  >
    <div class="sidebar-logo">
      <!-- quiero utilizar el logo de la aplicación lottery_fiends.png -->
      <img src="/lottery_friends.png" alt="Lottery Friends" class="sidebar-logo-img"/>

    </div>

    <v-divider></v-divider>

    <v-list density="comfortable" nav class="sidebar-nav">
      <v-list-item
        v-for="item in items"
        :key="item.route"
        :prepend-icon="item.icon"
        :title="item.title"
        :active="activeItem === item.route"
        @click="activeItem = item.route"
        :to="item.route"
        rounded="lg"
        class="nav-item"
      ></v-list-item>
    </v-list>
  </v-navigation-drawer>
</template>

<style scoped>
.sidebar {
  background: #ffffff !important;
  border-right: 1px solid #e5e7eb !important;
}

.sidebar-logo {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 56px;
  min-height: 56px;
}

.sidebar-nav {
  padding: 8px;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.nav-item {
  margin-bottom: 4px;
  width: 100%;
}

.nav-item.v-list-item--active {
  background: #eff6ff !important;
  color: #3b82f6 !important;
}

.sidebar-logo-img {
  max-height: 40px;
  max-width: 100%;
}
</style>
