<script setup>
import { ref, onMounted } from 'vue'
import { games as gamesService } from '../../services/api'
import { useAuth } from '../../composables/useAuth'
import GamesFormDialog from './GamesFormDialog.vue'

const { isAdmin } = useAuth()

const gameList = ref([])
const loading = ref(false)
const dialogVisible = ref(false)
const selectedGame = ref(null)
const snackbar = ref({ show: false, message: '', color: 'success' })

async function loadGames() {
  loading.value = true
  try {
    gameList.value = await gamesService.getAll()
  } catch (error) {
    showSnackbar(error.message, 'error')
  } finally {
    loading.value = false
  }
}

function openCreateDialog() {
  selectedGame.value = null
  dialogVisible.value = true
}

function openEditDialog(game) {
  selectedGame.value = game
  dialogVisible.value = true
}

async function handleSaved(data) {
  try {
    if (data.id) {
      await gamesService.update(data.id, {
        name: data.name,
        draw_days: data.draw_days,
        ticket_price: data.ticket_price,
        active: data.active
      })
      showSnackbar('Game updated', 'success')
    } else {
      await gamesService.create({
        name: data.name,
        draw_days: data.draw_days,
        ticket_price: data.ticket_price
      })
      showSnackbar('Game created', 'success')
    }
    dialogVisible.value = false
    await loadGames()
  } catch (error) {
    showSnackbar(error.message, 'error')
  }
}

async function toggleActive(game, active) {
  try {
    await gamesService.update(game.id, {
      name: game.name,
      draw_days: game.draw_days,
      ticket_price: game.ticket_price,
      active
    })
    showSnackbar(active ? 'Game activated' : 'Game deactivated', 'success')
  } catch (error) {
    showSnackbar(error.message, 'error')
    await loadGames()
  }
}

async function deleteGame(game) {
  try {
    await gamesService.delete(game.id)
    showSnackbar('Game deleted', 'success')
    await loadGames()
  } catch (error) {
    showSnackbar(error.message, 'error')
  }
}

function showSnackbar(message, color) {
  snackbar.value = { show: true, message, color }
}

onMounted(() => {
  loadGames()
})
</script>

<template>
  <div>
    <div class="table-header" v-if="isAdmin">
      <v-btn color="primary" @click="openCreateDialog">
        <v-icon start>mdi-plus</v-icon>
        Add Game
      </v-btn>
    </div>

    <v-card>
      <v-data-table
        :headers="[
          { title: 'ID', key: 'id', width: '80px' },
          { title: 'Name', key: 'name' },
          { title: 'Draw Days', key: 'draw_days' },
          { title: 'Ticket Price', key: 'ticket_price', width: '140px' },
          ...(isAdmin ? [
            { title: 'Status', key: 'active', width: '100px', sortable: true },
            { title: 'Actions', key: 'actions', width: '120px', sortable: false }
          ] : [])
        ]"
        :items="gameList"
        :loading="loading"
        striped-rows
      >
        <template v-slot:item.ticket_price="{ item }">
          {{ item.ticket_price.toFixed(2) }} €
        </template>

        <template v-slot:item.active="{ item }" v-if="isAdmin">
          <v-switch
            v-model="item.active"
            color="success"
            density="compact"
            hide-details
            @update:model-value="(val) => toggleActive(item, val)"
          />
        </template>

        <template v-slot:item.actions="{ item }" v-if="isAdmin">
          <v-btn icon="mdi-pencil" variant="text" size="small" @click="openEditDialog(item)" />
          <v-btn icon="mdi-delete" variant="text" size="small" color="error" @click="deleteGame(item)" />
        </template>
      </v-data-table>
    </v-card>

    <GamesFormDialog
      v-model:visible="dialogVisible"
      :game="selectedGame"
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
