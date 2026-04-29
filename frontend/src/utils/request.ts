import axios from 'axios'
import { useAuthStore } from '../stores/auth'

// 创建 axios 实例
const request = axios.create({
  baseURL: '/api/v1', // 使用 Vite 代理，指向我们后台 /api/v1
  timeout: 5000
})

// 请求拦截器
request.interceptors.request.use(
  config => {
    // 每次发送请求时，从 pinia 中提取 token，添加到请求头
    const authStore = useAuthStore()
    if (authStore.token) {
      config.headers['Authorization'] = `Bearer ${authStore.token}`
    }
    return config
  },
  error => {
    return Promise.reject(error)
  }
)

// 响应拦截器
request.interceptors.response.use(
  response => {
    return response.data
  },
  error => {
    if (error.response) {
      // 如果后端判定 401 无权限或者是过期 token
      if (error.response.status === 401) {
        const authStore = useAuthStore()
        authStore.logout() // 清除过期状态
      }
    }
    return Promise.reject(error)
  }
)

export default request
