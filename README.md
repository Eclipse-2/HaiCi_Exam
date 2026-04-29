# 🏥 简易在线挂号系统 (Hospital Registration System)

本项目为广州海鹚网络科技有限公司软件开发工程师岗位笔试题目，详细的笔试要求请见文件：软件开发工程师笔试题.docx

基于 **Vue 3 + TypeScript + Golang (Gin) + MySQL** 构建的前后端分离在线挂号平台。包含患者家属就诊预约、医生科室浏览以及管理员排班等核心挂号与院内资源管理功能。

## 🛠️ 技术栈

### 前端 (Frontend)
- **框架:** Vue 3 (Composition API) + TypeScript
- **构建工具:** Vite
- **状态管理:** Pinia
- **路由控制:** Vue Router
- **网络请求:** Axios

### 后端 (Backend)
- **语言:** Golang
- **Web 框架:** Gin
- **ORM:** GORM
- **数据库:** MySQL 8.0 (Docker 部署)
- **认证:** JWT (JSON Web Token) + Bcrypt 密码加密

---

## ✨ 核心功能

### 👤 普通用户端 (患者)
- 认证与安全：用户注册、登录、退出登录、JWT 会话保持。
- 就诊人管理：患者可以为自己或家人添加“就诊人”信息（实名制/身份证认证）。
- 科室与医生浏览：展示医院各科室以及该科室下的医生列表。
- 挂号预约：查看医生近期排班（上午/下午），确认号源充足后发起预约，支持高并发下的防超卖保护（乐观锁）。
- 订单管理：患者可以查看自己的挂号记录和预约详情。

### 👨‍💻 管理员端 (院方)
- 权限隔离：普通用户无法进入排班系统，严格的角色准入控制。
- 资源管理：可添加不同科室下的新医生资源。
- 排班发布：为指定医生发布未来日期的门诊排班，并设置每日上午/下午的具体号源池数量。

---

## 📂 项目结构

```text
HC_Code/
├── backend/                  # Golang 后端代码目录
│   ├── src/
│   │   ├── controllers/      # 控制器层 (API 逻辑)
│   │   ├── models/           # 数据库模型层与 GORM 映射
│   │   ├── routers/          # 路由注册与分发
│   │   ├── middlewares/      # 中间件 (JWT 拦截器等)
│   │   ├── utils/            # 工具类 (加密、Token生成等)
│   │   └── main.go           # 后端服务启动入口
│   ├── 数据库初始化脚本/     # SQL 迁移脚本 (init.sql)
│   ├── go.mod / go.sum       # Go 依赖管理
│   └── docker-compose.yml    # Docker MySQL 容器编排
└── frontend/                 # Vue3 前端代码目录
    ├── src/
    │   ├── views/            # 页面视图组件 (Home, Login, Admin 等)
    │   ├── stores/           # Pinia 状态管理库
    │   ├── router/           # 前端路由配置与路由守卫
    │   └── utils/            # 请求拦截器定义
    ├── package.json          # Node 依赖管理
    └── vite.config.ts        # Vite 构建配置
```

---

## 🚀 快速启动

### ⚠️ 环境要求
- [Docker Desktop](https://www.docker.com/) (用于运行 MySQL)
- [Go 1.20+](https://golang.org/) (用于编译运行后端)
- [Node.js 18+](https://nodejs.org/) (用于运行前端)

### 步骤 1：启动数据库 (MySQL)
确保 Docker 已运行，在后端目录下启动 Docker 容器。容器启动时会自动执行 `init.sql` 脚本进行完整数据库建表和数据初始化。
```bash
cd backend
docker-compose up -d
```
*(注：数据库默认用户: hospital_user，密码: hospital_password，端口: 3306)*

### 步骤 2：启动后端服务 (Golang)
开启一个新的终端窗口。
```bash
cd backend
# 解决可能的依赖包下载
go mod tidy
# 启动 Gin 服务 (默认监听 8080 端口)
go run src/main.go
```

### 步骤 3：启动前端服务 (Vue)
开启一个新的终端窗口。
```bash
cd frontend
npm install
npm run dev
```
启动成功后，浏览器访问终端中提示的本地地址 (通常是 `http://localhost:5173`) 即可体验系统。

*(Windows 用户若遇到 npm 执行策略报错，需先用管理员身份打开 PowerShell 并运行: `Set-ExecutionPolicy -Scope CurrentUser -ExecutionPolicy RemoteSigned`)*

---

## 🔒 账号与权限测试指南

- **普通用户测试**：可直接在前端登录页面点击“注册”，例如使用手机号 `13800138000` 注册体验全流程。
- **排班管理员测试的获取方式**：
  默认注册的用户角色为 `USER`。若需要体验 `[排班管理]` 功能，请修改数据库中对应账号的 `role` 字段为 `ADMIN`。
  
  使用 Docker 的快捷命令修改（以手机号为 13800138000 为例）：
  ```bash
  docker exec hospital_mysql mysql -uhospital_user -phospital_password hospital_registry -e "UPDATE users SET role = 'ADMIN' WHERE phone = '13800138000';"
  ```
  修改后**重新登录**该账号，首页顶部便会显示管理员专属入口！

---
