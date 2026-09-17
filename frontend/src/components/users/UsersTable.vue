<script setup>
import { ref, onMounted } from 'vue'
import { users as usersService } from '../../services/api'

const userList = ref([])
const loading = ref(false)
const toggling = ref(null)
const snackbar = ref({ show: false, message: '', color: 'success' })

const dialogVisible = ref(false)
const selectedUser = ref(null)
const form = ref({ name: '', email: '', role: 'user' })
const formRef = ref(null)

async function loadUsers() {
  loading.value = true
  try {
    userList.value = await usersService.getAll()
  } catch (error) {
    showSnackbar(error.message, 'error')
  } finally {
    loading.value = false
  }
}

function openEditDialog(user) {
  selectedUser.value = user
  form.value = { name: user.name, email: user.email, role: user.role }
  dialogVisible.value = true
}

function closeDialog() {
  dialogVisible.value = false
  selectedUser.value = null
}

async function saveUser() {
  const { valid } = await formRef.value.validate()
  if (!valid) return

  try {
    await usersService.update(selectedUser.value.id, form.value)
    showSnackbar('User updated', 'success')
    closeDialog()
    await loadUsers()
  } catch (error) {
    showSnackbar(error.message, 'error')
  }
}

async function toggleActive(user) {
  const previous = user.active
  user.active = !previous
  toggling.value = user.id
  try {
    if (previous) {
      await usersService.deactivate(user.id)
      showSnackbar(`${user.name} deactivated`, 'success')
    } else {
      await usersService.activate(user.id)
      showSnackbar(`${user.name} activated`, 'success')
    }
  } catch (error) {
    user.active = previous
    showSnackbar(error.message, 'error')
  } finally {
    toggling.value = null
  }
}

function showSnackbar(message, color) {
  snackbar.value = { show: true, message, color }
}

onMounted(() => {
  loadUsers()
})
</script>

<template>
  <div>
    <v-card>
      <div class="table-responsive">
        <v-data-table
        :headers="[
          { title: 'ID', key: 'id', width: '70px' },
          { title: 'Name', key: 'name' },
          { title: 'Email', key: 'email' },
          { title: 'Role', key: 'role', width: '100px' },
          { title: 'Status', key: 'active', width: '120px', sortable: true },
          { title: 'Joined', key: 'created_at', width: '140px' },
          { title: 'Actions', key: 'actions', width: '120px', sortable: false }
        ]"
        :items="userList"
        :loading="loading"
        striped-rows
      >
        <template v-slot:item.role="{ item }">
          <v-chip
            :color="item.role === 'admin' ? 'primary' : 'grey'"
            size="small"
            variant="tonal"
          >
            {{ item.role }}
          </v-chip>
        </template>

        <template v-slot:item.active="{ item }">
          <v-switch
            :model-value="item.active"
            :color="item.active ? 'success' : 'warning'"
            density="compact"
            hide-details
            :loading="toggling === item.id"
            @update:model-value="toggleActive(item)"
          />
        </template>

        <template v-slot:item.created_at="{ item }">
          {{ item.created_at ? item.created_at.split('T')[0] : '-' }}
        </template>

        <template v-slot:item.actions="{ item }">
          <v-btn icon="mdi-pencil" variant="text" size="small" @click="openEditDialog(item)" />
        </template>
      </v-data-table>
      </div>
    </v-card>

    <v-dialog v-model="dialogVisible" max-width="450" persistent>
      <v-card>
        <v-card-title class="text-h6">Edit User</v-card-title>
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
              type="email"
              variant="outlined"
              :rules="[v => !!v || 'Email is required']"
              class="mb-3"
            />
            <v-select
              v-model="form.role"
              :items="['user', 'admin']"
              label="Role"
              variant="outlined"
              :rules="[v => !!v || 'Role is required']"
            />
          </v-form>
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="text" @click="closeDialog">Cancel</v-btn>
          <v-btn color="primary" @click="saveUser">Save</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <v-snackbar v-model="snackbar.show" :color="snackbar.color" timeout="3000">
      {{ snackbar.message }}
    </v-snackbar>
  </div>
</template>

<style scoped>
.table-responsive {
  overflow-x: auto;
  -webkit-overflow-scrolling: touch;
}
</style>
