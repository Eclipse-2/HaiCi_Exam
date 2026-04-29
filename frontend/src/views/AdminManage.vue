<template>
  <div class="admin-manage">
    <h2>医生与排班管理 (管理员)</h2>

    <div class="admin-section">
      <h3>1. 添加医生</h3>
      <form @submit.prevent="handleAddDoctor" class="admin-form">
        <div class="form-group">
          <label>所属科室</label>
          <select v-model="doctorForm.dept_id" required>
            <option v-for="dept in departments" :key="dept.id" :value="dept.id">
              {{ dept.name }}
            </option>
          </select>
        </div>
        <div class="form-group">
          <label>医生姓名</label>
          <input v-model="doctorForm.name" required placeholder="如：华佗" />
        </div>
        <div class="form-group">
          <label>医生职称</label>
          <input v-model="doctorForm.title" required placeholder="如：主任医师" />
        </div>
        <div class="form-group">
          <label>专业特长</label>
          <input v-model="doctorForm.specialty" placeholder="如：心脏外科手术" />
        </div>
        <div class="form-group">
          <label>医生简介</label>
          <textarea v-model="doctorForm.description" placeholder="个人简介..."></textarea>
        </div>
        <button type="submit" class="admin-btn">添加医生</button>
      </form>
    </div>

    <div class="admin-section">
      <h3>2. 发布/更新排班号源</h3>
      <form @submit.prevent="handleAddSchedule" class="admin-form">
        <div class="form-group">
          <label>选择医生</label>
          <select v-model="scheduleForm.doctor_id" required>
            <option v-for="doc in allDoctors" :key="doc.id" :value="doc.id">
              {{ doc.name }} ({{ doc.title }})
            </option>
          </select>
        </div>
        <div class="form-group">
          <label>排班日期</label>
          <input type="date" v-model="scheduleForm.date" required :min="minDate" />
        </div>
        <div class="form-group">
          <label>班次 (上午/下午)</label>
          <select v-model="scheduleForm.session" required>
            <option value="MORNING">上午</option>
            <option value="AFTERNOON">下午</option>
          </select>
        </div>
        <div class="form-group">
          <label>号源数量</label>
          <input type="number" v-model.number="scheduleForm.total_slots" required min="1" max="200" />
        </div>
        <button type="submit" class="admin-btn">提交排班</button>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import request from '../utils/request';

// 数据定义
const departments = ref<any[]>([]);
const allDoctors = ref<any[]>([]);

const doctorForm = ref({
  dept_id: '',
  name: '',
  title: '',
  specialty: '',
  description: ''
});

const scheduleForm = ref({
  doctor_id: '',
  date: '',
  session: 'MORNING',
  total_slots: 20
});

// 计算最小允许选择的日期（今天或明天，视业务而定，我们用今天）
const minDate = computed(() => {
  const today = new Date();
  return today.toISOString().split('T')[0];
});

// 加载基础数据
const fetchDepartments = async () => {
  try {
    const res = await request.get('/departments');
    if (res.departments) {
      departments.value = res.departments;
    }
  } catch (error) {
    console.error('获取科室失败', error);
  }
};

const fetchAllDoctors = async () => {
  try {
    const res = await request.get('/admin/doctors');
    if (res.doctors) {
      allDoctors.value = res.doctors;
    }
  } catch (error) {
    console.error('获取医生列表失败', error);
  }
};

onMounted(() => {
  fetchDepartments();
  fetchAllDoctors();
});

const handleAddDoctor = async () => {
  try {
    const res = await request.post('/admin/doctors', doctorForm.value);
    alert('添加医生成功！');
    doctorForm.value = { dept_id: '', name: '', title: '', specialty: '', description: '' };
    fetchAllDoctors(); // 刷新医生列表
  } catch (error: any) {
    alert('添加医生失败: ' + error.response?.data?.error || error.message);
  }
};

const handleAddSchedule = async () => {
  try {
    const res = await request.post('/admin/schedules', scheduleForm.value);
    alert('排班发布成功！');
    scheduleForm.value.total_slots = 20; // 仅重置号源数量方便继续选
  } catch (error: any) {
    alert('排班发布失败: ' + error.response?.data?.error || error.message);
  }
};
</script>

<style scoped>
.admin-manage {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
}
.admin-section {
  background: white;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
  margin-bottom: 30px;
}
h3 {
  margin-bottom: 20px;
  border-bottom: 1px solid #eee;
  padding-bottom: 10px;
  color: #333;
}
.form-group {
  margin-bottom: 15px;
}
.form-group label {
  display: block;
  margin-bottom: 5px;
  font-weight: bold;
}
.form-group input,
.form-group select,
.form-group textarea {
  width: 100%;
  padding: 10px;
  border: 1px solid #ccc;
  border-radius: 4px;
}
.form-group textarea {
  resize: vertical;
  min-height: 80px;
}
.admin-btn {
  background-color: #2c3e50;
  color: white;
  padding: 10px 20px;
  border: none;
  cursor: pointer;
  border-radius: 4px;
  font-size: 16px;
  font-weight: bold;
}
.admin-btn:hover {
  background-color: #1a252f;
}
</style>
