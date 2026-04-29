package controllers

import (
	"log"
	"net/http"

	"hospital-api/src/models"
	"hospital-api/src/utils"

	"github.com/gin-gonic/gin"
)

// RegisterRequest 定义注册请求载荷
type RegisterRequest struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginRequest 定义登录请求载荷
type LoginRequest struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Register 处理用户注册
func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数验证失败: " + err.Error()})
		return
	}

	// 检查手机号是否已被注册
	var existingUser models.User
	if err := models.DB.Where("phone = ?", req.Phone).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "该手机号已注册"})
		return
	}

	// 对密码进行哈希处理
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "加密密码失败"})
		return
	}

	// 保存用户到数据库
	newUser := models.User{
		Phone:        req.Phone,
		PasswordHash: hashedPassword,
		Role:         "USER",
	}

	if err := models.DB.Create(&newUser).Error; err != nil {
		log.Printf("Failed to create user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "注册失败，请稍后再试"})
		return
	}

	// 注册时可顺便下发Token，或者让其重新登录，取决于业务
	c.JSON(http.StatusOK, gin.H{
		"message": "注册成功",
		"user": gin.H{
			"id":    newUser.ID,
			"phone": newUser.Phone,
			"role":  newUser.Role,
		},
	})
}

// Login 处理用户登录
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
		return
	}

	// 查找用户
	var user models.User
	if err := models.DB.Where("phone = ?", req.Phone).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 验证密码
	if match := utils.CheckPasswordHash(req.Password, user.PasswordHash); !match {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 密码正确，生成 JWT Token
	token, err := utils.GenerateToken(user.ID, user.Role)
	if err != nil {
		log.Printf("Failed to generate token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成Token失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "登录成功",
		"token":   token,
		"user": gin.H{
			"id":    user.ID,
			"phone": user.Phone,
			"role":  user.Role,
		},
	})
}
