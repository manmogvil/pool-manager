<template>
  <v-card elevation="0" rounded="lg" class="section-card">
    <v-card-title class="d-flex align-center pa-4 pb-0">
      <v-icon icon="mdi-chart-line" color="success" class="mr-2" />
      <span class="text-subtitle-1 font-weight-bold">Monthly Trend</span>
    </v-card-title>
    <v-card-text class="pa-4">
      <div v-if="months.length === 0" class="text-center text-grey py-4">
        No data available
      </div>
      <div v-else class="chart-container">
        <div class="chart-bars d-flex align-end justify-space-between" style="height: 120px;">
          <div
            v-for="m in months"
            :key="m.label + m.year"
            class="bar-wrapper d-flex flex-column align-center"
          >
            <div class="text-caption font-weight-medium mb-1">€{{ m.total.toFixed(0) }}</div>
            <div
              class="bar"
              :style="{ height: getBarHeight(m.total) + '%' }"
            />
            <div class="text-caption text-grey mt-1">{{ m.label }}</div>
          </div>
        </div>
      </div>
    </v-card-text>
  </v-card>
</template>

<script setup>
const props = defineProps({
  months: { type: Array, default: () => [] }
})

const getBarHeight = (total) => {
  if (props.months.length === 0) return 0
  const max = Math.max(...props.months.map(m => m.total))
  return max > 0 ? (total / max) * 100 : 0
}
</script>

<style scoped>
.section-card {
  border: 1px solid #e8e8e8;
}
.bar {
  width: 100%;
  max-width: 40px;
  background: linear-gradient(180deg, #3b82f6 0%, #60a5fa 100%);
  border-radius: 4px 4px 0 0;
  min-height: 4px;
  transition: height 0.3s ease;
}
.bar-wrapper {
  flex: 1;
  min-width: 40px;
}
@media (max-width: 599px) {
  .bar {
    max-width: 28px;
  }
  .bar-wrapper {
    min-width: 32px;
  }
}
</style>
