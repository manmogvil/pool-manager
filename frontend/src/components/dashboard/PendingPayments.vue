<template>
  <v-card elevation="0" rounded="lg" class="section-card">
    <v-card-title class="d-flex align-center pa-4 pb-0">
      <v-icon icon="mdi-alert-circle-outline" color="warning" class="mr-2" />
      <span class="text-subtitle-1 font-weight-bold">Pending Payments</span>
      <v-spacer />
      <v-chip v-if="payments.length" color="warning" size="small" variant="tonal">
        {{ payments.length }}
      </v-chip>
    </v-card-title>
    <v-card-text class="pa-0">
      <v-table v-if="payments.length" density="compact">
        <thead>
          <tr>
            <th>User</th>
            <th>Game</th>
            <th>Period</th>
            <th class="text-right">Amount</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(p, i) in payments" :key="i">
            <td>{{ p.user_name }}</td>
            <td>{{ p.game_name }}</td>
            <td>{{ monthNames[p.month - 1] }} {{ p.year }}</td>
            <td class="text-right font-weight-medium">€{{ p.amount.toFixed(2) }}</td>
          </tr>
        </tbody>
      </v-table>
      <div v-else class="pa-6 text-center text-grey">
        <v-icon icon="mdi-check-circle-outline" size="48" color="success" class="mb-2" />
        <div>All payments up to date!</div>
      </div>
    </v-card-text>
  </v-card>
</template>

<script setup>
defineProps({
  payments: { type: Array, default: () => [] }
})

const monthNames = ['January', 'February', 'March', 'April', 'May', 'June',
  'July', 'August', 'September', 'October', 'November', 'December']
</script>

<style scoped>
.section-card {
  border: 1px solid #e8e8e8;
}
</style>
