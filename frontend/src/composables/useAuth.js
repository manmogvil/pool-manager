import { ref, computed } from 'vue'
import { auth } from '../services/api'

const user = ref(auth.getUser())
const isAuthenticated = computed(() => !!user.value)

function login(token, userData) {
  auth.setToken(token)
  auth.setUser(userData)
  user.value = userData
}

function logout() {
  auth.logout()
  user.value = null
}

function refreshUser() {
  user.value = auth.getUser()
}

export function useAuth() {
  return {
    user,
    isAuthenticated,
    login,
    logout,
    refreshUser,
    isAdmin: computed(() => user.value?.role === 'admin'),
  }
}
