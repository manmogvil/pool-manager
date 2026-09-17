<template>
  <v-card class="stat-card" elevation="0" rounded="lg">
    <v-card-text class="d-flex align-center pa-4">
      <v-avatar :color="color" size="48" class="mr-4">
        <v-icon :icon="icon" color="white" size="24" />
      </v-avatar>
      <div class="stat-value-wrapper">
        <div class="text-caption text-grey mb-1">{{ title }}</div>
        <div class="stat-value text-h5 font-weight-bold">{{ value }}</div>
      </div>
    </v-card-text>
  </v-card>
</template>

<script setup>
defineProps({
  title: { type: String, required: true },
  value: { type: [String, Number], required: true },
  icon: { type: String, required: true },
  color: { type: String, default: 'primary' }
})
</script>

<style scoped>
.stat-card {
  border: 1px solid #e8e8e8;
  transition: box-shadow 0.2s;
}
.stat-card:hover {
  box-shadow: 0 2px 8px rgba(0,0,0,0.08);
}
.stat-value-wrapper {
  overflow: hidden;
  min-width: 0;
}
.stat-value {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
@media (max-width: 599px) {
  .stat-value {
    font-size: 1.1rem !important;
  }
}
</style>
