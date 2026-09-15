<script setup>
import { ref, onMounted } from 'vue'
import { draws as drawsService, games as gamesService } from '../../services/api'
import { useAuth } from '../../composables/useAuth'
import DrawsFormDialog from './DrawsFormDialog.vue'
import FetchResultsDialog from './FetchResultsDialog.vue'
import ConfirmDialog from '../layout/ConfirmDialog.vue'

const { isAdmin } = useAuth()

const drawList = ref([])
const gameList = ref([])
const loading = ref(false)
const dialogVisible = ref(false)
const selectedDraw = ref(null)
const snackbar = ref({ show: false, message: '', color: 'success' })

const fetchDialogVisible = ref(false)
const confirmVisible = ref(false)
const drawToDelete = ref(null)

async function loadData() {
  loading.value = true
  try {
    const [d, g] = await Promise.all([drawsService.getAll(), gamesService.getAll()])
    drawList.value = d || []
    gameList.value = g || []
  } catch (error) {
    showSnackbar(error.message, 'error')
  } finally {
    loading.value = false
  }
}

function getGameName(id) {
  const g = gameList.value.find(g => g.id === id)
  return g ? g.name : '-'
}

function formatDate(dateStr) {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  const day = String(d.getDate()).padStart(2, '0')
  const month = String(d.getMonth() + 1).padStart(2, '0')
  const year = d.getFullYear()
  return `${day}/${month}/${year}`
}

function openCreateDialog() {
  selectedDraw.value = null
  dialogVisible.value = true
}

function openEditDialog(draw) {
  selectedDraw.value = draw
  dialogVisible.value = true
}

function openFetchDialog() {
  fetchDialogVisible.value = true
}

function handleFetchSaved(result) {
  if (result.error) {
    showSnackbar(result.error, 'error')
  } else {
    showSnackbar(result.message || 'Results fetched', 'success')
    loadData()
  }
}

async function handleSaved(data) {
  try {
    if (data.id) {
      await drawsService.updateResults(data.id, {
        result_numbers: data.result_numbers || null,
        result_stars: data.result_stars || null
      })
      showSnackbar('Draw updated', 'success')
    } else {
      await drawsService.create({
        game_id: data.game_id || null,
        draw_date: data.draw_date
      })
      showSnackbar('Draw created', 'success')
    }
    dialogVisible.value = false
    await loadData()
  } catch (error) {
    showSnackbar(error.message, 'error')
  }
}

function confirmDelete(draw) {
  drawToDelete.value = draw
  confirmVisible.value = true
}

async function handleConfirmDelete() {
  try {
    await drawsService.delete(drawToDelete.value.id)
    showSnackbar('Draw deleted', 'success')
    await loadData()
  } catch (error) {
    showSnackbar(error.message, 'error')
  }
  confirmVisible.value = false
}

function showSnackbar(message, color) {
  snackbar.value = { show: true, message, color }
}

onMounted(() => {
  loadData()
})
</script>

<template>
  <div>
    <div class="table-header">
      <v-btn v-if="isAdmin" color="secondary" @click="openFetchDialog" class="mr-2">
        <v-icon start>mdi-cloud-download</v-icon>
        Fetch Results
      </v-btn>
      <v-btn color="primary" @click="openCreateDialog">
        <v-icon start>mdi-plus</v-icon>
        Add Draw
      </v-btn>
    </div>

    <v-card>
      <v-data-table
        :headers="[
          { title: 'ID', key: 'id', width: '70px' },
          { title: 'Game', key: 'game_id', width: '120px' },
          { title: 'Date', key: 'draw_date', width: '120px' },
          { title: 'Numbers', key: 'result_numbers', width: '150px' },
          { title: 'Stars', key: 'result_stars', width: '120px' },
          { title: 'Actions', key: 'actions', width: '120px', sortable: false }
        ]"
        :items="drawList"
        :loading="loading"
        striped-rows
      >
        <template v-slot:item.game_id="{ item }">
          {{ getGameName(item.game_id) }}
        </template>

        <template v-slot:item.draw_date="{ item }">
          {{ formatDate(item.draw_date) }}
        </template>

        <template v-slot:item.result_numbers="{ item }">
          <span v-if="item.result_numbers">{{ item.result_numbers }}</span>
          <span v-else class="text-grey">-</span>
        </template>

        <template v-slot:item.result_stars="{ item }">
          <span v-if="item.result_stars">{{ item.result_stars }}</span>
          <span v-else class="text-grey">-</span>
        </template>

        <template v-slot:item.actions="{ item }">
          <v-tooltip location="top" text="Edit draw results">
            <template v-slot:activator="{ props }">
              <v-btn icon="mdi-pencil" variant="text" size="small" v-bind="props" @click="openEditDialog(item)" />
            </template>
          </v-tooltip>
          <v-tooltip location="top" text="Delete draw">
            <template v-slot:activator="{ props }">
              <v-btn icon="mdi-delete" variant="text" size="small" color="error" v-bind="props" @click="confirmDelete(item)" />
            </template>
          </v-tooltip>
        </template>
      </v-data-table>
    </v-card>

    <DrawsFormDialog
      v-model:visible="dialogVisible"
      :draw="selectedDraw"
      :games="gameList"
      @saved="handleSaved"
    />

    <FetchResultsDialog
      v-model:visible="fetchDialogVisible"
      @saved="handleFetchSaved"
    />

    <ConfirmDialog
      v-model:visible="confirmVisible"
      title="Delete Draw"
      message="Are you sure you want to delete this draw?"
      @confirm="handleConfirmDelete"
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
