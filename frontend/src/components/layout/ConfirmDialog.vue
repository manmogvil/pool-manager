<script setup>
const props = defineProps({
  visible: Boolean,
  title: { type: String, default: 'Confirm' },
  message: { type: String, default: 'Are you sure?' },
  confirmText: { type: String, default: 'Delete' },
  color: { type: String, default: 'error' }
})

const emit = defineEmits(['update:visible', 'confirm'])

const iconMap = {
  error: 'mdi-alert-circle',
  warning: 'mdi-alert',
  info: 'mdi-information',
  success: 'mdi-check-circle'
}

function close() {
  emit('update:visible', false)
}
</script>

<template>
  <v-dialog :model-value="visible" @update:model-value="close" max-width="420" persistent>
    <v-card rounded="lg" class="confirm-card">
      <v-card-text class="confirm-body">
        <v-avatar :color="color" size="56" class="mb-4">
          <v-icon :icon="iconMap[color] || iconMap.error" size="32" color="white" />
        </v-avatar>
        <h3 class="text-h6 font-weight-bold mb-2">{{ title }}</h3>
        <p class="text-body-2 text-medium-emphasis">{{ message }}</p>
      </v-card-text>
      <v-divider />
      <v-card-actions class="confirm-actions">
        <v-btn variant="outlined" rounded="lg" @click="close">Cancel</v-btn>
        <v-btn :color="color" variant="flat" rounded="lg" @click="emit('confirm')">{{ confirmText }}</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<style scoped>
.confirm-card {
  overflow: hidden;
}
.confirm-body {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  padding: 32px 24px 24px;
}
.confirm-actions {
  padding: 12px 16px;
  justify-content: center;
  gap: 12px;
}
</style>
