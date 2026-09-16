<script setup>
import { ref, onMounted } from 'vue'
import { contributions as contributionsService, users as usersService, games as gamesService } from '../../services/api'
import { useAuth } from '../../composables/useAuth'
import ContributionsFormDialog from './ContributionsFormDialog.vue'
import ConfirmDialog from '../layout/ConfirmDialog.vue'

const { isAdmin, user: currentUser } = useAuth()

const contributionList = ref([])
const userList = ref([])
const gameList = ref([])
const loading = ref(false)
const dialogVisible = ref(false)
const selectedContribution = ref(null)
const confirmVisible = ref(false)
const contributionToDelete = ref(null)
const snackbar = ref({ show: false, message: '', color: 'success' })

const monthNames = [
  'January', 'February', 'March', 'April', 'May', 'June',
  'July', 'August', 'September', 'October', 'November', 'December'
]

async function loadData() {
  loading.value = true
  try {
    const [contributionData, gameData, namesData] = await Promise.all([
      contributionsService.getAll(),
      gamesService.getAll(),
      usersService.getNames()
    ])
    contributionList.value = contributionData
    gameList.value = gameData
    userList.value = namesData
  } catch (error) {
    showSnackbar(error.message, 'error')
  } finally {
    loading.value = false
  }
}

function openCreateDialog() {
  selectedContribution.value = null
  dialogVisible.value = true
}

function openEditDialog(contribution) {
  selectedContribution.value = contribution
  dialogVisible.value = true
}

async function handleSaved(data) {
  try {
    if (data.id) {
      await contributionsService.update(data.id, data)
      showSnackbar('Contribution updated', 'success')
    } else {
      await contributionsService.create(data)
      showSnackbar('Contribution created', 'success')
    }
    dialogVisible.value = false
    await loadData()
  } catch (error) {
    showSnackbar(error.message, 'error')
  }
}

function confirmDelete(contribution) {
  contributionToDelete.value = contribution
  confirmVisible.value = true
}

async function handleConfirmDelete() {
  try {
    await contributionsService.delete(contributionToDelete.value.id)
    showSnackbar('Contribution deleted', 'success')
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
        Add Contribution
      </v-btn>
    </div>

    <v-card>
      <v-data-table
        :headers="[
          { title: 'ID', key: 'id', width: '70px' },
          { title: 'User', key: 'user_name', width: '150px' },
          { title: 'Game', key: 'game_id', width: '120px' },
          { title: 'Period', key: 'month', width: '120px' },
          { title: 'Amount', key: 'amount', width: '100px' },
          { title: 'Paid', key: 'paid', width: '80px', sortable: true },
          { title: 'Method', key: 'payment_method', width: '100px' },
          { title: 'Date', key: 'payment_date', width: '120px' },
          ...(isAdmin ? [{ title: 'Actions', key: 'actions', width: '120px', sortable: false }] : [])
        ]"
        :items="contributionList"
        :loading="loading"
        striped-rows
      >
        <template v-slot:item.user_name="{ item }">
          {{ item.user_name || `User #${item.user_id}` }}
        </template>

        <template v-slot:item.game_id="{ item }">
          {{ gameList.find(g => g.id === item.game_id)?.name || `#${item.game_id}` }}
        </template>

        <template v-slot:item.month="{ item }">
          {{ monthNames[item.month - 1] }} {{ item.year }}
        </template>

        <template v-slot:item.amount="{ item }">
          {{ item.amount.toFixed(2) }} €
        </template>

        <template v-slot:item.paid="{ item }">
          <v-chip :color="item.paid ? 'success' : 'grey'" size="small">
            {{ item.paid ? 'Yes' : 'No' }}
          </v-chip>
        </template>

        <template v-slot:item.payment_method="{ item }">
          <v-chip v-if="item.payment_method" size="small" :color="item.payment_method === 'BIZUM' ? 'info' : 'warning'">
            {{ item.payment_method }}
          </v-chip>
          <span v-else class="text-grey">-</span>
        </template>

        <template v-slot:item.payment_date="{ item }">
          {{ item.payment_date ? item.payment_date.split('T')[0] : '-' }}
        </template>

        <template v-slot:item.actions="{ item }" v-if="isAdmin">
          <v-btn icon="mdi-pencil" variant="text" size="small" @click="openEditDialog(item)" />
          <v-btn icon="mdi-delete" variant="text" size="small" color="error" @click="confirmDelete(item)" />
        </template>
      </v-data-table>
    </v-card>

    <ContributionsFormDialog
      v-model:visible="dialogVisible"
      :contribution="selectedContribution"
      :users="userList"
      :games="gameList"
      :current-user="currentUser"
      @saved="handleSaved"
    />

    <ConfirmDialog
      v-model:visible="confirmVisible"
      title="Delete Contribution"
      message="Are you sure you want to delete this contribution?"
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
