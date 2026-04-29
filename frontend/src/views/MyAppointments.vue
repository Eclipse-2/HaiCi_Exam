<template>
  <div class="wrapper">
    <nav class="navbar">
      <div class="nav-left">
        <router-link to="/" class="back-link"> &lt; 返回首页 </router-link>
        <h1>我的挂号订单</h1>
      </div>
    </nav>

    <div class="content">
      <div v-if="loading" class="loading">加载中...</div>
      <div v-else-if="appointments.length === 0" class="empty">您还没有预约过任何号源。</div>
      <div v-else class="appt-list">
        <div 
          v-for="appt in appointments" 
          :key="appt.id" 
          :class="['appt-card', `status-${appt.status.toLowerCase()}`]"
        >
          <div class="appt-header">
            <span>流水号: <strong>{{ appt.appointment_no }}</strong></span>
            <span class="status-tag">{{ getStatusText(appt.status) }}</span>
          </div>
          <div class="appt-body">
            <p><strong>就诊人：</strong> {{ appt.patient?.name }}</p>
            <p><strong>就诊日期：</strong> {{ formatDate(appt.schedule?.date) }} {{ appt.schedule?.session === 'MORNING' ? '上午' : '下午' }}</p>
            <p class="meta-time">下单时间: {{ formatTime(appt.created_at) }}</p>
          </div>
          <div class="appt-footer" v-if="appt.status === 'PENDING'">
            <button class="cancel-btn" @click="handleCancel(appt.id)" :disabled="cancelingId === appt.id">
              {{ cancelingId === appt.id ? '取消中...' : '取消预约 (30分钟免费退)' }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import request from '../utils/request'

const appointments = ref<any[]>([])
const loading = ref(true)
const cancelingId = ref<number | null>(null)

const fetchAppointments = async () => {
  loading.value = true
  try {
    const res: any = await request.get('/appointments')
    appointments.value = res.appointments || []
  } catch (error) {
    console.error('获取预约订单失败', error)
  } finally {
    loading.value = false
  }
}

const getStatusText = (status: string) => {
  switch(status) {
    case 'PENDING': return '待就诊'
    case 'CANCELLED': return '已取消'
    case 'FINISHED': return '已完成'
    default: return status
  }
}

const formatDate = (dateStr: string) => dateStr ? dateStr.split('T')[0] : ''
const formatTime = (timeStr: string) => {
  if (!timeStr) return ''
  const d = new Date(timeStr)
  return d.toLocaleString()
}

const handleCancel = async (id: number) => {
  if (!confirm('您确定要取消这个预约号源吗？')) return
  cancelingId.value = id
  try {
    const res: any = await request.put(`/appointments/${id}/cancel`)
    alert(res.message)
    await fetchAppointments()
  } catch (error: any) {
    alert(error.response?.data?.error || '取消失败')
  } finally {
    cancelingId.value = null
  }
}

onMounted(() => {
  fetchAppointments()
})
</script>

<style scoped>
.navbar { background-color: #409eff; color: white; padding: 15px 30px; }
.nav-left { display: flex; align-items: center; }
.back-link { color: white; text-decoration: none; font-weight: bold; margin-right: 20px; }
.back-link:hover { text-decoration: underline; }
.content { max-width: 800px; margin: 30px auto; padding: 0 20px; }
.appt-list { display: flex; flex-direction: column; gap: 15px; }
.appt-card { background: white; padding: 20px; border-radius: 8px; box-shadow: 0 2px 8px rgba(0,0,0,0.1); border-left: 5px solid #ccc; }
.status-pending { border-left-color: #e6a23c; }
.status-cancelled { border-left-color: #909399; opacity: 0.7; }
.status-finished { border-left-color: #67c23a; }
.appt-header { display: flex; justify-content: space-between; border-bottom: 1px solid #ebeef5; padding-bottom: 10px; margin-bottom: 10px; }
.status-tag { font-weight: bold; }
.status-pending .status-tag { color: #e6a23c; }
.status-cancelled .status-tag { color: #909399; }
.status-finished .status-tag { color: #67c23a; }
.appt-body { color: #606266; line-height: 1.6; }
.appt-body p { margin: 5px 0; }
.meta-time { font-size: 12px; color: #a8abb2; margin-top: 10px !important; }
.appt-footer { margin-top: 15px; text-align: right; }
.cancel-btn { background-color: #f56c6c; color: white; border: none; padding: 8px 15px; border-radius: 4px; cursor: pointer; }
.cancel-btn:disabled { background-color: #fab6b6; cursor: not-allowed; }
.loading, .empty { text-align: center; color: #909399; padding: 40px; }
</style>