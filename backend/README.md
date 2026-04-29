# 后端启动说明

后端基于 Golang / Gin / GORM 构建，用于处理在线挂号业务逻辑。

## 环境依赖
- Go 1.20 或以上
- MySQL 5.7 或 8.0

## 运行步骤
1. 根据 \`../环境初始化脚本/init.sql\` 初始化数据库，并保证包含表结构与演示数据。
2. 配置好数据库连接信息 (暂可写在环境变量或 \`main.go\` 内)
3. 安装依赖：\`go mod tidy\`
4. 运行服务：\`go run src/main.go\`

## 项目架构
- \`src/main.go\` - 入口与路由
- 业务处理逻辑应可划分至 routers, controllers, services, models 中。
