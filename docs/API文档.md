# API接口文档 (RESTful)

## 1. 用户模块
- \`POST /api/v1/auth/register\` - 用户注册
- \`POST /api/v1/auth/login\` - 用户登录 (返回 JWT Token)

## 2. 就诊人模块
- \`GET /api/v1/patients\` - 获取当前账号下所有就诊人
- \`POST /api/v1/patients\` - 添加就诊人

## 3. 科室与医生模块
- \`GET /api/v1/departments\` - 获取所有科室
- \`GET /api/v1/departments/:id/doctors\` - 获取特定科室医生列表
- \`GET /api/v1/doctors/:id/schedules\` - 获取医生未来7天的排班

## 4. 预约挂号流程
- \`POST /api/v1/appointments\` - 提交预约挂号
  - 传参: \`schedule_id\`, \`patient_id\`
  - 业务验证: 限额、同一天同一科室限1次、并发安全性。
- \`GET /api/v1/appointments\` - 查询我的预约记录
- \`PUT /api/v1/appointments/:id/cancel\` - 取消预约
