<template>
  <v-container fluid class="view-container">
    <h1 class="text-h5 font-weight-bold view-title">Dashboard</h1>

    <v-progress-circular v-if="loading" indeterminate color="primary" class="d-block mx-auto my-8" />

    <template v-else>
      <!-- Summary Cards -->
      <v-row class="mb-4">
        <v-col cols="12" sm="6" md="3">
          <StatCard
            title="Active Users"
            :value="stats.active_users"
            icon="mdi-account-group"
            color="primary"
          />
        </v-col>
        <v-col cols="12" sm="6" md="3">
          <StatCard
            title="Total Collected"
            :value="'€' + stats.total_contributions.toFixed(2)"
            icon="mdi-cash-multiple"
            color="success"
          />
        </v-col>
        <v-col cols="12" sm="6" md="3">
          <StatCard
            title="Total Spent"
            :value="'€' + stats.total_spent.toFixed(2)"
            icon="mdi-cash-minus"
            color="error"
          />
        </v-col>
        <v-col cols="12" sm="6" md="3">
          <StatCard
            title="Prizes Won"
            :value="'€' + stats.total_prizes.toFixed(2)"
            icon="mdi-trophy"
            color="amber"
          />
        </v-col>
      </v-row>

      <v-row class="mb-4">
        <v-col cols="12" sm="6" md="3">
          <StatCard
            title="Active Games"
            :value="stats.active_games"
            icon="mdi-ticket"
            color="info"
          />
        </v-col>
        <v-col cols="12" sm="6" md="3">
          <StatCard
            title="Pending Draws"
            :value="stats.pending_draws"
            icon="mdi-calendar-clock"
            color="warning"
          />
        </v-col>
        <v-col cols="12" sm="6" md="3">
          <StatCard
            title="Total Tickets"
            :value="stats.total_tickets"
            icon="mdi-ticket-confirmation"
            color="teal"
          />
        </v-col>
        <v-col cols="12" sm="6" md="3">
          <StatCard
            title="Net Balance"
            :value="netBalance"
            icon="mdi-scale-balance"
            :color="netColor"
          />
        </v-col>
      </v-row>

      <!-- Pending Payments -->
      <v-row class="mb-4">
        <v-col cols="12">
          <PendingPayments :payments="stats.pending_payments" />
        </v-col>
      </v-row>

      <!-- Contributions by Game + Monthly Trend -->
      <v-row>
        <v-col cols="12" md="6">
          <ContributionsByGame :games="stats.contributions_by_game" />
        </v-col>
        <v-col cols="12" md="6">
          <MonthlyTrend :months="stats.monthly_trend" />
        </v-col>
      </v-row>
    </template>
  </v-container>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { dashboard } from '../services/api'
import StatCard from '../components/dashboard/StatCard.vue'
import PendingPayments from '../components/dashboard/PendingPayments.vue'
import ContributionsByGame from '../components/dashboard/ContributionsByGame.vue'
import MonthlyTrend from '../components/dashboard/MonthlyTrend.vue'

const loading = ref(true)
const stats = ref({
  total_users: 0,
  active_users: 0,
  total_games: 0,
  active_games: 0,
  total_draws: 0,
  pending_draws: 0,
  total_tickets: 0,
  total_contributions: 0,
  total_spent: 0,
  total_prizes: 0,
  pending_payments: [],
  contributions_by_game: [],
  monthly_trend: []
})

const netBalance = computed(() => {
  const balance = stats.value.total_prizes - stats.value.total_spent
  return (balance >= 0 ? '+' : '') + '€' + balance.toFixed(2)
})

const netColor = computed(() => {
  return stats.value.total_prizes >= stats.value.total_spent ? 'success' : 'error'
})

onMounted(async () => {
  try {
    const data = await dashboard.getStats()
    stats.value = data
  } catch (err) {
    console.error('Failed to load dashboard:', err)
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.view-title {
  margin-bottom: 24px;
}
</style>
