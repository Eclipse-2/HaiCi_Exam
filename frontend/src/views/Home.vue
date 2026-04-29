<template>
  <div class="home">
    <nav class="navbar">
      <h1>医院挂号主页</h1>
      <div v-if="authStore.token" class="user-info">
        欢迎您, {{ authStore.userInfo?.phone }}
        <router-link v-if="authStore.userInfo?.role === 'ADMIN'" to="/admin" class="nav-link" style="color: yellow; font-weight: bold;">[排班管理]</router-link>
        <router-link to="/my-appointments" class="nav-link">我的预约</router-link>
        <button class="logout-btn" @click="handleLogout">退出登录</button>
      </div>
      <div v-else>
        <router-link to="/login" class="nav-link">登录</router-link> | 
        <router-link to="/register" class="nav-link">注册</router-link>
      </div>
    </nav>

    <div class="content">
      <h2>请选择您要就诊的科室</h2>
      <div v-if="loading" class="loading">加载中...</div>
      <div v-else class="dept-grid">
        <div 
          class="dept-card" 
          v-for="dept in departments" 
          :key="dept.id"
          @click="goToDoctors(dept.id, dept.name)"
        >
          <h3>{{ dept.name }}</h3>
          <p>{{ dept.description || '暂无描述' }}</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import request from '../utils/request'

const authStore = useAuthStore()
const router = useRouter()

const departments = ref<any[]>([])
const loading = ref(true)

const fetchDepartments = async () => {
  try {
    const res: any = await request.get('/departments')
    departments.value = res.departments
  } catch (error) {
    console.error('获取科室列表失败', error)
    alert('网络或服务器异常，无法加载科室列表')
  } finally {
    loading.value = false
  }
}

const goToDoctors = (deptId: number, deptName: string) => {
  router.push({ path: `/department/${deptId}/doctors`, query: { name: deptName } })
}

const handleLogout = () => {
  authStore.logout()
}

onMounted(() => {
  fetchDepartments()
})
</script>

<style scoped>
.navbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background-color: #409eff;
  color: white;
  padding: 15px 30px;
}
.nav-link {
  color: white;
  text-decoration: none;
  margin: 0 10px;
}
.nav-link:hover {
  text-decoration: underline;
}
.logout-btn {
  margin-left: 15px;
  padding: 6px 12px;
  border: none;
  border-radius: 4px;
  background-color: #f56c6c;
  color: white;
  cursor: pointer;
}
.content {
  padding: 40px;
  max-width: 1000px;
  margin: 0 auto;
}
.dept-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 20px;
  margin-top: 20px;
}
.dept-card {
  background: white;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
  cursor: pointer;
  transition: transform 0.2s, box-shadow 0.2s;
}
.dept-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 4px 12px rgba(64, 158, 255, 0.3);
  border: 1px solid #409eff;
}
.dept-card h3 {
  margin-top: 0;
  color: #303133;
}
.dept-card p {
  color: #909399;
  font-size: 14px;
  line-height: 1.5;
}
.loading {
  text-align: center;
  color: #606266;
  padding: 40px;
}
</style>
