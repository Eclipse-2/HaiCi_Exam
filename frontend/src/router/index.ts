import { createRouter, createWebHistory } from 'vue-router'
import Home from '../views/Home.vue'
import Login from '../views/Login.vue'
import Register from '../views/Register.vue'
import DoctorList from '../views/DoctorList.vue'
import DoctorSchedule from '../views/DoctorSchedule.vue'
import PatientList from '../views/PatientList.vue'
import MyAppointments from '../views/MyAppointments.vue'
import AdminManage from '../views/AdminManage.vue'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: Home
  },
  {
    path: '/admin',
    name: 'AdminManage',
    component: AdminManage
  },
  {
    path: '/department/:id/doctors',
    name: 'DoctorList',
    component: DoctorList
  },
  {
    path: '/doctor/:id/schedule',
    name: 'DoctorSchedule',
    component: DoctorSchedule
  },
  {
    path: '/patients',
    name: 'PatientList',
    component: PatientList
  },
  {
    path: '/my-appointments',
    name: 'MyAppointments',
    component: MyAppointments
  },
  {
    path: '/login',
    name: 'Login',
    component: Login
  },
  {
    path: '/register',
    name: 'Register',
    component: Register
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 简单前置路由守卫
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  const userInfoStr = localStorage.getItem('userInfo')
  let userInfo = null
  if (userInfoStr) {
    try {
      userInfo = JSON.parse(userInfoStr)
    } catch (e) {}
  }
  const protectedRoutes = ['/patients', '/my-appointments', '/admin']

  if (protectedRoutes.includes(to.path) && !token) {
    alert('该页面需要登录才能访问，请先登录！')
    next('/login')
  } else if (to.path === '/admin' && userInfo?.role !== 'ADMIN') {
    alert('权限不足：只有管理员可进入排班管理系统。')
    next('/')
  } else {
    next()
  }
})

export default router
