package models

import "time"

// Patient 就诊人信息模型
type Patient struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"column:user_id;not null"` // 关联的账号ID
	Name      string    `json:"name" gorm:"column:name;not null"`
	IDCard    string    `json:"id_card" gorm:"column:id_card;unique;not null"` // 身份证号
	Phone     string    `json:"phone" gorm:"column:phone;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
}

// TableName 指定数据表名
func (Patient) TableName() string {
	return "patients"
}
