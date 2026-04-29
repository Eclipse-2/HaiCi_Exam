<template>
  <div class="wrapper">
    <nav class="navbar">
      <div class="nav-left">
        <router-link to="/" class="back-link"> &lt; 返回首页 </router-link>
        <h1>就诊人管理</h1>
      </div>
    </nav>
    <div class="content">
      <div class="section-title">
        <h2>我的就诊人列表</h2>
        <button class="add-btn" @click="showForm = !showForm">{{ showForm ? '取消添加' : '+ 添加就诊人' }}</button>
      </div>

      <div v-if="loading" class="loading">加载中...</div>
      <div v-else-if="patients.length === 0" class="empty">您还没有添加任何就诊人</div>
      <div v-else class="patient-grid">
        <div class="patient-card" v-for="p in patients" :key="p.id">
          <div class="info-row"><strong>姓名:</strong> {{ p.name }}</div>
          <div class="info-row"><strong>身份证号:</strong> {{ p.id_card }}</div>
          <div class="info-row"><strong>手机号:</strong> {{ p.phone }}</div>
        </div>
      </div>

      <!-- 添加就诊人表单 -->
      <div v-if="showForm" class="add-form-box">
        <h3>添加新就诊人</h3>
        <form @submit.prevent="handleAddPatient">
          <div class="input-group">
            <label>姓名</label>
            <input type="text" v-model="form.name" required placeholder="请输入真实姓名">
          </div>
          <div class="input-group">
            <label>身份证号</label>
            <input type="text" v-model="form.id_card" required maxlength="18" placeholder="请输入18位身份证号">
          </div>
          <div class="input-group">
            <label>手机号</label>
            <input type="text" v-model="form.phone" required maxlength="11" placeholder="请输入联系手机号">
          </div>
          <button type="submit" class="submit-btn" :disabled="submitting">确认添加</button>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import request from '../utils/request'

const patients = ref<any[]>([])
const loading = ref(true)
const showForm = ref(false)
const submitting = ref(false)

const form = ref({
  name: '',
  id_card: '',
  phone: ''
})

const fetchPatients = async () => {
  loading.value = true
  try {
    const res: any = await request.get('/patients')
    patients.value = res.patients
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

const handleAddPatient = async () => {
  if (form.value.id_card.length !== 18) {
    alert('请输入正确的18位身份证号')
    return
  }
  submitting.value = true
  try {
    await request.post('/patients', form.value)
    alert('添加成功')
    form.value = { name: '', id_card: '', phone: '' }
    showForm.value = false
    fetchPatients() // 重新获取列表
  } catch (error: any) {
    alert(error.response?.data?.error || '添加失败')
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  fetchPatients()
})
</script>

<style scoped>
.navbar { background-color: #409eff; color: white; padding: 15px 30px; }
.nav-left { display: flex; align-items: center; }
.back-link { color: white; text-decoration: none; font-weight: bold; margin-right: 20px; }
.back-link:hover { text-decoration: underline; }
.content { max-width: 800px; margin: 30px auto; padding: 0 20px; }
.section-title { display: flex; justify-content: space-between; align-items: center; border-bottom: 1px solid #ebeef5; padding-bottom: 10px; margin-bottom: 20px; }
.add-btn { background-color: #67c23a; color: white; border: none; padding: 8px 15px; border-radius: 4px; cursor: pointer; }
.patient-grid { display: grid; gap: 15px; grid-template-columns: repeat(auto-fill, minmax(250px, 1fr)); margin-bottom: 30px; }
.patient-card { background: white; padding: 15px; border-radius: 8px; box-shadow: 0 2px 8px rgba(0,0,0,0.1); border-left: 4px solid #409eff; }
.info-row { margin-bottom: 8px; font-size: 14px; color: #606266; }
.add-form-box { background: #fdfdfd; padding: 20px; border-radius: 8px; border: 1px solid #e4e7ed; box-shadow: 0 2px 12px 0 rgba(0,0,0,0.05); }
.input-group { margin-bottom: 15px; }
.input-group label { display: block; margin-bottom: 5px; color: #606266; font-size: 14px; }
.input-group input { width: 100%; padding: 10px; border: 1px solid #dcdfe6; border-radius: 4px; box-sizing: border-box; }
.submit-btn { background-color: #409eff; color: white; border: none; padding: 10px 20px; border-radius: 4px; cursor: pointer; width: 100%; }
.submit-btn:disabled { background-color: #a0cfff; }
.loading, .empty { text-align: center; color: #909399; padding: 40px; }
</style>