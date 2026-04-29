import { defineStore } from 'pinia'
import { ref } from 'vue'
import router from '../router'

export const useAuthStore = defineStore('auth', () => {
  // state
  const token = ref(localStorage.getItem('token') || '')
  const userInfo = ref(JSON.parse(localStorage.getItem('userInfo') || 'null'))

  // actions
  const setAuth = (newToken: string, user: any) => {
    token.value = newToken
    userInfo.value = user
    localStorage.setItem('token', newToken)
    localStorage.setItem('userInfo', JSON.stringify(user))
  }

  const logout = () => {
    token.value = ''
    userInfo.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('userInfo')
    router.push('/login')
  }

  return { token, userInfo, setAuth, logout }
})
