<template>
  <div>
    <nav class="navbar">
      <div class="nav-left">
        <a @click.prevent="goBack" class="back-link" style="cursor:pointer"> &lt; 返回上一页</a>
        <h1>{{ docName }} - 排班表</h1>
      </div>
    </nav>

    <div class="content">
      <div v-if="loading" class="loading">排班加载中...</div>
      <div v-else-if="schedules.length === 0" class="empty">
        该医生近期暂无排班
      </div>
      <div v-else class="schedule-grid">
        <div 
          v-for="sched in schedules" 
          :key="sched.id" 
          :class="['schedule-card', { 'sold-out': sched.available_slots === 0 || sched.status !== 'AVAILABLE' }]"
        >
          <div class="date-info">
            <span class="date">{{ sched.date }}</span>
            <span class="session">{{ sched.session === 'MORNING' ? '上午' : '下午' }}</span>
          </div>
          <div class="status-info">
            <span v-if="sched.status !== 'AVAILABLE'" class="bad-status">停诊</span>
            <span v-else-if="sched.available_slots === 0" class="bad-status">已满诊</span>
            <span v-else class="good-status">余号: {{ sched.available_slots }}</span>
          </div>
          <button 
            class="submit-btn" 
            :disabled="sched.available_slots === 0 || sched.status !== 'AVAILABLE'"
            @click="openBookingModal(sched)"
          >
            {{ (sched.available_slots === 0 || sched.status !== 'AVAILABLE') ? '不可预约' : '预约挂号' }}
          </button>
        </div>
      </div>
    </div>

    <!-- 选择就诊人专属弹窗 -->
    <div v-if="showModal" class="modal-overlay" @click.self="closeModal">
      <div class="modal">
        <h3>确认挂号信息</h3>
        <p>医生：<strong>{{ docName }}</strong></p>
        <p>排班：<strong>{{ selectedSched?.date }} {{ selectedSched?.session === 'MORNING' ? '上午' : '下午' }}</strong></p>
        <hr/>
        <div class="patient-select">
          <h4>请选择就诊人：</h4>
          <div v-if="loadingPatients">正在加载您的就诊人...</div>
          <div v-else-if="myPatients.length === 0">
            <p>您还没有添加任何就诊人。</p>
            <button class="nav-btn" @click="goToPatients">去添加就诊人</button>
          </div>
          <div v-else class="patient-options">
            <label v-for="p in myPatients" :key="p.id" class="patient-radio">
              <input type="radio" :value="p.id" v-model="selectedPatientId" />
              <span>{{ p.name }} ({{ p.id_card.replace(/^(.{4})(?:\d+)(.{4})$/, "$1****$2") }})</span>
            </label>
          </div>
        </div>

        <div class="modal-actions">
          <button class="cancel-btn" @click="closeModal">取消</button>
          <button 
            class="confirm-btn" 
            :disabled="!selectedPatientId || submitting" 
            @click="confirmBooking"
          >
            {{ submitting ? '提交中...' : '提交预约' }}
          </button>
        </div>
      </div>
    </div>

  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import request from '../utils/request'
import { useAuthStore } from '../stores/auth'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const docId = route.params.id
const docName = route.query.name || '医生'
const schedules = ref<any[]>([])
const loading = ref(true)

// 挂号弹窗与参数
const showModal = ref(false)
const selectedSched = ref<any>(null)
const myPatients = ref<any[]>([])
const selectedPatientId = ref<number | null>(null)
const loadingPatients = ref(false)
const submitting = ref(false)

const fetchSchedules = async () => {
  try {
    const res: any = await request.get(`/doctors/${docId}/schedules`)
    schedules.value = res.schedules.map((s: any) => {
      s.date = s.date.split('T')[0]
      return s
    })
  } catch (error) {
    console.error('获取排班失败', error)
  } finally {
    loading.value = false
  }
}

const goBack = () => {
  router.back()
}

const goToPatients = () => {
  router.push('/patients')
}

const fetchMyPatients = async () => {
  loadingPatients.value = true
  try {
    const res: any = await request.get('/patients')
    myPatients.value = res.patients || []
  } catch (error) {
    console.error(error)
  } finally {
    loadingPatients.value = false
  }
}

// 点击预约，开启弹窗
const openBookingModal = async (sched: any) => {
  if (!authStore.token) {
    alert('请先登录后再进行挂号操作')
    router.push('/login')
    return
  }
  selectedSched.value = sched
  showModal.value = true
  await fetchMyPatients()
}

const closeModal = () => {
  showModal.value = false
  selectedPatientId.value = null
  selectedSched.value = null
}

const confirmBooking = async () => {
  if (!selectedPatientId.value || !selectedSched.value) return
  submitting.value = true
  try {
    const res: any = await request.post('/appointments', {
      schedule_id: selectedSched.value.id,
      patient_id: selectedPatientId.value
    })
    alert('🎉预约成功！流水号: ' + res.appointment.appointment_no)
    closeModal()
    router.push('/my-appointments') // 去我的预约列表
  } catch (error: any) {
    alert(error.response?.data?.error || '抱歉，挂号失败，请重试')
    // 如果是因为号源已空，建议刷新该排班页面
    fetchSchedules()
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  fetchSchedules()
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
  font-weight: bold;
  margin-right: 20px;
}
.back-link:hover {
  text-decoration: underline;
}
.content {
  max-width: 900px;
  margin: 30px auto;
  padding: 0 20px;
}
.schedule-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 20px;
}
.schedule-card {
  background: white;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
  display: flex;
  flex-direction: column;
  align-items: center;
  border: 1px solid transparent;
}
.schedule-card.sold-out {
  background: #fdfdfd;
  opacity: 0.8;
}
.date-info {
  margin-bottom: 10px;
  font-size: 18px;
  font-weight: 500;
  color: #303133;
}
.session {
  margin-left: 10px;
  font-size: 14px;
  color: #909399;
}
.status-info {
  margin-bottom: 20px;
}
.bad-status {
  color: #f56c6c;
  background: #fef0f0;
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 14px;
}
.good-status {
  color: #67c23a;
  background: #f0f9eb;
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 14px;
  font-weight: bold;
}
.submit-btn {
  width: 100%;
  background-color: #409eff;
  color: white;
  border: none;
  padding: 10px;
  border-radius: 4px;
  cursor: pointer;
}
.submit-btn:disabled {
  background-color: #c0c4cc;
  cursor: not-allowed;
}
.loading, .empty {
  text-align: center;
  padding: 40px;
  color: #909399;
}
/* 模态框样式 */
.modal-overlay {
  position: fixed;
  top: 0; left: 0; right: 0; bottom: 0;
  background: rgba(0,0,0,0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 100;
}
.modal {
  background: white;
  padding: 30px;
  border-radius: 8px;
  width: 90%;
  max-width: 400px;
}
.modal h3 {
  margin-top: 0;
}
.patient-select {
  margin: 20px 0;
}
.patient-options {
  display: flex;
  flex-direction: column;
  gap: 10px;
}
.patient-radio {
  display: flex;
  align-items: center;
  gap: 10px;
  background: #f4f4f5;
  padding: 10px;
  border-radius: 4px;
  cursor: pointer;
}
.patient-radio input {
  margin: 0;
}
.nav-btn {
  background: #e6a23c; color: white; padding: 6px 12px; border: none; border-radius: 4px; cursor: pointer;
}
.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 15px;
}
.cancel-btn {
  background: #909399; color: white; padding: 8px 16px; border: none; border-radius: 4px; cursor: pointer;
}
.confirm-btn {
  background: #409eff; color: white; padding: 8px 16px; border: none; border-radius: 4px; cursor: pointer;
}
.confirm-btn:disabled {
  background: #a0cfff; cursor: not-allowed;
}
</style>
