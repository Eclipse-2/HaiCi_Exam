package models

import "time"

// User 与数据库 users 表映射的模型
type User struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Phone        string    `json:"phone" gorm:"unique;not null;column:phone"`
	PasswordHash string    `json:"-" gorm:"not null;column:password_hash"` // 加上 json:"-" 确保返回数据时不暴露密码
	Role         string    `json:"role" gorm:"default:'USER';column:role"`
	CreatedAt    time.Time `json:"created_at" gorm:"column:created_at"`
}

// 定义表名，GORM默认会加s，明确指定为数据库中的表名
func (User) TableName() string {
	return "users"
}
