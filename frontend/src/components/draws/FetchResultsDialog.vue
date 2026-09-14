<script setup>
import { ref } from 'vue'
import { draws as drawsService } from '../../services/api'

const props = defineProps({
  visible: Boolean
})

const emit = defineEmits(['update:visible', 'saved'])

const loading = ref(false)
const form = ref({ game_slug: '', from: '', to: '' })
const lotteryGamesAPI = [
  { title: 'EuroMillones'             , value: 'euromillones' },
  { title: 'La Primitiva'             , value: 'primitiva'    },
  { title: 'Bonoloto'                 , value: 'bonoloto'     },
  { title: 'El Gordo de la Primitiva' , value: 'gordo'        },
  { title: 'Lotería Nacional'         , value: 'nacional'     },
  { title: 'Eurodreams'               , value: 'eurodreams'   }
]


function close() {
  emit('update:visible', false)
}

async function handleFetch() {
  if (!form.value.game_slug || !form.value.from || !form.value.to) return
  loading.value = true
  try {
    const result = await drawsService.fetchResults(form.value)
    emit('saved', result)
    close()
  } catch (error) {
    emit('saved', { error: error.message })
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <v-dialog :model-value="visible" @update:model-value="emit('update:visible', $event)" max-width="450">
    <v-card>
      <v-card-title class="text-h6">Fetch Results from API</v-card-title>
      <v-card-text>
        <v-select
          v-model="form.game_slug"
          :items="lotteryGamesAPI"
          label="Game"
          class="mb-2"
        />
        <v-text-field
          v-model="form.from"
          label="From"
          type="date"
          class="mb-2"
          :formatted-value="form.from ? new Date(form.from + 'T00:00:00').toLocaleDateString('es-ES') : ''"
        />
        <v-text-field
          v-model="form.to"
          label="To"
          type="date"
          :formatted-value="form.to ? new Date(form.to + 'T00:00:00').toLocaleDateString('es-ES') : ''"
        />
      </v-card-text>
      <v-card-actions>
        <v-spacer />
        <v-btn @click="close">Cancel</v-btn>
        <v-btn color="primary" @click="handleFetch" :loading="loading">Fetch</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>
