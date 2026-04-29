<template>
  <div>
    <nav class="navbar">
      <div class="nav-left">
        <router-link to="/" class="back-link"> &lt; 返回首页 </router-link>
        <h1>{{ deptName }} - 医生列表</h1>
      </div>
    </nav>
    <div class="content">
      <div v-if="loading" class="loading">加载中...</div>
      <div v-else-if="doctors.length === 0" class="empty">
        当前科室下暂无医生排班信息
      </div>
      <div v-else class="doctor-list">
        <div class="doctor-card" v-for="doc in doctors" :key="doc.id">
          <div class="doc-info">
            <h3>{{ doc.name }}</h3>
            <span class="title">{{ doc.title }}</span>
            <p v-if="doc.specialty" class="specialty"><strong>专长:</strong> {{ doc.specialty }}</p>
            <p class="desc">{{ doc.description || '暂无简介' }}</p>
          </div>
          <button class="book-btn" @click="goToSchedule(doc.id, doc.name)">查看排班</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import request from '../utils/request'

const route = useRoute()
const router = useRouter()

const deptId = route.params.id
const deptName = route.query.name || '科室'
const doctors = ref<any[]>([])
const loading = ref(true)

const fetchDoctors = async () => {
  try {
    const res: any = await request.get(`/departments/${deptId}/doctors`)
    doctors.value = res.doctors
  } catch (error) {
    console.error('获取医生列表失败', error)
  } finally {
    loading.value = false
  }
}

const goToSchedule = (docId: number, docName: string) => {
  router.push({ path: `/doctor/${docId}/schedule`, query: { name: docName } })
}

onMounted(() => {
  fetchDoctors()
})
</script>

<style scoped>
.navbar {
  background-color: #409eff;
  color: white;
  padding: 15px 30px;
}
.nav-left {
  display: flex;
  align-items: center;
}
.back-link {
  color: white;
  text-decoration: none;
  margin-right: 20px;
  font-weight: bold;
}
.back-link:hover {
  text-decoration: underline;
}
.content {
  max-width: 800px;
  margin: 30px auto;
  padding: 0 20px;
}
.doctor-list {
  display: flex;
  flex-direction: column;
  gap: 15px;
}
.doctor-card {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: white;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
}
.doc-info h3 {
  margin: 0;
  display: inline-block;
  margin-right: 10px;
  color: #303133;
}
.title {
  background: #e1f3d8;
  color: #67c23a;
  padding: 2px 8px;
  border-radius: 12px;
  font-size: 12px;
}
.specialty, .desc {
  color: #606266;
  font-size: 14px;
  margin: 8px 0 0 0;
}
.book-btn {
  background-color: #409eff;
  color: white;
  border: none;
  padding: 10px 20px;
  border-radius: 4px;
  cursor: pointer;
  white-space: nowrap;
}
.book-btn:hover {
  background-color: #66b1ff;
}
.loading, .empty {
  text-align: center;
  padding: 40px;
  color: #909399;
}
</style>
