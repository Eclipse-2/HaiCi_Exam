package controllers

import (
	"log"
	"net/http"

	"hospital-api/src/models"

	"github.com/gin-gonic/gin"
)

// AddPatientRequest 添加就诊人的请求载荷
type AddPatientRequest struct {
	Name   string `json:"name" binding:"required"`
	IDCard string `json:"id_card" binding:"required,len=18"` // 简单验证十八位
	Phone  string `json:"phone" binding:"required,len=11"`   // 简单验证手机号长度
}

// GetPatients 获取当前登录用户的所有就诊人列表
func GetPatients(c *gin.Context) {
	// 从上下文中获取 JWTAuth 中间件写入的 userID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未获取到用户身份"})
		return
	}

	var patients []models.Patient
	// 只查询该 user_id 下的就诊人
	if err := models.DB.Where("user_id = ?", userID).Find(&patients).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询就诊人失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "获取成功",
		"patients": patients,
	})
}

// AddPatient 为当前账号添加一个新的就诊人
func AddPatient(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未获取到用户身份"})
		return
	}

	var req AddPatientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
		return
	}

	// 检查身份证是否已存在 (同一系统内一般不允许相同的身份证由于安全及数据原因被多次注册)
	var existingPatient models.Patient
	if err := models.DB.Where("id_card = ?", req.IDCard).First(&existingPatient).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "该身份证号已被添加过"})
		return
	}

	// 保存新就诊人
	newPatient := models.Patient{
		UserID: userID.(uint),
		Name:   req.Name,
		IDCard: req.IDCard,
		Phone:  req.Phone,
	}

	if err := models.DB.Create(&newPatient).Error; err != nil {
		log.Printf("Failed to create patient: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "添加就诊人失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "添加就诊人成功",
		"patient": newPatient,
	})
}
