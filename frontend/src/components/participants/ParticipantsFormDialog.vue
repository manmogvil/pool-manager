<script setup>
import { computed, ref, watch } from 'vue'

const props = defineProps({
  visible: Boolean,
  participant: Object
})

const emit = defineEmits(['update:visible', 'saved'])

const isEditing = computed(() => !!props.participant)

const form = ref({ name: '', email: '' })
const formRef = ref(null)

watch(() => props.visible, (val) => {
  if (val) {
    form.value = props.participant
      ? { name: props.participant.name, email: props.participant.email }
      : { name: '', email: '' }
  }
})

function close() {
  emit('update:visible', false)
}

async function save() {
  const { valid } = await formRef.value.validate()
  if (!valid) return
  emit('saved', { ...form.value, id: props.participant?.id })
}
</script>

<template>
  <v-dialog :model-value="visible" @update:model-value="close" max-width="400" persistent>
    <v-card>
      <v-card-title class="text-h6">
        {{ isEditing ? 'Edit Participant' : 'New Participant' }}
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
            v-model="form.email"
            label="Email"
            variant="outlined"
            :rules="[
              v => !!v || 'Email is required',
              v => /.+@.+\..+/.test(v) || 'Please enter a valid email'
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
