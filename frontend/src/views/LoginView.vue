<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { auth } from '../services/api'
import { useAuth } from '../composables/useAuth'

const router = useRouter()
const { login: doLogin } = useAuth()
const tab = ref('login')
const email = ref('')
const password = ref('')
const name = ref('')
const loading = ref(false)
const error = ref('')
const snackbar = ref({ show: false, message: '', color: 'success' })

async function handleLogin() {
  error.value = ''
  loading.value = true
  try {
    const result = await auth.login(email.value, password.value)
    doLogin(result.token, result.user)
    router.push('/')
  } catch (e) {
    error.value = e.message
  } finally {
    loading.value = false
  }
}

async function handleRegister() {
  error.value = ''
  loading.value = true
  try {
    await auth.register(name.value, email.value, password.value)
    snackbar.value = { show: true, message: 'Registration successful. Wait for admin approval.', color: 'success' }
    tab.value = 'login'
    password.value = ''
  } catch (e) {
    error.value = e.message
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <v-main class="login-bg">
    <v-container fill-height>
      <v-row justify="center">
        <v-col cols="12" sm="8" md="4">
          <v-card class="login-card" elevation="4">
            <v-card-title class="text-h5 text-center pa-6">
              <v-icon color="primary" size="large" class="mr-2">mdi-ticket</v-icon>
              Lottery Pool Manager
            </v-card-title>

            <v-tabs v-model="tab" color="primary" grow>
              <v-tab value="login">Login</v-tab>
              <v-tab value="register">Register</v-tab>
            </v-tabs>

            <v-card-text class="pa-6">
              <v-alert v-if="error" type="error" variant="tonal" density="compact" class="mb-4">
                {{ error }}
              </v-alert>

              <v-form @submit.prevent="tab === 'login' ? handleLogin() : handleRegister()">
                <v-text-field
                  v-if="tab === 'register'"
                  v-model="name"
                  label="Name"
                  variant="outlined"
                  prepend-inner-icon="mdi-account"
                  class="mb-3"
                  :rules="[v => !!v || 'Name is required']"
                />
                <v-text-field
                  v-model="email"
                  label="Email"
                  type="email"
                  variant="outlined"
                  prepend-inner-icon="mdi-email"
                  class="mb-3"
                  :rules="[v => !!v || 'Email is required']"
                />
                <v-text-field
                  v-model="password"
                  label="Password"
                  type="password"
                  variant="outlined"
                  prepend-inner-icon="mdi-lock"
                  class="mb-4"
                  :rules="[v => !!v || 'Password is required']"
                />
                <v-btn
                  type="submit"
                  color="primary"
                  block
                  size="large"
                  :loading="loading"
                >
                  {{ tab === 'login' ? 'Login' : 'Register' }}
                </v-btn>
              </v-form>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>
    </v-container>
  </v-main>

  <v-snackbar v-model="snackbar.show" :color="snackbar.color" timeout="5000">
    {{ snackbar.message }}
  </v-snackbar>
</template>

<style scoped>
.login-bg {
  background: #f8f9fa;
  min-height: 100vh;
  min-height: 100dvh;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow-y: auto;
}
.login-card {
  border-radius: 12px;
  margin: 16px;
}
</style>
