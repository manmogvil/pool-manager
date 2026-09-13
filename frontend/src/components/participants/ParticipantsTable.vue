<script setup>
import { ref, onMounted } from 'vue'
import { participants as participantsService } from '../../services/api'
import ParticipantsFormDialog from './ParticipantsFormDialog.vue'

const participantList = ref([])
const loading = ref(false)
const dialogVisible = ref(false)
const selectedParticipant = ref(null)
const snackbar = ref({ show: false, message: '', color: 'success' })

async function loadParticipants() {
  loading.value = true
  try {
    participantList.value = await participantsService.getAll()
  } catch (error) {
    showSnackbar(error.message, 'error')
  } finally {
    loading.value = false
  }
}

function openCreateDialog() {
  selectedParticipant.value = null
  dialogVisible.value = true
}

function openEditDialog(participant) {
  selectedParticipant.value = participant
  dialogVisible.value = true
}

async function handleSaved(data) {
  try {
    if (data.id) {
      await participantsService.update(data.id, { name: data.name, email: data.email })
      showSnackbar('Participant updated', 'success')
    } else {
      await participantsService.create({ name: data.name, email: data.email })
      showSnackbar('Participant created', 'success')
    }
    dialogVisible.value = false
    await loadParticipants()
  } catch (error) {
    showSnackbar(error.message, 'error')
  }
}

async function toggleActive(participant, active) {
  try {
    if (active) {
      await participantsService.activate(participant.id)
      showSnackbar('Participant activated', 'success')
    } else {
      await participantsService.delete(participant.id)
      showSnackbar('Participant deactivated', 'success')
    }
  } catch (error) {
    showSnackbar(error.message, 'error')
    await loadParticipants()
  }
}

function showSnackbar(message, color) {
  snackbar.value = { show: true, message, color }
}

onMounted(() => {
  loadParticipants()
})
</script>

<template>
  <div>
    <div class="table-header">
      <v-btn color="primary" @click="openCreateDialog">
        <v-icon start>mdi-plus</v-icon>
        Add Participant
      </v-btn>
    </div>

    <v-card>
      <v-data-table
        :headers="[
          { title: 'ID', key: 'id', width: '80px' },
          { title: 'Name', key: 'name' },
          { title: 'Email', key: 'email' },
          { title: 'Status', key: 'active', width: '100px', sortable: true },
          { title: 'Actions', key: 'actions', width: '80px', sortable: false }
        ]"
        :items="participantList"
        :loading="loading"
        striped-rows
      >
        <template v-slot:item.active="{ item }">
          <v-switch
            v-model="item.active"
            color="success"
            density="compact"
            hide-details
            @update:model-value="(val) => toggleActive(item, val)"
          />
        </template>

        <template v-slot:item.actions="{ item }">
          <v-btn icon="mdi-pencil" variant="text" size="small" @click="openEditDialog(item)" />
        </template>
      </v-data-table>
    </v-card>

    <ParticipantsFormDialog
      v-model:visible="dialogVisible"
      :participant="selectedParticipant"
      @saved="handleSaved"
    />

    <v-snackbar v-model="snackbar.show" :color="snackbar.color" timeout="3000">
      {{ snackbar.message }}
    </v-snackbar>
  </div>
</template>

<style scoped>
.table-header {
  display: flex;
  justify-content: flex-end;
  margin-bottom: 16px;
}
</style>
