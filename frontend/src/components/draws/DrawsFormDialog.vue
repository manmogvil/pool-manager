<script setup>
import { computed, ref, watch } from 'vue'

const props = defineProps({
  visible: Boolean,
  draw: Object,
  games: { type: Array, default: () => [] }
})

const emit = defineEmits(['update:visible', 'saved'])

const isEditing = computed(() => !!props.draw)

const form = ref({})
const formRef = ref(null)
const dateMenu = ref(false)
const datePickerValue = ref(null)

function formatToISO(dateStr) {
  if (!dateStr) return ''
  const d = new Date(dateStr)
  const day = String(d.getDate()).padStart(2, '0')
  const month = String(d.getMonth() + 1).padStart(2, '0')
  const year = d.getFullYear()
  return `${year}-${month}-${day}`
}

function formatDisplay(dateStr) {
  if (!dateStr) return ''
  const parts = dateStr.split('-')
  return `${parts[2]}/${parts[1]}/${parts[0]}`
}

watch(() => props.visible, (val) => {
  if (val) {
    if (props.draw) {
      const isoDate = formatToISO(props.draw.draw_date)
      form.value = { ...props.draw, draw_date: isoDate }
      datePickerValue.value = isoDate ? new Date(isoDate + 'T00:00:00') : null
    } else {
      const today = formatToISO(new Date().toISOString())
      form.value = { game_id: null, draw_date: today }
      datePickerValue.value = new Date()
    }
  }
})

function onDateSelect(date) {
  const d = new Date(date)
  const day = String(d.getDate()).padStart(2, '0')
  const month = String(d.getMonth() + 1).padStart(2, '0')
  const year = d.getFullYear()
  form.value.draw_date = `${year}-${month}-${day}`
  dateMenu.value = false
}

function close() {
  emit('update:visible', false)
}

async function save() {
  const { valid } = await formRef.value.validate()
  if (!valid) return
  emit('saved', {
    ...form.value,
    id: props.draw?.id
  })
}
</script>

<template>
  <v-dialog :model-value="visible" @update:model-value="close" max-width="450" persistent>
    <v-card>
      <v-card-title class="text-h6">
        {{ isEditing ? 'Edit Draw' : 'New Draw' }}
      </v-card-title>
      <v-card-text>
        <v-form ref="formRef">
          <v-select
            v-model="form.game_id"
            :items="games"
            item-title="name"
            item-value="id"
            label="Game"
            variant="outlined"
            :rules="[v => !!v || 'Selecting a game is required']"
            class="mb-3"
          />
          <v-menu v-model="dateMenu" :close-on-content-click="false" location="bottom">
            <template v-slot:activator="{ props }">
              <v-text-field
                :model-value="formatDisplay(form.draw_date)"
                label="Draw Date"
                variant="outlined"
                readonly
                v-bind="props"
              />
            </template>
            <v-date-picker
              v-model="datePickerValue"
              @update:model-value="onDateSelect"
              color="primary"
            />
          </v-menu>
        </v-form>
      </v-card-text>
      <v-card-actions>
        <v-spacer />
        <v-btn variant="text" @click="close">Cancel</v-btn>
        <v-btn color="primary" @click="save">Save</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>
