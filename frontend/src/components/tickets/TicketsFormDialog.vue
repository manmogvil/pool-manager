<script setup>
import { computed, ref, watch } from 'vue'

const props = defineProps({
  visible: Boolean,
  ticket: Object,
  draws: { type: Array, default: () => [] },
  games: { type: Array, default: () => [] }
})

const emit = defineEmits(['update:visible', 'saved'])

const isEditing = computed(() => !!props.ticket)

const form = ref({})
const formRef = ref(null)

function formatDate(dateStr) {
  if (!dateStr) return ''
  const d = new Date(dateStr)
  const day = String(d.getDate()).padStart(2, '0')
  const month = String(d.getMonth() + 1).padStart(2, '0')
  const year = d.getFullYear()
  return `${day}/${month}/${year}`
}

function getDrawLabel(draw) {
  const game = props.games.find(g => g.id === draw.game_id)
  const gameName = game ? game.name : 'Unknown'
  return `${gameName} #${draw.id} - ${formatDate(draw.draw_date)}`
}

const selectedGame = computed(() => {
  if (!form.value.draw_id) return null
  const draw = props.draws.find(d => d.id === form.value.draw_id)
  if (!draw) return null
  return props.games.find(g => g.id === draw.game_id) || null
})

const ticketCost = computed(() => {
  return selectedGame.value ? selectedGame.value.ticket_price : 0
})

watch(() => props.visible, (val) => {
  if (val) {
    if (props.ticket) {
      form.value = { ...props.ticket }
    } else {
      form.value = { draw_id: null, numbers: '', stars: '' }
    }
  }
})

function close() {
  emit('update:visible', false)
}

async function save() {
  const { valid } = await formRef.value.validate()
  if (!valid) return
  emit('saved', {
    ...form.value,
    cost: ticketCost.value,
    id: props.ticket?.id
  })
}
</script>

<template>
  <v-dialog :model-value="visible" @update:model-value="close" max-width="450" persistent>
    <v-card>
      <v-card-title class="text-h6">
        {{ isEditing ? 'Edit Ticket' : 'New Ticket' }}
      </v-card-title>
      <v-card-text>
        <v-form ref="formRef">
          <v-select
            v-model="form.draw_id"
            :items="draws"
            :item-title="(item) => getDrawLabel(item)"
            item-value="id"
            label="Draw"
            variant="outlined"
            :rules="[v => !!v || 'Selecting a draw is required']"
            class="mb-3"
          />
          <v-text-field
            v-model="form.numbers"
            label="Numbers"
            placeholder="e.g. 5,12,23,34,45"
            variant="outlined"
            :rules="[v => !!v || 'Numbers are required']"
            class="mb-3"
          />
          <v-text-field
            v-model="form.stars"
            label="Stars"
            placeholder="e.g. 3,7"
            variant="outlined"
            class="mb-3"
          />
          <v-text-field
            :model-value="ticketCost > 0 ? ticketCost.toFixed(2) + ' €' : 'Select a draw first'"
            label="Cost"
            variant="outlined"
            readonly
            :disabled="!form.draw_id"
            :hint="selectedGame ? `Price from ${selectedGame.name}` : ''"
            persistent-hint
          />
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
