<script setup>
import { computed, ref, watch } from 'vue'

const props = defineProps({
  visible: Boolean,
  game: Object
})

const emit = defineEmits(['update:visible', 'saved'])

const isEditing = computed(() => !!props.game)

const form = ref({ name: '', draw_days: '', ticket_price: null })
const formRef = ref(null)

watch(() => props.visible, (val) => {
  if (val) {
    form.value = props.game
      ? { name: props.game.name, draw_days: props.game.draw_days, ticket_price: props.game.ticket_price }
      : { name: '', draw_days: '', ticket_price: null }
  }
})

function close() {
  emit('update:visible', false)
}

async function save() {
  const { valid } = await formRef.value.validate()
  if (!valid) return
  emit('saved', { ...form.value, id: props.game?.id, active: props.game?.active })
}
</script>

<template>
  <v-dialog :model-value="visible" @update:model-value="close" max-width="450" persistent>
    <v-card>
      <v-card-title class="text-h6">
        {{ isEditing ? 'Edit Game' : 'New Game' }}
      </v-card-title>
      <v-card-text>
        <v-form ref="formRef">
          <v-text-field
            v-model="form.name"
            label="Name"
            variant="outlined"
            :rules="[v => !!v || 'Name is required']"
            class="mb-3"
          />
          <v-text-field
            v-model="form.draw_days"
            label="Draw Days"
            placeholder="e.g. Tuesday, Friday"
            variant="outlined"
            :rules="[v => !!v || 'Draw days are required']"
            class="mb-3"
          />
          <v-text-field
            v-model.number="form.ticket_price"
            label="Ticket Price"
            type="number"
            prefix="€"
            variant="outlined"
            :rules="[
              v => !!v || 'Ticket price is required',
              v => v > 0 || 'Price must be greater than 0'
            ]"
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
