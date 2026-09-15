<script setup>
import { ref, onMounted } from 'vue'
import { tickets as ticketsService, draws as drawsService, games } from '../../services/api'
import TicketsFormDialog from './TicketsFormDialog.vue'
import ConfirmDialog from '../layout/ConfirmDialog.vue'

const ticketList = ref([])
const drawList = ref([])
const gameList = ref([])
const loading = ref(false)
const checkingId = ref(null)
const dialogVisible = ref(false)
const selectedTicket = ref(null)
const confirmVisible = ref(false)
const ticketToDelete = ref(null)
const snackbar = ref({ show: false, message: '', color: 'success' })

async function loadData() {
  loading.value = true
  try {
    const [t, d, g] = await Promise.all([ticketsService.getAll(), drawsService.getAll(), games.getAll()])
    ticketList.value = t || []
    drawList.value = d || []
    gameList.value = g || []
  } catch (error) {
    showSnackbar(error.message, 'error')
  } finally {
    loading.value = false
  }
}

function getDrawInfo(drawId) {
  const d = drawList.value.find(d => d.id === drawId)
  if (!d) return `#${drawId}`
  const game = gameList.value.find(g => g.id === d.game_id)
  const gameName = game ? game.name : 'Unknown'
  const date = new Date(d.draw_date)
  const day = String(date.getDate()).padStart(2, '0')
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const year = date.getFullYear()
  return `${gameName} #${d.id} (${day}/${month}/${year})`
}

function formatCurrency(value) {
  if (value === null || value === undefined) return '-'
  return `${value.toFixed(2)} €`
}

function openCreateDialog() {
  selectedTicket.value = null
  dialogVisible.value = true
}

function openEditDialog(ticket) {
  selectedTicket.value = ticket
  dialogVisible.value = true
}

async function checkTicket(ticket) {
  checkingId.value = ticket.id
  try {
    const result = await ticketsService.check(ticket.id)
    showSnackbar(result.message || 'Ticket checked', 'success')
    await loadData()
  } catch (error) {
    showSnackbar(error.message, 'error')
  } finally {
    checkingId.value = null
  }
}

async function handleSaved(data) {
  try {
    if (data.id) {
      await ticketsService.updatePrize(data.id, {
        prize_tier: data.prize_tier || null,
        prize_amount: data.prize_amount || null,
        matched_numbers: data.matched_numbers || null,
        matched_stars: data.matched_stars || null
      })
      showSnackbar('Ticket updated', 'success')
    } else {
      await ticketsService.create({
        draw_id: data.draw_id,
        numbers: data.numbers,
        stars: data.stars,
        cost: data.cost
      })
      showSnackbar('Ticket created', 'success')
    }
    dialogVisible.value = false
    await loadData()
  } catch (error) {
    showSnackbar(error.message, 'error')
  }
}

function confirmDelete(ticket) {
  ticketToDelete.value = ticket
  confirmVisible.value = true
}

async function handleConfirmDelete() {
  try {
    await ticketsService.delete(ticketToDelete.value.id)
    showSnackbar('Ticket deleted', 'success')
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
      <v-btn color="primary" @click="openCreateDialog">
        <v-icon start>mdi-plus</v-icon>
        Add Ticket
      </v-btn>
    </div>

    <v-card>
      <v-data-table
        :headers="[
          { title: 'ID', key: 'id', width: '60px' },
          { title: 'Draw', key: 'draw_id', width: '160px' },
          { title: 'Numbers', key: 'numbers', width: '140px' },
          { title: 'Stars', key: 'stars', width: '100px' },
          { title: 'Cost', key: 'cost', width: '90px' },
          { title: 'Prize Tier', key: 'prize_tier', width: '110px' },
          { title: 'Prize', key: 'prize_amount', width: '100px' },
          { title: 'Matched', key: 'matched_numbers', width: '90px' },
          { title: 'Actions', key: 'actions', width: '160px', sortable: false }
        ]"
        :items="ticketList"
        :loading="loading"
        striped-rows
      >
        <template v-slot:item.draw_id="{ item }">
          {{ getDrawInfo(item.draw_id) }}
        </template>

        <template v-slot:item.cost="{ item }">
          {{ formatCurrency(item.cost) }}
        </template>

        <template v-slot:item.prize_tier="{ item }">
          <v-chip v-if="item.prize_tier" size="small" color="success">{{ item.prize_tier }}</v-chip>
          <span v-else class="text-grey">-</span>
        </template>

        <template v-slot:item.prize_amount="{ item }">
          {{ formatCurrency(item.prize_amount) }}
        </template>

        <template v-slot:item.matched_numbers="{ item }">
          <span v-if="item.matched_numbers !== null">{{ item.matched_numbers }}</span>
          <span v-else class="text-grey">-</span>
        </template>

        <template v-slot:item.actions="{ item }">
          <v-tooltip v-if="!item.prize_tier" location="top" text="Check prize against API">
            <template v-slot:activator="{ props }">
              <v-btn
                icon="mdi-magnify"
                variant="text"
                size="small"
                color="primary"
                :loading="checkingId === item.id"
                v-bind="props"
                @click="checkTicket(item)"
              />
            </template>
          </v-tooltip>
          <v-tooltip location="top" text="Edit ticket">
            <template v-slot:activator="{ props }">
              <v-btn icon="mdi-pencil" variant="text" size="small" v-bind="props" @click="openEditDialog(item)" />
            </template>
          </v-tooltip>
          <v-tooltip location="top" text="Delete ticket">
            <template v-slot:activator="{ props }">
              <v-btn icon="mdi-delete" variant="text" size="small" color="error" v-bind="props" @click="confirmDelete(item)" />
            </template>
          </v-tooltip>
        </template>
      </v-data-table>
    </v-card>

    <TicketsFormDialog
      v-model:visible="dialogVisible"
      :ticket="selectedTicket"
      :draws="drawList"
      :games="gameList"
      @saved="handleSaved"
    />

    <ConfirmDialog
      v-model:visible="confirmVisible"
      title="Delete Ticket"
      message="Are you sure you want to delete this ticket?"
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
