package main

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type kvs map[string]string

var _ sql.Scanner = (*kvs)(nil)

func (k *kvs) Scan(val any) error {
	if val == nil {
		*k = make(kvs)
		return nil
	}

	var ba []byte
	switch v := val.(type) {
	case []byte:
		ba = v
	case string:
		ba = []byte(v)
	default:
		return fmt.Errorf("failed to unmarshal JSONB value: %w", val)
	}

	return json.Unmarshal(ba, k)
}

// GormDataType gorm common data type
func (kvs) GormDataType() string {
	return "json"
}

// GormDBDataType gorm db data type
func (kvs) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case "sqlite":
		return "JSON"
	case "mysql":
		return "JSON"
	case "postgres":
		return "JSONB"
	}
	return ""
}

// func (k kvs) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
// 	if len(k) == 0 {
// 		return gorm.Expr("NULL")
// 	}
//
// 	data, _ := json.Marshal(k)
//
// 	switch db.Dialector.Name() {
// 	case "mysql":
// 		if v, ok := db.Dialector.(*mysql.Dialector); ok && !strings.Contains(v.ServerVersion, "MariaDB") {
// 			return gorm.Expr("CAST(? AS JSON)", string(data))
// 		}
// 	}
//
// 	return gorm.Expr("?", string(data))
// }
