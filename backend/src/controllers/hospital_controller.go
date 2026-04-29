package controllers

import (
	"net/http"
	"time"

	"hospital-api/src/models"

	"github.com/gin-gonic/gin"
)

// GetDepartments 获取所有科室列表
func GetDepartments(c *gin.Context) {
	var departments []models.Department
	if err := models.DB.Find(&departments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询科室失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "获取成功",
		"departments": departments,
	})
}

// GetDoctorsByDept 根据科室ID获取医生列表
func GetDoctorsByDept(c *gin.Context) {
	deptID := c.Param("id")

	var doctors []models.Doctor
	if err := models.DB.Where("dept_id = ?", deptID).Find(&doctors).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询医生失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "获取成功",
		"doctors": doctors,
	})
}

// GetDoctorSchedules 获取某医生的排班
func GetDoctorSchedules(c *gin.Context) {
	doctorID := c.Param("id")

	today := time.Now().Format("2006-01-02")
	sevenDaysLater := time.Now().AddDate(0, 0, 7).Format("2006-01-02")

	var schedules []models.Schedule
	err := models.DB.Where("doctor_id = ? AND date BETWEEN ? AND ?", doctorID, today, sevenDaysLater).
		Order("date ASC, session ASC").
		Find(&schedules).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询排班失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "获取成功",
		"schedules": schedules,
	})
}

// 以下为管理员API ------------------------------

type AdminAddDoctorRequest struct {
	DeptID      uint   `json:"dept_id" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Title       string `json:"title" binding:"required"`
	Specialty   string `json:"specialty"`
	Description string `json:"description"`
}

func AddDoctor(c *gin.Context) {
	var req AdminAddDoctorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误", "details": err.Error()})
		return
	}

	doc := models.Doctor{
		DeptID:      req.DeptID,
		Name:        req.Name,
		Title:       req.Title,
		Specialty:   req.Specialty,
		Description: req.Description,
	}
	if err := models.DB.Create(&doc).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "添加医生失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "添加成功", "doctor": doc})
}

func GetAllDoctors(c *gin.Context) {
	var docs []models.Doctor
	if err := models.DB.Find(&docs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询医生失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "获取成功", "doctors": docs})
}

type AdminAddScheduleRequest struct {
	DoctorID   uint   `json:"doctor_id" binding:"required"`
	Date       string `json:"date" binding:"required"`
	Session    string `json:"session" binding:"required"`
	TotalSlots int    `json:"total_slots" binding:"required"`
}

func AddSchedule(c *gin.Context) {
	var req AdminAddScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误", "details": err.Error()})
		return
	}

	// 检查是否已存在当天同时段的排班
	var existing models.Schedule
	err := models.DB.Where("doctor_id = ? AND date = ? AND session = ?", req.DoctorID, req.Date, req.Session).First(&existing).Error
	if err == nil {
		// 已存在，可以更新
		existing.TotalSlots = req.TotalSlots
		existing.AvailableSlots = req.TotalSlots // 简单起见，重置
		models.DB.Save(&existing)
		c.JSON(http.StatusOK, gin.H{"message": "更新排班成功", "schedule": existing})
		return
	}

	schedule := models.Schedule{
		DoctorID:       req.DoctorID,
		Date:           req.Date,
		Session:        req.Session,
		TotalSlots:     req.TotalSlots,
		AvailableSlots: req.TotalSlots,
		Status:         "AVAILABLE",
	}

	if err := models.DB.Create(&schedule).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "添加排班失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "添加排班成功", "schedule": schedule})
}
