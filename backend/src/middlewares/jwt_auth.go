package middlewares

import (
	"net/http"
	"strings"

	"hospital-api/src/utils"

	"github.com/gin-gonic/gin"
)

// JWTAuth 验证并解析 JWT Token
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 尝试从 Authorization Header 中获取 Token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "请求未携带Token，无权访问"})
			c.Abort()
			return
		}

		// Token 一般以 Bearer 开头: "Bearer eyJhbGci..."
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token 格式错误"})
			c.Abort()
			return
		}

		tokenStr := parts[1]
		claims, err := utils.ParseToken(tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的或已过期的Token"})
			c.Abort()
			return
		}

		// 将用户信息存入上下文，供后续业务处理使用
		c.Set("userID", claims.UserID)
		c.Set("userRole", claims.Role)

		c.Next()
	}
}
