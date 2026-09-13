<script setup>
import { ref, onMounted } from 'vue'
import { draws as drawsService, games as gamesService } from '../../services/api'
import DrawsFormDialog from './DrawsFormDialog.vue'

const drawList = ref([])
const gameList = ref([])
const loading = ref(false)
const dialogVisible = ref(false)
const selectedDraw = ref(null)
const snackbar = ref({ show: false, message: '', color: 'success' })

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

async function markAsProcessed(draw) {
  try {
    await draws.markAsProcessed(draw.id)
    showSnackbar('Draw marked as processed', 'success')
    await loadData()
  } catch (error) {
    showSnackbar(error.message, 'error')
  }
}

async function deleteDraw(draw) {
  try {
    await drawsService.delete(draw.id)
    showSnackbar('Draw deleted', 'success')
    await loadData()
  } catch (error) {
    showSnackbar(error.message, 'error')
  }
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
          { title: 'Processed', key: 'processed', width: '110px', sortable: true },
          { title: 'Actions', key: 'actions', width: '150px', sortable: false }
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

        <template v-slot:item.processed="{ item }">
          <v-chip :color="item.processed ? 'success' : 'warning'" size="small">
            {{ item.processed ? 'Yes' : 'No' }}
          </v-chip>
        </template>

        <template v-slot:item.actions="{ item }">
          <v-btn icon="mdi-pencil" variant="text" size="small" @click="openEditDialog(item)" />
          <v-btn
            v-if="!item.processed"
            icon="mdi-check-circle"
            variant="text"
            size="small"
            color="success"
            @click="markAsProcessed(item)"
          />
          <v-btn icon="mdi-delete" variant="text" size="small" color="error" @click="deleteDraw(item)" />
        </template>
      </v-data-table>
    </v-card>

    <DrawsFormDialog
      v-model:visible="dialogVisible"
      :draw="selectedDraw"
      :games="gameList"
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
