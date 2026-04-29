package models

// Department 科室模型
type Department struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	Name        string `json:"name" gorm:"column:name;not null"`
	Description string `json:"description" gorm:"column:description"`
}

// TableName 指定数据表名
func (Department) TableName() string {
	return "departments"
}

// Doctor 医生模型
type Doctor struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	DeptID      uint   `json:"dept_id" gorm:"column:dept_id;not null"`
	Name        string `json:"name" gorm:"column:name;not null"`
	Title       string `json:"title" gorm:"column:title;not null"`
	Specialty   string `json:"specialty" gorm:"column:specialty"`
	Description string `json:"description" gorm:"column:description"`
}

func (Doctor) TableName() string {
	return "doctors"
}

// Schedule 医生排班模型
type Schedule struct {
	ID             uint   `json:"id" gorm:"primaryKey"`
	DoctorID       uint   `json:"doctor_id" gorm:"column:doctor_id;not null"`
	Date           string `json:"date" gorm:"column:date;type:date;not null"` // 简化处理，使用 string 存储 YYYY-MM-DD
	Session        string `json:"session" gorm:"column:session;not null"`     // MORNING, AFTERNOON
	TotalSlots     int    `json:"total_slots" gorm:"column:total_slots;not null"`
	AvailableSlots int    `json:"available_slots" gorm:"column:available_slots;not null"`
	Status         string `json:"status" gorm:"column:status;default:'AVAILABLE'"` // AVAILABLE, FULL, SUSPENDED
}

func (Schedule) TableName() string {
	return "schedules"
}
