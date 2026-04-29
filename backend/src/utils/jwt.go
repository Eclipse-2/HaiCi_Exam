package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// 定义全局的 Jwt 秘钥 (实际项目中应从环境变量或配置文件读取)
var jwtKey = []byte("my_hospital_super_secret_key_12345")

// Claims 定义要在 JWT 中保存的负载数据
type Claims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken 生成具有指定UserID和Role的JWT Token
func GenerateToken(userID uint, role string) (string, error) {
	// 设置为 24 小时过期
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	return tokenString, err
}

// ParseToken 验证和解析 JWT Token
func ParseToken(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	return claims, nil
}
