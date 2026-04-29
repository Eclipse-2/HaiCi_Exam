<template>
  <div class="login-container">
    <div class="form-box">
      <h2>登录挂号系统</h2>
      <form @submit.prevent="handleLogin">
        <div class="input-group">
          <label>手机号</label>
          <input type="text" v-model="form.phone" required placeholder="请输入手机号" />
        </div>
        <div class="input-group">
          <label>密码</label>
          <input type="password" v-model="form.password" required placeholder="请输入密码" />
        </div>
        <button type="submit" class="submit-btn" :disabled="loading">登录</button>
        <p v-if="errorMsg" class="error-msg">{{ errorMsg }}</p>
        <div class="links">
          <span>还没有账号？ <router-link to="/register">立即注册</router-link></span>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import request from '../utils/request'

const router = useRouter()
const authStore = useAuthStore()

const form = ref({
  phone: '',
  password: ''
})

const loading = ref(false)
const errorMsg = ref('')

const handleLogin = async () => {
  loading.value = true
  errorMsg.value = ''
  try {
    const res: any = await request.post('/auth/login', form.value)
    // 登录成功存入 token 和用户信息到 Pinia 和 localStorage
    authStore.setAuth(res.token, res.user)
    router.push('/') // 回到首页
  } catch (error: any) {
    if (error.response && error.response.data) {
      errorMsg.value = error.response.data.error || '登录失败'
    } else {
      errorMsg.value = '服务器连接异常'
    }
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  background-color: #f5f7fa;
}
.form-box {
  background: white;
  padding: 40px;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  width: 100%;
  max-width: 400px;
}
h2 {
  text-align: center;
  margin-bottom: 20px;
  color: #333;
}
.input-group {
  margin-bottom: 15px;
}
.input-group label {
  display: block;
  margin-bottom: 5px;
  color: #666;
}
.input-group input {
  width: 100%;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
  box-sizing: border-box; /* 防止边距超出 */
}
.submit-btn {
  width: 100%;
  padding: 12px;
  background-color: #409eff;
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 16px;
  cursor: pointer;
  margin-top: 10px;
}
.submit-btn:disabled {
  background-color: #a0cfff;
}
.error-msg {
  color: red;
  margin-top: 10px;
  text-align: center;
}
.links {
  margin-top: 15px;
  text-align: center;
  font-size: 14px;
}
</style>
