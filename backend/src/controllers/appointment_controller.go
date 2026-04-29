package controllers

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"hospital-api/src/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type BookAppointmentRequest struct {
	ScheduleID uint `json:"schedule_id" binding:"required"`
	PatientID  uint `json:"patient_id" binding:"required"`
}

// BookAppointment 提交预约请求
// 会使用事务（Transaction）确保号源扣减操作原子性，防超卖。
// 并进行业务规则判断（同一就诊人+同科室+同一天只能预约1次）。
func BookAppointment(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	var req BookAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	// 1. 验证就诊人是否属于当前用户
	var patient models.Patient
	if err := models.DB.Where("id = ? AND user_id = ?", req.PatientID, userID).First(&patient).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "并非该账号下的有效就诊人"})
		return
	}

	// 开启事务 (GORM)
	tx := models.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 2. 读取要预约的那条排班，并联合查询出它背后的医生和科室信息
	var sched models.Schedule
	if err := tx.First(&sched, req.ScheduleID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "排班不存在"})
		return
	}

	if sched.Status != "AVAILABLE" || sched.AvailableSlots <= 0 {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "号源已派完或已停诊"})
		return
	}

	// 获取此排班对应的医生所在科室
	var doc models.Doctor
	tx.First(&doc, sched.DoctorID)

	// 3. 业务规则校验：同一天同科室同就诊人只限一次 (且状态不能是已取消)
	var duplicateCount int64
	tx.Table("appointments").
		Joins("JOIN schedules ON appointments.schedule_id = schedules.id").
		Joins("JOIN doctors ON schedules.doctor_id = doctors.id").
		Where("appointments.patient_id = ? AND appointments.status != 'CANCELLED' AND schedules.date = ? AND doctors.dept_id = ?",
			req.PatientID, sched.Date, doc.DeptID).
		Count(&duplicateCount)

	if duplicateCount > 0 {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "系统检测到该就诊人在同一科室同一天内已有预约，不可重复挂号"})
		return
	}

	// 4. 扣减号源 (采用乐观检查，在排版中 WHERE 剩余可预约 > 0 才 Update)
	res := tx.Model(&models.Schedule{}).
		Where("id = ? AND available_slots > 0", req.ScheduleID).
		UpdateColumn("available_slots", gorm.Expr("available_slots - 1"))

	if res.RowsAffected == 0 {
		// Update failed，说明在并发下刚刚被别人抢空了
		tx.Rollback()
		c.JSON(http.StatusConflict, gin.H{"error": "十分抱歉，手慢了一步，号源已空"})
		return
	}

	// 5. 生成预约号: AP + 时间戳 + 三位随机数
	appointmentNo := fmt.Sprintf("AP%s%03d", time.Now().Format("060102150405"), rand.Intn(1000))

	// 6. 创建预约记录
	appt := models.Appointment{
		AppointmentNo: appointmentNo,
		UserID:        userID.(uint),
		PatientID:     req.PatientID,
		ScheduleID:    req.ScheduleID,
		Status:        "PENDING",
	}

	if err := tx.Create(&appt).Error; err != nil {
		tx.Rollback()
		log.Printf("预约创建失败 %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误，预约失败"})
		return
	}

	// 🎉 一切正常，提交完整事务
	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"message":     "预约成功",
		"appointment": appt,
	})
}

// GetMyAppointments 查看当前用户的所有挂号记录
func GetMyAppointments(c *gin.Context) {
	userID, _ := c.Get("userID")

	var appts []models.Appointment
	// 联表预加载（Preload）让最后查出来的数据顺带携带排期与就诊人信息
	err := models.DB.Preload("Patient").Preload("Schedule").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&appts).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询预约记录失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "获取成功",
		"appointments": appts,
	})
}

// CancelAppointment 取消挂号 (30分钟免费退号源限制)
func CancelAppointment(c *gin.Context) {
	userID, _ := c.Get("userID")
	apptID := c.Param("id")

	tx := models.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 查找该预约必须是自己的且为PENDING状态
	var appt models.Appointment
	if err := tx.Where("id = ? AND user_id = ? AND status = ?", apptID, userID, "PENDING").First(&appt).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "未找到待就诊的合法预约"})
		return
	}

	// 判断是否超过 30 分钟
	// 如果超过限制不可免费取消（根据医院规则，这里演示直接阻断）
	if time.Since(appt.CreatedAt) > 30*time.Minute {
		tx.Rollback()
		c.JSON(http.StatusForbidden, gin.H{"error": "无法取消：预约成功超过30分钟，号源已被锁定。"})
		return
	}

	// 可取消 -> 恢复排版的可用号源
	if err := tx.Model(&models.Schedule{}).
		Where("id = ?", appt.ScheduleID).
		UpdateColumn("available_slots", gorm.Expr("available_slots + 1")).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "恢复号源失败"})
		return
	}

	// 更改状态为已取消
	if err := tx.Model(&appt).Update("status", "CANCELLED").Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新预约状态失败"})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"message": "非常痛心：您的预约已成功取消",
	})
}
