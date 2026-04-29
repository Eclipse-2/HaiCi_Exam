package models

import "time"

// Appointment 预约记录模型
type Appointment struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	AppointmentNo string    `json:"appointment_no" gorm:"column:appointment_no;unique;not null"`
	UserID        uint      `json:"user_id" gorm:"column:user_id;not null"`
	PatientID     uint      `json:"patient_id" gorm:"column:patient_id;not null"`
	ScheduleID    uint      `json:"schedule_id" gorm:"column:schedule_id;not null"`
	Status        string    `json:"status" gorm:"column:status;default:'PENDING'"` // PENDING, FINISHED, CANCELLED
	CreatedAt     time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"column:updated_at"`

	// 便于后续联合查询时获取级联对象的关联结构
	Patient  *Patient  `json:"patient,omitempty" gorm:"foreignKey:PatientID"`
	Schedule *Schedule `json:"schedule,omitempty" gorm:"foreignKey:ScheduleID"`
}

func (Appointment) TableName() string {
	return "appointments"
}
