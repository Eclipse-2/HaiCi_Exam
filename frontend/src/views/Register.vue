<template>
  <div class="login-container">
    <div class="form-box">
      <h2>注册挂号系统</h2>
      <form @submit.prevent="handleRegister">
        <div class="input-group">
          <label>手机号</label>
          <input type="text" v-model="form.phone" required placeholder="请输入手机号 (如 138xxxx)" />
        </div>
        <div class="input-group">
          <label>密码</label>
          <input type="password" v-model="form.password" required minlength="6" placeholder="至少6位密码" />
        </div>
        <button type="submit" class="submit-btn" :disabled="loading">注册</button>
        <p v-if="errorMsg" class="error-msg">{{ errorMsg }}</p>
        <div class="links">
          <span>已有账号？ <router-link to="/login">马上登录</router-link></span>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import request from '../utils/request'

const router = useRouter()

const form = ref({
  phone: '',
  password: ''
})

const loading = ref(false)
const errorMsg = ref('')

const handleRegister = async () => {
  loading.value = true
  errorMsg.value = ''
  try {
    await request.post('/auth/register', form.value)
    alert('注册成功，请去登录！')
    router.push('/login') // 跳转到登录页
  } catch (error: any) {
    if (error.response && error.response.data) {
      errorMsg.value = error.response.data.error || '注册失败'
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
  box-sizing: border-box;
}
.submit-btn {
  width: 100%;
  padding: 12px;
  background-color: #67c23a;
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 16px;
  cursor: pointer;
  margin-top: 10px;
}
.submit-btn:disabled {
  background-color: #b3e19d;
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
