<template>
  <v-card elevation="0" rounded="lg" class="section-card">
    <v-card-title class="d-flex align-center pa-4 pb-0">
      <v-icon icon="mdi-chart-bar" color="primary" class="mr-2" />
      <span class="text-subtitle-1 font-weight-bold">Contributions by Game</span>
    </v-card-title>
    <v-card-text class="pa-4">
      <div v-if="games.length === 0" class="text-center text-grey py-4">
        No data available
      </div>
      <div v-else>
        <div v-for="game in games" :key="game.game_id" class="game-row mb-3">
          <div class="d-flex justify-space-between align-center mb-1">
            <span class="text-body-2 font-weight-medium">{{ game.game_name }}</span>
            <span class="text-body-2 text-grey">€{{ game.total.toFixed(2) }}</span>
          </div>
          <v-progress-linear
            :model-value="getPercentage(game.total)"
            color="primary"
            rounded
            height="8"
          />
          <div class="text-caption text-grey mt-1">{{ game.count }} contributions</div>
        </div>
      </div>
    </v-card-text>
  </v-card>
</template>

<script setup>
const props = defineProps({
  games: { type: Array, default: () => [] }
})

const getPercentage = (total) => {
  if (props.games.length === 0) return 0
  const max = Math.max(...props.games.map(g => g.total))
  return max > 0 ? (total / max) * 100 : 0
}
</script>

<style scoped>
.section-card {
  border: 1px solid #e8e8e8;
}
.game-row {
  padding-bottom: 8px;
  border-bottom: 1px solid #f0f0f0;
}
.game-row:last-child {
  border-bottom: none;
  padding-bottom: 0;
}
</style>
