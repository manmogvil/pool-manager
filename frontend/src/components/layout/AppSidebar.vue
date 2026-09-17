<script setup>
import { ref, computed, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useDisplay } from 'vuetify'
import { useAuth } from '../../composables/useAuth'

const router = useRouter()
const { user, isAdmin, logout: doLogout } = useAuth()
const { mobile } = useDisplay()

const drawer = ref(false)

const allItems = computed(() => [
  { title: 'Dashboard',     icon: 'mdi-home',           route: '/',              visible: true },
  { title: 'Users',         icon: 'mdi-account-group',  route: '/users',         visible: isAdmin.value },
  { title: 'Contributions', icon: 'mdi-wallet',         route: '/contributions', visible: true },
  { title: 'Games',         icon: 'mdi-ticket',         route: '/games',         visible: true },
  { title: 'Draws',         icon: 'mdi-calendar',       route: '/draws',         visible: true },
  { title: 'Tickets',       icon: 'mdi-book',           route: '/tickets',       visible: true },
])

const items = computed(() => allItems.value.filter(item => item.visible))

function logout() {
  doLogout()
  router.push('/login')
}

watch(() => router.currentRoute.value.path, () => {
  if (mobile.value) drawer.value = false
})
</script>

<template>
  <v-btn
    v-if="mobile"
    icon="mdi-menu"
    variant="text"
    size="small"
    class="menu-btn"
    @click="drawer = !drawer"
  />

  <v-navigation-drawer
    v-model="drawer"
    :rail="!mobile"
    :temporary="mobile"
    :permanent="!mobile"
    expand-on-hover
    width="250"
    class="sidebar"
  >
    <div class="sidebar-logo">
      <img src="/lottery_friends.png" alt="Lottery Friends" class="sidebar-logo-img"/>
    </div>

    <v-divider></v-divider>

    <v-list density="comfortable" nav class="sidebar-nav">
      <v-list-item
        v-for="item in items"
        :key="item.route"
        :prepend-icon="item.icon"
        :title="item.title"
        :to="item.route"
        rounded="lg"
        class="nav-item"
      ></v-list-item>
    </v-list>

    <template v-slot:append>
      <v-divider class="mb-2"></v-divider>
      <v-list density="comfortable" nav class="sidebar-nav">
        <v-list-item
          v-if="user"
          prepend-icon="mdi-account"
          :title="user.name"
          :subtitle="user.role"
          class="nav-item user-info"
        />
        <v-list-item
          prepend-icon="mdi-logout"
          title="Logout"
          @click="logout"
          class="nav-item logout-item"
        />
      </v-list>
    </template>
  </v-navigation-drawer>
</template>

<style scoped>
.sidebar {
  background: #ffffff !important;
  border-right: 1px solid #e5e7eb !important;
}

.menu-btn {
  position: fixed;
  top: 12px;
  left: 12px;
  z-index: 1000;
  background: white;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
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

.user-info {
  opacity: 0.8;
}

.logout-item {
  color: #ef4444 !important;
}
</style>
