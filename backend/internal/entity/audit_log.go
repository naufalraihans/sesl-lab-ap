package entity

import "time"

type AuditLog struct {
	ID          int       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      *int      `gorm:"index" json:"user_id"`
	NIM         string    `gorm:"type:varchar(50);index" json:"nim"`
	Nama        string    `gorm:"type:varchar(150)" json:"nama"`
	Role        string    `gorm:"type:varchar(50)" json:"role"`
	Action      string    `gorm:"type:varchar(100);not null;index" json:"action"`
	Description string    `gorm:"type:text;not null" json:"description"`
	IPAddress   string    `gorm:"type:varchar(45)" json:"ip_address"`
	UserAgent   string    `gorm:"type:text" json:"user_agent"`
	CreatedAt   time.Time `json:"created_at"`

	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (AuditLog) TableName() string { return "audit_logs" }
