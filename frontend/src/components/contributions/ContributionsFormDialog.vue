<script setup>
import { computed, ref, watch } from 'vue'

const props = defineProps({
  visible: Boolean,
  contribution: Object,
  participants: { type: Array, default: () => [] },
  games: { type: Array, default: () => [] }
})

const emit = defineEmits(['update:visible', 'saved'])

const isEditing = computed(() => !!props.contribution)

const monthNames = [
  'January', 'February', 'March', 'April', 'May', 'June',
  'July', 'August', 'September', 'October', 'November', 'December'
]

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
    if (props.contribution) {
      const isoDate = formatToISO(props.contribution.payment_date)
      form.value = { ...props.contribution, payment_date: isoDate }
      datePickerValue.value = isoDate ? new Date(isoDate + 'T00:00:00') : null
    } else {
      const today = formatToISO(new Date().toISOString())
      form.value = { participant_id: null, game_id: null, month: new Date().getMonth() + 1, year: new Date().getFullYear(), amount: null, paid: true, payment_date: today, payment_method: '', comments: '' }
      datePickerValue.value = new Date()
    }
  }
})

function onDateSelect(date) {
  const d = new Date(date)
  const day = String(d.getDate()).padStart(2, '0')
  const month = String(d.getMonth() + 1).padStart(2, '0')
  const year = d.getFullYear()
  form.value.payment_date = `${year}-${month}-${day}`
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
    id: props.contribution?.id,
    payment_date: form.value.payment_date || null,
    payment_method: form.value.payment_method || null
  })
}
</script>

<template>
  <v-dialog :model-value="visible" @update:model-value="close" max-width="550" persistent>
    <v-card>
      <v-card-title class="text-h6">
        {{ isEditing ? 'Edit Contribution' : 'New Contribution' }}
      </v-card-title>
      <v-card-text>
        <v-form ref="formRef">
          <v-select
            v-model="form.participant_id"
            :items="participants"
            item-title="name"
            item-value="id"
            label="Participant"
            variant="outlined"
            :rules="[v => !!v || 'Selecting a participant is required']"
            class="mb-3"
          />
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
          <v-row>
            <v-col cols="6">
              <v-select
                v-model="form.month"
                :items="monthNames.map((name, i) => ({ title: name, value: i + 1 }))"
                item-title="title"
                item-value="value"
                label="Month"
                variant="outlined"
                :rules="[v => !!v || 'Month is required']"
              />
            </v-col>
            <v-col cols="6">
              <v-text-field
                v-model.number="form.year"
                label="Year"
                type="number"
                variant="outlined"
                :rules="[
                  v => !!v || 'Year is required',
                  v => v >= 2020 || 'Year must be 2020 or later'
                ]"
              />
            </v-col>
          </v-row>
          <v-text-field
            v-model.number="form.amount"
            label="Amount"
            type="number"
            prefix="€"
            variant="outlined"
            :rules="[
              v => !!v || 'Amount is required',
              v => v > 0 || 'Amount must be greater than 0'
            ]"
            class="mb-3"
          />
          <v-select
            v-model="form.payment_method"
            :items="['CASH', 'BIZUM']"
            label="Payment Method"
            variant="outlined"
            clearable
            class="mb-3"
          />
          <v-menu v-model="dateMenu" :close-on-content-click="false" location="bottom">
            <template v-slot:activator="{ props }">
              <v-text-field
                :model-value="formatDisplay(form.payment_date)"
                label="Payment Date"
                variant="outlined"
                readonly
                class="mb-3"
                v-bind="props"
              />
            </template>
            <v-date-picker
              v-model="datePickerValue"
              @update:model-value="onDateSelect"
              color="primary"
            />
          </v-menu>
          <v-textarea
            v-model="form.comments"
            label="Comments"
            variant="outlined"
            rows="2"
            auto-grow
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
