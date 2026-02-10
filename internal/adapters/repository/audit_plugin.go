package repository

import (
	"bulls-lab-be/internal/core/domain"
	"encoding/json"
	"fmt"
	"reflect"

	"gorm.io/gorm"
)

func RegisterAuditCallbacks(db *gorm.DB) {
	// 1. Create callback
	db.Callback().Create().After("gorm:create").Register("audit:create", func(tx *gorm.DB) {
		if tx.Statement.Table == "audit_logs" {
			return
		}
		saveAudit(tx, "CREATE", nil)
	})

	// 2. Update callback
	db.Callback().Update().After("gorm:update").Register("audit:update", func(tx *gorm.DB) {
		if tx.Statement.Table == "audit_logs" {
			return
		}
		saveAudit(tx, "UPDATE", nil)
	})

	// 3. Delete callback
	db.Callback().Delete().After("gorm:delete").Register("audit:delete", func(tx *gorm.DB) {
		if tx.Statement.Table == "audit_logs" {
			return
		}
		saveAudit(tx, "DELETE", nil)
	})
}

func saveAudit(tx *gorm.DB, action string, oldValues interface{}) {
	newValuesJSON, _ := json.Marshal(tx.Statement.Dest)
	oldValuesJSON := "{}"
	if oldValues != nil {
		ov, _ := json.Marshal(oldValues)
		oldValuesJSON = string(ov)
	}

	// Try to get RecordID
	var recordID int
	if tx.Statement.ReflectValue.Kind() == reflect.Struct {
		if idField := tx.Statement.ReflectValue.FieldByName("ID"); idField.IsValid() {
			recordID = int(idField.Int())
		}
	}

	audit := domain.AuditLog{
		Action:    action,
		TableName: tx.Statement.Table,
		RecordID:  recordID,
		OldValues: oldValuesJSON,
		NewValues: string(newValuesJSON),
	}

	// We use a new session to avoid recursion or transaction issues
	tx.Session(&gorm.Session{NewDB: true}).Create(&audit)
	fmt.Printf("🛡️ Audit Log: %s on %s (ID: %d)\n", action, tx.Statement.Table, recordID)
}
