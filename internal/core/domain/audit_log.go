package domain

import "time"

type AuditLog struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Action    string    `json:"action"`     // CREATE, UPDATE, DELETE
	TableName string    `json:"table_name"` // e.g., "users"
	RecordID  int       `json:"record_id"`
	OldValues string    `json:"old_values" gorm:"type:jsonb"`
	NewValues string    `json:"new_values" gorm:"type:jsonb"`
	CreatedAt time.Time `json:"created_at"`
}
