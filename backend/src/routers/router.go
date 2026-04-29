package routers

import (
	"hospital-api/src/controllers"
	"hospital-api/src/middlewares"

	"github.com/gin-gonic/gin"
)

// SetupRouter 配置所有的 API 路由
func SetupRouter(r *gin.Engine) {
	api := r.Group("/api/v1")

	// 公开路由 (无需鉴权)
	// 科室、医生、排期查询
	api.GET("/departments", controllers.GetDepartments)
	api.GET("/departments/:id/doctors", controllers.GetDoctorsByDept)
	api.GET("/doctors/:id/schedules", controllers.GetDoctorSchedules)

	// 鉴权相关路由
	auth := api.Group("/auth")
	{
		auth.POST("/register", controllers.Register) // 注册
		auth.POST("/login", controllers.Login)       // 登录
	}

	// 需要鉴权的接口组
	authorized := api.Group("/")
	authorized.Use(middlewares.JWTAuth())
	{
		// 管理员模块 (简化：与需要鉴权的服务合在一起，实际中应有单独角色校验)
		admin := authorized.Group("/admin")
		{
			admin.POST("/doctors", controllers.AddDoctor)
			admin.GET("/doctors", controllers.GetAllDoctors)
			admin.POST("/schedules", controllers.AddSchedule)
		}

		// 就诊人管理模块
		patients := authorized.Group("/patients")
		{
			patients.GET("", controllers.GetPatients)
			patients.POST("", controllers.AddPatient)
		}

		// 预约挂号管理模块
		appointments := authorized.Group("/appointments")
		{
			appointments.POST("", controllers.BookAppointment)             // 提交预约
			appointments.GET("", controllers.GetMyAppointments)            // 查看预约记录
			appointments.PUT("/:id/cancel", controllers.CancelAppointment) // 取消预约
		}
	}
}
